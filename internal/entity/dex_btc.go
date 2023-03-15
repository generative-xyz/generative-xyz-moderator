package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type DexBtcProjectInfo struct {
	ProjectID string `bson:"project_id"`
}

type DexBTCListing struct {
	BaseEntity    `bson:",inline"`
	RawPSBT       string `bson:"raw_psbt"`
	SplitTx       string `bson:"split_tx"`
	InscriptionID string `bson:"inscription_id"`
	Amount        uint64 `bson:"amount"`
	// InscriptionOutputValue uint64     `bson:"inscription_output_value"`
	SellerAddress string `bson:"seller_address"`
	Verified      bool   `bson:"verified"`
	// IsValid       bool       `bson:"is_valid"`
	CancelAt       *time.Time `bson:"cancel_at"`
	Cancelled      bool       `bson:"cancelled"`
	CancelTx       string     `bson:"cancel_tx"`
	Inputs         []string   `bson:"inputs"`
	Matched        bool       `bson:"matched"`
	MatchedTx      string     `bson:"matched_tx"`
	MatchAt        *time.Time `bson:"matched_at"`
	Buyer          string     `bson:"buyer"`
	InvalidMatch   bool       `bson:"invalid_match"`
	InvalidMatchTx string     `bson:"invalid_match_tx"`

	CreatedVerifiedActivity bool `bson:"created_verified_activity"`
	CreatedCancelledActivity bool `bson:"created_cancelled_activity"`
	CreatedMatchedActivity bool `bson:"created_matched_activity"`
}

type DexBtcListingWithProjectInfo struct {
	DexBTCListing
	ProjectInfo []DexBtcProjectInfo `bson:"project_info"`
}

type GetDexBtcListingWithProjectInfoReq struct {
	Page  int64
	Limit int64
}

func (u DexBTCListing) TableName() string {
	return utils.COLLECTION_DEX_BTC_LISTING
}

func (u DexBTCListing) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}

type DexBTCBuyWithETH struct {
	BaseEntity     `bson:",inline"`
	OrderID        string             `bson:"order_id" json:"order_id"`
	InscriptionID  string             `bson:"inscription_id" json:"inscription_id"`
	AmountBTC      uint64             `bson:"amount_btc" json:"amount_btc"`
	Confirmation   int                `bson:"confirmation" json:"confirmation" `
	AmountETH      string             `bson:"amount_eth" json:"amount_eth"`
	UserID         string             `bson:"user_id" json:"user_id"`
	ReceiveAddress string             `bson:"receive_address" json:"receive_address"`
	RefundAddress  string             `bson:"refund_address" json:"refund_address"`
	ExpiredAt      time.Time          `bson:"expired_at" json:"expired_at"`
	BuyTx          string             `bson:"buy_tx" json:"buy_tx"`
	RefundTx       string             `bson:"refund_tx" json:"refund_tx"`
	MasterTx       string             `bson:"master_tx" json:"master_tx"`
	FeeRate        uint64             `bson:"fee_rate" json:"fee_rate"`
	Status         DexBTCETHBuyStatus `bson:"status" json:"status"`
	ETHKey         string             `bson:"eth_key" json:"eth_key"`
}

func (u DexBTCBuyWithETH) TableName() string {
	return utils.COLLECTION_DEX_BTC_BUY_ETH
}

func (u DexBTCBuyWithETH) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}

type DexBTCETHBuyStatus int

const (
	StatusDEXBuy_Pending DexBTCETHBuyStatus = iota // 0: pending: waiting for fund
	StatusDEXBuy_ReceivedFund
	StatusDEXBuy_Buying
	StatusDEXBuy_Bought
	StatusDEXBuy_WaitingToRefund
	StatusDEXBuy_Refunding
	StatusDEXBuy_Refunded
	StatusDEXBuy_SendingMaster
	StatusDEXBuy_SentMaster
	StatusDEXBuy_Expired
)

var StatusDexBTCETHToText = map[DexBTCETHBuyStatus]string{
	StatusDEXBuy_Pending:         "Waiting for payment",
	StatusDEXBuy_ReceivedFund:    "Received payment",
	StatusDEXBuy_Buying:          "Buying",
	StatusDEXBuy_Bought:          "Bought",
	StatusDEXBuy_WaitingToRefund: "Waiting to refund",
	StatusDEXBuy_Refunding:       "Refunding",
	StatusDEXBuy_Refunded:        "Refunded",
	StatusDEXBuy_SendingMaster:   "Sending to master",
	StatusDEXBuy_SentMaster:      "Sent to master",
	StatusDEXBuy_Expired:         "Expired",
}
