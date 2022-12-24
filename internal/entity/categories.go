package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type Categories struct {
	BaseEntity`bson:",inline"`
	Name string `bson:"name"`
	
}

type FilterCategories struct {
	BaseFilters
	Name *string
	ID *string
}

func (u Categories) TableName() string { 
	return utils.COLLECTION_CATEGORIES
}

func (u Categories) ToBson()  (*bson.D, error) { 
	return helpers.ToDoc(u)
}