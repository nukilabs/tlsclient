package tlsclient

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	tls "github.com/refraction-networking/utls"
	http "github.com/sparkaio/fhttp"
	"github.com/sparkaio/fhttp/http2"
	"github.com/sparkaio/tlsclient/bandwidth"
	"github.com/sparkaio/tlsclient/profiles"
	"golang.org/x/net/proxy"
)

type RoundTripper struct {
	sync.Mutex
	profile profiles.ClientProfile
	dialer  proxy.ContextDialer
	pinner  *Pinner
	tracker bandwidth.Tracker

	clientSessionCache tls.ClientSessionCache
	serverNameOverride string
	insecureSkipVerify bool
	disableKeepAlives  bool
	idleConnTimeout    time.Duration
	disableIPV4        bool
	disableIPV6        bool

	transportLock sync.Mutex
	transports    map[string]http.RoundTripper
	connections   map[string]net.Conn
}

func NewRoundTripper(profile profiles.ClientProfile, dialer proxy.ContextDialer, pinner *Pinner, tracker bandwidth.Tracker, opts *TransportOptions) *RoundTripper {
	var clientSessionCache tls.ClientSessionCache
	if supportsSessionResumption(profile.ClientHelloSpec()) {
		clientSessionCache = tls.NewLRUClientSessionCache(32)
	}
	var serverNameOverride string
	var insecureSkipVerify, disableKeepAlives bool
	var idleConnTimeout time.Duration = 90 * time.Second
	var disableIPV4, disableIPV6 bool
	if opts != nil {
		serverNameOverride = opts.ServerNameOverride
		insecureSkipVerify = opts.InsecureSkipVerify
		disableKeepAlives = opts.DisableKeepAlives
		if opts.IdleConnTimeout != 0 {
			idleConnTimeout = opts.IdleConnTimeout
		}
		disableIPV4 = opts.DisableIPV4
		disableIPV6 = opts.DisableIPV6
	}
	return &RoundTripper{
		profile: profile,
		dialer:  dialer,
		pinner:  pinner,
		tracker: tracker,

		clientSessionCache: clientSessionCache,
		serverNameOverride: serverNameOverride,
		insecureSkipVerify: insecureSkipVerify,
		disableKeepAlives:  disableKeepAlives,
		idleConnTimeout:    idleConnTimeout,
		disableIPV4:        disableIPV4,
		disableIPV6:        disableIPV6,

		transports:  make(map[string]http.RoundTripper),
		connections: make(map[string]net.Conn),
	}
}

func (rt *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var addr string
	host, port, err := net.SplitHostPort(req.URL.Host)
	if err != nil {
		addr = net.JoinHostPort(req.URL.Host, "443")
	} else {
		addr = net.JoinHostPort(host, port)
	}

	transport, err := rt.getTransport(req.Context(), req.URL.Scheme, addr)
	if err != nil {
		return nil, err
	}

	return transport.RoundTrip(req)
}

func (rt *RoundTripper) getTransport(ctx context.Context, scheme, addr string) (http.RoundTripper, error) {
	rt.transportLock.Lock()
	defer rt.transportLock.Unlock()

	if t, ok := rt.transports[addr]; ok {
		return t, nil
	}

	switch scheme {
	case "http":
		rt.transports[addr] = rt.buildHttp1Transport()
	case "https":
		if _, err := rt.dialTLSContext(ctx, "tcp", addr); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported scheme: %s", scheme)
	}
	return rt.transports[addr], nil
}

func (rt *RoundTripper) buildHttp1Transport() http.RoundTripper {
	return &http.Transport{
		DialContext:        rt.dialContext,
		DialTLSContext:     rt.dialTLSContext,
		DisableCompression: true,
		TLSClientConfig: &tls.Config{
			ClientSessionCache: rt.clientSessionCache,
			InsecureSkipVerify: rt.insecureSkipVerify,
			OmitEmptyPsk:       true,
		},
		DisableKeepAlives: rt.disableKeepAlives,
		IdleConnTimeout:   rt.idleConnTimeout,
	}
}

func (rt *RoundTripper) buildHttp2Transport() http.RoundTripper {
	return &http2.Transport{
		DialTLSContext:     rt.dialTLSContextHTTP2,
		DisableCompression: true,
		TLSClientConfig: &tls.Config{
			ClientSessionCache: rt.clientSessionCache,
			InsecureSkipVerify: rt.insecureSkipVerify,
			OmitEmptyPsk:       true,
		},
		ConnectionFlow:    rt.profile.ConnectionFlow,
		Settings:          rt.profile.Settings,
		Priorities:        rt.profile.Priorities,
		HeaderPriority:    rt.profile.HeaderPriority,
		PseudoHeaderOrder: rt.profile.PseudoHeaderOrder,
		IdleConnTimeout:   rt.idleConnTimeout,
	}
}

func (rt *RoundTripper) dialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	if network == "tcp" && (rt.disableIPV6) {
		network = "tcp4"
	}
	return rt.dialer.DialContext(ctx, network, addr)
}

func (rt *RoundTripper) dialTLSContext(ctx context.Context, network, addr string) (net.Conn, error) {
	rt.Lock()
	defer rt.Unlock()

	if conn := rt.connections[addr]; conn != nil {
		delete(rt.connections, addr)
		return conn, nil
	}

	if network == "tcp" && (rt.disableIPV6) {
		network = "tcp4"
	}
	if network == "tcp" && (rt.disableIPV4) {
		network = "tcp6"
	}

	rawConn, err := rt.dialer.DialContext(ctx, network, addr)
	if err != nil {
		return nil, err
	}
	rawConn = bandwidth.NewTrackedConn(rawConn, rt.tracker)

	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	tlsConfig := &tls.Config{
		ClientSessionCache: rt.clientSessionCache,
		ServerName:         host,
		InsecureSkipVerify: rt.insecureSkipVerify,
		OmitEmptyPsk:       true,
	}
	conn := tls.UClient(rawConn, tlsConfig, tls.HelloCustom)
	if err := conn.ApplyPreset(rt.profile.ClientHelloSpec()); err != nil {
		conn.Close()
		return nil, err
	}

	if err := conn.HandshakeContext(ctx); err != nil {
		conn.Close()
		return nil, err
	}

	if err := rt.pinner.Pin(conn, host); err != nil {
		conn.Close()
		return nil, err
	}

	if rt.transports[addr] != nil {
		return conn, nil
	}

	state := conn.ConnectionState()
	switch state.NegotiatedProtocol {
	case http2.NextProtoTLS:
		rt.transports[addr] = rt.buildHttp2Transport()
	default:
		rt.transports[addr] = rt.buildHttp1Transport()
	}

	rt.connections[addr] = conn

	return nil, nil
}

func (rt *RoundTripper) dialTLSContextHTTP2(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
	return rt.dialTLSContext(ctx, network, addr)
}
