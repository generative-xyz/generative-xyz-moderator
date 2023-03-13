package request

type CreateDexBTCListing struct {
	RawPSBT       string `json:"raw_psbt"`
	SplitTx       string `json:"split_tx"`
	InscriptionID string `json:"inscription_id"`
}

type CancelDexBTCListing struct {
	Txhash        string `json:"txhash"`
	InscriptionID string `json:"inscription_id"`
	OrderID       string `json:"order_id"`
}

type SubmitDexBTCBuy struct {
	Txhash        string `json:"txhash"`
	InscriptionID string `json:"inscription_id"`
	OrderID       string `json:"order_id"`
}

// type SubmitDexBTCBuyETH struct {
// 	Txhash  string `json:"txhash"`
// 	OrderID string `json:"order_id"`
// 	FeeRate uint64 `json:"fee_rate"`
// }

type GenDexBTCBuyETH struct {
	Amount  uint64 `json:"amount"`
	OrderID string `json:"order_id"`
	FeeRate uint64 `json:"fee_rate"`
}
