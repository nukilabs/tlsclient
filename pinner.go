package tlsclient

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"net"
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
		pins: make(map[string][]string),
	}
}

func (p *Pinner) AddPins(hostname string, pins []string) {
	p.Lock()
	defer p.Unlock()
	p.pins[hostname] = pins
}

func (p *Pinner) Pin(conn *tls.UConn, hostname string) error {
	p.RLock()
	defer p.RUnlock()

	if _, ok := p.pins[hostname]; !ok && p.auto {
		p.AutoPin(hostname)
	} else if !ok {
		return nil
	}

	state := conn.ConnectionState()
	for _, cert := range state.PeerCertificates {
		fingerprint := p.Fingerprint(cert)
		for _, pin := range p.pins[hostname] {
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

func (p *Pinner) AutoPin(hostname string) {
	p.Lock()
	defer p.Unlock()

	var addr string
	host, port, err := net.SplitHostPort(hostname)
	if err == nil {
		addr = net.JoinHostPort(host, port)
	} else {
		addr = net.JoinHostPort(hostname, "443")
	}

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

	p.pins[hostname] = pins
}
