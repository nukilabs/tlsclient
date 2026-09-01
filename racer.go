package tlsclient

import (
	"context"
	"errors"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nukilabs/http"
	"github.com/nukilabs/quic-go"
	"github.com/nukilabs/quic-go/http3"
	tls "github.com/nukilabs/utls"
)

// The racer implements Chrome's "alternative service" model for HTTP/3:
//
//   - It never speaks QUIC to an origin it has no prior knowledge of. HTTP/3 is
//     attempted only once a host advertises it via Alt-Svc, or when the client
//     is pinned to HTTP/3 (forced).
//   - On learning a hint it warms a QUIC connection in the background, the way
//     Chrome pre-connects an alternative service. A request then prefers an
//     established QUIC session but only waits a short head start for one to come
//     up before falling back to TCP, so a black-holed UDP path costs one head
//     start rather than a full handshake timeout.
//   - A request is sent on exactly one transport. Once it is committed to a live
//     QUIC connection its result stands: a failure there is a real error, never
//     a reason to replay a request that may already have reached the server.
//   - Repeated dial failures put an origin in exponential backoff, so a black
//     hole is not re-dialed on every request.
//
// errUseTCP is the sentinel returned by connection when HTTP/3 was not used and
// the caller should fall back to TCP. It never escapes RoundTripper.
var errUseTCP = errors.New("tlsclient: no usable HTTP/3 connection")

const (
	defaultH3RaceDelay  = 300 * time.Millisecond
	h3DialTimeout       = 10 * time.Second
	h3BackoffBase       = 1 * time.Second
	h3BackoffMax        = 60 * time.Second
	h3MaxEntries        = 512
	altSvcDefaultMaxAge = 24 * time.Hour
)

// quicDialer dials a QUIC connection to addr. It matches RoundTripper.dialQuic.
type quicDialer func(ctx context.Context, addr string, tlsConf *tls.Config, cfg *quic.Config) (*quic.Conn, error)

// racer owns all HTTP/3 state for a RoundTripper: the fingerprinted transport,
// the Alt-Svc hint cache, and the per-origin QUIC connection pool. It is safe
// for concurrent use.
type racer struct {
	transport *http3.Transport
	dial      quicDialer
	headStart time.Duration
	forceH3   bool // every https request goes over HTTP/3, no TCP fallback

	mu    sync.Mutex
	hints map[string]time.Time // addr -> Alt-Svc h3 hint expiry
	conns map[string]*h3conn   // addr -> QUIC connection state
}

func newRacer(transport *http3.Transport, dial quicDialer, headStart time.Duration, forceH3 bool) *racer {
	if headStart <= 0 {
		headStart = defaultH3RaceDelay
	}
	return &racer{
		transport: transport,
		dial:      dial,
		headStart: headStart,
		forceH3:   forceH3,
		hints:     make(map[string]time.Time),
		conns:     make(map[string]*h3conn),
	}
}

// h3conn is the state of one origin's QUIC connection. It doubles as a future:
// establish fills cc/err (and, on failure, retryAt) before closing ready, after
// which those fields are safe to read.
type h3conn struct {
	ready   chan struct{}
	cc      *http3.ClientConn
	err     error
	fails   int
	retryAt time.Time // when a failed origin may be dialed again
}

// live reports whether a resolved entry holds a usable connection.
func (c *h3conn) live() bool {
	return c.err == nil && c.cc != nil && c.cc.Context().Err() == nil
}

// connection returns a QUIC connection to serve req over HTTP/3. It returns
// errUseTCP when HTTP/3 should not be used and the request — still unsent —
// should go over TCP instead. A forced request that cannot reach HTTP/3 returns
// its dial error instead, since it has no TCP fallback.
func (r *racer) connection(req *http.Request, addr string) (*http3.ClientConn, error) {
	forced := r.forceH3 || forceHTTP3FromContext(req.Context())
	if !forced && !r.hinted(addr) {
		return nil, errUseTCP
	}

	c := r.acquire(addr, forced)
	if c == nil {
		return nil, errUseTCP // in backoff
	}

	if err := r.await(req.Context(), c, forced); err != nil {
		if forced {
			return nil, err
		}
		return nil, errUseTCP // still warming, or the dial failed
	}
	return c.cc, nil
}

