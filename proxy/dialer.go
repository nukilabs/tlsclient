package proxy

import (
	"net"
	"net/url"
	"sync"
	"time"

	"github.com/nukilabs/http/http2"
	"github.com/nukilabs/quic-go"
	"github.com/nukilabs/quic-go/http3"
	tls "github.com/nukilabs/utls"
	"github.com/yosida95/uritemplate/v3"
)

type Dialer struct {
	proxyURL   *url.URL
	template   *uritemplate.Template
	authHeader string
	timeout    time.Duration
	tlsConf    *tls.Config

	h2DialLock   sync.Mutex
	h2Conn       net.Conn
	h2ClientConn *http2.ClientConn

	h3DialOnce   sync.Once
	h3DialErr    error
	h3Conn       *quic.Conn
	h3ClientConn *http3.ClientConn
}

func (d *Dialer) expandTemplate(addr string) (*url.URL, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	str, err := d.template.Expand(uritemplate.Values{
		uriTemplateTargetHost: uritemplate.String(host),
		uriTemplateTargetPort: uritemplate.String(port),
	})
	if err != nil {
		return nil, err
	}
	return url.Parse(str)
}
