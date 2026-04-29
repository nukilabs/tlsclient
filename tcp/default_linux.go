package tcp

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"syscall"
)

// Process-wide singleton. The first call to defaultManager initializes
// the manager bound to either the explicit interface set via
// SetInterface, the TCP_IFACE environment variable, or
// the auto-detected default-route interface.
var (
	defaultMu       sync.Mutex
	defaultMgr      *manager
	defaultErr      error
	defaultIfaceCfg string // explicit override; "" means auto-detect
)

// SetInterface overrides the auto-detected default-route interface
// used by the singleton. Must be called before the first dial through
// tlsclient.WithTCP (which triggers the lazy load). After the singleton
// is initialized, this is a no-op.
func SetInterface(iface string) {
	defaultMu.Lock()
	defer defaultMu.Unlock()
	if defaultMgr != nil {
		return
	}
	defaultIfaceCfg = iface
}

// defaultManager returns the process-wide singleton, initializing it on
// first call. If init fails, the error is cached and returned on every
// subsequent call so the failure is observable at the dial site.
func defaultManager() (*manager, error) {
	defaultMu.Lock()
	defer defaultMu.Unlock()
	if defaultMgr != nil || defaultErr != nil {
		return defaultMgr, defaultErr
	}
	iface := defaultIfaceCfg
	if iface == "" {
		iface = os.Getenv("TCP_IFACE")
	}
	if iface == "" {
		var err error
		iface, err = detectDefaultInterface()
		if err != nil {
			defaultErr = fmt.Errorf("tcp: detect default interface: %w", err)
			return nil, defaultErr
		}
	}
	mgr, err := newManager(iface)
	if err != nil {
		defaultErr = err
		return nil, err
	}
	defaultMgr = mgr
	return mgr, nil
}

// MarkControl returns a syscall.RawConn Control function that registers p
// (idempotently) on the singleton and tags each socket with the matching
// SO_MARK. Used by tlsclient.WithTCP. Errors during init or registration
// surface as dial errors.
func MarkControl(p Profile) func(network, address string, c syscall.RawConn) error {
	return func(network, address string, c syscall.RawConn) error {
		mgr, err := defaultManager()
		if err != nil {
			return err
		}
		mark, err := mgr.register(p)
		if err != nil {
			return err
		}
		return setMark(c, mark)
	}
}

// detectDefaultInterface reads /proc/net/route and returns the interface
// owning the first 0.0.0.0 route entry.
func detectDefaultInterface() (string, error) {
	f, err := os.Open("/proc/net/route")
	if err != nil {
		return "", err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	if !s.Scan() {
		return "", errors.New("empty /proc/net/route")
	}
	for s.Scan() {
		fields := strings.Fields(s.Text())
		// Iface Destination Gateway Flags RefCnt Use Metric Mask MTU Window IRTT
		if len(fields) < 8 {
			continue
		}
		if fields[1] == "00000000" {
			return fields[0], nil
		}
	}
	if err := s.Err(); err != nil {
		return "", err
	}
	return "", errors.New("no default route in /proc/net/route")
}
