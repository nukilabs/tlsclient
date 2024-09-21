package tlsclient

import (
	"errors"
	"fmt"
)

var ErrCertificatePinningFailed = errors.New("tlsclient: certificate pinning failed")

type ErrProxy struct {
	Message string
}

func (e *ErrProxy) Error() string {
	return fmt.Sprintf("Proxy error: %s", e.Message)
}

func NewErrProxy(message string) *ErrProxy {
	return &ErrProxy{Message: message}
}
