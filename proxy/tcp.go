package proxy

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/nukilabs/http"
	"github.com/nukilabs/http/http2"
	tls "github.com/nukilabs/utls"
)

func (d *Dialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	proxy, dst := opAddr(d.proxyURL.Host), opAddr(addr)

	req := &http.Request{
		Method: http.MethodConnect,
		URL:    &url.URL{Host: addr},
		Header: make(http.Header),
		Host:   addr,
	}
	if d.authHeader != "" {
		req.Header.Set("Proxy-Authorization", d.authHeader)
	}

	// Reuse an existing HTTP/2 session to the proxy when possible.
	d.h2DialLock.Lock()
	h2Conn, h2ClientConn := d.h2Conn, d.h2ClientConn
	d.h2DialLock.Unlock()
	if h2ClientConn != nil && h2ClientConn.CanTakeNewRequest() {
		c, err := d.connectHttp2(ctx, req, h2Conn, h2ClientConn)
		if err != nil {
			return nil, &net.OpError{Op: "connect", Net: network, Source: proxy, Addr: dst, Err: err}
		}
		return c, nil
	}

	switch d.proxyURL.Scheme {
	case "http":
		dd := net.Dialer{Timeout: d.timeout}
		conn, err := dd.DialContext(ctx, network, d.proxyURL.Host)
		if err != nil {
			return nil, &net.OpError{Op: "connect", Net: network, Source: proxy, Addr: dst, Err: fmt.Errorf("dialing proxy failed: %w", err)}
		}
		c, err := d.connectHttp1(ctx, req, conn)
		if err != nil {
			return nil, &net.OpError{Op: "connect", Net: network, Source: proxy, Addr: dst, Err: err}
		}
		return c, nil
	case "https":
		tlsConf := d.tlsConf.Clone()
		tlsConf.ServerName = d.proxyURL.Hostname()
		tlsConf.NextProtos = []string{"http/1.1", "h2"}
		dd := net.Dialer{Timeout: d.timeout}
		rawConn, err := dd.DialContext(ctx, network, d.proxyURL.Host)
		if err != nil {
			return nil, &net.OpError{Op: "connect", Net: network, Source: proxy, Addr: dst, Err: fmt.Errorf("dialing proxy failed: %w", err)}
		}
		conn := tls.Client(rawConn, tlsConf)
		if deadline, ok := d.deadline(ctx); ok {
			if err := conn.SetDeadline(deadline); err != nil {
				conn.Close()
				return nil, &net.OpError{Op: "connect", Net: network, Source: proxy, Addr: dst, Err: fmt.Errorf("setting connection deadline failed: %w", err)}
			}
		}
		if err := conn.HandshakeContext(ctx); err != nil {
			conn.Close()
			return nil, &net.OpError{Op: "connect", Net: network, Source: proxy, Addr: dst, Err: fmt.Errorf("tls handshake failed: %w", err)}
		}
		if err := conn.SetDeadline(noDeadline); err != nil {
			conn.Close()
			return nil, &net.OpError{Op: "connect", Net: network, Source: proxy, Addr: dst, Err: fmt.Errorf("clearing connection deadline failed: %w", err)}
		}

		state := conn.ConnectionState()
		switch state.NegotiatedProtocol {
		case "http/1.1":
			c, err := d.connectHttp1(ctx, req, conn)
			if err != nil {
				return nil, &net.OpError{Op: "connect", Net: network, Source: proxy, Addr: dst, Err: err}
			}
			return c, nil
		case "h2":
			tr := &http2.Transport{}
			clientConn, err := tr.NewClientConn(conn)
			if err != nil {
				conn.Close()
				return nil, &net.OpError{Op: "connect", Net: network, Source: proxy, Addr: dst, Err: fmt.Errorf("dialing h2 client connection failed: %w", err)}
			}
			d.h2DialLock.Lock()
			d.h2Conn = conn
			d.h2ClientConn = clientConn
			d.h2DialLock.Unlock()

			c, err := d.connectHttp2(ctx, req, conn, clientConn)
			if err != nil {
				return nil, &net.OpError{Op: "connect", Net: network, Source: proxy, Addr: dst, Err: err}
			}
			return c, nil
		default:
			conn.Close()
			return nil, &net.OpError{Op: "connect", Net: network, Source: proxy, Addr: dst, Err: errors.New("unsupported protocol: " + state.NegotiatedProtocol)}
		}
	default:
		return nil, &net.OpError{Op: "connect", Net: network, Source: proxy, Addr: dst, Err: errors.New("unsupported proxy scheme: " + d.proxyURL.Scheme)}
	}
}

func (d *Dialer) connectHttp1(ctx context.Context, req *http.Request, conn net.Conn) (net.Conn, error) {
	if deadline, ok := d.deadline(ctx); ok {
		if err := conn.SetDeadline(deadline); err != nil {
			conn.Close()
			return nil, fmt.Errorf("setting connection deadline failed: %w", err)
		}
	}
	defer context.AfterFunc(ctx, func() { conn.SetDeadline(aLongTimeAgo) })()
	if err := req.Write(conn); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	res, err := http.ReadResponse(bufio.NewReader(conn), req)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		conn.Close()
		return nil, fmt.Errorf("server responded with %d", res.StatusCode)
	}
	if err := conn.SetDeadline(noDeadline); err != nil {
		conn.Close()
		return nil, fmt.Errorf("clearing connection deadline failed: %w", err)
	}
	return conn, nil
}

// connectHttp2 opens a CONNECT stream on the shared HTTP/2 session. It never
// closes the session conn itself: other tunnels may be multiplexed over it,
// and a dead session is detected via CanTakeNewRequest on the next dial.
func (d *Dialer) connectHttp2(ctx context.Context, req *http.Request, conn net.Conn, clientConn *http2.ClientConn) (net.Conn, error) {
	upR, upW := net.Pipe()
	req.Body = upR

	streamCtx, cancel := context.WithCancel(context.Background())
	defer context.AfterFunc(ctx, cancel)()
	if d.timeout > 0 {
		defer time.AfterFunc(d.timeout, cancel).Stop()
	}

	res, err := clientConn.RoundTrip(req.WithContext(streamCtx))
	if err != nil {
		cancel()
		upW.Close()
		return nil, fmt.Errorf("failed to round trip request: %w", err)
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		cancel()
		upW.Close()
		res.Body.Close()
		return nil, fmt.Errorf("server responded with %d", res.StatusCode)
	}

	return newH2Conn(conn, upW, res.Body, cancel), nil
}
