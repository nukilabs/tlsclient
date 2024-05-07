package tlsclient

import (
	tls "github.com/refraction-networking/utls"
	http "github.com/sparkaio/fhttp"
)

func supportsSessionResumption(spec *tls.ClientHelloSpec) bool {
	if spec == nil {
		return false
	}
	for _, ext := range spec.Extensions {
		if _, ok := ext.(*tls.UtlsPreSharedKeyExtension); ok {
			return true
		}
	}

	return false
}

var defaultRedirectFunc = func(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}
