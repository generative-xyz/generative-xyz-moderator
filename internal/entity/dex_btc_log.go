package entity

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils/helpers"
)

type DexBTCLog struct {
	BaseEntity `bson:",inline"`
	Data       map[string]interface{} `bson:"data"`
	Function   string                 `bson:"function"`
}

func (u DexBTCLog) ToJsonString() string {
	dataBytes, _ := json.Marshal(u)
	return string(dataBytes)
}

func (u DexBTCLog) TableName() string {
	return "dex_btc_log"
}

func (u DexBTCLog) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
