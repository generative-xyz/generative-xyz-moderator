package response

type DexBTCListingOrderInfo struct {
	RawPSBT string `json:"raw_psbt"`
}

type DexBTCHistoryListing struct {
	OrderID       string `json:"order_id"`
	Type          string `json:"type"`
	Timestamp     uint64 `json:"timestamp"`
	InscriptionID string `json:"inscription_id"`
	Txhash        string `json:"txhash"`
	Amount        string `json:"amount"`
}
