package profiles

import (
	"github.com/nukilabs/fhttp/http2"
	tls "github.com/refraction-networking/utls"
)

type ClientProfile struct {
	ClientHelloSpec   func() *tls.ClientHelloSpec
	Settings          []http2.Setting
	ConnectionFlow    uint32
	Priorities        []http2.Priority
	HeaderPriority    *http2.PriorityParam
	PseudoHeaderOrder []string
}
