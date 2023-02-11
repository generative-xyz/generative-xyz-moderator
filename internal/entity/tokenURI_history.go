package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type TokenHistoryType string
const (
	MINT TokenHistoryType = "mint"
	SENT TokenHistoryType = "sent"
)

type FilterTokenUriHistory struct {
	BaseFilters
	TokenID *string
	MinterAddress *string
	Owner *string
}

type TokenUriHistories struct {
	BaseEntity `bson:",inline"`
	TokenID string `bson:"token_id" json:"token_id"`
	MinterAddress string `bson:"minter_address"`
	Owner string `bson:"owner"`
	Action TokenHistoryType `bson:"action"`
	Commit string `bson:"commit"`
	Reveal string `bson:"reveal"`
	Fees int `bson:"reveal"`
}

func (u TokenUriHistories) TableName() string { 
	return utils.COLLECTION_TOKEN_URI_HISTORIES
}

func (u TokenUriHistories) ToBson()  (*bson.D, error) { 
	return helpers.ToDoc(u)
}
