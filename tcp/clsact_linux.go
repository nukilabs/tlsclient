package tcp

// clsact_linux.go: fallback attach path for kernels older than 6.6 (TCX
// requires Linux >= 6.6). Adds the clsact qdisc on the interface (if not
// already present) and attaches the BPF program as a TC classifier in
// direct-action mode on the egress hook.
//
// Implemented as raw rtnetlink rather than pulling in vishvananda/netlink.

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

const (
	tcaKind    uint16 = 1
	tcaOptions uint16 = 2

	tcaBpfFd    uint16 = 6 // <linux/pkt_cls.h> TCA_BPF_FD
	tcaBpfName  uint16 = 7 // TCA_BPF_NAME
	tcaBpfFlags uint16 = 8 // TCA_BPF_FLAGS

	tcaBpfFlagActDirect uint32 = 1

	// TC_H_CLSACT and the egress minor handle (linux/pkt_sched.h).
	tcHClsact     uint32 = 0xFFFFFFF1
	tcHMajClsact  uint32 = 0xFFFF0000
	tcHMinEgress  uint32 = 0xFFF3
	tcQdiscHandle uint32 = 0xFFFF0000

	ethPAll uint16 = 0x0003
)

// attachClsact ensures a clsact qdisc on ifindex and attaches a direct-
// action BPF classifier (progFD) to its egress hook. Returns a detach
// function that removes the filter; the qdisc is left in place since
// other programs may share it.
func attachClsact(ifindex int, progFD int, progName string) (func() error, error) {
	if err := ensureClsactQdisc(ifindex); err != nil {
		return nil, err
	}
	if err := addBPFFilter(ifindex, progFD, progName); err != nil {
		return nil, err
	}
	return func() error { return removeBPFFilter(ifindex) }, nil
}

// ----- netlink message construction --------------------------------------

// tcmsgBytes builds the 20-byte struct tcmsg payload.
func tcmsgBytes(family uint8, ifindex int32, handle, parent, info uint32) []byte {
	b := make([]byte, 20)
	b[0] = family
	binary.LittleEndian.PutUint32(b[4:8], uint32(ifindex))
	binary.LittleEndian.PutUint32(b[8:12], handle)
	binary.LittleEndian.PutUint32(b[12:16], parent)
	binary.LittleEndian.PutUint32(b[16:20], info)
	return b
}

// nlAttr encodes a netlink attribute with 4-byte alignment.
func nlAttr(typ uint16, value []byte) []byte {
	length := 4 + len(value)
	aligned := (length + 3) &^ 3
	out := make([]byte, aligned)
	binary.LittleEndian.PutUint16(out[0:2], uint16(length))
	binary.LittleEndian.PutUint16(out[2:4], typ)
	copy(out[4:], value)
	return out
}

func nlAttrU32(typ uint16, val uint32) []byte {
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], val)
	return nlAttr(typ, b[:])
}

func nlAttrCString(typ uint16, val string) []byte {
	b := make([]byte, len(val)+1)
	copy(b, val)
	return nlAttr(typ, b)
}

// ----- netlink transport -------------------------------------------------

