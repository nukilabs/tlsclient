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
