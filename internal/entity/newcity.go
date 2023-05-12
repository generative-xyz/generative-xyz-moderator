package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils/helpers"
)

type NewCityGm struct {
	BaseEntity `bson:",inline"`

	UserAddress string `bson:"user_address" json:"userAddress"` // to faucet ...
	ENS         string `bson:"ens" json:"ens"`
	Avatar      string `bson:"avatar" json:"avatar"`

	Address    string `bson:"address" json:"address"`
	PrivateKey string `bson:"private_key" json:"-"`
	KeyVersion int    `bson:"key_version" json:"-"`

	Status int    `bson:"status" json:"status"`
	Type   string `bson:"type" json:"type"`

	NativeAmount []string `bson:"native_amounts" json:"nativeAmounts"`

	TokenAmounts []TokenAmounts `bson:"token_amounts" json:"tokenAmounts"`

	TxNatives []string `bson:"tx_natives" json:"txNatives"`

	TxTokens []string `bson:"tx_tokens" json:"txTokens"`
}

type TokenAmounts struct {
	Token  string `bson:"token" json:"token"`
	Amount string `bson:"amount" json:"amount"`
}

func (u NewCityGm) TableName() string {
	return "new_city_gm"
}

func (u NewCityGm) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
