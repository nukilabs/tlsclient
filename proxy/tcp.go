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
	req := &http.Request{
		Method: http.MethodConnect,
		URL:    &url.URL{Host: addr},
		Header: make(http.Header),
		Host:   addr,
	}
	req.WithContext(ctx)
	if d.authHeader != "" {
		req.Header.Set("Proxy-Authorization", d.authHeader)
	}

	d.h2DialLock.Lock()
	if d.h2ClientConn != nil {
		if d.h2ClientConn.CanTakeNewRequest() {
			d.h2DialLock.Unlock()
			return d.connectHttp2(req, d.h2Conn, d.h2ClientConn)
		}
	}
	d.h2DialLock.Unlock()

	switch d.proxyURL.Scheme {
	case "http":
		var dd net.Dialer
		conn, err := dd.DialContext(ctx, network, d.proxyURL.Host)
		if err != nil {
			return nil, fmt.Errorf("dialing proxy failed: %w", err)
		}
		return d.connectHttp1(req, conn)
	case "https":
		tlsConf := d.tlsConf.Clone()
		tlsConf.ServerName = d.proxyURL.Hostname()
		tlsConf.NextProtos = []string{"http/1.1", "h2"}
		conn, err := tls.Dial(network, d.proxyURL.Host, tlsConf)
		if err != nil {
			return nil, fmt.Errorf("dialing tls connection failed: %w", err)
		}
		if err := conn.HandshakeContext(ctx); err != nil {
			conn.Close()
			return nil, fmt.Errorf("tls handshake failed: %w", err)
		}

		state := conn.ConnectionState()
		switch state.NegotiatedProtocol {
		case "http/1.1":
			return d.connectHttp1(req, conn)
		case "h2":
			d.h2DialLock.Lock()
			defer d.h2DialLock.Unlock()

			tr := &http2.Transport{}
			clientConn, err := tr.NewClientConn(conn)
			if err != nil {
				conn.Close()
				return nil, fmt.Errorf("dialing h2 client connection failed: %w", err)
			}
			d.h2Conn = conn
			d.h2ClientConn = clientConn

			return d.connectHttp2(req, conn, clientConn)
		default:
			conn.Close()
			return nil, fmt.Errorf("unsupported protocol: %s", state.NegotiatedProtocol)
		}
	default:
		return nil, errors.New("unsupported proxy scheme: " + d.proxyURL.Scheme)
	}
}

func (d *Dialer) connectHttp1(req *http.Request, conn net.Conn) (net.Conn, error) {
	deadline := time.Now().Add(d.timeout)
	if err := conn.SetDeadline(deadline); err != nil {
		conn.Close()
		return nil, fmt.Errorf("setting connection deadline failed: %w", err)
	}
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
	if err := conn.SetDeadline(time.Time{}); err != nil {
		conn.Close()
		return nil, fmt.Errorf("clearing connection deadline failed: %w", err)
	}
	return conn, nil
}

func (d *Dialer) connectHttp2(req *http.Request, conn net.Conn, clientConn *http2.ClientConn) (net.Conn, error) {
	pr, pw := net.Pipe()
	req.Body = pr

	res, err := clientConn.RoundTrip(req)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to round trip request: %w", err)
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		conn.Close()
		return nil, fmt.Errorf("server responded with %d", res.StatusCode)
	}

	return newH2Conn(conn, pw, res.Body), nil
}
