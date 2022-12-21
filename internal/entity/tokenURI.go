package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type TokenUri struct {
	BaseEntity `bson:",inline"`
	TokenID string `bson:"token_id" json:"token_id"`
	ContractAddress string `bson:"contract_address" json:"contract_address"`
	Name string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
	Image string `bson:"image" json:"image"`
	AnimationURL string `bson:"animation_url" json:"animation_url"`
	Attributes string `bson:"attributes" json:"attributes"`
	ParsedAttributes []TokenUriAttr `bson:"parsed_attributes" json:"parsed_attributes"`
}

type TokenUriAttr struct {
	TraitType string `bson:"trait_type" json:"trait_type"`
	Value interface{} `bson:"value" json:"value"`
}

func (u TokenUri) TableName() string { 
	return utils.COLLECTION_TOKEN_URI
}

func (u TokenUri) ToBson()  (*bson.D, error) { 
	return helpers.ToDoc(u)
}