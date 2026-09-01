package tests

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
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

func verify(t *testing.T, data H3ImpersonateData) {
	if len(data.HTTP3.Settings) != 5 {
		t.Errorf("Expected 5 settings, got %d", len(data.HTTP3.Settings))
	}
	var settings []string
	for _, setting := range data.HTTP3.Settings {
		if setting.Name == "GREASE" {
			settings = append(settings, "GREASE")
		} else {
			settings = append(settings, strconv.Itoa(setting.ID)+":"+strconv.Itoa(setting.Value))
		}
	}

	if len(data.HTTP3.HeaderOrder) < 4 {
		t.Fatalf("Expected at least 4 header order, got %d", len(data.HTTP3.HeaderOrder))
	}
	headerOrder := data.HTTP3.HeaderOrder[:4]
	if strings.Join(settings, ";") != "1:65536;6:262144;7:100;51:1;GREASE" {
		t.Errorf("Expected settings %s, got %s", "1:65536;6:262144;7:100;51:1;GREASE", strings.Join(settings, ";"))
	}
	if strings.Join(headerOrder, ";") != ":method;:authority;:scheme;:path" {
		t.Errorf("Expected header order %s, got %s", ":method;:authority;:scheme;:path", strings.Join(headerOrder, ";"))
	}
}

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
	c := tlsclient.New(profiles.Chrome138, tlsclient.WithTransportOptions(tlsclient.TransportOptions{
		ForceHTTP3: true,
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
	if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	verify(t, data)
}

func TestH3SocksProxy(t *testing.T) {
	server := socks.NewServer()
	server.Authentication = socks.UserPass("user", "password")

	go server.ListenAndServe("tcp", ":1080")

	c := tlsclient.New(profiles.Chrome152, tlsclient.WithTransportOptions(tlsclient.TransportOptions{
		ForceHTTP3: true,
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
	verify(t, data)
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

	c := tlsclient.New(profiles.Chrome138,
		tlsclient.WithTLSConfig(&tls.Config{RootCAs: certPool}),
		tlsclient.WithTransportOptions(tlsclient.TransportOptions{ForceHTTP3: true}),
	)
	c.SetProxy(proxyURL)

	c.Get("https://fp.impersonate.pro/api/http3")
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
	verify(t, data)
}
