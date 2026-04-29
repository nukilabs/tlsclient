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
// egress hook via TCX (kernel >= 6.6). It returns the loaded object set,
// the attachment link, and any error encountered.
//
// Requires CAP_NET_ADMIN + CAP_BPF (or root). On older kernels TCX returns
// ENOSYS / EOPNOTSUPP at attach time; in that case a clsact-qdisc-based
// fallback would be needed (not implemented here — modern deployments
// should run kernels with TCX support).
func loadAndAttach(iface string) (*synrewriteObjects, link.Link, error) {
	netif, err := net.InterfaceByName(iface)
	if err != nil {
		return nil, nil, fmt.Errorf("tcp: interface %q: %w", iface, err)
	}

	// On kernels < 5.11 BPF programs are charged against RLIMIT_MEMLOCK.
	// Newer kernels use a per-cgroup memcg accounting and ignore the rlimit,
	// but bumping it here keeps older kernels happy and is a no-op elsewhere.
	if err := rlimit.RemoveMemlock(); err != nil {
		return nil, nil, wrapPermErr(err, "raise RLIMIT_MEMLOCK")
	}

	objs := &synrewriteObjects{}
	if err := loadSynrewriteObjects(objs, nil); err != nil {
		return nil, nil, wrapPermErr(err, "load BPF object")
	}

	l, err := link.AttachTCX(link.TCXOptions{
		Program:   objs.SynRewrite,
		Attach:    ebpf.AttachTCXEgress,
		Interface: netif.Index,
	})
	if err != nil {
		objs.Close()
		if errors.Is(err, syscall.ENOSYS) || errors.Is(err, syscall.EOPNOTSUPP) {
			return nil, nil, fmt.Errorf("tcp: TCX egress attach not supported on this kernel (need >= 6.6) for %s: %w", iface, err)
		}
		return nil, nil, wrapPermErr(err, fmt.Sprintf("attach TCX egress on %s", iface))
	}

	return objs, l, nil
}

// wrapPermErr augments EPERM/EACCES with a hint about the required caps.
func wrapPermErr(err error, op string) error {
	if errors.Is(err, syscall.EPERM) || errors.Is(err, syscall.EACCES) || errors.Is(err, os.ErrPermission) {
		return fmt.Errorf("tcp: %s: %w (need CAP_NET_ADMIN + CAP_BPF, or root)", op, err)
	}
	return fmt.Errorf("tcp: %s: %w", op, err)
}
