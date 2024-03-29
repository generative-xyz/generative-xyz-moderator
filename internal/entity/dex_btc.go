package entity

import (
	"encoding/json"
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

	CreatedVerifiedActivity  bool `bson:"created_verified_activity"`
	CreatedCancelledActivity bool `bson:"created_cancelled_activity"`
	CreatedMatchedActivity   bool `bson:"created_matched_activity"`

	IsTimeSeriesData bool `json:"is_time_series_data"`

	FromOtherMkp bool `bson:"from_other_mkp"`
}

type DexBtcListingWithProjectInfo struct {
	DexBTCListing `bson:",inline"`
	ProjectInfo   []DexBtcProjectInfo `bson:"project_info"`
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
	BaseEntity      `bson:",inline"`
	OrderID         string             `bson:"order_id" json:"order_id"`
	InscriptionID   string             `bson:"inscription_id" json:"inscription_id"`
	AmountBTC       uint64             `bson:"amount_btc" json:"amount_btc"`
	Confirmation    int                `bson:"confirmation" json:"confirmation" `
	AmountETH       string             `bson:"amount_eth" json:"amount_eth"`
	UserID          string             `bson:"user_id" json:"user_id"`
	ReceiveAddress  string             `bson:"receive_address" json:"receive_address"`
	RefundAddress   string             `bson:"refund_address" json:"refund_address"`
	ExpiredAt       time.Time          `bson:"expired_at" json:"expired_at"`
	BuyTx           string             `bson:"buy_tx" json:"buy_tx"`
	RefundTx        string             `bson:"refund_tx" json:"refund_tx"`
	MasterTx        string             `bson:"master_tx" json:"master_tx"`
	SplitTx         string             `bson:"split_tx" json:"split_tx"`
	FeeRate         uint64             `bson:"fee_rate" json:"fee_rate"`
	Status          DexBTCETHBuyStatus `bson:"status" json:"status"`
	ETHKey          string             `bson:"eth_key" json:"eth_key"`
	ETHAddress      string             `bson:"eth_address" json:"eth_address"`
	IsMultiBuy      bool               `bson:"multi_buy" json:"multi_buy"`
	InscriptionList []string           `bson:"inscription_list" json:"inscription_list"`
	SellOrderList   []string           `bson:"sell_order_list" json:"order_list"`
}

func (u DexBTCBuyWithETH) ToJsonString() string {
	dataBytes, _ := json.Marshal(u)
	return string(dataBytes)
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
	StatusDEXBuy_ReceivedFund:    "Payment received",
	StatusDEXBuy_Buying:          "Buying",
	StatusDEXBuy_Bought:          "Bought",
	StatusDEXBuy_WaitingToRefund: "Waiting to refund",
	StatusDEXBuy_Refunding:       "Refunding",
	StatusDEXBuy_Refunded:        "Refunded",
	StatusDEXBuy_SendingMaster:   "Sending to master",
	StatusDEXBuy_SentMaster:      "Sent to master",
	StatusDEXBuy_Expired:         "Expired",
}

type DexBTCTrackingInternal struct {
	BaseEntity `bson:",inline"`
	// OrderID       string                       `bson:"order_id" json:"order_id"`
	// BuyEthOrder   string                       `bson:"buy_eth_order" json:"buy_eth_order"`
	InscriptionList []string                     `bson:"inscription_list" json:"inscription_list"`
	Txhash          string                       `bson:"txhash" json:"txhash"`
	Status          DexBTCTrackingInternalStatus `bson:"status" json:"status"`
}

func (u DexBTCTrackingInternal) TableName() string {
	return utils.COLLECTION_DEX_BTC_TRACKING_INTERNAL
}

func (u DexBTCTrackingInternal) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}

type DexBTCTrackingInternalStatus int

const (
	StatusDEXBTCTracking_Pending DexBTCTrackingInternalStatus = iota // 0: pending
	StatusDEXBTCTracking_Success                                     // 1: successful
	StatusDEXBTCTracking_Failed                                      // 2: failed
)

type Report2ndSale struct {
	Amount           float64 `bson:"total_amount" json:"total_amount"`
	AmountUSD        float64 `bson:"-" json:"amountUSD"`
	WalletAddressBTC string  `bson:"walletAddressBtc" json:"walletAddressBtc"`
	WalletAddress    string  `bson:"-" json:"walletAddress"`
}
