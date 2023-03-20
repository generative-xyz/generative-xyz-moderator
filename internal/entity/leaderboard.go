package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type LeaderBoard struct {
	BaseEntity        `bson:",inline"`
	UserWalletAddress string `bson:"walletAddress" json:"walletAddress"`
	AlphaBoost        int    `bson:"alphaBoost" json:"alphaBoost"`
	Points            int    `bson:"points" json:"points"`
}

func (u LeaderBoard) TableName() string {
	return utils.COLLECTION_LEADERBOARD
}

func (u LeaderBoard) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
