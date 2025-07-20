package tests

import (
	"encoding/json"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/nukilabs/tlsclient"
	"github.com/nukilabs/tlsclient/profiles"
)

func TestHTTP3(t *testing.T) {
	c := tlsclient.New(profiles.Chrome138)

	res, err := c.Get("https://http3.is/")
	if err != nil {
		t.Fatal(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(body), "it does support HTTP/3!") {
		t.Fatal("Response did not contain HTTP3 result")
	}
}

func TestH3SettingsOrder(t *testing.T) {
	c := tlsclient.New(profiles.Chrome138)
	c.Get("https://fp.impersonate.pro/api/http3")
	time.Sleep(1 * time.Second) // Allow time for the connection to establish
	res, err := c.Get("https://fp.impersonate.pro/api/http3")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var data H3ImpersonateData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.HTTP3.PerkHash != "e1d11ee6f2f4c7b1f11bfaaf4dbbc211" {
		t.Errorf("Expected peetprint hash e1d11ee6f2f4c7b1f11bfaaf4dbbc211, got %s", data.HTTP3.PerkHash)
	}
}
