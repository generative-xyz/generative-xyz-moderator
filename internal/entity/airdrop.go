package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type Airdrop struct {
	BaseEntity                `bson:",inline"`
	Tx                        string `bson:"tx"`
	File                      string `bson:"file"`
	Receiver                  string `bson:"receiver"`
	ReceiverBtcAddressTaproot string `bson:"receiverBtcAddressTaproot"`
	Type                      int    `bson:"type"`   // 0 artist, 1 collector
	Status                    int    `bson:"status"` // 0 pending, 1 success, 2 fail
}

func (u Airdrop) TableName() string {
	return utils.COLLECTION_AIRDROP
}

func (u Airdrop) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
