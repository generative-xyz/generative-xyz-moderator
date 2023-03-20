package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type BTCTransactionSubmit struct {
	BaseEntity `bson:",inline"`
	Txhash     string `bson:"txhash" json:"txhash"`
	Raw        string `bson:"raw" json:"raw"`
	Status     int    `bson:"status" json:"status"`
}

func (u BTCTransactionSubmit) TableName() string {
	return utils.COLLECTION_BTC_TX_SUBMIT
}

func (u BTCTransactionSubmit) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
