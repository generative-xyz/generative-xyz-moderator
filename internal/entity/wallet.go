package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type WalletTrackTx struct {
	BaseEntity            `bson:",inline"`
	Txhash                string   `bson:"txhash"`
	Address               string   `bson:"address"`
	Status                string   `bson:"status"`
	Type                  string   `bson:"type"`
	InscriptionID         string   `bson:"inscription_id"`
	InscriptionNumber     uint64   `bson:"inscription_number"`
	InscriptionList       []string `json:"inscription_list"`
	InscriptionNumberList []uint64 `json:"inscription_number_list"`
	Amount                uint64   `bson:"amount"`
	Receiver              string   `bson:"receiver"`
}

func (u WalletTrackTx) TableName() string {
	return utils.WALLET_TRACK_TX
}

func (u WalletTrackTx) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
