package tcp

import (
	"errors"
	"fmt"
	"net"
	"os"
	"syscall"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
)

// loadAndAttach opens the embedded BPF object and attaches it to iface's
// egress hook.
//
// The attach path is TCX on Linux >= 6.6, falling back to clsact qdisc +
// direct-action TC classifier (via raw rtnetlink) on older kernels.
//
// Returns the loaded object set, a detach function that cleans up the
// attachment, and any error. Requires CAP_NET_ADMIN + CAP_BPF (or root).
func loadAndAttach(iface string) (*synrewriteObjects, func() error, error) {
	netif, err := net.InterfaceByName(iface)
	if err != nil {
		return nil, nil, fmt.Errorf("tcp: interface %q: %w", iface, err)
	}

	// On kernels < 5.11 BPF programs are charged against RLIMIT_MEMLOCK.
	// Newer kernels use a per-cgroup memcg accounting and ignore the rlimit;
	// bumping it here keeps older kernels happy and is a no-op elsewhere.
	if err := rlimit.RemoveMemlock(); err != nil {
		return nil, nil, wrapPermErr(err, "raise RLIMIT_MEMLOCK")
	}

	objs := &synrewriteObjects{}
	if err := loadSynrewriteObjects(objs, nil); err != nil {
		return nil, nil, wrapPermErr(err, "load BPF object")
	}

	// Try TCX first.
	l, err := link.AttachTCX(link.TCXOptions{
		Program:   objs.SynRewrite,
		Attach:    ebpf.AttachTCXEgress,
		Interface: netif.Index,
	})
	if err == nil {
		return objs, func() error { return l.Close() }, nil
	}
	if !isTCXUnsupported(err) {
		objs.Close()
		return nil, nil, wrapPermErr(err, fmt.Sprintf("attach TCX egress on %s", iface))
	}

	// Fall back to clsact + tc filter via netlink.
	detach, err := attachClsact(netif.Index, objs.SynRewrite.FD(), "syn_rewrite")
	if err != nil {
		objs.Close()
		return nil, nil, wrapPermErr(err, fmt.Sprintf("attach clsact egress on %s", iface))
	}
	return objs, detach, nil
}

// isTCXUnsupported recognizes the various ways the kernel reports that
// TCX (BPF_PROG_TYPE_SCHED_CLS via the link API) isn't available — kernel
// < 6.6 returns ENOSYS / EOPNOTSUPP / EINVAL depending on minor version,
// and cilium/ebpf wraps unsupported probes in ebpf.ErrNotSupported.
func isTCXUnsupported(err error) bool {
	return errors.Is(err, ebpf.ErrNotSupported) ||
		errors.Is(err, syscall.ENOSYS) ||
		errors.Is(err, syscall.EOPNOTSUPP) ||
		errors.Is(err, syscall.EINVAL)
}

// wrapPermErr augments EPERM/EACCES with a hint about the required caps.
func wrapPermErr(err error, op string) error {
	if errors.Is(err, syscall.EPERM) || errors.Is(err, syscall.EACCES) || errors.Is(err, os.ErrPermission) {
		return fmt.Errorf("tcp: %s: %w (need CAP_NET_ADMIN + CAP_BPF, or root)", op, err)
	}
	return fmt.Errorf("tcp: %s: %w", op, err)
}