// nlSendRecv opens a NETLINK_ROUTE socket, sends one request, and waits
// for the matching ACK or NLMSG_ERROR. Returns the syscall error reported
// by the kernel (or nil for a clean ACK).
func nlSendRecv(typ uint16, flags uint16, body []byte) error {
	fd, err := unix.Socket(unix.AF_NETLINK, unix.SOCK_RAW|unix.SOCK_CLOEXEC, unix.NETLINK_ROUTE)
	if err != nil {
		return fmt.Errorf("netlink socket: %w", err)
	}
	defer unix.Close(fd)
	if err := unix.Bind(fd, &unix.SockaddrNetlink{Family: unix.AF_NETLINK}); err != nil {
		return fmt.Errorf("netlink bind: %w", err)
	}

	seq := uint32(os.Getpid()) ^ 0xC0DE
	hdr := make([]byte, 16)
	total := uint32(16 + len(body))
	binary.LittleEndian.PutUint32(hdr[0:4], total)
	binary.LittleEndian.PutUint16(hdr[4:6], typ)
	binary.LittleEndian.PutUint16(hdr[6:8], flags|unix.NLM_F_REQUEST|unix.NLM_F_ACK)
	binary.LittleEndian.PutUint32(hdr[8:12], seq)
	// pid 0 = let kernel use ours

	msg := append(hdr, body...)
	if err := unix.Sendto(fd, msg, 0, &unix.SockaddrNetlink{Family: unix.AF_NETLINK}); err != nil {
		return fmt.Errorf("netlink sendto: %w", err)
	}

	buf := make([]byte, 65536)
	n, err := unix.Read(fd, buf)
	if err != nil {
		return fmt.Errorf("netlink read: %w", err)
	}
	if n < 16 {
		return errors.New("netlink: short response")
	}
	msgType := binary.LittleEndian.Uint16(buf[4:6])
	if msgType != unix.NLMSG_ERROR {
		return fmt.Errorf("netlink: unexpected response type %d", msgType)
	}
	if n < 20 {
		return errors.New("netlink: truncated NLMSG_ERROR")
	}
	errCode := int32(binary.LittleEndian.Uint32(buf[16:20]))
	if errCode == 0 {
		return nil
	}
	return unix.Errno(-errCode)
}

// ----- qdisc + filter operations -----------------------------------------

func ensureClsactQdisc(ifindex int) error {
	body := append(
		tcmsgBytes(unix.AF_UNSPEC, int32(ifindex), tcQdiscHandle, tcHClsact, 0),
		nlAttrCString(tcaKind, "clsact")...,
	)
	err := nlSendRecv(unix.RTM_NEWQDISC,
		unix.NLM_F_CREATE|unix.NLM_F_EXCL, body)
	if err == nil || errors.Is(err, unix.EEXIST) {
		return nil
	}
	return fmt.Errorf("tcp: add clsact qdisc on ifindex %d: %w", ifindex, err)
}

func addBPFFilter(ifindex int, progFD int, progName string) error {
	parent := tcHMajClsact | tcHMinEgress           // 0xFFFFFFF3
	info := uint32(1)<<16 | uint32(ethPAll)<<8 // prio=1 << 16 | htons(ETH_P_ALL)
	// Note: ethPAll<<8 produces 0x0300 which is htons(0x0003).

	bpfOpts := nlAttrU32(tcaBpfFd, uint32(progFD))
	bpfOpts = append(bpfOpts, nlAttrCString(tcaBpfName, progName)...)
	bpfOpts = append(bpfOpts, nlAttrU32(tcaBpfFlags, tcaBpfFlagActDirect)...)

	body := tcmsgBytes(unix.AF_UNSPEC, int32(ifindex), 0, parent, info)
	body = append(body, nlAttrCString(tcaKind, "bpf")...)
	body = append(body, nlAttr(tcaOptions, bpfOpts)...)

	if err := nlSendRecv(unix.RTM_NEWTFILTER, unix.NLM_F_CREATE, body); err != nil {
		return fmt.Errorf("tcp: add bpf tc filter on ifindex %d: %w", ifindex, err)
	}
	return nil
}

func removeBPFFilter(ifindex int) error {
	parent := tcHMajClsact | tcHMinEgress
	info := uint32(1)<<16 | uint32(ethPAll)<<8

	body := tcmsgBytes(unix.AF_UNSPEC, int32(ifindex), 0, parent, info)
	body = append(body, nlAttrCString(tcaKind, "bpf")...)

	if err := nlSendRecv(unix.RTM_DELTFILTER, 0, body); err != nil {
		return fmt.Errorf("tcp: remove bpf tc filter on ifindex %d: %w", ifindex, err)
	}
	return nil
}
