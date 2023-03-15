package response

type DexBTCListingOrderInfo struct {
	RawPSBT      string `json:"raw_psbt"`
	Buyable      bool   `json:"buyable"`
	SellVerified bool   `json:"sell_verified"`
	PriceBTC     uint64 `json:"priceBTC"`
	PriceETH     string `json:"priceETH"`
	OrderID      string `json:"orderID"`
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
	OrderID         string `json:"order_id"`
	ETHAddress      string `json:"eth_address"`
	ETHAmount       string `json:"eth_amount"`
	ExpiredAt       int64  `json:"expired_at"`
	ETHAmountOrigin string `json:"eth_amount_origin"`
	ETHFee          string `json:"eth_fee"`
}

type DEXBuyEthHistory struct {
	ID             string `bson:"id" json:"id"`
	CreatedAt      int64  `bson:"created_at" json:"created_at"`
	OrderID        string `bson:"order_id" json:"order_id"`
	InscriptionID  string `bson:"inscription_id" json:"inscription_id"`
	AmountBTC      uint64 `bson:"amount_btc" json:"amount_btc"`
	AmountETH      string `bson:"amount_eth" json:"amount_eth"`
	UserID         string `bson:"user_id" json:"user_id"`
	ReceiveAddress string `bson:"receive_address" json:"receive_address"`
	RefundAddress  string `bson:"refund_address" json:"refund_address"`
	ExpiredAt      int64  `bson:"expired_at" json:"expired_at"`
	BuyTx          string `bson:"buy_tx" json:"buy_tx"`
	RefundTx       string `bson:"refund_tx" json:"refund_tx"`
	FeeRate        uint64 `bson:"fee_rate" json:"fee_rate"`
	Status         string `bson:"status" json:"status"`
}
