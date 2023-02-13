package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type BTCWalletAddressV2 struct {
	BaseEntity`bson:",inline"`
	UserAddress string `bson:"user_address"` //user's wallet address from FE
	OriginUserAddress string `bson:"origin_user_address"` //user's wallet address from FE
	Amount string `bson:"amount"`
	MintFee string `bson:"mint_fee"`
	SentTokenFee string `bson:"sent_token_fee"`
	OrdAddress string `bson:"ordAddress"` // address is generated from ORD service, which receive all amount
	FileURI string `bson:"fileURI"` // FileURI will be mount if OrdAddress get all amount
	IsConfirm bool  `bson:"isConfirm"` //default: false, if OrdAddress get all amount it will be set true
	InscriptionID string `bson:"inscriptionID"` // tokenID in ETH
	Mnemonic string `bson:"mnemonic"` 
	IsMinted bool  `bson:"isMinted"`//default: false. If InscriptionID exist which means token is minted, it's true
	MintResponse MintStdoputResponse `bson:"mintResponse"` // after token has been mint
	Balance string `bson:"balance"` // balance after check
	FeeRate int32	`bson:"fee_rate"`
}

func (u BTCWalletAddressV2) TableName() string { 
	return utils.COLLECTION_BTC_WALLET_ADDRESS_V2
}

func (u BTCWalletAddressV2) ToBson()  (*bson.D, error) { 
	return helpers.ToDoc(u)
}
