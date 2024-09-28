package tlsclient

import (
	"time"

	"github.com/sparkaio/tlsclient/bandwidth"
)

type TransportOptions struct {
	ServerNameOverride string
	InsecureSkipVerify bool
	DisableKeepAlives  bool
	IdleConnTimeout    time.Duration
	DisableIPV4        bool
	DisableIPV6        bool
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

func WithTransportOptions(opts TransportOptions) Option {
	return func(c *Client) {
		c.opts = &opts
	}
}
