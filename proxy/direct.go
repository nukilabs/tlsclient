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
	family   string // "4", "6", or "" (any)
}

// Direct creates a new direct dialer with optional local IP for binding.
func Direct(ip net.IP, timeout time.Duration) *direct {
	d := &direct{
		dialer: net.Dialer{Timeout: timeout},
	}
	if ip != nil {
		if ip.To4() != nil {
			d.family = "4"
			ctrl := control(ip, nil)
			d.dialer.Control = ctrl
			d.listener.Control = ctrl
		} else {
			d.family = "6"
			ctrl := control(nil, ip)
			d.dialer.Control = ctrl
			d.listener.Control = ctrl
		}
	}
	return d
}

// DirectDualStack creates a new direct dialer with both IPv4 and IPv6 local addresses.
func DirectDualStack(ipv4, ipv6 net.IP, timeout time.Duration) *direct {
	ctrl := control(ipv4, ipv6)
	return &direct{
		dialer: net.Dialer{
			Timeout:       timeout,
			FallbackDelay: time.Second,
			Control:       ctrl,
		},
		listener: net.ListenConfig{
			Control: ctrl,
		},
	}
}

// DialContext dials the address using the configured dialer.
func (d *direct) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	return d.dialer.DialContext(ctx, d.forceFamily(network), addr)
}

// ListenPacket creates a packet connection for QUIC/UDP using the configured listener.
func (d *direct) ListenPacket(ctx context.Context, network, addr string) (net.PacketConn, error) {
	return d.listener.ListenPacket(ctx, d.forceFamily(network), ":0")
}

// SupportHTTP3 indicates that the direct dialer supports HTTP/3.
func (d *direct) SupportHTTP3() bool {
	return true
}

// forceFamily overrides the network to match the bound IP family.
func (d *direct) forceFamily(network string) string {
	if d.family == "" {
		return network
	}
	switch network {
	case "tcp", "tcp4", "tcp6":
		return "tcp" + d.family
	case "udp", "udp4", "udp6":
		return "udp" + d.family
	default:
		return network
	}
}
