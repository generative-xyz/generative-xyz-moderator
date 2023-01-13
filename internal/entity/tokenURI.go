package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type FilterTokenUris struct {
	BaseFilters
	ContractAddress *string
	OwnerAddr *string
	CreatorAddr *string
	GenNFTAddr *string
	Keyword *string
	CollectionIDs []string
	TokenIDs []string
}

type TokenUri struct {
	BaseEntityNoID `bson:",inline"`
	TokenID string `bson:"token_id" json:"token_id"`
	TokenIDInt int `bson:"token_id_int" json:"token_id_int"`
	TokenIDMini *int `bson:"token_id_mini" json:"token_id_mini"`
	ContractAddress string `bson:"contract_address" json:"contract_address"`
	Name string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
	Image string `bson:"image" json:"image"`
	ParsedImage *string `bson:"parsed_image" json:"parsed_image"`
	AnimationURL string `bson:"animation_url" json:"animation_url"`
	Attributes string `bson:"attributes" json:"attributes"`
	ParsedAttributes []TokenUriAttr `bson:"parsed_attributes" json:"parsed_attributes"`
	ProjectID string `bson:"project_id" json:"project_id"`
	ProjectIDInt int64 `bson:"project_id_int" json:"project_id_int"`
	BlockNumberMinted *string `bson:"block_number_minted" json:"block_number_minted"`
	MintedTime *time.Time `bson:"minted_time" json:"minted_time"` 
	GenNFTAddr string `bson:"gen_nft_addrress"`
	Thumbnail string `bson:"thumbnail"`

	OwnerAddr string `bson:"owner_addrress"`
	CreatorAddr string `bson:"creator_address"`
	Priority *int `bson:"priority"`

	//accept duplicated data to query more faster
	Owner *Users `bson:"owner"`
	Project *Projects `bson:"project"`
	Creator *Users `bson:"creator"`
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
