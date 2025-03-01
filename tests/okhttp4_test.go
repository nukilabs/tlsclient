package tests

import (
	"encoding/json"
	"testing"

	"github.com/nukilabs/tlsclient"
	"github.com/nukilabs/tlsclient/profiles"
)

func TestOkhttp4Android21(t *testing.T) {
	c := tlsclient.New(profiles.Okhttp4Android(21))
	res, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var data PeetsApiCleanData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.PeetprintHash != "eb19c6e2b63eca32210d31a8b086cd44" {
		t.Errorf("Expected peetprint hash eb19c6e2b63eca32210d31a8b086cd44, got %s", data.PeetprintHash)
	}
	if data.AkamaiHash != "605a1154008045d7e3cb3c6fb062c0ce" {
		t.Errorf("Expected akamai hash 605a1154008045d7e3cb3c6fb062c0ce, got %s", data.AkamaiHash)
	}
}

func TestOkhttp4Android22(t *testing.T) {
	c := tlsclient.New(profiles.Okhttp4Android(22))
	res, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var data PeetsApiCleanData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.PeetprintHash != "eb19c6e2b63eca32210d31a8b086cd44" {
		t.Errorf("Expected peetprint hash eb19c6e2b63eca32210d31a8b086cd44, got %s", data.PeetprintHash)
	}
	if data.AkamaiHash != "605a1154008045d7e3cb3c6fb062c0ce" {
		t.Errorf("Expected akamai hash 605a1154008045d7e3cb3c6fb062c0ce, got %s", data.AkamaiHash)
	}
}

func TestOkhttp4Android23(t *testing.T) {
	c := tlsclient.New(profiles.Okhttp4Android(23))
	res, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var data PeetsApiCleanData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.PeetprintHash != "ed2770a897ccdd7f03daa812229c065e" {
		t.Errorf("Expected peetprint hash ed2770a897ccdd7f03daa812229c065e, got %s", data.PeetprintHash)
	}
	if data.AkamaiHash != "605a1154008045d7e3cb3c6fb062c0ce" {
		t.Errorf("Expected akamai hash 605a1154008045d7e3cb3c6fb062c0ce, got %s", data.AkamaiHash)
	}
}

func TestOkhttp4Android24(t *testing.T) {
	c := tlsclient.New(profiles.Okhttp4Android(24))
	res, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var data PeetsApiCleanData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.PeetprintHash != "0fdf9274db1be95c335c873e900a68ed" {
		t.Errorf("Expected peetprint hash 0fdf9274db1be95c335c873e900a68ed, got %s", data.PeetprintHash)
	}
	if data.AkamaiHash != "605a1154008045d7e3cb3c6fb062c0ce" {
		t.Errorf("Expected akamai hash 605a1154008045d7e3cb3c6fb062c0ce, got %s", data.AkamaiHash)
	}
}

func TestOkhttp4Android25(t *testing.T) {
	c := tlsclient.New(profiles.Okhttp4Android(25))
	res, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var data PeetsApiCleanData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.PeetprintHash != "730a96a208d85d7fdb8a2668780076f5" {
		t.Errorf("Expected peetprint hash 730a96a208d85d7fdb8a2668780076f5, got %s", data.PeetprintHash)
	}
	if data.AkamaiHash != "605a1154008045d7e3cb3c6fb062c0ce" {
		t.Errorf("Expected akamai hash 605a1154008045d7e3cb3c6fb062c0ce, got %s", data.AkamaiHash)
	}
}

func TestOkhttp4Android26(t *testing.T) {
	c := tlsclient.New(profiles.Okhttp4Android(26))
	res, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var data PeetsApiCleanData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.PeetprintHash != "4e947cd8ae385704ba0effe9fd805cd6" {
		t.Errorf("Expected peetprint hash 4e947cd8ae385704ba0effe9fd805cd6, got %s", data.PeetprintHash)
	}
	if data.AkamaiHash != "605a1154008045d7e3cb3c6fb062c0ce" {
		t.Errorf("Expected akamai hash 605a1154008045d7e3cb3c6fb062c0ce, got %s", data.AkamaiHash)
	}
}

func TestOkhttp4Android27(t *testing.T) {
	c := tlsclient.New(profiles.Okhttp4Android(27))
	res, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var data PeetsApiCleanData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.PeetprintHash != "4e947cd8ae385704ba0effe9fd805cd6" {
		t.Errorf("Expected peetprint hash 4e947cd8ae385704ba0effe9fd805cd6, got %s", data.PeetprintHash)
	}
	if data.AkamaiHash != "605a1154008045d7e3cb3c6fb062c0ce" {
		t.Errorf("Expected akamai hash 605a1154008045d7e3cb3c6fb062c0ce, got %s", data.AkamaiHash)
	}
}

