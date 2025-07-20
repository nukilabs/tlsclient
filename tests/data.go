package tests

type PeetsApiCleanData struct {
	Ja3           string `json:"ja3"`
	Ja3Hash       string `json:"ja3_hash"`
	Akamai        string `json:"akamai"`
	AkamaiHash    string `json:"akamai_hash"`
	Peetprint     string `json:"peetprint"`
	PeetprintHash string `json:"peetprint_hash"`
}

type H3ImpersonateData struct {
	HTTP3 struct {
		PerkText string `json:"perk_text"`
		PerkHash string `json:"perk_hash"`
	} `json:"http3"`
}
