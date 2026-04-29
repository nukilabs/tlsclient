package tlsclient

import (
	"syscall"

	"github.com/nukilabs/tlsclient/tcp"
)

// WithTCP enables Linux TC eBPF SYN fingerprint spoofing using the
// process-wide singleton manager (auto-attaches on first use to the
// default-route interface, overridable via tcp.SetInterface or
// the TCP_IFACE environment variable).
//
// On non-Linux this option is a silent no-op: dials proceed with the
// host kernel's normal TCP fingerprint.
func WithTCP(p tcp.Profile) Option {
	ctrl := tcp.MarkControl(p)
	return func(c *Client) {
		c.tcpControl = ctrl
	}
}

// markControl is the Control function signature stored on the Client.
type markControl = func(network, address string, c syscall.RawConn) error
