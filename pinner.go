package tlsclient

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"sync"

	tls "github.com/nukilabs/utls"
)

type Pinner struct {
	sync.RWMutex
	auto bool
	pins map[string][]string
}

func NewPinner(auto bool) *Pinner {
	return &Pinner{
		auto: auto,
		pins: make(map[string][]string),
	}
}

func (p *Pinner) AddPins(hostname string, pins []string) {
	p.Lock()
	defer p.Unlock()
	p.pins[hostname] = pins
}

func (p *Pinner) Pin(certs []*x509.Certificate, addr string) error {
	p.RLock()
	defer p.RUnlock()

	if _, ok := p.pins[addr]; !ok && p.auto {
		p.RUnlock()
		p.AutoPin(addr)
		p.RLock()
	} else if !ok {
		return nil
	}

	for _, cert := range certs {
		fingerprint := p.Fingerprint(cert)
		for _, pin := range p.pins[addr] {
			if fingerprint == pin {
				return nil
			}
		}
	}

	return ErrCertificatePinningFailed
}

func (p *Pinner) Fingerprint(cert *x509.Certificate) string {
	digest := sha256.Sum256(cert.RawSubjectPublicKeyInfo)
	return base64.StdEncoding.EncodeToString(digest[:])
}

func (p *Pinner) AutoPin(addr string) {
	p.Lock()
	defer p.Unlock()

	conn, err := tls.Dial("tcp", addr, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return
	}

	state := conn.ConnectionState()
	conn.Close()

	pins := make([]string, 0, len(state.PeerCertificates))

	for _, cert := range state.PeerCertificates {
		pins = append(pins, p.Fingerprint(cert))
	}

	p.pins[addr] = pins
}
