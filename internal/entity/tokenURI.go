package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type TokenUri struct {
	BaseEntity `bson:",inline"`
	TokenID string `bson:"token_id" json:"token_id"`
	TokenIDInt int `bson:"token_id_int" json:"token_id_int"`
	ContractAddress string `bson:"contract_address" json:"contract_address"`
	Name string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
	Image string `bson:"image" json:"image"`
	ParsedImage *string `bson:"parsed_image" json:"parsed_image"`
	AnimationURL string `bson:"animation_url" json:"animation_url"`
	Attributes string `bson:"attributes" json:"attributes"`
	ParsedAttributes []TokenUriAttr `bson:"parsed_attributes" json:"parsed_attributes"`
	ProjectID string `bson:"project_id" json:"project_id"`
	BlockNumberMinted *string `bson:"block_number_minted" json:"block_number_minted"`
	MintedTime *time.Time `bson:"minted_time" json:"minted_time"` 
	GenNFTAddr string `bson:"gen_nft_addrress"`

	OwnerAddr string `bson:"owner_addrress"`
	CreatorAddr string `bson:"creator_address"`
	Owner *Users `bson:"-"`
	Project *Projects `bson:"-"`
	Creator *Users `bson:"-"`
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
