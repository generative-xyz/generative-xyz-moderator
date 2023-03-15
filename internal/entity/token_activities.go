package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type TokenActivityType int
const (
	TokenMint TokenActivityType = 0
	TokenListing TokenActivityType = 1
	TokenCancelListing TokenActivityType = 2
	TokenMatched TokenActivityType = 3
)

type TokenActivity struct {
	BaseEntity    `bson:"base_entity" json:"base_entity"`
	Type          TokenActivityType `bson:"type" json:"type"`
	Title         string            `bson:"title" json:"title"`
	UserAAddress  string            `bson:"user_a_address" json:"user_a_address"`
	UserA         *Users            `bson:"-" json:"user_a"`
	UserBAddress  string            `bson:"user_b_address" json:"user_b_address"`
	UserB         *Users            `bson:"-" json:"user_b"`
	Amount        int64             `bson:"amount" json:"amount"`
	Time          *time.Time        `bson:"time" json:"time"`
	InscriptionID string            `bson:"inscription_id" json:"inscription_id"`
	ProjectID     string            `bson:"project_id" json:"project_id"`
}

func (u TokenActivity) TableName() string { 
	return utils.COLLECTION_TOKEN_ACTIVITY
}

func (u TokenActivity) ToBson()  (*bson.D, error) { 
	return helpers.ToDoc(u)
}
