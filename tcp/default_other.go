//go:build !linux

package tcp

import "syscall"

// SetInterface is a no-op on non-Linux.
func SetInterface(iface string) {}

// MarkControl on non-Linux returns nil — the tlsclient option silently
// proceeds without SYN spoofing instead of failing every dial.
func MarkControl(p Profile) func(network, address string, c syscall.RawConn) error {
	return nil
}
