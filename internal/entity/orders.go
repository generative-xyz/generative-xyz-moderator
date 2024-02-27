package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type OrderStatus int

const (
	Order_Pending OrderStatus = iota // 0: pending: waiting for payment
	Order_Paid                       // 1:paid
	Order_Cancel                     // 2: cancel

)

type OrdersAddress struct {
	BaseEntityNoID `bson:",inline"`
	Status         OrderStatus `bson:"status" json:"status"`
	Address        string      `bson:"address" json:"address"`          //payment address
	AddressType    string      `bson:"address_type" json:"addressType"` //payment address_type
	PrivateKey     string      `bson:"private_key" json:"-"`            //payment private_key
	OrderID        string      `bson:"order_id" json:"orderID"`
	Amount         string      `bson:"-" json:"-"`
}

func (u OrdersAddress) TableName() string {
	return utils.ORDERS
}

func (u OrdersAddress) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
