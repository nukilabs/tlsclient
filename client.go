package tlsclient

import (
	"io"
	"net/url"
	"strings"
	"sync/atomic"
	"time"

	"github.com/nukilabs/http"
	"github.com/nukilabs/http/cookiejar"
	"github.com/nukilabs/quic-go"
	"github.com/nukilabs/tlsclient/bandwidth"
	"github.com/nukilabs/tlsclient/profiles"
	"github.com/nukilabs/tlsclient/proxy"
	tls "github.com/nukilabs/utls"
)

type Client struct {
	http.Client
	profile  profiles.ClientProfile
	pinner   *Pinner
	tracker  bandwidth.Tracker
	proxyURL *url.URL
	hooks    []HookFunc
	inHook   atomic.Bool
	redirect func(req *http.Request, via []*http.Request) error
	tlsConf  *tls.Config
	quicConf *quic.Config
	opts     *TransportOptions

	AutoDecompress bool
}

type HookFunc func(*Client, *http.Response) (*http.Response, error)

func New(profile profiles.ClientProfile, options ...Option) *Client {
	jar, _ := cookiejar.New(nil)
	client := &Client{
		Client: http.Client{
			Timeout:       30 * time.Second,
			Jar:           jar,
			CheckRedirect: nil,
		},
		profile: profile,
		tlsConf: &tls.Config{},
		quicConf: &quic.Config{
			KeepAlivePeriod: 30 * time.Second,
			EnableDatagrams: true,
		},
		AutoDecompress: true,
		tracker:        bandwidth.NewNoopTracker(),
	}
	for _, option := range options {
		option(client)
	}
	if client.pinner == nil {
		client.pinner = NewPinner(false)
	}
	dialer := proxy.Direct(client.Timeout)
	client.Transport = NewRoundTripper(profile, dialer, client.pinner, client.tracker, client.tlsConf, client.quicConf, client.opts)
	return client
}

func (c *Client) GetProxy() *url.URL {
	return c.proxyURL
}

func (c *Client) SetProxy(proxyUrl *url.URL) error {
	oldProxy := c.proxyURL
	c.proxyURL = proxyUrl
	if err := c.applyProxy(); err != nil {
		c.proxyURL = oldProxy
		if err := c.applyProxy(); err != nil {
			c.proxyURL = nil
			return c.applyProxy()
		}
	}
	return nil
}

func (c *Client) applyProxy() error {
	dialer, err := proxy.New(c.proxyURL, c.Timeout, c.tlsConf)
	if err != nil {
		return err
	}
	c.Transport = NewRoundTripper(c.profile, dialer, c.pinner, c.tracker, c.tlsConf, c.quicConf, c.opts)
	return nil
}

func (c *Client) SetHooks(hooks ...HookFunc) {
	c.hooks = hooks
}

func (c *Client) AddHooks(hooks ...HookFunc) {
	c.hooks = append(c.hooks, hooks...)
}

func (c *Client) SetCookieJar(jar http.CookieJar) {
	c.Client.Jar = jar
}

func (c *Client) GetCookieJar() http.CookieJar {
	return c.Client.Jar
}

func (c *Client) GetCookies(u *url.URL) []*http.Cookie {
	return c.Client.Jar.Cookies(u)
}

func (c *Client) SetCookies(u *url.URL, cookies []*http.Cookie) {
	c.Client.Jar.SetCookies(u, cookies)
}

func (c *Client) CloseIdleConnections() {
	c.Client.CloseIdleConnections()
}

func (c *Client) SetRedirectFunc(f func(req *http.Request, via []*http.Request) error) {
	c.Client.CheckRedirect = f
	c.redirect = f
}

func (c *Client) SetFollowRedirects(follow bool) {
	if follow && c.redirect != nil {
		c.Client.CheckRedirect = c.redirect
	} else if follow {
		c.Client.CheckRedirect = nil
	} else {
		c.Client.CheckRedirect = defaultRedirectFunc
	}
}

func (c *Client) ResetInHook() {
	c.inHook.Store(false)
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if c.AutoDecompress {
		DecompressBody(res)
	}
	for _, hook := range c.hooks {
		if c.inHook.CompareAndSwap(false, true) {
			res, err = hook(c, res)
			c.inHook.Store(false)
			if err != nil {
				return nil, err
			}
		}
	}
	return res, nil
}

func (c *Client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) Head(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return c.Do(req)
}

func (c *Client) PostForm(url string, data url.Values) (*http.Response, error) {
	return c.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}
