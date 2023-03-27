package request

type CreateDexBTCListing struct {
	RawPSBT       string `json:"raw_psbt"`
	SplitTx       string `json:"split_tx"`
	InscriptionID string `json:"inscription_id"`
	Address       string `json:"address"`
}

type CancelDexBTCListing struct {
	Txhash        string `json:"txhash"`
	InscriptionID string `json:"inscription_id"`
	OrderID       string `json:"order_id"`
	Address       string `json:"address"`
}

// type SubmitDexBTCBuy struct {
// 	Txhash        string `json:"txhash"`
// 	InscriptionID string `json:"inscription_id"`
// 	OrderID       string `json:"order_id"`
// }

// type SubmitDexBTCBuyETH struct {
// 	Txhash  string `json:"txhash"`
// 	OrderID string `json:"order_id"`
// 	FeeRate uint64 `json:"fee_rate"`
// }

type GenDexBTCBuyETH struct {
	OrderID        string   `json:"order_id"`
	FeeRate        uint64   `json:"fee_rate"`
	ReceiveAddress string   `json:"receive_address"`
	RefundAddress  string   `json:"refund_address"`
	IsEstimate     bool     `json:"is_estimate"`
	OrderIDList    []string `json:"order_list"`
}

type UpdateDexBTCBuyETHTx struct {
	Txhash  string `json:"txhash"`
	OrderID string `json:"order_id"`
}

type RetrieveBTCListingOrdersInfo struct {
	OrderList []string `json:"order_list"`
}

type SubmitOWPurchaseTx struct {
	Address        string `json:"address"`
	InscriptionID  string `json:"inscription_id"`
	PurchaseSigned string `json:"purchase_raw`
	SetupSigned    string `json:"setup_raw`
}
