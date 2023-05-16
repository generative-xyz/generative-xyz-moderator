package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type CachedGMDashBoard struct {
	BaseEntity `bson:",inline"`
	OldValue   interface{} `bson:"old_value"`
	Value      interface{} `bson:"value"`
	Key        string      `bson:"key"`
}

func (u CachedGMDashBoard) TableName() string {
	return utils.COLLECTION_CACHED_GM_DASHBOARD
}

func (u CachedGMDashBoard) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
