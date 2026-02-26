//go:build windows

package proxy

import (
	"context"
	"net"
	"syscall"
)

func control(ipv4, ipv6 net.IP) func(context.Context, string, string, syscall.RawConn) error {
	return func(ctx context.Context, network, address string, c syscall.RawConn) error {
		var err error
		c.Control(func(fd uintptr) {
			switch network {
			case "tcp4", "udp4":
				if ipv4 != nil {
					sa := &syscall.SockaddrInet4{}
					copy(sa.Addr[:], ipv4.To4())
					err = syscall.Bind(syscall.Handle(fd), sa)
				}
			case "tcp6", "udp6":
				if ipv6 != nil {
					sa := &syscall.SockaddrInet6{}
					copy(sa.Addr[:], ipv6.To16())
					err = syscall.Bind(syscall.Handle(fd), sa)
				}
			}
		})
		return err
	}
}
