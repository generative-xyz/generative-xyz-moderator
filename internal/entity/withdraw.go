package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

const (
	StatusWithdraw_Pending        = iota  // 0: pending: waiting for approve
	StatusWithdraw_Approve
	StatusWithdraw_Done
	StatusWithdraw_Reject
)


type FilterWithdraw struct {
	BaseFilters
	PaymentType *string
	ProjectID *string
	ProjectIDs []string
	WalletAddress *string
	Status *int
}

type Withdraw struct {
	BaseEntity              `bson:",inline" json:"-"`
	Amount       string `bson:"amount" json:"amount"`
	PayType  string `bson:"payType" json:"payType"`
	ProjectID string `bson:"projectID"  json:"projectID"`
	Status int `bson:"status" json:"status"`
	WalletAddress string `bson:"walletAddress" json:"walletAddress"`
	WithdrawFrom string `bson:"withdrawFrom" json:"withdrawFrom"`
	EarningReferal string `bson:"earningReferal" json:"earningReferal"`
	EarningVolume string `bson:"earningVolume" json:"earningVolume"`
	TotalEarnings string `bson:"totalEarnings" json:"totalEarnings"`
}

func (u Withdraw) TableName() string {
	return utils.COLLECTION_WITHDRAW
}

func (u Withdraw) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
