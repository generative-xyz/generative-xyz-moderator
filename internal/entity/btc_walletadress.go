package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type BTCWalletAddress struct {
	BaseEntity`bson:",inline"`
	UserAddress string `bson:"user_address"` //user's wallet address from FE
	Amount uint64 `bson:"amount"`
	OrdAddress string `bson:"ordAddress"` // address is generated from ORD service, which receive all amount
	FileURI string `bson:"fileURI"` // FileURI will be mount if OrdAddress get all amount
	IsConfirm bool  `bson:"isConfirm"` //default: false, if OrdAddress get all amount it will be set true
	InscriptionID string `bson:"inscriptionID"` // tokenID in ETH
	
}

type FilterBTCWalletAddress struct {
	BaseFilters
	UserAddress *string
	Amount *string
	OrdAddress *string
	IsConfirm *string
	InscriptionID *string
}

func (u BTCWalletAddress) TableName() string { 
	return utils.COLLECTION_BTC_WALLET_ADDRESS
}

func (u BTCWalletAddress) ToBson()  (*bson.D, error) { 
	return helpers.ToDoc(u)
}