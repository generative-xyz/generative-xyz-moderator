package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type ActionType int
type ObjectType int

const (
	DISLIKE ActionType = 0
	LIKE    ActionType = 1
)

const (
	TokenURI ObjectType = 0
	Project  ObjectType = 1
)

type Actions struct {
	BaseEntity `bson:",inline"`
	ObjectID   string     `bson:"object_id" json:"object_id"`
	ObjectType ObjectType `bson:"object_type" json:"object_type"`
	Action     ActionType `bson:"action" json:"action"`
	CreatedBy  string     `bson:"created_by" json:"created_by"`
}

type FilterActions struct {
	BaseFilters
	Parent     *string
	ObjectID   *string
	ObjectType *ObjectType
	Action     *ActionType
}

func (u Actions) TableName() string {
	return utils.COLLECTION_ACTIONS
}

func (u Actions) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
