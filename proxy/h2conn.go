package proxy

import (
	"context"
	"io"
	"net"
	"sync"
	"time"
)

// h2Conn adapts a CONNECT stream multiplexed over a shared HTTP/2 session to
// the proxy into a net.Conn. All I/O and deadlines operate on the stream's
// pipe ends only — never on the shared session conn, which other tunnels are
// multiplexed over as well.
type h2Conn struct {
	session net.Conn      // shared HTTP/2 session, used for addresses only
	up      net.Conn      // pipe end whose writes become request-body DATA frames
	down    net.Conn      // pipe end fed by the response body
	body    io.ReadCloser // response body carrying the tunneled bytes
	cancel  context.CancelFunc
	once    sync.Once
}

var _ net.Conn = &h2Conn{}

func newH2Conn(session net.Conn, up net.Conn, body io.ReadCloser, cancel context.CancelFunc) *h2Conn {
	downR, downW := net.Pipe()
	c := &h2Conn{
		session: session,
		up:      up,
		down:    downR,
		body:    body,
		cancel:  cancel,
	}
	go func() {
		io.Copy(downW, body)
		downW.Close()
	}()
	return c
}

func (c *h2Conn) Read(p []byte) (n int, err error) {
	return c.down.Read(p)
}

func (c *h2Conn) Write(p []byte) (n int, err error) {
	return c.up.Write(p)
}

func (c *h2Conn) Close() error {
	c.once.Do(func() {
		c.cancel()
		c.up.Close()
		c.body.Close()
		c.down.Close()
	})
	return nil
}

func (c *h2Conn) LocalAddr() net.Addr {
	return c.session.LocalAddr()
}

func (c *h2Conn) RemoteAddr() net.Addr {
	return c.session.RemoteAddr()
}

func (c *h2Conn) SetDeadline(t time.Time) error {
	err := c.up.SetWriteDeadline(t)
	if rerr := c.down.SetReadDeadline(t); rerr != nil {
		err = rerr
	}
	return err
}

func (c *h2Conn) SetReadDeadline(t time.Time) error {
	return c.down.SetReadDeadline(t)
}

func (c *h2Conn) SetWriteDeadline(t time.Time) error {
	return c.up.SetWriteDeadline(t)
}
