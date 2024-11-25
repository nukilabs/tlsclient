package tests

import (
	"encoding/json"
	"testing"

	"github.com/nukilabs/tlsclient"
	"github.com/nukilabs/tlsclient/profiles"
)

func TestChrome(t *testing.T) {
	c := tlsclient.New(profiles.Chrome_124)
	c.Get("https://tls.peet.ws/api/clean")
	res, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var data PeetsApiCleanData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.PeetprintHash != "9cb72b909981b498e833d0f5e5929c70" {
		t.Errorf("Expected peetprint hash 9cb72b909981b498e833d0f5e5929c70, got %s", data.PeetprintHash)
	}
	if data.AkamaiHash != "90224459f8bf70b7d0a8797eb916dbc9" {
		t.Errorf("Expected akamai hash 605a1154008045d7e3cb3c6fb062c0ce, got %s", data.AkamaiHash)
	}
}

func TestChrome131(t *testing.T) {
	c := tlsclient.New(profiles.Chrome_131)
	c.Get("https://tls.peet.ws/api/clean")
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
