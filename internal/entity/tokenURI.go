package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type TokenPaidType string

const (
	ETH TokenPaidType = "eth"
	BIT TokenPaidType = "btc"
)

type InscriptionDetail struct {
	Chain         string                 `json:"chain"`
	GenesisFee    int64                  `json:"genesis_fee"`
	GenesisHeight int64                  `json:"genesis_height"`
	Address       string                 `json:"address"`
	ContentType   string                 `json:"content_type"`
	InscriptionId string                 `json:"inscription_id"`
	Next          string                 `json:"next"`
	Number        int64                  `json:"number"`
	Output        map[string]interface{} `json:"output"`
	Previous      string                 `json:"previous"`
	Sat           int64                  `json:"sat"`
	Satpoint      string                 `json:"satpoint"`
	Timestamp     string                 `json:"timestamp"`
}

type TokenUriAttrFilter struct {
	TraitType string
	Values    []string
}

type TokenUriListingPage struct {
	TotalData  []TokenUriListingFilter `bson:"totalData" json:"totalData"`
	TotalCount []struct {
		Count int64 `bson:"count" json:"count"`
	} `bson:"totalCount" json:"totalCount"`
}

type TokenUriListingVolume struct {
	TotalAmount uint64 `bson:"totalAmount" json:"totalAmount"`
}

type FilterTokenUris struct {
	BaseFilters
	ContractAddress  *string
	OwnerAddr        *string
	CreatorAddr      *string
	GenNFTAddr       *string
	Keyword          *string
	Search           *string
	CollectionIDs    []string
	TokenIDs         []string
	Attributes       []TokenUriAttrFilter
	RarityAttributes []TokenUriAttrFilter
	HasPrice         *bool
	FromPrice        *int64
	ToPrice          *int64
	Ids              []string
	IsBuynow         *bool
}

type TokenStats struct {
	PriceInt *int64 `bson:"price_int,omitempty" json:"price_int,omitempty"`
}

type TokenUri struct {
	BaseEntityNoID      `bson:",inline"`
	TokenID             string            `bson:"token_id" json:"token_id"`
	TokenIDInt          int               `bson:"token_id_int" json:"token_id_int"`
	TokenIDMini         *int              `bson:"token_id_mini" json:"token_id_mini"`
	ContractAddress     string            `bson:"contract_address" json:"contract_address"`
	Name                string            `bson:"name" json:"name"`
	Description         string            `bson:"description" json:"description"`
	Image               string            `bson:"image" json:"image"`
	ParsedImage         *string           `bson:"parsed_image" json:"parsed_image"`
	AnimationURL        string            `bson:"animation_url" json:"animation_url"`
	AnimationHtml       *string           `bson:"animation_html"`
	Attributes          string            `bson:"attributes" json:"attributes"`
	ParsedAttributes    []TokenUriAttr    `bson:"parsed_attributes" json:"parsed_attributes"`
	ParsedAttributesStr []TokenUriAttrStr `bson:"parsed_attributes_str" json:"parsed_attributes_str"`
	ProjectID           string            `bson:"project_id" json:"project_id"`
	ProjectIDInt        int64             `bson:"project_id_int" json:"project_id_int"`
	BlockNumberMinted   *string           `bson:"block_number_minted" json:"block_number_minted"`
	MintedTime          *time.Time        `bson:"minted_time" json:"minted_time"`
	GenNFTAddr          string            `bson:"gen_nft_addrress"`
	Thumbnail           string            `bson:"thumbnail"`
	ThumbnailCapturedAt *time.Time        `bson:"thumbnailCapturedAt"`

	Stats TokenStats `bson:"stats" json:"stats"`

	OwnerAddr     string  `bson:"owner_addrress"`
	CreatorAddr   string  `bson:"creator_address"`
	Priority      *int    `bson:"priority"`
	MinterAddress *string `bson:"minter_address"`
	//accept duplicated data to query more faster
	Owner                          *Users        `bson:"owner"`
	Project                        *Projects     `bson:"project"`
	Creator                        *Users        `bson:"creator"`
	PaidType                       TokenPaidType `bson:"paidType"`
	IsOnchain                      bool          `bson:"isOnchain"`
	InscriptionIndex               string        `bson:"inscription_index"`
	OrderInscriptionIndex          int           `bson:"order_inscription_index" json:"order_inscription_index"`
	SyncedInscriptionInfo          bool          `bson:"synced_inscription_info"`
	CreatedByCollectionInscription bool          `bson:"created_by_collection_inscription"`
	Source                         string        `bson:"source" json:"source"`
	NftTokenId                     string        `bson:"nftTokenId"`
	InscribedBy                    string        `bson:"inscribedBy"`
	OriginalInscribedBy            string        `bson:"originalInscribedBy"`
	CreatedMintActivity            bool          `bson:"created_mint_activity"`
	CreatedTokenTx                 bool          `bson:"created_token_tx"`
}

