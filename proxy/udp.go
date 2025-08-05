package proxy

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

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
	u, err := d.expandTemplate(addr)
	if err != nil {
		return nil, err
	}
	d.h3DialOnce.Do(func() {
		tlsConf := d.tlsConf.Clone()
		tlsConf.NextProtos = []string{http3.NextProtoH3}
		conn, err := quic.DialAddr(ctx, u.Host, tlsConf, &quic.Config{
			EnableDatagrams:   true,
			InitialPacketSize: 1350,
		})
		if err != nil {
			d.h3DialErr = fmt.Errorf("dialing quic connection failed: %w", err)
			return
		}
		d.h3Conn = conn
		tr := &http3.Transport{EnableDatagrams: true}
		d.h3ClientConn = tr.NewClientConn(conn)
	})
	if d.h3DialErr != nil {
		return nil, d.h3DialErr
	}
	select {
	case <-ctx.Done():
		return nil, context.Cause(ctx)
	case <-d.h3ClientConn.Context().Done():
		return nil, context.Cause(d.h3ClientConn.Context())
	case <-d.h3ClientConn.ReceivedSettings():
	}
	settings := d.h3ClientConn.Settings()
	if !settings.EnableExtendedConnect {
		return nil, errors.New("server didn't enable extended connect")
	}
	if !settings.EnableDatagrams {
		return nil, errors.New("server didn't enable datagrams")
	}

	rstr, err := d.h3ClientConn.OpenRequestStream(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to open request stream: %w", err)
	}
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
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	res, err := rstr.ReadResponse()
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, fmt.Errorf("server responded with %d", res.StatusCode)
	}
	return newH3Conn(rstr, d.h3Conn.LocalAddr()), nil
}

func (d *Dialer) SupportHTTP3() bool {
	raw := d.proxyURL.String()
	return strings.Contains(raw, uriTemplateTargetHost) && strings.Contains(raw, uriTemplateTargetPort)
}
