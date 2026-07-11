package proxy

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/nukilabs/http"
	"github.com/nukilabs/quic-go"
	"github.com/nukilabs/quic-go/http3"
	"github.com/nukilabs/quic-go/quicvarint"
)

var contextIDZero = quicvarint.Append([]byte{}, 0)

const (
	uriTemplateTargetHost = "target_host"
	uriTemplateTargetPort = "target_port"
)

func (d *Dialer) ListenPacket(ctx context.Context, network, addr string) (net.PacketConn, error) {
	proxy, dst := opAddr(d.proxyURL.Host), opAddr(addr)

	u, err := d.expandTemplate(addr)
	if err != nil {
		return nil, &net.OpError{Op: "listen", Net: network, Source: proxy, Addr: dst, Err: err}
	}
	d.h3DialLock.Lock()
	if d.h3Conn == nil || d.h3Conn.Context().Err() != nil {
		tlsConf := d.tlsConf.Clone()
		tlsConf.NextProtos = []string{http3.NextProtoH3}
		conn, err := quic.DialAddr(ctx, u.Host, tlsConf, &quic.Config{
			EnableDatagrams:      true,
			InitialPacketSize:    1350,
			HandshakeIdleTimeout: d.timeout,
		})
		if err != nil {
			d.h3DialLock.Unlock()
			return nil, &net.OpError{Op: "listen", Net: network, Source: proxy, Addr: dst, Err: fmt.Errorf("dialing quic connection failed: %w", err)}
		}
		d.h3Conn = conn
		tr := &http3.Transport{EnableDatagrams: true}
		d.h3ClientConn = tr.NewClientConn(conn)
	}
	h3Conn, clientConn := d.h3Conn, d.h3ClientConn
	d.h3DialLock.Unlock()

	var timeoutCh <-chan time.Time
	if d.timeout > 0 {
		timer := time.NewTimer(d.timeout)
		defer timer.Stop()
		timeoutCh = timer.C
	}
	select {
	case <-ctx.Done():
		return nil, &net.OpError{Op: "listen", Net: network, Source: proxy, Addr: dst, Err: context.Cause(ctx)}
	case <-timeoutCh:
		return nil, &net.OpError{Op: "listen", Net: network, Source: proxy, Addr: dst, Err: os.ErrDeadlineExceeded}
	case <-clientConn.Context().Done():
		return nil, &net.OpError{Op: "listen", Net: network, Source: proxy, Addr: dst, Err: context.Cause(clientConn.Context())}
	case <-clientConn.ReceivedSettings():
	}
	settings := clientConn.Settings()
	if !settings.EnableExtendedConnect {
		return nil, &net.OpError{Op: "listen", Net: network, Source: proxy, Addr: dst, Err: errors.New("server didn't enable extended connect")}
	}
	if !settings.EnableDatagrams {
		return nil, &net.OpError{Op: "listen", Net: network, Source: proxy, Addr: dst, Err: errors.New("server didn't enable datagrams")}
	}

	rstr, err := clientConn.OpenRequestStream(ctx)
	if err != nil {
		return nil, &net.OpError{Op: "listen", Net: network, Source: proxy, Addr: dst, Err: fmt.Errorf("failed to open request stream: %w", err)}
	}
	if deadline, ok := d.deadline(ctx); ok {
		if err := rstr.SetDeadline(deadline); err != nil {
			return nil, &net.OpError{Op: "listen", Net: network, Source: proxy, Addr: dst, Err: fmt.Errorf("setting stream deadline failed: %w", err)}
		}
	}
	defer context.AfterFunc(ctx, func() { rstr.SetDeadline(aLongTimeAgo) })()
	if err := rstr.SendRequestHeader(&http.Request{
		Method: http.MethodConnect,
		Proto:  "connect-udp",
		Host:   u.Host,
		Header: http.Header{
			http3.CapsuleProtocolHeader: {"?1"},
			"Proxy-Authorization":       {d.authHeader},
		},
		URL: u,
	}); err != nil {
		return nil, &net.OpError{Op: "listen", Net: network, Source: proxy, Addr: dst, Err: fmt.Errorf("failed to send request: %w", err)}
	}
	res, err := rstr.ReadResponse()
	if err != nil {
		return nil, &net.OpError{Op: "listen", Net: network, Source: proxy, Addr: dst, Err: fmt.Errorf("failed to read response: %w", err)}
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, &net.OpError{Op: "listen", Net: network, Source: proxy, Addr: dst, Err: fmt.Errorf("server responded with %d", res.StatusCode)}
	}
	if err := rstr.SetDeadline(noDeadline); err != nil {
		return nil, &net.OpError{Op: "listen", Net: network, Source: proxy, Addr: dst, Err: fmt.Errorf("clearing stream deadline failed: %w", err)}
	}
	return newH3Conn(rstr, h3Conn.LocalAddr()), nil
}

func (d *Dialer) SupportHTTP3() bool {
	var host, port bool
	for _, name := range d.template.Varnames() {
		switch name {
		case uriTemplateTargetHost:
			host = true
		case uriTemplateTargetPort:
			port = true
		}
	}
	return host && port
}
