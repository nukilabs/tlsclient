package tlsclient

import (
	"time"

	"github.com/nukilabs/http"
	"github.com/nukilabs/quic-go"
	"github.com/nukilabs/tlsclient/bandwidth"
	tls "github.com/nukilabs/utls"
)

type TransportOptions struct {
	DisableKeepAlives bool
	IdleConnTimeout   time.Duration
	DisableIPV4       bool
	DisableIPV6       bool
	DisableHTTP3      bool
	// ForceHTTP3 routes every https request over HTTP/3 with no TCP fallback,
	// instead of waiting for an Alt-Svc hint. It has no effect when HTTP/3 is
	// unavailable (no h3 profile, DisableHTTP3, or a dialer without h3 support).
	ForceHTTP3 bool
	// HTTP3RaceDelay is the head start given to a QUIC connection before a
	// request that has an Alt-Svc h3 hint falls back to TCP. Zero uses the
	// default (300ms), matching Chrome's alternative-service race.
	HTTP3RaceDelay time.Duration
}

type Option func(*Client)

func WithAutoPinning() Option {
	return func(c *Client) {
		c.pinner = NewPinner(true)
	}
}

func WithNoAutoDecompress() Option {
	return func(c *Client) {
		c.AutoDecompress = false
	}
}

func WithCookieJar(jar http.CookieJar) Option {
	return func(c *Client) {
		c.Client.Jar = jar
	}
}

func WithNoCookieJar() Option {
	return func(c *Client) {
		c.Client.Jar = nil
	}
}

func WithNoFollowRedirects() Option {
	return func(c *Client) {
		c.Client.CheckRedirect = nil
	}
}

func WithPinner(pinner *Pinner) Option {
	return func(c *Client) {
		c.pinner = pinner
	}
}

func WithTracker(tracker bandwidth.Tracker) Option {
	return func(c *Client) {
		c.tracker = tracker
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.Timeout = timeout
	}
}

func WithTLSConfig(tlsConf *tls.Config) Option {
	return func(c *Client) {
		c.tlsConf = tlsConf
	}
}

func WithQUICConfig(quicConf *quic.Config) Option {
	return func(c *Client) {
		c.quicConf = quicConf
	}
}

func WithTransportOptions(opts TransportOptions) Option {
	return func(c *Client) {
		c.opts = &opts
	}
}
