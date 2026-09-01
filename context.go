package tlsclient

import "context"

type contextkey int

const forceHTTP3Key contextkey = iota

// WithForceHTTP3 returns a child context that forces the request it is attached
// to over HTTP/3, with no TCP fallback. It is the per-request counterpart to
// TransportOptions.ForceHTTP3, which forces every request. It has no effect when
// HTTP/3 is unavailable (no h3 profile, DisableHTTP3, or a dialer without h3
// support).
func WithForceHTTP3(ctx context.Context) context.Context {
	return context.WithValue(ctx, forceHTTP3Key, true)
}

// forceHTTP3FromContext reports whether ctx was marked by WithForceHTTP3.
func forceHTTP3FromContext(ctx context.Context) bool {
	forced, _ := ctx.Value(forceHTTP3Key).(bool)
	return forced
}
