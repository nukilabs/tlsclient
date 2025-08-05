package bandwidth

import (
	"errors"
	"net"
)

type TrackedConn struct {
	net.Conn
	tracker Tracker
}

func (c *TrackedConn) Read(p []byte) (n int, err error) {
	n, err = c.Conn.Read(p)
	c.tracker.AddReadBytes(int64(n))
	return
}

func (c *TrackedConn) Write(p []byte) (n int, err error) {
	n, err = c.Conn.Write(p)
	c.tracker.AddWriteBytes(int64(n))
	return
}

func NewTrackedConn(conn net.Conn, tracker Tracker) net.Conn {
	return &TrackedConn{
		Conn:    conn,
		tracker: tracker,
	}
}

type PacketConn interface {
	net.PacketConn
	SetReadBuffer(bytes int) error
	SetWriteBuffer(bytes int) error
}

type TrackedPacketConn struct {
	PacketConn
	tracker Tracker
}

func (c *TrackedPacketConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	n, addr, err = c.PacketConn.ReadFrom(p)
	c.tracker.AddReadBytes(int64(n))
	return
}

func (c *TrackedPacketConn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	n, err = c.PacketConn.WriteTo(p, addr)
	c.tracker.AddWriteBytes(int64(n))
	return
}

func (c *TrackedPacketConn) SetReadBuffer(bytes int) error {
	return c.PacketConn.SetReadBuffer(bytes)
}

func (c *TrackedPacketConn) SetWriteBuffer(bytes int) error {
	return c.PacketConn.SetWriteBuffer(bytes)
}

func NewTrackedPacketConn(conn net.PacketConn, tracker Tracker) (net.PacketConn, error) {
	pconn, ok := conn.(PacketConn)
	if !ok {
		return nil, errors.New("connection does not implement set read/write buffer methods")
	}
	return &TrackedPacketConn{
		PacketConn: pconn,
		tracker:    tracker,
	}, nil
}
