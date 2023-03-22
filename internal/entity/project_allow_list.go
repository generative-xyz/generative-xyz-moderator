package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type AllowedByType string

const (
	ERC20  AllowedByType = "erc20"
	ERC721 AllowedByType = "erc721"
)

type ProjectAllowList struct {
	BaseEntity                  `bson:",inline" json:"-"`
	ProjectID                   string        `bson:"projectID" json:"projectID"`
	UserWalletAddress           string        `bson:"userWalletAddress" json:"userWalletAddress"`
	UserWalletAddressBTCTaproot string        `bson:"userWalletAddressBTCTaproot" json:"userWalletAddressBTCTaproot"`
	AllowedBy                   AllowedByType `bson:"allowedBy" json:"allowedBy"` //to capture that the type allowed this user
}

func (u ProjectAllowList) TableName() string {
	return utils.COLLECTION_PROJECT_ALLOW_LIST
}

func (u ProjectAllowList) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
