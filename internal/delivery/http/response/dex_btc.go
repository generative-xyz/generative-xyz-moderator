package response

type DexBTCListingOrderInfo struct {
	RawPSBT string `json:"raw_psbt"`
}

type DexBTCHistoryListing struct {
	OrderID       string `json:"order_id"`
	Type          string `json:"type"`
	Timestamp     int64  `json:"timestamp"`
	InscriptionID string `json:"inscription_id"`
	Txhash        string `json:"txhash"`
	Amount        string `json:"amount"`
}

type GenDexBTCBuyETH struct {
	OrderID    string `json:"order_id"`
	ETHAddress string `json:"eth_address"`
	ETHAmount  string `json:"eth_amount"`
	ExpiredAt  int64  `json:"expired_at"`
}
