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
	listenV4 string // local listen address for udp4 (e.g. "1.2.3.4:0")
	listenV6 string // local listen address for udp6 (e.g. "[::1]:0")
}

// Direct creates a new direct dialer with optional local IP for binding.
func Direct(ip net.IP, timeout time.Duration) *direct {
	d := &direct{
		dialer: net.Dialer{Timeout: timeout},
	}
	if ip != nil {
		if ip.To4() != nil {
			d.family = "4"
			d.listenV4 = ip.String() + ":0"
			d.dialer.Control = control(ip, nil)
		} else {
			d.family = "6"
			d.listenV6 = "[" + ip.String() + "]:0"
			d.dialer.Control = control(nil, ip)
		}
	}
	return d
}

// DirectDualStack creates a new direct dialer with both IPv4 and IPv6 local addresses.
func DirectDualStack(ipv4, ipv6 net.IP, timeout time.Duration) *direct {
	return &direct{
		dialer: net.Dialer{
			Timeout:       timeout,
			FallbackDelay: time.Second,
			Control:       control(ipv4, ipv6),
		},
		listenV4: ipv4.String() + ":0",
		listenV6: "[" + ipv6.String() + "]:0",
	}
}

// DialContext dials the address using the configured dialer.
func (d *direct) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	return d.dialer.DialContext(ctx, d.forceFamily(network), addr)
}

// ListenPacket creates a packet connection for QUIC/UDP using the configured listener.
func (d *direct) ListenPacket(ctx context.Context, network, addr string) (net.PacketConn, error) {
	network = d.forceFamily(network)
	return d.listener.ListenPacket(ctx, network, d.listenAddr(network))
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

// listenAddr returns the local address to bind the UDP listener to.
func (d *direct) listenAddr(network string) string {
	switch network {
	case "udp4":
		if d.listenV4 != "" {
			return d.listenV4
		}
	case "udp6":
		if d.listenV6 != "" {
			return d.listenV6
		}
	}
	return ":0"
}
