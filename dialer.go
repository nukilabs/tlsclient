package tlsclient

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"
	"sync"
	"time"

	http "github.com/sparkaio/fhttp"
	"github.com/sparkaio/fhttp/http2"
	"golang.org/x/net/proxy"
)

type ConnectDialer struct {
	ProxyURL      *url.URL
	Dialer        net.Dialer
	DefaultHeader http.Header
	Timeout       time.Duration

	cacheH2Lock       sync.Mutex
	cacheH2ClientConn *http2.ClientConn
	cacheH2RawConn    net.Conn
}

func NewConnectDialer(proxyUrl *url.URL, timeout time.Duration) (*ConnectDialer, error) {
	switch proxyUrl.Scheme {
	case "http":
		if proxyUrl.Port() == "" {
			proxyUrl.Host = net.JoinHostPort(proxyUrl.Hostname(), "80")
		}
	case "https":
		if proxyUrl.Port() == "" {
			proxyUrl.Host = net.JoinHostPort(proxyUrl.Hostname(), "443")
		}
	case "socks5", "socks5h":
		if proxyUrl.Port() == "" {
			proxyUrl.Host = net.JoinHostPort(proxyUrl.Hostname(), "1080")
		}
	}

	dialer := &ConnectDialer{
		ProxyURL:      proxyUrl,
		DefaultHeader: make(http.Header),
		Timeout:       timeout,
	}

	if proxyUrl.User != nil {
		if username := proxyUrl.User.Username(); username != "" {
			password, _ := proxyUrl.User.Password()
			dialer.DefaultHeader.Set("Proxy-Authorization", fmt.Sprintf("Basic %s",
				base64.StdEncoding.EncodeToString([]byte(username+":"+password))))
		}
	}

	return dialer, nil
}

func (d *ConnectDialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	if strings.HasPrefix(d.ProxyURL.Scheme, "socks") {
		dial, err := proxy.FromURL(d.ProxyURL, proxy.Direct)
		if err != nil {
			return nil, err
		}
		return dial.(proxy.ContextDialer).DialContext(ctx, network, addr)
	}

	req := &http.Request{
		Method: "CONNECT",
		URL:    &url.URL{Host: addr},
		Header: make(http.Header),
		Host:   addr,
	}
	req.WithContext(ctx)

	for k, v := range d.DefaultHeader {
		req.Header[k] = v
	}

	d.cacheH2Lock.Lock()
	if d.cacheH2ClientConn != nil && d.cacheH2RawConn != nil {
		if d.cacheH2ClientConn.CanTakeNewRequest() {
			d.cacheH2Lock.Unlock()
			return d.connectHttp2(req, d.cacheH2RawConn, d.cacheH2ClientConn)
		}
	}
	d.cacheH2Lock.Unlock()

	switch d.ProxyURL.Scheme {
	case "http":
		rawConn, err := d.Dialer.DialContext(ctx, network, d.ProxyURL.Host)
		if err != nil {
			return nil, err
		}
		return d.connectHttp1(req, rawConn)
	case "https":
		rawConn, err := tls.Dial(network, d.ProxyURL.Host, &tls.Config{
			ServerName: d.ProxyURL.Hostname(),
			NextProtos: []string{"h2", "http/1.1"},
		})
		if err != nil {
			return nil, err
		}
		if err := rawConn.HandshakeContext(ctx); err != nil {
			return nil, err
		}

		state := rawConn.ConnectionState()
		switch state.NegotiatedProtocol {
		case "http/1.1":
			return d.connectHttp1(req, rawConn)
		case "h2":
			tr := http2.Transport{}
			h2clientConn, err := tr.NewClientConn(rawConn)
			if err != nil {
				rawConn.Close()
				return nil, err
			}

			proxyConn, err := d.connectHttp2(req, rawConn, h2clientConn)
			if err != nil {
				rawConn.Close()
				return nil, err
			}

			d.cacheH2Lock.Lock()
			d.cacheH2ClientConn = h2clientConn
			d.cacheH2RawConn = rawConn
			d.cacheH2Lock.Unlock()

			return proxyConn, nil
		default:
			return nil, errors.New("Unsupported negotiated protocol: " + state.NegotiatedProtocol)
		}
	default:
		return nil, errors.New("Unsupported proxy scheme: " + d.ProxyURL.Scheme)
	}
}

func (d *ConnectDialer) connectHttp1(req *http.Request, rawConn net.Conn) (net.Conn, error) {
	req.Proto = "HTTP/1.1"
	req.ProtoMajor = 1
	req.ProtoMinor = 1

	// Set a deadline for the CONNECT request
	deadline := time.Now().Add(d.Timeout)
	if err := rawConn.SetDeadline(deadline); err != nil {
		rawConn.Close()
		return nil, err
	}

	if err := req.Write(rawConn); err != nil {
		rawConn.Close()
		return nil, err
	}

	res, err := http.ReadResponse(bufio.NewReader(rawConn), req)
	if err != nil {
		rawConn.Close()
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		rawConn.Close()
		return nil, errors.New("Proxy responded with non 200 code: " + res.Status)
	}

	// Reset the deadline
	if err := rawConn.SetDeadline(time.Time{}); err != nil {
		rawConn.Close()
		return nil, err
	}

	return rawConn, nil
}

func (d *ConnectDialer) connectHttp2(req *http.Request, rawConn net.Conn, h2clientConn *http2.ClientConn) (net.Conn, error) {
	req.Proto = "HTTP/2.0"
	req.ProtoMajor = 2
	req.ProtoMinor = 0

	pr, pw := io.Pipe()
	req.Body = pr

	res, err := h2clientConn.RoundTrip(req)
	if err != nil {
		rawConn.Close()
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		rawConn.Close()
		return nil, errors.New("Proxy responded with non 200 code: " + res.Status)
	}

	return newHttp2Conn(rawConn, pw, res.Body), nil
}

func newHttp2Conn(c net.Conn, pipedReqBody *io.PipeWriter, respBody io.ReadCloser) net.Conn {
	return &http2Conn{Conn: c, in: pipedReqBody, out: respBody}
}

type http2Conn struct {
	net.Conn
	in  *io.PipeWriter
	out io.ReadCloser
}

func (h *http2Conn) Read(p []byte) (n int, err error) {
	return h.out.Read(p)
}

func (h *http2Conn) Write(p []byte) (n int, err error) {
	return h.in.Write(p)
}

func (h *http2Conn) Close() error {
	var retErr error = nil
	if err := h.in.Close(); err != nil {
		retErr = err
	}
	if err := h.out.Close(); err != nil {
		retErr = err
	}
	return retErr
}

func (h *http2Conn) CloseConn() error {
	return h.Conn.Close()
}

func (h *http2Conn) CloseWrite() error {
	return h.in.Close()
}

func (h *http2Conn) CloseRead() error {
	return h.out.Close()
}