func TestOkhttp4Android28(t *testing.T) {
	c := tlsclient.New(profiles.Okhttp4Android(28))
	res, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var data PeetsApiCleanData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.PeetprintHash != "56b4a9c0c62272cd3a9257311697e41f" {
		t.Errorf("Expected peetprint hash 56b4a9c0c62272cd3a9257311697e41f, got %s", data.PeetprintHash)
	}
	if data.AkamaiHash != "605a1154008045d7e3cb3c6fb062c0ce" {
		t.Errorf("Expected akamai hash 605a1154008045d7e3cb3c6fb062c0ce, got %s", data.AkamaiHash)
	}
}

func TestOkhttp4Android29(t *testing.T) {
	c := tlsclient.New(profiles.Okhttp4Android(29))
	res, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var data PeetsApiCleanData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.PeetprintHash != "f5a2547a1c4e1c66ba5ba13314a7c111" {
		t.Errorf("Expected peetprint hash f5a2547a1c4e1c66ba5ba13314a7c111, got %s", data.PeetprintHash)
	}
	if data.AkamaiHash != "605a1154008045d7e3cb3c6fb062c0ce" {
		t.Errorf("Expected akamai hash 605a1154008045d7e3cb3c6fb062c0ce, got %s", data.AkamaiHash)
	}
}

func TestOkhttp4Android30(t *testing.T) {
	c := tlsclient.New(profiles.Okhttp4Android(30))
	res, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var data PeetsApiCleanData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.PeetprintHash != "f5a2547a1c4e1c66ba5ba13314a7c111" {
		t.Errorf("Expected peetprint hash f5a2547a1c4e1c66ba5ba13314a7c111, got %s", data.PeetprintHash)
	}
	if data.AkamaiHash != "605a1154008045d7e3cb3c6fb062c0ce" {
		t.Errorf("Expected akamai hash 605a1154008045d7e3cb3c6fb062c0ce, got %s", data.AkamaiHash)
	}
}

func TestOkhttp4Android31(t *testing.T) {
	c := tlsclient.New(profiles.Okhttp4Android(31))
	res, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var data PeetsApiCleanData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.PeetprintHash != "f5a2547a1c4e1c66ba5ba13314a7c111" {
		t.Errorf("Expected peetprint hash f5a2547a1c4e1c66ba5ba13314a7c111, got %s", data.PeetprintHash)
	}
	if data.AkamaiHash != "605a1154008045d7e3cb3c6fb062c0ce" {
		t.Errorf("Expected akamai hash 605a1154008045d7e3cb3c6fb062c0ce, got %s", data.AkamaiHash)
	}
}

func TestOkhttp4Android32(t *testing.T) {
	c := tlsclient.New(profiles.Okhttp4Android(32))
	res, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var data PeetsApiCleanData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.PeetprintHash != "f5a2547a1c4e1c66ba5ba13314a7c111" {
		t.Errorf("Expected peetprint hash f5a2547a1c4e1c66ba5ba13314a7c111, got %s", data.PeetprintHash)
	}
	if data.AkamaiHash != "605a1154008045d7e3cb3c6fb062c0ce" {
		t.Errorf("Expected akamai hash 605a1154008045d7e3cb3c6fb062c0ce, got %s", data.AkamaiHash)
	}
}

func TestOkhttp4Android33(t *testing.T) {
	c := tlsclient.New(profiles.Okhttp4Android(33))
	res, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var data PeetsApiCleanData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.PeetprintHash != "f5a2547a1c4e1c66ba5ba13314a7c111" {
		t.Errorf("Expected peetprint hash f5a2547a1c4e1c66ba5ba13314a7c111, got %s", data.PeetprintHash)
	}
	if data.AkamaiHash != "605a1154008045d7e3cb3c6fb062c0ce" {
		t.Errorf("Expected akamai hash 605a1154008045d7e3cb3c6fb062c0ce, got %s", data.AkamaiHash)
	}
}

func TestOkhttp4Android34(t *testing.T) {
	c := tlsclient.New(profiles.Okhttp4Android(34))
	res, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var data PeetsApiCleanData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.PeetprintHash != "f5a2547a1c4e1c66ba5ba13314a7c111" {
		t.Errorf("Expected peetprint hash f5a2547a1c4e1c66ba5ba13314a7c111, got %s", data.PeetprintHash)
	}
	if data.AkamaiHash != "605a1154008045d7e3cb3c6fb062c0ce" {
		t.Errorf("Expected akamai hash 605a1154008045d7e3cb3c6fb062c0ce, got %s", data.AkamaiHash)
	}
}

func TestOkhttp4Android35(t *testing.T) {
	c := tlsclient.New(profiles.Okhttp4Android(35))
	res, err := c.Get("https://tls.peet.ws/api/clean")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var data PeetsApiCleanData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if data.PeetprintHash != "f5a2547a1c4e1c66ba5ba13314a7c111" {
		t.Errorf("Expected peetprint hash f5a2547a1c4e1c66ba5ba13314a7c111, got %s", data.PeetprintHash)
	}
	if data.AkamaiHash != "605a1154008045d7e3cb3c6fb062c0ce" {
		t.Errorf("Expected akamai hash 605a1154008045d7e3cb3c6fb062c0ce, got %s", data.AkamaiHash)
	}
}
