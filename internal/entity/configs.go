package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type Configs struct {
	BaseEntity`bson:",inline"`
	Key string `bson:"key"`
	Value string `bson:"value"`
	
}

type FilterConfigs struct {
	BaseFilters
	Name *string
	UploadedBy *string
}

func (u Configs) TableName() string { 
	return utils.COLLECTION_CONFIGS
}

func (u Configs) ToBson()  (*bson.D, error) { 
	return helpers.ToDoc(u)
}