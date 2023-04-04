package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils/helpers"
)

type Faucet struct {
	BaseEntity `bson:",inline"`
	Tx         string `bson:"tx"`
	Address    string `bson:"address"`
	Status     int    `bson:"status"` // 0 pending, 1 success, 2 fail, -1: init
	Amount     string `bson:"amount"`
}

func (u Faucet) TableName() string {
	return "faucets"
}

func (u Faucet) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
