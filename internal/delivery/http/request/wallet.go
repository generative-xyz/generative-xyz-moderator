package request

type InscriptionByOutput struct {
	Outputs []string `json:"outputs"`
}

type TrackTx struct {
	Txhash                string   `json:"txhash"`
	Type                  string   `json:"type"`
	InscriptionID         string   `json:"inscription_id"`
	InscriptionNumber     uint64   `json:"inscription_number"`
	InscriptionList       []string `json:"inscription_list"`
	InscriptionNumberList []uint64 `json:"inscription_number_list"`
	Amount                uint64   `json:"send_amount"`
	Address               string   `json:"address"`
	Receiver              string   `json:"receiver"`
}

type SubmitTx struct {
	Txs map[string]string `json:"txs"`
}
