package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type Projects struct {
	BaseEntity`bson:",inline"`
	ContractAddress string `bson:"contractAddress"`
	TokenID string `bson:"tokenID"`
	Tags []string `bson:"tags"`
	Categories []string `bson:"categories"`
}

type FilterProjects struct {
	BaseFilters
}

func (u Projects) TableName() string { 
	return utils.COLLECTION_PROJECTS
}

func (u Projects) ToBson()  (*bson.D, error) { 
	return helpers.ToDoc(u)
}