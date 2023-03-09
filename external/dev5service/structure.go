package dev5service

type InscriptionsResp struct {
	Inscriptions []string `json:"inscriptions"`
	Prev int `json:"prev"`
	Next int `json:"next"`
}

type InscriptionResp struct {
	Chain              string `json:"chain"`
	GenesisFee         int `json:"genesis_fee"`
	GenesisHeight            int    `json:"genesis_height"`
	Address             string   `json:"address"`
	ContentType  bool   `json:"content_type"`
	InscriptionID string `json:"inscription_id"`
	Next          string `json:"next"`
}
