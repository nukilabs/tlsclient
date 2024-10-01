package tlsclient

import (
	"errors"
	"fmt"
)

var ErrCertificatePinningFailed = errors.New("tlsclient: certificate pinning failed")

type ProxyError struct {
	Message string
}

func (e ProxyError) Error() string {
	return fmt.Sprintf("Proxy error: %s", e.Message)
}

func proxyError(message string) ProxyError {
	return ProxyError{Message: message}
}
