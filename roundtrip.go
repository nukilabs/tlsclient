package tlsclient

import (
	"context"
	"fmt"
	"net"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/nukilabs/http"
	"github.com/nukilabs/http/http2"
	"github.com/nukilabs/quic-go"
	"github.com/nukilabs/quic-go/http3"
	"github.com/nukilabs/tlsclient/bandwidth"
	"github.com/nukilabs/tlsclient/profiles"
	tls "github.com/nukilabs/utls"
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

	maxUploadBufferPerConnection int32
	maxReadFrameSize             uint32
	maxHeaderListSize            uint32
	maxHeaderTableSize           uint32

	transportLock sync.Mutex
	transports    map[string]http.RoundTripper
	connections   map[string]net.Conn

	altsvc sync.Map
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
	var maxHeaderTableSize, maxReadFrameSize, maxHeaderListSize uint32
	if idx := slices.IndexFunc(profile.Settings, func(s http2.Setting) bool {
		return s.ID == http2.SettingHeaderTableSize
	}); idx != -1 {
		maxHeaderListSize = profile.Settings[idx].Val
	}
	if idx := slices.IndexFunc(profile.Settings, func(s http2.Setting) bool {
		return s.ID == http2.SettingMaxFrameSize
	}); idx != -1 {
		maxReadFrameSize = profile.Settings[idx].Val
	}
	if idx := slices.IndexFunc(profile.Settings, func(s http2.Setting) bool {
		return s.ID == http2.SettingMaxHeaderListSize
	}); idx != -1 {
		maxHeaderTableSize = profile.Settings[idx].Val
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

		maxUploadBufferPerConnection: int32(profile.ConnectionFlow),
		maxReadFrameSize:             maxReadFrameSize,
		maxHeaderListSize:            maxHeaderListSize,
		maxHeaderTableSize:           maxHeaderTableSize,

		transports:  make(map[string]http.RoundTripper),
		connections: make(map[string]net.Conn),
	}
}

func (rt *RoundTripper) CloseIdleConnections() {
	rt.transportLock.Lock()
	defer rt.transportLock.Unlock()

	for addr, transport := range rt.transports {
		if t, ok := transport.(*http3.Transport); ok {
			t.CloseIdleConnections()
		} else if t, ok := transport.(*http2.Transport); ok {
			t.CloseIdleConnections()
		} else if t, ok := transport.(*http.Transport); ok {
			t.CloseIdleConnections()
		}
		delete(rt.transports, addr)
	}

	rt.Lock()
	defer rt.Unlock()

	for addr, conn := range rt.connections {
		conn.Close()
		delete(rt.connections, addr)
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

	res, err := transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	altsvc := res.Header.Get("Alt-Svc")
	if strings.HasPrefix(altsvc, "h3") {
		rt.altsvc.Store(addr, true)
	}

	return res, nil
}

func (rt *RoundTripper) getTransport(ctx context.Context, scheme, addr string) (http.RoundTripper, error) {
	rt.transportLock.Lock()
	defer rt.transportLock.Unlock()

	if _, ok := rt.altsvc.Load(addr); ok {
		rt.transports[addr] = rt.buildHttp3Transport()
		return rt.transports[addr], nil
	}

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
		MaxUploadBufferPerConnection: rt.maxUploadBufferPerConnection,
		Settings:                     rt.profile.Settings,
		Priorities:                   rt.profile.Priorities,
		HeaderPriority:               rt.profile.HeaderPriority,
		PseudoHeaderOrder:            rt.profile.PseudoHeaderOrder,
		MaxReadFrameSize:             rt.maxReadFrameSize,
		MaxHeaderListSize:            rt.maxHeaderListSize,
		MaxDecoderHeaderTableSize:    rt.maxHeaderTableSize,
		IdleConnTimeout:              rt.idleConnTimeout,
	}
}

func (rt *RoundTripper) buildHttp3Transport() http.RoundTripper {
	settings := make(map[uint64]uint64)
	order := make([]uint64, 0, len(rt.profile.H3Settings))
	for _, setting := range rt.profile.H3Settings {
		settings[setting.ID] = setting.Val
		order = append(order, setting.ID)
	}
	return &http3.Transport{
		DisableCompression: true,
		TLSClientConfig: &tls.Config{
			ClientSessionCache: rt.clientSessionCache,
			InsecureSkipVerify: rt.insecureSkipVerify,
			OmitEmptyPsk:       true,
		},
		QUICConfig: &quic.Config{
			MaxIdleTimeout:                 90 * time.Second,
			KeepAlivePeriod:                30 * time.Second,
			InitialStreamReceiveWindow:     512 * 1024,
			InitialConnectionReceiveWindow: 1024 * 1024,
			EnableDatagrams:                true,
		},
		AdditionalSettings:      settings,
		AdditionalSettingsOrder: order,
		PseudoHeaderOrder:       rt.profile.PseudoHeaderOrder,
		Dial:                    rt.dialQuic,
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
	case http3.NextProtoH3:
		rt.transports[addr] = rt.buildHttp3Transport()
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

func (rt *RoundTripper) dialQuic(ctx context.Context, addr string, tlscfg *tls.Config, cfg *quic.Config) (*quic.Conn, error) {
	udpaddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	udpConn, err := net.ListenUDP("udp", nil)
	if err != nil {
		return nil, err
	}
	trackedUdpConn := bandwidth.NewTrackedUDPConn(udpConn, rt.tracker)
	transport := &quic.Transport{
		Conn: trackedUdpConn,
	}
	conn, err := transport.DialEarly(ctx, udpaddr, tlscfg, cfg)
	if err != nil {
		udpConn.Close()
		return nil, err
	}
	return conn, nil
}
