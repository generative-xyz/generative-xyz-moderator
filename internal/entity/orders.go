package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type OrdersAddress struct {
	BaseEntityNoID `bson:",inline"`
	Address        string `bson:"address" json:"address"`
	AddressType    string `bson:"address_type" json:"addressType"`
	PrivateKey     string `bson:"private_key" json:"-"`
	OrderID        string `bson:"order_id" json:"orderID"`
}

func (u OrdersAddress) TableName() string {
	return utils.ORDERS
}

func (u OrdersAddress) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
