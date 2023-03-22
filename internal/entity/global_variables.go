package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)


type GlobalVariable struct {
	BaseEntity`bson:",inline"`
	Key        string `bson:"key"`
	Value      interface{} `bson:"value"`
}

func (u GlobalVariable) TableName() string { 
	return utils.COLLECTION_GLOBAL_VARIABLE
}

func (u GlobalVariable) ToBson()  (*bson.D, error) { 
	return helpers.ToDoc(u)
}
