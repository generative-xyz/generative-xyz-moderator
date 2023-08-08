package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

// Convert key for the specific projectID
type DiscordCustomTypes struct {
	BaseEntity  `bson:",inline"`
	ProjectID   int               `json:"project_id" bson:"project_id"`
	CustomTypes map[string]string `json:"custom_types" bson:"custom_types"`
}

func (u DiscordCustomTypes) TableName() string {
	return utils.COLLECTION_DISCORD_CUTOM_TYPES
}

func (u DiscordCustomTypes) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
