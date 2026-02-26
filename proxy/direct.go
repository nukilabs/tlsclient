package proxy

import (
	"context"
	"net"
	"time"
)

// direct implements Dialer by making network connections directly using a net.Dialer and net.ListenConfig.
type direct struct {
	dialer   net.Dialer
	listener net.ListenConfig
}

// Direct creates a new direct dialer with optional local address for binding.
func Direct(addr net.Addr, timeout time.Duration) *direct {
	return &direct{
		dialer: net.Dialer{
			Timeout:   timeout,
			LocalAddr: addr,
		},
	}
}

// DirectDualStack creates a new direct dialer with both IPv4 and IPv6 local addresses.
func DirectDualStack(ipv4, ipv6 net.IP, timeout time.Duration) *direct {
	return &direct{
		dialer: net.Dialer{
			Timeout:        timeout,
			FallbackDelay:  time.Second,
			ControlContext: control(ipv4, ipv6),
		},
	}
}

// DialContext dials the address using the configured dialer.
func (d *direct) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	return d.dialer.DialContext(ctx, network, addr)
}

// ListenPacket instantiates a net.ListenConfig and invokes its ListenPacket receiver for packet connections.
func (d *direct) ListenPacket(ctx context.Context, network, addr string) (net.PacketConn, error) {
	return d.listener.ListenPacket(ctx, network, ":0")
}

// SupportHTTP3 indicates that the direct dialer supports HTTP/3.
func (d *direct) SupportHTTP3() bool {
	return true
}
