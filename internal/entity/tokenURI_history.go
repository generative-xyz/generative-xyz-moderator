package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type TokenHistoryType string
const (
	BLANCE TokenHistoryType = "balance"
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
	ProccessID string `bson:"processID"` //map to eth_wallet_adress or btc_wallet_address
	TokenID string `bson:"token_id" json:"tokenID"`
	MinterAddress string `bson:"minter_address"`
	Owner string `bson:"owner"`
	Action TokenHistoryType `bson:"action"`
	Commit string `bson:"commit"`
	Reveal string `bson:"reveal"`
	Fees int `bson:"fees"`
	ProjectID string `bson:"projectID"`
	Type TokenPaidType `bson:"type"`
	TraceID string `bson:"traceID"`
	Balance string `bson:"balance"`
	Amount string `bson:"amount"`
}

type TokenUriHistoriesArr struct {
	BaseEntity `bson:",inline"`
	ProccessID string `bson:"processID"` //map to eth_wallet_adress or btc_wallet_address
	TokenID string `bson:"token_id" json:"tokenID"`
	MinterAddress string `bson:"minter_address"`
	Owner string `bson:"owner"`
	Action TokenHistoryType `bson:"action"`
	Commit string `bson:"commit"`
	Reveal string `bson:"reveal"`
	Fees int `bson:"fees"`
	ProjectID string `bson:"projectID"`
	Type TokenPaidType `bson:"type"`
	TraceID string `bson:"traceID"`
	Balance string `bson:"balance"`
	Token []TokenUri `bson:"token"`
}

func (u TokenUriHistories) TableName() string { 
	return utils.COLLECTION_TOKEN_URI_HISTORIES
}

func (u TokenUriHistories) ToBson()  (*bson.D, error) { 
	return helpers.ToDoc(u)
}
