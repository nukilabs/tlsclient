package tests

import (
	"encoding/json"
	"testing"

	"github.com/nukilabs/tlsclient"
	"github.com/nukilabs/tlsclient/profiles"
)

func TestHttpUrlConnectionAndroid21(t *testing.T) {
	c := tlsclient.New(profiles.HttpUrlConnectionAndroid(21))
	res, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var data PeetsApiCleanData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.PeetprintHash != "62766c449998184019a8a0bb607ea739" {
		t.Errorf("Expected peetprint hash 62766c449998184019a8a0bb607ea739, got %s", data.PeetprintHash)
	}
}

func TestHttpUrlConnectionAndroid22(t *testing.T) {
	c := tlsclient.New(profiles.HttpUrlConnectionAndroid(22))
	res, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var data PeetsApiCleanData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.PeetprintHash != "62766c449998184019a8a0bb607ea739" {
		t.Errorf("Expected peetprint hash 62766c449998184019a8a0bb607ea739, got %s", data.PeetprintHash)
	}
}

func TestHttpUrlConnectionAndroid23(t *testing.T) {
	c := tlsclient.New(profiles.HttpUrlConnectionAndroid(23))
	res, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var data PeetsApiCleanData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.PeetprintHash != "3dbba614799048f4cba3a71412538415" {
		t.Errorf("Expected peetprint hash 3dbba614799048f4cba3a71412538415, got %s", data.PeetprintHash)
	}
}

func TestHttpUrlConnectionAndroid24(t *testing.T) {
	c := tlsclient.New(profiles.HttpUrlConnectionAndroid(24))
	res, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var data PeetsApiCleanData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.PeetprintHash != "4f5e8f0498b9ca9b990212e23053ed09" {
		t.Errorf("Expected peetprint hash 4f5e8f0498b9ca9b990212e23053ed09, got %s", data.PeetprintHash)
	}
}

func TestHttpUrlConnectionAndroid25(t *testing.T) {
	c := tlsclient.New(profiles.HttpUrlConnectionAndroid(25))
	res, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var data PeetsApiCleanData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.PeetprintHash != "f08087e49108105745ed6350d5b9b369" {
		t.Errorf("Expected peetprint hash f08087e49108105745ed6350d5b9b369, got %s", data.PeetprintHash)
	}
}

func TestHttpUrlConnectionAndroid26(t *testing.T) {
	c := tlsclient.New(profiles.HttpUrlConnectionAndroid(26))
	res, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var data PeetsApiCleanData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.PeetprintHash != "1bd6e81e447473180b7357172a31313e" {
		t.Errorf("Expected peetprint hash 1bd6e81e447473180b7357172a31313e, got %s", data.PeetprintHash)
	}
}

func TestHttpUrlConnectionAndroid27(t *testing.T) {
	c := tlsclient.New(profiles.HttpUrlConnectionAndroid(27))
	res, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var data PeetsApiCleanData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.PeetprintHash != "1bd6e81e447473180b7357172a31313e" {
		t.Errorf("Expected peetprint hash 1bd6e81e447473180b7357172a31313e, got %s", data.PeetprintHash)
	}
}

func TestHttpUrlConnectionAndroid28(t *testing.T) {
	c := tlsclient.New(profiles.HttpUrlConnectionAndroid(28))
	res, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var data PeetsApiCleanData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.PeetprintHash != "aa2d2a21efdcc185bc1266c7cfd212ed" {
		t.Errorf("Expected peetprint hash aa2d2a21efdcc185bc1266c7cfd212ed, got %s", data.PeetprintHash)
	}
}

func TestHttpUrlConnectionAndroid29(t *testing.T) {
	c := tlsclient.New(profiles.HttpUrlConnectionAndroid(29))
	res, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var data PeetsApiCleanData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.PeetprintHash != "7664b85b73761fcb1c0d0709bc7a00d7" {
		t.Errorf("Expected peetprint hash 7664b85b73761fcb1c0d0709bc7a00d7, got %s", data.PeetprintHash)
	}
}

func TestHttpUrlConnectionAndroid30(t *testing.T) {
	c := tlsclient.New(profiles.HttpUrlConnectionAndroid(30))
	res, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var data PeetsApiCleanData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.PeetprintHash != "7664b85b73761fcb1c0d0709bc7a00d7" {
		t.Errorf("Expected peetprint hash 7664b85b73761fcb1c0d0709bc7a00d7, got %s", data.PeetprintHash)
	}
}

func TestHttpUrlConnectionAndroid31(t *testing.T) {
	c := tlsclient.New(profiles.HttpUrlConnectionAndroid(31))
	res, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var data PeetsApiCleanData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.PeetprintHash != "7664b85b73761fcb1c0d0709bc7a00d7" {
		t.Errorf("Expected peetprint hash 7664b85b73761fcb1c0d0709bc7a00d7, got %s", data.PeetprintHash)
	}
}

func TestHttpUrlConnectionAndroid32(t *testing.T) {
	c := tlsclient.New(profiles.HttpUrlConnectionAndroid(32))
	res, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var data PeetsApiCleanData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.PeetprintHash != "7664b85b73761fcb1c0d0709bc7a00d7" {
		t.Errorf("Expected peetprint hash 7664b85b73761fcb1c0d0709bc7a00d7, got %s", data.PeetprintHash)
	}
}

func TestHttpUrlConnectionAndroid33(t *testing.T) {
	c := tlsclient.New(profiles.HttpUrlConnectionAndroid(33))
	res, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var data PeetsApiCleanData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.PeetprintHash != "7664b85b73761fcb1c0d0709bc7a00d7" {
		t.Errorf("Expected peetprint hash 7664b85b73761fcb1c0d0709bc7a00d7, got %s", data.PeetprintHash)
	}
}

func TestHttpUrlConnectionAndroid34(t *testing.T) {
	c := tlsclient.New(profiles.HttpUrlConnectionAndroid(34))
	res, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var data PeetsApiCleanData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.PeetprintHash != "7664b85b73761fcb1c0d0709bc7a00d7" {
		t.Errorf("Expected peetprint hash 7664b85b73761fcb1c0d0709bc7a00d7, got %s", data.PeetprintHash)
	}
}

func TestHttpUrlConnectionAndroid35(t *testing.T) {
	c := tlsclient.New(profiles.HttpUrlConnectionAndroid(34))
	res, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var data PeetsApiCleanData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.PeetprintHash != "7664b85b73761fcb1c0d0709bc7a00d7" {
		t.Errorf("Expected peetprint hash 7664b85b73761fcb1c0d0709bc7a00d7, got %s", data.PeetprintHash)
	}
}
