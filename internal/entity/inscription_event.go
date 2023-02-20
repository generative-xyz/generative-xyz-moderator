package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type InscriptionEvent struct {
	BaseEntity  `bson:",inline"`
	Inscription string `bson:"inscription" json:"inscription"`
	BlockHeight int    `bson:"block_height" json:"block_height"`
	EventType   string `bson:"event_type" json:"event_type"`
}

func (u InscriptionEvent) TableName() string {
	return utils.INSCRIPTION_EVENT
}

func (u InscriptionEvent) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
