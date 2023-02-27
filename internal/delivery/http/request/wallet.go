package request

type InscriptionByOutput struct {
	Outputs []string `json:"outputs"`
}

type TrackTx struct {
	Txhash            string `json:"txhash"`
	Type              string `json:"type"`
	InscriptionID     string `json:"inscription_id"`
	InscriptionNumber uint64 `json:"inscription_number"`
	Amount            uint64 `json:"amount"`
	Address           string `json:"address"`
	Receiver          string `json:"receiver"`
}
