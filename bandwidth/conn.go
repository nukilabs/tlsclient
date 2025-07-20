package bandwidth

import "net"

type TrackedConn struct {
	net.Conn
	tracker Tracker
}

func (t *TrackedConn) Read(p []byte) (n int, err error) {
	n, err = t.Conn.Read(p)
	t.tracker.AddReadBytes(int64(n))
	return
}

func (t *TrackedConn) Write(p []byte) (n int, err error) {
	n, err = t.Conn.Write(p)
	t.tracker.AddWriteBytes(int64(n))
	return
}

func NewTrackedConn(conn net.Conn, tracker Tracker) net.Conn {
	return &TrackedConn{
		Conn:    conn,
		tracker: tracker,
	}
}

type TrackedUDPConn struct {
	net.PacketConn
	tracker Tracker
}

func (t *TrackedUDPConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	n, addr, err = t.PacketConn.ReadFrom(p)
	t.tracker.AddReadBytes(int64(n))
	return
}

func (t *TrackedUDPConn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	n, err = t.PacketConn.WriteTo(p, addr)
	t.tracker.AddWriteBytes(int64(n))
	return
}

func NewTrackedUDPConn(udpConn net.PacketConn, tracker Tracker) net.PacketConn {
	return &TrackedUDPConn{
		PacketConn: udpConn,
		tracker:    tracker,
	}
}
