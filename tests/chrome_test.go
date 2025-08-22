package tests

import (
	"encoding/json"
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

func TestChrome134(t *testing.T) {
	c := tlsclient.New(profiles.Chrome134)
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

func TestChrome135(t *testing.T) {
	c := tlsclient.New(profiles.Chrome135)
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

func TestChrome136(t *testing.T) {
	c := tlsclient.New(profiles.Chrome136)
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

func TestChrome137(t *testing.T) {
	c := tlsclient.New(profiles.Chrome137)
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

func TestChrome138(t *testing.T) {
	c := tlsclient.New(profiles.Chrome138)
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

func TestChrome139(t *testing.T) {
	c := tlsclient.New(profiles.Chrome138)
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
