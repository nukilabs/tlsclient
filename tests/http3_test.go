package tests

import (
	"io"
	"strings"
	"testing"

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
