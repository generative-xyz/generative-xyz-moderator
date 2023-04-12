package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type DiscordPartner struct {
	BaseEntity            `bson:",inline"`
	Webhooks              map[string]string `bson:"webhooks"`
	ProjectIDs            []string          `bson:"project_ids"`
	Name                  string            `bson:"name"`
	AmountGreaterThanZero bool              `bson:"greater_than_zero"`
	Categories            []string          `bson:"categories"`
}

func (u DiscordPartner) MatchProject(t string) bool {
	if len(u.ProjectIDs) == 0 {
		return true
	}

	for _, p := range u.ProjectIDs {
		if p == t {
			return true
		}
	}

	return false
}

func (u DiscordPartner) MatchCategory(t string) bool {
	if len(u.Categories) == 0 {
		return true
	}

	for _, p := range u.Categories {
		if p == t {
			return true
		}
	}

	return false
}

func (u DiscordPartner) MatchAmountGreaterThanZero(amount uint64) bool {
	if !u.AmountGreaterThanZero {
		return true
	}

	if u.AmountGreaterThanZero && amount > 0 {
		return true
	}

	return false
}

func (u DiscordPartner) TableName() string {
	return utils.COLLECTION_DISCORD_PARTNER
}

func (u DiscordPartner) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
