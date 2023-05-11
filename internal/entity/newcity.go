package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils/helpers"
)

type NewCityGm struct {
	BaseEntity `bson:",inline"`

	UserAddress string `bson:"user_address" json:"userAddress"` // to faucet ...

	Address    string `bson:"address" json:"address"`
	PrivateKey string `bson:"private_key" json:"-"`

	Status int    `bson:"status" json:"status"`
	Type   string `bson:"type" json:"type"`
}

func (u NewCityGm) TableName() string {
	return "new_city_gm"
}

func (u NewCityGm) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
