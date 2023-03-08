package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type TokenMetadataStatus string
const (
	METADATA_UPLOADED TokenHistoryType = "balance"
	METADATA_UPDATED TokenHistoryType = "mint"
)

type FilterTokenUriMetadata struct {
	BaseFilters
	ProjectID []string
	TokenID []string
}

type TokenUriMetadata struct {
	BaseEntity`bson:",inline"`
	ProjectID string `bson:"projectID"`
	UploadedFile string `bson:"uploadedFile"`
	Content []TokenTraits `bson:"content"`
}

type TokenTraits struct {
	ID string `bson:"id" json:"id"`
	Atrributes []TraitAttribute `bson:"atrributes" json:"atrributes"`
}

type TraitAttribute struct {
	TraitType string `bson:"trait_type" json:"trait_type"`
	Value string `bson:"value"  json:"value"`
}

func (u TokenUriMetadata) TableName() string { 
	return utils.COLLECTION_TOKEN_URI_METADATA
}

func (u TokenUriMetadata) ToBson()  (*bson.D, error) { 
	return helpers.ToDoc(u)
}
