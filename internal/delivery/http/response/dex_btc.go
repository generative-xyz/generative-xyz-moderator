package response

import (
	"rederinghub.io/internal/entity"
)

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

type DEXBuyEthHistory struct {
	ID             string                    `bson:"id" json:"id"`
	OrderID        string                    `bson:"order_id" json:"order_id"`
	AmountETH      string                    `bson:"amount_eth" json:"amount_eth"`
	UserID         string                    `bson:"user_id" json:"user_id"`
	ReceiveAddress string                    `bson:"receive_address" json:"receive_address"`
	RefundAddress  string                    `bson:"refund_address" json:"refund_address"`
	ExpiredAt      int64                     `bson:"expired_at" json:"expired_at"`
	BuyTx          string                    `bson:"buy_tx" json:"buy_tx"`
	RefundTx       string                    `bson:"refund_tx" json:"refund_tx"`
	FeeRate        uint64                    `bson:"fee_rate" json:"fee_rate"`
	Status         entity.DexBTCETHBuyStatus `bson:"status" json:"status"`
}
