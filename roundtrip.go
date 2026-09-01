package tlsclient

import (
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
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

// RoundTripper routes requests to one of three transports:
//
//   - t1 serves HTTP/1.1 and cleartext requests. Its TLS dials go through
//     dialTLSContext, which performs the fingerprinted uTLS handshake. When a
//     connection negotiates "h2" via ALPN, t1 hands it off to t2 through
//     TLSNextProto, so the protocol is decided per connection at handshake
//     time and a host may switch protocols between reconnects.
//   - t2 is the fingerprinted HTTP/2 transport fed by t1's handoff. It never
//     dials on its own.
//   - HTTP/3 is handled by the racer, used when HTTP/3 is requested explicitly
//     or advertised by the host via Alt-Svc. Unlike t1/t2 it is independent: it
//     dials its own QUIC connections and shares nothing. The racer is nil when
//     HTTP/3 is unavailable (no h3 profile, disabled, or the dialer can't).
type RoundTripper struct {
	profile profiles.ClientProfile
	dialer  proxy.ContextDialer
	pinner  *Pinner
	tracker bandwidth.Tracker

	tlsConf  *tls.Config
	quicConf *quic.Config

	clientSessionCache tls.ClientSessionCache
	disableIPV4        bool
	disableIPV6        bool

	t1    *http.Transport
	t2    *http2.Transport
	racer *racer
}

func NewRoundTripper(profile profiles.ClientProfile, dialer proxy.ContextDialer, pinner *Pinner, tracker bandwidth.Tracker, tlsConf *tls.Config, quicConf *quic.Config, opts *TransportOptions) *RoundTripper {
	var clientSessionCache tls.ClientSessionCache
	if supportsSessionResumption(profile.ClientHelloSpec()) {
		clientSessionCache = tls.NewLRUClientSessionCache(32)
	}

	var disableKeepAlives bool
	idleConnTimeout := 90 * time.Second
	h3RaceDelay := defaultH3RaceDelay
	var disableIPV4, disableIPV6, disableHTTP3, forceHTTP3 bool
	if opts != nil {
		disableKeepAlives = opts.DisableKeepAlives
		if opts.IdleConnTimeout != 0 {
			idleConnTimeout = opts.IdleConnTimeout
		}
		if opts.HTTP3RaceDelay != 0 {
			h3RaceDelay = opts.HTTP3RaceDelay
		}
		disableIPV4 = opts.DisableIPV4
		disableIPV6 = opts.DisableIPV6
		disableHTTP3 = opts.DisableHTTP3
		forceHTTP3 = opts.ForceHTTP3
	}

	rt := &RoundTripper{
		profile: profile,
		dialer:  dialer,
		pinner:  pinner,
		tracker: tracker,

		tlsConf:  tlsConf,
		quicConf: quicConf,

		clientSessionCache: clientSessionCache,
		disableIPV4:        disableIPV4,
		disableIPV6:        disableIPV6,
	}

	rt.t1 = &http.Transport{
		DialContext:         rt.dialContext,
		DialTLSContext:      rt.dialTLSContext,
		DisableCompression:  true,
		DisableKeepAlives:   disableKeepAlives,
		IdleConnTimeout:     idleConnTimeout,
		MaxIdleConnsPerHost: 6,
	}

	if profile.H2 != nil {
		if t2, err := http2.ConfigureTransports(rt.t1); err == nil {
			t2.DisableCompression = true
			t2.IdleConnTimeout = idleConnTimeout
			t2.MaxUploadBufferPerConnection = int32(profile.H2.ConnectionFlow)
			t2.Settings = profile.H2.Settings
			t2.Priorities = profile.H2.Priorities
			t2.HeaderPriority = profile.H2.HeaderPriority
			t2.PseudoHeaderOrder = profile.PseudoHeaderOrder
			t2.ReadIdleTimeout = profile.H2.ReadIdleTimeout
			t2.InflowTimeout = profile.H2.InflowTimeout
			t2.PrefacePing = profile.H2.PrefacePing
			for _, s := range profile.H2.Settings {
				switch s.ID {
				case http2.SettingHeaderTableSize:
					t2.MaxDecoderHeaderTableSize = s.Val
				case http2.SettingMaxFrameSize:
					t2.MaxReadFrameSize = s.Val
				case http2.SettingMaxHeaderListSize:
					t2.MaxHeaderListSize = s.Val
				}
			}
			rt.t2 = t2
		}
	}

	if profile.H3 != nil && !disableHTTP3 && dialer.SupportHTTP3() {
		settings := make(map[uint64]uint64, len(profile.H3.Settings))
		order := make([]uint64, 0, len(profile.H3.Settings))
		for _, setting := range profile.H3.Settings {
			settings[setting.ID] = setting.Val
			order = append(order, setting.ID)
		}
		h3TLSConf := tlsConf.Clone()
		h3TLSConf.ClientSessionCache = clientSessionCache
		h3TLSConf.OmitEmptyPsk = true
		h3QUICConf := quicConf.Clone()
		h3QUICConf.MaxIdleTimeout = idleConnTimeout
		h3QUICConf.EnableDatagrams = true
		t3 := &http3.Transport{
			DisableCompression:      true,
			TLSClientConfig:         h3TLSConf,
			QUICConfig:              h3QUICConf,
			AdditionalSettings:      settings,
			AdditionalSettingsOrder: order,
			PseudoHeaderOrder:       profile.PseudoHeaderOrder,
			Dial:                    rt.dialQuic,
			EnableDatagrams:         true,
		}
		rt.racer = newRacer(t3, rt.dialQuic, h3RaceDelay, forceHTTP3)
	}

	return rt
}

func (rt *RoundTripper) CloseIdleConnections() {
	rt.t1.CloseIdleConnections()
	if rt.t2 != nil {
		rt.t2.CloseIdleConnections()
	}
	if rt.racer != nil {
		rt.racer.close()
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
	switch req.URL.Scheme {
	case "http", "https":
	default:
		return nil, fmt.Errorf("unsupported scheme: %s", req.URL.Scheme)
	}
	addr := net.JoinHostPort(host, port)

	if rt.racer != nil && req.URL.Scheme == "https" {
		res, err := rt.roundTripHTTP3(req, addr)
		if !errors.Is(err, errUseTCP) {
			return res, err
		}
		// h3 was declined before the request was sent: fall through to TCP with
		// the original, untouched request.
	}

	res, err := rt.t1.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	if rt.racer != nil {
		rt.racer.recordAltSvc(addr, res.Header.Get("Alt-Svc"))
	}
	return res, nil
}

// roundTripHTTP3 serves the request over an HTTP/3 connection from the racer. It
// returns errUseTCP when HTTP/3 was declined before the request was sent and the
// caller should fall back to TCP; any other error is a genuine, committed
// round-trip failure.
func (rt *RoundTripper) roundTripHTTP3(req *http.Request, addr string) (*http.Response, error) {
	cc, err := rt.racer.connection(req, addr)
	if err != nil {
		return nil, err
	}
	res, err := cc.RoundTrip(req)
	if err != nil {
		rt.racer.forget(addr, cc)
		return nil, err
	}
	rt.racer.recordAltSvc(addr, res.Header.Get("Alt-Svc"))
	return res, nil
}

// restrict narrows "tcp"/"udp" to a single address family when one is
// disabled.
func (rt *RoundTripper) restrict(network string) string {
	switch network {
	case "tcp", "udp":
		if rt.disableIPV6 {
			return network + "4"
		}
		if rt.disableIPV4 {
			return network + "6"
		}
	}
	return network
}

func (rt *RoundTripper) dialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	conn, err := rt.dialer.DialContext(ctx, rt.restrict(network), addr)
	if err != nil {
		return nil, err
	}
	return bandwidth.NewTrackedConn(conn, rt.tracker), nil
}

func (rt *RoundTripper) dialTLSContext(ctx context.Context, network, addr string) (net.Conn, error) {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}

	rawConn, err := rt.dialer.DialContext(ctx, rt.restrict(network), addr)
	if err != nil {
		return nil, err
	}
	rawConn = bandwidth.NewTrackedConn(rawConn, rt.tracker)

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

	// Return the embedded *tls.Conn: with the handshake complete it behaves
	// identically to the UConn, and it lets the transport read the negotiated
	// ALPN protocol and hand h2 connections off through TLSNextProto.
	return conn.Conn, nil
}

func (rt *RoundTripper) dialQuic(ctx context.Context, addr string, tlscfg *tls.Config, cfg *quic.Config) (*quic.Conn, error) {
	udpaddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	pconn, err := rt.dialer.ListenPacket(ctx, rt.restrict("udp"), udpaddr.String())
	if err != nil {
		return nil, err
	}
	trackedPconn, err := bandwidth.NewTrackedPacketConn(pconn, rt.tracker)
	if err != nil {
		pconn.Close()
		return nil, err
	}

	tlscfg = tlscfg.Clone()
	verify := tlscfg.VerifyPeerCertificate
	tlscfg.VerifyPeerCertificate = func(rawCerts [][]byte, chains [][]*x509.Certificate) error {
		if verify != nil {
			if err := verify(rawCerts, chains); err != nil {
				return err
			}
		}
		certs := make([]*x509.Certificate, 0, len(rawCerts))
		for _, raw := range rawCerts {
			cert, err := x509.ParseCertificate(raw)
			if err != nil {
				return err
			}
			certs = append(certs, cert)
		}
		return rt.pinner.Pin(certs, addr)
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
