package profiles

import (
	"github.com/nukilabs/http/http2"
	"github.com/nukilabs/quic-go/http3"
	tls "github.com/nukilabs/utls"
)

type ClientProfile struct {
	ClientHelloSpec   func() *tls.ClientHelloSpec
	Settings          []http2.Setting
	H3Settings        []http3.Setting
	ConnectionFlow    uint32
	Priorities        []http2.Priority
	HeaderPriority    http2.PriorityParam
	PseudoHeaderOrder []string
}
