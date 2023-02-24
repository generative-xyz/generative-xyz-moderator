package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type WalletTrackTx struct {
	BaseEntity `bson:",inline"`
	Txhash     string `json:"txhash"`
	Address    string `json:"address"`
	Status     string `json:"status"`
}

func (u WalletTrackTx) TableName() string {
	return utils.WALLET_TRACK_TX
}

func (u WalletTrackTx) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
