package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type ActivityType int
const (
	View ActivityType = 0
	Mint ActivityType = 1
	Buy ActivityType = 2
)

type Activity struct {
	BaseEntity`bson:",inline"`
	Type       ActivityType `bson:"type" json:"type"`   // Type of activity, now it is ViewProject, MintProject and BuyAProject.
	Value      int64        `bson:"value" json:"value"` // Value is used for calculating the mint volumn and marketplace volumn of a project.
	ProjectID  string				`bson:"project_id"`         // ProjectUUID of the activity
	Reference  string				`bson:"reference"`          // Uuid of reference object's uuid (currently, it can only be token's inscriptionID).
}

func (u Activity) TableName() string { 
	return utils.COLLECTION_ACTIVITIES
}

func (u Activity) ToBson()  (*bson.D, error) { 
	return helpers.ToDoc(u)
}