// await blocks until the connection is ready to use. A forced request waits for
// the dial to resolve; an opportunistic one waits only the head start before
// giving up so the request can proceed over TCP.
func (r *racer) await(ctx context.Context, c *h3conn, forced bool) error {
	if !forced {
		timer := time.NewTimer(r.headStart)
		defer timer.Stop()
		select {
		case <-c.ready:
		case <-timer.C:
			return errUseTCP
		case <-ctx.Done():
			return ctx.Err()
		}
	} else {
		select {
		case <-c.ready:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	if !c.live() {
		return c.err
	}
	return nil
}

// acquire returns the QUIC connection entry for addr, starting a dial if none
// is in flight. It returns nil when the origin is in failure backoff (unless
// the request forces HTTP/3, which has no fallback and so ignores backoff).
func (r *racer) acquire(addr string, forced bool) *h3conn {
	r.mu.Lock()
	defer r.mu.Unlock()

	fails := 0
	if c, ok := r.conns[addr]; ok {
		select {
		case <-c.ready:
			if c.live() {
				return c
			}
			if !forced && time.Now().Before(c.retryAt) {
				return nil
			}
			fails = c.fails // carry the count forward so backoff escalates
			delete(r.conns, addr)
		default:
			return c // dial in flight
		}
	}

	r.pruneLocked()
	c := &h3conn{ready: make(chan struct{}), fails: fails}
	r.conns[addr] = c
	go r.establish(addr, c)
	return c
}

func (r *racer) establish(addr string, c *h3conn) {
	defer close(c.ready)

	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		r.markFailed(c, err)
		return
	}

	tlsConf := r.transport.TLSClientConfig.Clone()
	tlsConf.ServerName = host
	tlsConf.NextProtos = []string{http3.NextProtoH3}

	ctx, cancel := context.WithTimeout(context.Background(), h3DialTimeout)
	defer cancel()

	conn, err := r.dial(ctx, addr, tlsConf, r.transport.QUICConfig)
	if err != nil {
		r.markFailed(c, err)
		return
	}
	select {
	case <-conn.HandshakeComplete():
	case <-ctx.Done():
		conn.CloseWithError(quic.ApplicationErrorCode(http3.ErrCodeNoError), "")
		r.markFailed(c, ctx.Err())
		return
	}

	c.cc = r.transport.NewClientConn(conn)
}

// markFailed records a dial failure and schedules the backoff before which the
// origin will not be dialed again. Safe without the lock: the fields are read
// only after ready is closed.
func (r *racer) markFailed(c *h3conn, err error) {
	c.err = err
	c.fails++
	backoff := min(h3BackoffBase<<min(c.fails-1, 6), h3BackoffMax)
	c.retryAt = time.Now().Add(backoff)
}

// forget drops a connection that failed mid-request and closes it. cc is the
// connection handed out by connection; the entry is only removed if it still
// holds that same connection.
func (r *racer) forget(addr string, cc *http3.ClientConn) {
	r.mu.Lock()
	if c, ok := r.conns[addr]; ok && c.cc == cc {
		delete(r.conns, addr)
	}
	r.mu.Unlock()

	cc.CloseWithError(quic.ApplicationErrorCode(http3.ErrCodeNoError), "")
}

// hinted reports whether addr has a live Alt-Svc HTTP/3 hint.
func (r *racer) hinted(addr string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	exp, ok := r.hints[addr]
	if !ok {
		return false
	}
	if time.Now().After(exp) {
		delete(r.hints, addr)
		return false
	}
	return true
}

// recordAltSvc parses an Alt-Svc response header, records whether the origin
// offers HTTP/3, and on first discovery warms a QUIC connection in the
// background so the next request finds the session ready.
func (r *racer) recordAltSvc(addr, header string) {
	if header == "" {
		return
	}
	if strings.EqualFold(strings.TrimSpace(header), "clear") {
		r.mu.Lock()
		delete(r.hints, addr)
		r.mu.Unlock()
		return
	}

	maxAge, ok := parseH3AltSvc(header)
	if !ok {
		return
	}

	r.mu.Lock()
	_, known := r.hints[addr]
	r.hints[addr] = time.Now().Add(maxAge)
	r.pruneLocked()
	r.mu.Unlock()

	if !known {
		r.acquire(addr, false) // warm the connection for the next request
	}
}

// parseH3AltSvc returns the max-age of an Alt-Svc header's HTTP/3 entry, if any.
func parseH3AltSvc(header string) (time.Duration, bool) {
	for entry := range strings.SplitSeq(header, ",") {
		entry = strings.TrimSpace(entry)
		proto, _, ok := strings.Cut(entry, "=")
		if !ok || (proto != "h3" && !strings.HasPrefix(proto, "h3-")) {
			continue
		}
		maxAge := altSvcDefaultMaxAge
		for _, param := range strings.Split(entry, ";")[1:] {
			key, val, ok := strings.Cut(strings.TrimSpace(param), "=")
			if ok && strings.EqualFold(strings.TrimSpace(key), "ma") {
				if secs, err := strconv.Atoi(strings.TrimSpace(val)); err == nil && secs > 0 {
					maxAge = time.Duration(secs) * time.Second
				}
			}
		}
		return maxAge, true
	}
	return 0, false
}

// pruneLocked bounds the hint and connection maps by dropping expired hints and
// dead connections once either grows past its limit. Callers must hold mu.
func (r *racer) pruneLocked() {
	now := time.Now()
	if len(r.hints) > h3MaxEntries {
		for addr, exp := range r.hints {
			if now.After(exp) {
				delete(r.hints, addr)
			}
		}
	}
	if len(r.conns) > h3MaxEntries {
		for addr, c := range r.conns {
			select {
			case <-c.ready:
				if !c.live() && now.After(c.retryAt) {
					delete(r.conns, addr)
				}
			default:
			}
		}
	}
}

// close tears down every cached QUIC connection and clears all state.
func (r *racer) close() {
	r.transport.CloseIdleConnections()

	r.mu.Lock()
	conns := r.conns
	r.conns = make(map[string]*h3conn)
	r.hints = make(map[string]time.Time)
	r.mu.Unlock()

	for _, c := range conns {
		if c.cc != nil {
			c.cc.CloseWithError(quic.ApplicationErrorCode(http3.ErrCodeNoError), "")
		}
	}
}
