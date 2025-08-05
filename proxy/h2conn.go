package proxy

import (
	"io"
	"net"
)

type h2Conn struct {
	net.Conn
	in  io.WriteCloser
	out io.ReadCloser
}

var _ net.Conn = &h2Conn{}

func newH2Conn(conn net.Conn, in io.WriteCloser, out io.ReadCloser) *h2Conn {
	return &h2Conn{
		Conn: conn,
		in:   in,
		out:  out,
	}
}

func (h *h2Conn) Read(p []byte) (n int, err error) {
	return h.out.Read(p)
}

func (h *h2Conn) Write(p []byte) (n int, err error) {
	return h.in.Write(p)
}

func (h *h2Conn) Close() error {
	var retErr error = nil
	if err := h.in.Close(); err != nil {
		retErr = err
	}
	if err := h.out.Close(); err != nil {
		retErr = err
	}
	return retErr
}
