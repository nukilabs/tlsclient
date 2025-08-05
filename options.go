package tlsclient

import (
	"time"

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
