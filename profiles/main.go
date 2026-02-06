package profiles

import (
	"time"

	"github.com/nukilabs/http"
	"github.com/nukilabs/http/http2"
	"github.com/nukilabs/quic-go/http3"
	tls "github.com/nukilabs/utls"
)

type ClientProfile struct {
	ClientHelloSpec   func() *tls.ClientHelloSpec
	PseudoHeaderOrder []string
	H2                *H2ClientProfile
	H3                *H3ClientProfile
}

type H2ClientProfile struct {
	Settings        []http2.Setting
	ConnectionFlow  uint32
	Priorities      []http2.Priority
	HeaderPriority  func(*http.Request) http2.PriorityParam
	InflowTimeout   time.Duration
	ReadIdleTimeout time.Duration
	PrefacePing     bool
}

type H3ClientProfile struct {
	Settings []http3.Setting
}
