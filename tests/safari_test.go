package tests

import (
	"encoding/json"
	"testing"

	"github.com/nukilabs/tlsclient"
	"github.com/nukilabs/tlsclient/profiles"
)

func TestSafari18(t *testing.T) {
	c := tlsclient.New(profiles.Safari18)
	res, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var data PeetsApiCleanData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.PeetprintHash != "fdf2c64009327d63a456cbab56a7bdde" {
		t.Errorf("Expected peetprint hash fdf2c64009327d63a456cbab56a7bdde, got %s", data.PeetprintHash)
	}
	if data.AkamaiHash != "c52879e43202aeb92740be6e8c86ea96" {
		t.Errorf("Expected akamai hash c52879e43202aeb92740be6e8c86ea96, got %s", data.AkamaiHash)
	}
}

func TestSafari26(t *testing.T) {
	c := tlsclient.New(profiles.Safari26)
	res, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var data PeetsApiCleanData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.PeetprintHash != "62b834de729e78a9f0ebd1dd099314a7" {
		t.Errorf("Expected peetprint hash 62b834de729e78a9f0ebd1dd099314a7, got %s", data.PeetprintHash)
	}
	if data.AkamaiHash != "c52879e43202aeb92740be6e8c86ea96" {
		t.Errorf("Expected akamai hash c52879e43202aeb92740be6e8c86ea96, got %s", data.AkamaiHash)
	}
}
