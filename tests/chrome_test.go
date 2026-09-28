package tests

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"slices"
	"strings"
	"testing"

	"github.com/nukilabs/tlsclient"
	"github.com/nukilabs/tlsclient/profiles"
)

func TestChrome131(t *testing.T) {
	c := tlsclient.New(profiles.Chrome131)
	res, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var data PeetsApiCleanData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.PeetprintHash != "e7eaab546c55bcccf568e89056ffbd70" {
		t.Errorf("Expected peetprint hash 9cb72b909981b498e833d0f5e5929c70, got %s", data.PeetprintHash)
	}
	if data.AkamaiHash != "52d84b11737d980aef856699f885ca86" {
		t.Errorf("Expected akamai hash 605a1154008045d7e3cb3c6fb062c0ce, got %s", data.AkamaiHash)
	}
}

func TestChrome133(t *testing.T) {
	c := tlsclient.New(profiles.Chrome133)
	res1, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res1.Body.Close()
	var data1 PeetsApiCleanData
	if err := json.NewDecoder(res1.Body).Decode(&data1); err != nil {
		t.Fatal(err)
	}
	if data1.PeetprintHash != "1d4ffe9b0e34acac0bd883fa7f79d7b5" {
		t.Errorf("Expected peetprint hash 1d4ffe9b0e34acac0bd883fa7f79d7b5, got %s", data1.PeetprintHash)
	}
	if data1.AkamaiHash != "52d84b11737d980aef856699f885ca86" {
		t.Errorf("Expected akamai hash 52d84b11737d980aef856699f885ca86, got %s", data1.AkamaiHash)
	}
	res2, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res2.Body.Close()
	var data2 PeetsApiCleanData
	if err := json.NewDecoder(res2.Body).Decode(&data2); err != nil {
		t.Fatal(err)
	}
	if data2.PeetprintHash != "d44d68f0fce54cd423d6792272a242b8" {
		t.Errorf("Expected peetprint hash d44d68f0fce54cd423d6792272a242b8, got %s", data2.PeetprintHash)
	}
	if data2.AkamaiHash != "52d84b11737d980aef856699f885ca86" {
		t.Errorf("Expected akamai hash 52d84b11737d980aef856699f885ca86, got %s", data2.AkamaiHash)
	}
}

func TestChrome150(t *testing.T) {
	c := tlsclient.New(profiles.Chrome150)
	res1, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res1.Body.Close()
	var data1 PeetsApiCleanData
	if err := json.NewDecoder(res1.Body).Decode(&data1); err != nil {
		t.Fatal(err)
	}
	if data1.PeetprintHash != "67c3e9111bed9e7f03d2f21d6d88994b" {
		t.Errorf("Expected peetprint hash 67c3e9111bed9e7f03d2f21d6d88994b, got %s", data1.PeetprintHash)
	}
	if data1.AkamaiHash != "52d84b11737d980aef856699f885ca86" {
		t.Errorf("Expected akamai hash 52d84b11737d980aef856699f885ca86, got %s", data1.AkamaiHash)
	}
	res2, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res2.Body.Close()
	var data2 PeetsApiCleanData
	if err := json.NewDecoder(res2.Body).Decode(&data2); err != nil {
		t.Fatal(err)
	}
	if data2.PeetprintHash != "35fc5e864929e3b01e9ba9eb41bc1360" {
		t.Errorf("Expected peetprint hash 35fc5e864929e3b01e9ba9eb41bc1360, got %s", data2.PeetprintHash)
	}
	if data2.AkamaiHash != "52d84b11737d980aef856699f885ca86" {
		t.Errorf("Expected akamai hash 52d84b11737d980aef856699f885ca86, got %s", data2.AkamaiHash)
	}
}

func TestChrome152(t *testing.T) {
	c := tlsclient.New(profiles.Chrome152)
	res1, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res1.Body.Close()
	var data1 PeetsApiCleanData
	if err := json.NewDecoder(res1.Body).Decode(&data1); err != nil {
		t.Fatal(err)
	}
	if data1.PeetprintHash != "fc97c1cdfb1409c9a9326c1b726d1dee" {
		t.Errorf("Expected peetprint hash 67c3e9111bed9e7f03d2f21d6d88994b, got %s", data1.PeetprintHash)
	}
	if data1.AkamaiHash != "52d84b11737d980aef856699f885ca86" {
		t.Errorf("Expected akamai hash 52d84b11737d980aef856699f885ca86, got %s", data1.AkamaiHash)
	}
	res2, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res2.Body.Close()
	var data2 PeetsApiCleanData
	if err := json.NewDecoder(res2.Body).Decode(&data2); err != nil {
		t.Fatal(err)
	}
	if data2.PeetprintHash != "5fa343c29062ede7d0e28fd46c1052a7" {
		t.Errorf("Expected peetprint hash 35fc5e864929e3b01e9ba9eb41bc1360, got %s", data2.PeetprintHash)
	}
	if data2.AkamaiHash != "52d84b11737d980aef856699f885ca86" {
		t.Errorf("Expected akamai hash 52d84b11737d980aef856699f885ca86, got %s", data2.AkamaiHash)
	}
}

func TestChrome154(t *testing.T) {
	c := tlsclient.New(profiles.Chrome154)
	// The second request resumes the session, which adds pre_shared_key.
	peetprints := []string{"fc97c1cdfb1409c9a9326c1b726d1dee", "5fa343c29062ede7d0e28fd46c1052a7"}
	var first [][]byte
	for i, peetprint := range peetprints {
		res, err := c.Get("https://tls.peet.ws/api/all")
		if err != nil {
			t.Fatal(err)
		}
		var data PeetsApiAllData
		err = json.NewDecoder(res.Body).Decode(&data)
		res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
		if data.TLS.PeetprintHash != peetprint {
			t.Errorf("Expected peetprint hash %s, got %s", peetprint, data.TLS.PeetprintHash)
		}
		if data.HTTP2.AkamaiFingerprintHash != "52d84b11737d980aef856699f885ca86" {
			t.Errorf("Expected akamai hash 52d84b11737d980aef856699f885ca86, got %s", data.HTTP2.AkamaiFingerprintHash)
		}
		anchors := trustAnchors(t, data)
		if len(anchors) == 0 {
			t.Fatal("trust_anchors extension missing")
		}
		if !slices.IsSortedFunc(anchors, bytes.Compare) {
			t.Errorf("Expected trust anchors sorted, got %x", anchors)
		}
		if i == 0 {
			first = anchors
		} else if !slices.EqualFunc(first, anchors, bytes.Equal) {
			t.Errorf("Expected the same trust anchor order on every connection")
		}
	}
}

// trustAnchors decodes the identifiers in the trust_anchors (0xca34) extension.
func trustAnchors(t *testing.T, data PeetsApiAllData) [][]byte {
	for _, ext := range data.TLS.Extensions {
		if !strings.Contains(ext.Name, "51764") {
			continue
		}
		b, err := hex.DecodeString(ext.Data)
		if err != nil || len(b) < 2 {
			t.Fatalf("bad trust_anchors data %q", ext.Data)
		}
		b = b[2:]
		var anchors [][]byte
		for len(b) > 0 {
			n := int(b[0])
			if len(b) < 1+n {
				t.Fatalf("truncated trust_anchors data %q", ext.Data)
			}
			anchors = append(anchors, b[1:1+n])
			b = b[1+n:]
		}
		return anchors
	}
	return nil
}
