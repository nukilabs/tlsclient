package tlsclient

import (
	"io"
	"net/url"
	"strings"
	"time"

	http "github.com/sparkaio/fhttp"
	"github.com/sparkaio/fhttp/cookiejar"
	"github.com/sparkaio/tlsclient/bandwidth"
	"github.com/sparkaio/tlsclient/profiles"
	"golang.org/x/net/proxy"
)

type Client struct {
	http.Client
	profile  profiles.ClientProfile
	pinner   *Pinner
	tracker  bandwidth.Tracker
	proxyURL *url.URL
	hooks    []HookFunc
	redirect func(req *http.Request, via []*http.Request) error

	AutoDecompress bool
}

type HookFunc func(*Client, *http.Response) (*http.Response, error)

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

func New(profile profiles.ClientProfile, options ...Option) *Client {
	jar, _ := cookiejar.New(nil)
	client := &Client{
		Client: http.Client{
			Timeout:       30 * time.Second,
			Jar:           jar,
			CheckRedirect: nil,
		},
		profile: profile,

		AutoDecompress: true,
		tracker:        bandwidth.NewNoopTracker(),
	}
	for _, option := range options {
		option(client)
	}
	if client.pinner == nil {
		client.pinner = NewPinner(false)
	}
	client.Transport = NewRoundTripper(profile, proxy.Direct, client.pinner, client.tracker)
	return client
}

func (c *Client) GetProxy() *url.URL {
	return c.proxyURL
}

func (c *Client) SetProxy(proxyUrl *url.URL) error {
	currentProxy := c.proxyURL
	c.proxyURL = proxyUrl
	if err := c.applyProxy(); err != nil {
		c.proxyURL = currentProxy
		if err := c.applyProxy(); err != nil {
			c.proxyURL = nil
			return c.applyProxy()
		}
	}
	return nil
}

func (c *Client) applyProxy() error {
	var dialer proxy.ContextDialer = proxy.Direct
	if c.proxyURL != nil {
		proxyDialer, err := NewConnectDialer(c.proxyURL, c.Timeout)
		if err != nil {
			return err
		}
		dialer = proxyDialer
	}
	c.Transport = NewRoundTripper(c.profile, dialer, c.pinner, c.tracker)
	return nil
}

func (c *Client) RemoveProxy() {
	c.proxyURL = nil
	c.Transport = NewRoundTripper(c.profile, proxy.Direct, c.pinner, c.tracker)
}

func (c *Client) AddHooks(hooks ...HookFunc) {
	c.hooks = append(c.hooks, hooks...)
}

func (c *Client) RemoveHooks() {
	c.hooks = nil
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

func (c *Client) SetCustomRedirectFunc(f func(req *http.Request, via []*http.Request) error) {
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

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if c.AutoDecompress {
		DecompressBody(res)
	}
	for _, hook := range c.hooks {
		res, err = hook(c, res)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (c *Client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
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
