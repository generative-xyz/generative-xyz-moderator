package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

const (
	StatusWithdraw_Available        = iota  // 0: pending: waiting for approve
	StatusWithdraw_Pending       
	StatusWithdraw_Approve
	StatusWithdraw_Reject
)

type Withdrawtype string 

const (
	WithDrawProject Withdrawtype = Withdrawtype(`project`)
	WithDrawReferal Withdrawtype =  Withdrawtype(`referal`)
)


type FilterWithdraw struct {
	BaseFilters
	PaymentType *string
	WithdrawItemID *string
	WithdrawItemIDs []string
	WalletAddress *string
	WithdrawType *string
	Status *int
	Statuses []int
}

type Withdraw struct {
	BaseEntity              `bson:",inline" json:"-"`
	Amount       string `bson:"amount" json:"amount"`
	PayType  string `bson:"payType" json:"payType"`
	Status int `bson:"status" json:"status"`
	WalletAddress string `bson:"walletAddress" json:"walletAddress"`
	WithdrawFrom string `bson:"withdrawFrom" json:"withdrawFrom"`
	EarningReferal string `bson:"earningReferal" json:"earningReferal"`
	EarningVolume string `bson:"earningVolume" json:"earningVolume"`
	TotalEarnings string `bson:"totalEarnings" json:"totalEarnings"`
	AvailableBalance string `bson:"availableBalance" json:"availableBalance"`
	WithdrawType Withdrawtype `bson:"withdrawType" json:"withdrawType"`
	WithdrawItemID string `bson:"withdrawItemID"  json:"withdrawItemID"`
	User  WithdrawUserInfo `bson:"user"`
}

type WithdrawUserInfo struct {
	WalletAddress *string `bson:"walletAddress"`
	WalletAddressPayment *string `bson:"walletAddressPayment"`
	WalletAddressBTC *string `bson:"walletAddressBTC"`
	DisplayName *string  `bson:"displayName"`
	Avatar   *string  `bson:"avatar"`
}


func (u Withdraw) TableName() string {
	return utils.COLLECTION_WITHDRAW
}

func (u Withdraw) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
