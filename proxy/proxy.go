package proxy

import (
	"context"
	"encoding/base64"
	"errors"
	"net"
	"net/url"
	"time"

	"github.com/nukilabs/socks"
	tls "github.com/nukilabs/utls"
	"github.com/yosida95/uritemplate/v3"
)

type ContextDialer interface {
	DialContext(ctx context.Context, network, addr string) (net.Conn, error)
	ListenPacket(ctx context.Context, network, addr string) (net.PacketConn, error)
	SupportHTTP3() bool
}

func New(proxyURL *url.URL, timeout time.Duration, tlsConf *tls.Config) (ContextDialer, error) {
	if proxyURL == nil {
		return Direct(nil, timeout), nil
	}

	switch proxyURL.Scheme {
	case "":
		ip := net.ParseIP(proxyURL.Host)
		if ip == nil {
			return nil, errors.New("invalid ip address for direct connection: " + proxyURL.Host)
		}
		return Direct(ip, timeout), nil
	case "socks5", "socks5h":
		dialer, err := socks.NewDialer(proxyURL)
		if err != nil {
			return nil, err
		}
		dialer.Timeout = timeout
		return dialer, nil
	case "http", "https":
		var authHeader string
		if proxyURL.User != nil {
			data := []byte(proxyURL.User.String())
			authHeader = "Basic " + base64.StdEncoding.EncodeToString(data)
		}
		template, err := uritemplate.New(unescapeBraces(proxyURL.String()))
		if err != nil {
			return nil, err
		}
		return &Dialer{
			proxyURL:   proxyURL,
			authHeader: authHeader,
			template:   template,
			timeout:    timeout,
			tlsConf:    tlsConf,
		}, nil
	default:
		return nil, errors.New("unsupported proxy scheme: " + proxyURL.Scheme)
	}
}
