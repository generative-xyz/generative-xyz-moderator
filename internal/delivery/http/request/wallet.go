package request

type InscriptionByOutput struct {
	Outputs []string `json:"outputs"`
}

type TrackTx struct {
	Txhash  string `json:"txhash"`
	Address string `json:"address"`
}
