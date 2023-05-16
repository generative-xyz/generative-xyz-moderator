package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
	"time"
)

type CachedGMDashBoard struct {
	BaseEntity `bson:",inline"`
	OldValue   interface{} `bson:"old_value"`
	Value      interface{} `bson:"value"`
	Key        string      `bson:"key"`
}

type AggregatedGMDashBoard struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt    *time.Time         `bson:"created_at" json:"created_at"`
	Usdt         float64            `bson:"usdt" json:"usdt"`
	Contributors int64              `bson:"contributors" json:"contributors"`
}

func (u CachedGMDashBoard) TableName() string {
	return utils.COLLECTION_CACHED_GM_DASHBOARD
}

func (u CachedGMDashBoard) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
