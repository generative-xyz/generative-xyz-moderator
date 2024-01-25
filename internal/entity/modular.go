package entity

import (
	"rederinghub.io/utils"
	"time"
)

type Total struct {
	Total int64 `json:"total" bson:"total"`
}

type ModularTokenUri struct {
	BaseEntityNoID      `bson:",inline"`
	TokenID             string            `bson:"token_id" json:"token_id"`
	Name                string            `bson:"name" json:"name"`
	Description         string            `bson:"description" json:"description"`
	Image               string            `bson:"image" json:"image"`
	ParsedImage         *string           `bson:"parsed_image" json:"parsed_image"`
	AnimationURL        string            `bson:"animation_url" json:"animation_url"`
	ParsedAttributesStr []TokenUriAttrStr `bson:"parsed_attributes_str" json:"attributes"`
	ProjectID           string            `bson:"project_id" json:"project_id"`
	BlockNumberMinted   *string           `bson:"block_number_minted" json:"block_number_minted"`
	MintedTime          *time.Time        `bson:"minted_time" json:"minted_time"`
	Thumbnail           string            `bson:"thumbnail" json:"thumbnail"`

	OwnerAddr string `bson:"owner_addrress" json:"owner_address"`
	//accept duplicated data to query more faster
	Owner   *ModularUsers    `bson:"owner" json:"owner"`
	Project *ModularProjects `bson:"project" json:"project"`
}

type ModularProjects struct {
	BaseEntity      `bson:",inline" json:"-"`
	TokenID         string `bson:"tokenid" json:"token_id"`
	Name            string `bson:"name" json:"name"`
	CreatorName     string `bson:"creatorName" json:"creator_name"`
	CreatorAddrr    string `bson:"creatorAddress" json:"creator_address"`
	CreatorAddrrBTC string `bson:"creatorAddrrBTC" json:"creator_address_btc"`
	Description     string `bson:"description" json:"description"`
	Thumbnail       string `bson:"thumbnail" json:"thumbnail"`
	GenNFTAddr      string `bson:"genNFTAddr" json:"gen_nft_address"`
}

type ModularUsers struct {
	BaseEntity `bson:",inline" json:"-"`
	// ID                      string        `bson:"id" json:"id,omitempty"`
	WalletAddress           string `bson:"wallet_address" json:"wallet_address,omitempty"`                         // eth wallet define user in platform by connect wallet and sign
	WalletAddressPayment    string `bson:"wallet_address_payment" json:"wallet_address_payment,omitempty"`         // eth wallet artist receive royalty
	WalletAddressBTC        string `bson:"wallet_address_btc" json:"wallet_address_btc,omitempty"`                 // btc wallet artist receive royalty
	WalletAddressBTCTaproot string `bson:"wallet_address_btc_taproot" json:"wallet_address_btc_taproot,omitempty"` // btc wallet receive minted nft
	DisplayName             string `bson:"display_name" json:"display_name,omitempty"`
	Bio                     string `bson:"bio" json:"bio,omitempty"`
	Avatar                  string `bson:"avatar" json:"avatar"`
	Banner                  string `bson:"banner" json:"banner"`
	IsUpdatedAvatar         *bool  `bson:"is_updated_avatar" json:"is_updated_avatar,omitempty"`
}

type ModularInscription struct {
	BaseEntityNoID `bson:",inline"`
	IsCreatedToken bool   `bson:"is_created_token" json:"is_created_token"`
	InscriptionID  string `bson:"inscription_id" json:"inscription_id"`
	BlockHeight    uint64 `bson:"block_height" json:"block_height"`
}

func (g *ModularInscription) TableName() string {
	return utils.COLLECTION_MODULAR_INSCRIPTION
}
