package request

type InscriptionByOutput struct {
	Outputs []string `json:"outputs"`
}

type TrackTx struct {
	Txhash        string `json:"txhash"`
	Type          string `json:"type"`
	InscriptionID string `json:"inscription_id"`
	Amount        uint64 `json:"amount"`
	Address       string `json:"address"`
}
