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
	"github.com/nukilabs/tlsclient/proxy"
	tls "github.com/nukilabs/utls"
)

type RoundTripper struct {
	sync.Mutex
	profile profiles.ClientProfile
	dialer  proxy.ContextDialer
	pinner  *Pinner
	tracker bandwidth.Tracker

	tlsConf  *tls.Config
	quicConf *quic.Config

	clientSessionCache tls.ClientSessionCache
	disableKeepAlives  bool
	idleConnTimeout    time.Duration
	disableIPV4        bool
	disableIPV6        bool
	disableHTTP3       bool

	maxUploadBufferPerConnection int32
	maxReadFrameSize             uint32
	maxHeaderListSize            uint32
	maxHeaderTableSize           uint32

	transportLock sync.Mutex
	transports    map[string]http.RoundTripper
	connections   map[string]net.Conn

	altsvc sync.Map
}

func NewRoundTripper(profile profiles.ClientProfile, dialer proxy.ContextDialer, pinner *Pinner, tracker bandwidth.Tracker, tlsConf *tls.Config, quicConf *quic.Config, opts *TransportOptions) *RoundTripper {
	var clientSessionCache tls.ClientSessionCache
	if supportsSessionResumption(profile.ClientHelloSpec()) {
		clientSessionCache = tls.NewLRUClientSessionCache(32)
	}
	var disableKeepAlives bool
	var idleConnTimeout time.Duration = 90 * time.Second
	var disableIPV4, disableIPV6, disableHTTP3 bool
	if opts != nil {
		disableKeepAlives = opts.DisableKeepAlives
		if opts.IdleConnTimeout != 0 {
			idleConnTimeout = opts.IdleConnTimeout
		}
		disableIPV4 = opts.DisableIPV4
		disableIPV6 = opts.DisableIPV6
		disableHTTP3 = opts.DisableHTTP3
	}
	var maxHeaderTableSize, maxReadFrameSize, maxHeaderListSize uint32
	if profile.H2 != nil {
		if idx := slices.IndexFunc(profile.H2.Settings, func(s http2.Setting) bool {
			return s.ID == http2.SettingHeaderTableSize
		}); idx != -1 {
			maxHeaderListSize = profile.H2.Settings[idx].Val
		}
		if idx := slices.IndexFunc(profile.H2.Settings, func(s http2.Setting) bool {
			return s.ID == http2.SettingMaxFrameSize
		}); idx != -1 {
			maxReadFrameSize = profile.H2.Settings[idx].Val
		}
		if idx := slices.IndexFunc(profile.H2.Settings, func(s http2.Setting) bool {
			return s.ID == http2.SettingMaxHeaderListSize
		}); idx != -1 {
			maxHeaderTableSize = profile.H2.Settings[idx].Val
		}
	}

	var maxUploadBufferPerConnection int32
	if profile.H2 != nil {
		maxUploadBufferPerConnection = int32(profile.H2.ConnectionFlow)
	}

	return &RoundTripper{
		profile: profile,
		dialer:  dialer,
		pinner:  pinner,
		tracker: tracker,

		tlsConf:  tlsConf,
		quicConf: quicConf,

		clientSessionCache: clientSessionCache,
		disableKeepAlives:  disableKeepAlives,
		idleConnTimeout:    idleConnTimeout,
		disableIPV4:        disableIPV4,
		disableIPV6:        disableIPV6,
		disableHTTP3:       disableHTTP3,

		maxUploadBufferPerConnection: maxUploadBufferPerConnection,
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
	host := req.URL.Hostname()
	port := req.URL.Port()
	if port == "" {
		switch req.URL.Scheme {
		case "http":
			port = "80"
		case "https":
			port = "443"
		}
	}
	addr := net.JoinHostPort(host, port)

	transport, err := rt.getTransport(req.Context(), req.Proto, req.URL.Scheme, addr)
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

func (rt *RoundTripper) forceHTTP3(proto, addr string) bool {
	if rt.profile.H3 == nil || rt.disableHTTP3 || !rt.dialer.SupportHTTP3() {
		return false
	}
	if strings.EqualFold(proto, "HTTP/1.0") || strings.EqualFold(proto, "HTTP/1.1") {
		return false
	}
	if strings.EqualFold(proto, "HTTP/2.0") || strings.EqualFold(proto, "h2") {
		return false
	}
	if _, ok := rt.altsvc.Load(addr); ok {
		return true
	}
	if len(rt.tlsConf.NextProtos) == 1 && rt.tlsConf.NextProtos[0] == http3.NextProtoH3 {
		return true
	}
	if strings.EqualFold(proto, "HTTP/3.0") || strings.EqualFold(proto, "h3") {
		return true
	}
	return false
}

func (rt *RoundTripper) getTransport(ctx context.Context, proto, scheme, addr string) (http.RoundTripper, error) {
	rt.transportLock.Lock()
	defer rt.transportLock.Unlock()

	if rt.forceHTTP3(proto, addr) {
		if t, ok := rt.transports[addr]; ok {
			if _, ok := t.(*http3.Transport); ok {
				return t, nil
			}
		}
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
	tlsConf := rt.tlsConf.Clone()
	tlsConf.ClientSessionCache = rt.clientSessionCache
	tlsConf.OmitEmptyPsk = true
	return &http.Transport{
		DialContext:        rt.dialContext,
		DialTLSContext:     rt.dialTLSContext,
		DisableCompression: true,
		TLSClientConfig:    tlsConf,
		DisableKeepAlives:  rt.disableKeepAlives,
		IdleConnTimeout:    rt.idleConnTimeout,
	}
}

func (rt *RoundTripper) buildHttp2Transport() http.RoundTripper {
	tlsConf := rt.tlsConf.Clone()
	tlsConf.ClientSessionCache = rt.clientSessionCache
	tlsConf.OmitEmptyPsk = true
	return &http2.Transport{
		DialTLSContext: func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
			return rt.dialTLSContext(ctx, network, addr)
		},
		DisableCompression:           true,
		TLSClientConfig:              tlsConf,
		MaxUploadBufferPerConnection: rt.maxUploadBufferPerConnection,
		Settings:                     rt.profile.H2.Settings,
		Priorities:                   rt.profile.H2.Priorities,
		HeaderPriority:               rt.profile.H2.HeaderPriority,
		PseudoHeaderOrder:            rt.profile.PseudoHeaderOrder,
		MaxReadFrameSize:             rt.maxReadFrameSize,
		MaxHeaderListSize:            rt.maxHeaderListSize,
		MaxDecoderHeaderTableSize:    rt.maxHeaderTableSize,
		IdleConnTimeout:              rt.idleConnTimeout,
		ReadIdleTimeout:              rt.profile.H2.ReadIdleTimeout,
		InflowTimeout:                rt.profile.H2.InflowTimeout,
		PrefacePing:                  rt.profile.H2.PrefacePing,
	}
}

func (rt *RoundTripper) buildHttp3Transport() http.RoundTripper {
	settings := make(map[uint64]uint64)
	order := make([]uint64, 0, len(rt.profile.H3.Settings))
	for _, setting := range rt.profile.H3.Settings {
		settings[setting.ID] = setting.Val
		order = append(order, setting.ID)
	}
	tlsConf := rt.tlsConf.Clone()
	tlsConf.ClientSessionCache = rt.clientSessionCache
	tlsConf.OmitEmptyPsk = true
	quicConf := rt.quicConf.Clone()
	quicConf.MaxIdleTimeout = rt.idleConnTimeout
	quicConf.EnableDatagrams = true
	return &http3.Transport{
		DisableCompression:      true,
		TLSClientConfig:         tlsConf,
		QUICConfig:              quicConf,
		AdditionalSettings:      settings,
		AdditionalSettingsOrder: order,
		PseudoHeaderOrder:       rt.profile.PseudoHeaderOrder,
		Dial:                    rt.dialQuic,
		EnableDatagrams:         true,
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
	tlsConf := rt.tlsConf.Clone()
	tlsConf.ServerName = host
	tlsConf.ClientSessionCache = rt.clientSessionCache
	tlsConf.OmitEmptyPsk = true
	conn := tls.UClient(rawConn, tlsConf, tls.HelloCustom)
	if err := conn.ApplyPreset(rt.profile.ClientHelloSpec()); err != nil {
		conn.Close()
		return nil, err
	}

	if err := conn.HandshakeContext(ctx); err != nil {
		conn.Close()
		return nil, err
	}

	state := conn.ConnectionState()
	if err := rt.pinner.Pin(state.PeerCertificates, addr); err != nil {
		conn.Close()
		return nil, err
	}

	if rt.transports[addr] != nil {
		return conn, nil
	}

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

func (rt *RoundTripper) dialQuic(ctx context.Context, addr string, tlscfg *tls.Config, cfg *quic.Config) (*quic.Conn, error) {
	udpaddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	network := "udp"
	if rt.disableIPV6 {
		network = "udp4"
	}
	if rt.disableIPV4 {
		network = "udp6"
	}

	pconn, err := rt.dialer.ListenPacket(ctx, network, udpaddr.String())
	if err != nil {
		return nil, err
	}
	trackedPconn, err := bandwidth.NewTrackedPacketConn(pconn, rt.tracker)
	if err != nil {
		pconn.Close()
		return nil, err
	}

	cfg = cfg.Clone()
	cfg.DisablePathMTUDiscovery = true
	conn, err := quic.DialEarly(ctx, trackedPconn, udpaddr, tlscfg, cfg)
	if err != nil {
		pconn.Close()
		return nil, err
	}

	return conn, nil
}
