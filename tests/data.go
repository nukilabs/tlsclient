package tests

type PeetsApiCleanData struct {
	Ja3           string `json:"ja3"`
	Ja3Hash       string `json:"ja3_hash"`
	Akamai        string `json:"akamai"`
	AkamaiHash    string `json:"akamai_hash"`
	Peetprint     string `json:"peetprint"`
	PeetprintHash string `json:"peetprint_hash"`
}

type PeetsApiAllData struct {
	TLS struct {
		PeetprintHash string `json:"peetprint_hash"`
		Extensions    []struct {
			Name string `json:"name"`
			Data string `json:"data"`
		} `json:"extensions"`
	} `json:"tls"`
	HTTP2 struct {
		AkamaiFingerprintHash string `json:"akamai_fingerprint_hash"`
	} `json:"http2"`
}

type H3ImpersonateData struct {
	HTTP3 struct {
		Settings []struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Value int    `json:"value"`
		} `json:"settings"`
		HeaderOrder []string `json:"header_order"`
	} `json:"http3"`
}
