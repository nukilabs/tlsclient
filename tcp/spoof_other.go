//go:build !linux

package tcp

// manager is a no-op on non-Linux. MarkControl returns nil so dials
// silently proceed without spoofing; this type is only here to keep
// cross-platform compilation working.
type manager struct{}
