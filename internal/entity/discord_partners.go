package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type DiscordPartner struct {
	BaseEntity    `bson:",inline"`
	Webhooks   map[string]string `bson:"webhooks"`
	ProjectIDs []string          `bson:"project_ids"`
	Name       string            `bson:"name"`
}

func (u DiscordPartner) TableName() string {
	return utils.COLLECTION_DISCORD_PARTNER
}

func (u DiscordPartner) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
