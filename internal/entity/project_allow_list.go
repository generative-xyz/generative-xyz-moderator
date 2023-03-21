package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type ProjectAllowList struct {
	BaseEntity        `bson:",inline" json:"-"`
	ProjectID         string `bson:"projectID" json:"projectID"`
	UserWalletAddress string `bson:"userWalletAddress" json:"userWalletAddress"`
}

func (u ProjectAllowList) TableName() string {
	return utils.COLLECTION_PROJECT_ALLOW_LIST
}

func (u ProjectAllowList) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
