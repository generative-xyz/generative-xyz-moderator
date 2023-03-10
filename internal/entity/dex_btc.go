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
	CancelAt  *time.Time `bson:"cancel_at"`
	Cancelled bool       `bson:"cancelled"`
	CancelTx  string     `bson:"cancel_tx"`
	Inputs    []string   `bson:"inputs"`
	Matched   bool       `bson:"matched"`
	MatchedTx string     `bson:"matched_tx"`
	MatchAt   *time.Time `bson:"matched_at"`
	Buyer     string     `bson:"buyer"`
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
	BaseEntity `bson:",inline"`
	OrderID    string             `bson:"order_id" json:"order_id"`
	AmountBTC  uint64             `bson:"amount_btc" json:"amount_btc"`
	AmountETH  uint64             `bson:"amount_eth" json:"amount_eth"`
	UserID     string             `bson:"user_id" json:"user_id"`
	Txhash     string             `bson:"txhash" json:"txhash"`
	BTCTx      string             `bson:"btc_tx" json:"btc_tx"`
	UXTOList   []string           `bson:"uxto_list" json:"uxto_list"`
	BuyTx      string             `bson:"buy_tx" json:"buy_tx"`
	RefundTx   string             `bson:"refund_tx" json:"refund"`
	FeeRate    uint64             `bson:"fee_rate" json:"fee_rate"`
	Status     DexBTCETHBuyStatus `bson:"status" json:"status"`
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
)
