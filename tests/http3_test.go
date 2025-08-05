package tests

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/nukilabs/http"
	"github.com/nukilabs/masque-go"
	"github.com/nukilabs/quic-go/http3"
	"github.com/nukilabs/socks"
	"github.com/nukilabs/tlsclient"
	"github.com/nukilabs/tlsclient/profiles"
	tls "github.com/nukilabs/utls"
	"github.com/yosida95/uritemplate/v3"
)

func TestHTTP3(t *testing.T) {
	c := tlsclient.New(profiles.Chrome138)
	c.Get("https://http3.is/")

	res, err := c.Get("https://http3.is/")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.Proto != "HTTP/3.0" {
		t.Fatalf("Expected HTTP/3.0, got %s", res.Proto)
	}

	res2, err := c.Get("https://http3.is/")
	if err != nil {
		t.Fatal(err)
	}
	defer res2.Body.Close()

	if res2.Proto != "HTTP/3.0" {
		t.Fatalf("Expected HTTP/3.0, got %s", res.Proto)
	}
}

func TestH3SettingsOrder(t *testing.T) {
	c := tlsclient.New(profiles.Chrome138, tlsclient.WithTLSConfig(&tls.Config{
		NextProtos: []string{http3.NextProtoH3},
	}))
	res, err := c.Get("https://fp.impersonate.pro/api/http3")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.Proto != "HTTP/3.0" {
		t.Fatalf("Expected HTTP/3.0, got %s", res.Proto)
	}

	var data H3ImpersonateData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.HTTP3.PerkHash != "e1d11ee6f2f4c7b1f11bfaaf4dbbc211" {
		t.Errorf("Expected h3 hash e1d11ee6f2f4c7b1f11bfaaf4dbbc211, got %s", data.HTTP3.PerkHash)
	}
}

func TestH3SocksProxy(t *testing.T) {
	server := socks.NewServer()
	server.Authentication = socks.UserPass("user", "password")

	go server.ListenAndServe("tcp", ":1080")

	c := tlsclient.New(profiles.Chrome138, tlsclient.WithTLSConfig(&tls.Config{
		NextProtos: []string{http3.NextProtoH3},
	}))
	c.SetProxy(&url.URL{Scheme: "socks5h", Host: "localhost:1080", User: url.UserPassword("user", "password")})

	res, err := c.Get("https://fp.impersonate.pro/api/http3")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.Proto != "HTTP/3.0" {
		t.Fatalf("Expected HTTP/3.0, got %s", res.Proto)
	}

	var data H3ImpersonateData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.HTTP3.PerkHash != "e1d11ee6f2f4c7b1f11bfaaf4dbbc211" {
		t.Errorf("Expected h3 hash e1d11ee6f2f4c7b1f11bfaaf4dbbc211, got %s", data.HTTP3.PerkHash)
	}
}

func TestH3HttpProxy(t *testing.T) {
	template := uritemplate.MustNew("https://localhost:4443/masque?h={target_host}&p={target_port}")
	proxyURL, err := url.Parse("https://user:password@localhost:4443/masque?h={target_host}&p={target_port}")
	if err != nil {
		t.Fatal(err)
	}

	mux := http.NewServeMux()
	server := http3.Server{
		Addr:            ":4443",
		TLSConfig:       http3.ConfigureTLSConfig(tlsConf),
		EnableDatagrams: true,
		Handler:         mux,
	}

	var proxy masque.Proxy
	mux.HandleFunc("/masque", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Proxy-Authorization")
		if authHeader != "Basic dXNlcjpwYXNzd29yZA==" { // user:password
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		req, err := masque.ParseRequest(r, template)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		proxy.Proxy(w, req)
	})
	go server.ListenAndServe()

	c := tlsclient.New(profiles.Chrome138, tlsclient.WithTLSConfig(&tls.Config{
		RootCAs:    certPool,
		NextProtos: []string{http3.NextProtoH3},
	}))
	c.SetProxy(proxyURL)

	req, err := http.NewRequest(http.MethodGet, "https://fp.impersonate.pro/api/http3", nil)
	if err != nil {
		t.Fatal(err)
	}

	res, err := c.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if res.Proto != "HTTP/3.0" {
		t.Fatalf("Expected HTTP/3.0, got %s", res.Proto)
	}

	var data H3ImpersonateData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.HTTP3.PerkHash != "e1d11ee6f2f4c7b1f11bfaaf4dbbc211" {
		t.Errorf("Expected h3 hash e1d11ee6f2f4c7b1f11bfaaf4dbbc211, got %s", data.HTTP3.PerkHash)
	}
}
