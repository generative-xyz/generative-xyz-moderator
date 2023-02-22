package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type CollectionInscription struct {
	BaseEntity`bson:",inline"`
	ID   string `json:"id" bson:"id"`
	Meta struct {
		Name       string `bson:"name" json:"name"`
		Status     string `bson:"status" json:"status"`
		Rank       int    `bson:"rank" json:"rank"`
		Attributes []struct {
			TraitType string `bson:"trait_type" json:"trait_type"`
			Value     string `bson:"value" json:"value"`
			Status    string `bson:"status" json:"status"`
			Percent   string `bson:"percent" json:"percent"`
		} `bson:"attributes" json:"attributes"`
	} `bson:"meta" json:"meta"`
	CollectionInscriptionIcon string `bson:"collection_inscription_icon"`
	TokenCreated bool `bson:"token_created"`
}

func (u CollectionInscription) TableName() string { 
	return utils.COLLECTION_COLLECTION_INSCRIPTION
}

func (u CollectionInscription) ToBson()  (*bson.D, error) { 
	return helpers.ToDoc(u)
}