type AggregateTokenUriTraits struct {
	AggregateTokenUriTraitsID `bson:"_id"`
	ParsedAttributes          []TokenUriAttr    `bson:"parsed_attributes" json:"parsed_attributes"`
	ParsedAttributesStr       []TokenUriAttrStr `bson:"parsed_attributes_str" json:"parsed_attributes_str"`
}

type AggregateTokenUriTraitsID struct {
	ProjectID string `bson:"project_id"`
	TokenID   string `bson:"token_id"`
}

type TokenUriListingFilter struct {
	ID                    primitive.ObjectID   `bson:"_id" json:"_id"`
	TokenID               string               `bson:"token_id" json:"tokenID"`
	Name                  string               `bson:"name" json:"name"`
	Image                 string               `bson:"image" json:"image"`
	ContractAddress       string               `bson:"contract_address" json:"contract_address"`
	AnimationURL          string               `bson:"animation_url" json:"animation_url"`
	AnimationHtml         *string              `bson:"animation_html"`
	ProjectID             string               `bson:"project_id" json:"projectID"`
	MintedTime            *time.Time           `bson:"minted_time" json:"minted_time"`
	GenNFTAddr            string               `bson:"gen_nft_addrress" json:"genNFTAddr"`
	Thumbnail             string               `bson:"thumbnail" json:"thumbnail"`
	InscriptionIndex      string               `bson:"inscription_index" json:"inscriptionIndex"`
	OrderInscriptionIndex int                  `bson:"order_inscription_index" json:"orderInscriptionIndex"`
	OrderID               primitive.ObjectID   `bson:"orderID" json:"orderID"`
	Price                 int64                `bson:"priceBTC" json:"priceBTC"`
	PriceETH              string               `bson:"priceETH" json:"priceETH"`
	Buyable               bool                 `bson:"buyable" json:"buyable"`
	SellVerified          bool                 `bson:"sell_verified" json:"sell_verified"`
	OwnerAddr             bool                 `bson:"ownerAddr" json:"ownerAddr"`
	OwnerAddress          string               `bson:"owner_addrress" json:"owner_addrress"`
	Owner                 TokenURIListingOwner `bson:"owner_object" json:"owner"`
	Project               struct {
		TokenID string `bson:"tokenid" json:"tokenID"`
		Royalty int64  `bson:"royalty" json:"royalty"`
	} `bson:"project" json:"project"`
	NftTokenID string `bson:"nftTokenId" json:"nftTokenId"`
}

type TokenURIListingOwner struct {
	WalletAddress           string `bson:"wallet_address" json:"walletAddress"`                       // eth wallet define user in platform by connect wallet and sign     // btc wallet artist receive royalty
	WalletAddressBTCTaproot string `bson:"wallet_address_btc_taproot" json:"walletAddressBtcTaproot"` // btc wallet receive minted nft
	DisplayName             string `bson:"display_name" json:"displayName"`
	Avatar                  string `bson:"avatar" json:"avatar"`
}

type TokenUriAttr struct {
	TraitType string      `bson:"trait_type" json:"trait_type"`
	Value     interface{} `bson:"value" json:"value"`
}

type TokenUriAttrStr struct {
	TraitType string `bson:"trait_type" json:"trait_type"`
	Value     string `bson:"value" json:"value"`
}

func (u TokenUri) TableName() string {
	return utils.COLLECTION_TOKEN_URI
}

func (u TokenUri) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
