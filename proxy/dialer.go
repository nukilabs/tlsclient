package proxy

import (
	"context"
	"net"
	"net/url"
	"strings"
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

	h3DialLock   sync.Mutex
	h3Conn       *quic.Conn
	h3ClientConn *http3.ClientConn
}

var (
	// noDeadline is the zero time, used to clear deadlines once
	// connection establishment has completed.
	noDeadline = time.Time{}
	// aLongTimeAgo is a non-zero time far in the past, used for immediate
	// deadlines when the dial context is canceled.
	aLongTimeAgo = time.Unix(1, 0)
)

// deadline returns the earliest of the context's deadline and
// now plus the dialer's timeout.
func (d *Dialer) deadline(ctx context.Context) (time.Time, bool) {
	deadline, ok := ctx.Deadline()
	if d.timeout > 0 {
		if t := time.Now().Add(d.timeout); !ok || t.Before(deadline) {
			return t, true
		}
	}
	return deadline, ok
}

type opAddr string

func (a opAddr) Network() string { return "" }
func (a opAddr) String() string  { return string(a) }

// unescapeBraces restores the RFC 6570 template braces that url.URL.String()
// percent-encodes when they appear in the path (e.g. {target_host} becomes
// %7Btarget_host%7D). Without this, uritemplate can't recognize the variables
// and expansion leaves the placeholders literal. Braces in the query are left
// untouched by String(), so this is a no-op for query-form templates.
var braceUnescaper = strings.NewReplacer("%7B", "{", "%7b", "{", "%7D", "}", "%7d", "}")

func unescapeBraces(raw string) string {
	return braceUnescaper.Replace(raw)
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
