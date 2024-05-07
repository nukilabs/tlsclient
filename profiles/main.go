package profiles

import (
	tls "github.com/refraction-networking/utls"
	"github.com/sparkaio/fhttp/http2"
)

type ClientProfile struct {
	ClientHelloSpec   func() *tls.ClientHelloSpec
	Settings          []http2.Setting
	ConnectionFlow    uint32
	Priorities        []http2.Priority
	HeaderPriority    *http2.PriorityParam
	PseudoHeaderOrder []string
}
