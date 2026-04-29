package tcp

import (
	"errors"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/cilium/ebpf/link"
)

// manager loads and attaches the SYN-rewriting BPF program and tracks
// per-profile fwmarks. There is one manager per process; access goes
// through the package singleton (defaultManager). Use SetInterface
// or TCP_IFACE to pick the egress interface.
type manager struct {
	iface  string
	objs   *synrewriteObjects
	link   link.Link
	closed atomic.Bool

	mu       sync.Mutex
	nextMark uint32
	byEntry  map[synrewriteTcpProfile]uint32
}

// newManager loads the BPF program and attaches it to the egress hook of
// iface. Requires CAP_NET_ADMIN + CAP_BPF (or root).
func newManager(iface string) (*manager, error) {
	objs, l, err := loadAndAttach(iface)
	if err != nil {
		return nil, err
	}
	return &manager{
		iface: iface,
		objs:  objs,
		link:  l,
	}, nil
}

// register installs p into the BPF map and returns the fwmark to set on
// connections that should use it. Registering a Profile whose value has
// already been registered returns the existing mark.
func (m *manager) register(p Profile) (uint32, error) {
	if m.closed.Load() {
		return 0, errors.New("tcp: manager is closed")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.registerProfileLocked(p)
}

// close detaches the program and releases the BPF object set. Safe to
// call multiple times. Currently only invoked indirectly from tests.
func (m *manager) close() error {
	if m.closed.Swap(true) {
		return nil
	}
	var errs []error
	if m.link != nil {
		if err := m.link.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if m.objs != nil {
		if err := m.objs.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// setMark applies SO_MARK to the raw socket. Called from MarkControl.
func setMark(c syscall.RawConn, mark uint32) error {
	var setErr error
	if cerr := c.Control(func(fd uintptr) {
		setErr = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_MARK, int(mark))
	}); cerr != nil {
		return cerr
	}
	return setErr
}
