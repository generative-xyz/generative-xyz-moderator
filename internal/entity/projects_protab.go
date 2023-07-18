package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

// Protab
type ProjectsProtab struct {
	BaseEntityNoID    `bson:",inline" json:"-"`
	ContractAddress   string  `bson:"contractAddress" json:"contractAddress"`
	TokenID           string  `bson:"tokenid" json:"tokenID"`
	TokenIDInt        int64   `bson:"tokenIDInt" json:"tokenIDInt"`
	Name              string  `bson:"name" json:"name"`
	CreatorAddress    string  `bson:"creatorAddress" json:"creatorAddrr"`
	CreatorAddressBTC string  `bson:"creatorAddrrBTC" json:"creatorAddrrBTC"`
	Thumbnail         string  `bson:"thumbnail" json:"thumbnail"`
	MaxSupply         int64   `bson:"maxSupply" json:"maxSupply"`
	MintVolume        float64 `bson:"mint_volume" json:"mint_volume"`
	Volume            float64 `bson:"volume" json:"volume"`
	FloorPrice        float64 `bson:"floor_price" json:"floor_price"`
	CEXVolume         float64 `bson:"cex_volume" json:"cex_volume"`
	Listed            int     `bson:"listed" json:"listed"`
	UniqueOwners      int     `bson:"unique_owners" json:"unique_owners"`
	IsBuyable         bool    `bson:"is_buyable" json:"is_buyable"`
}

type UpdateProjectsProtab struct {
	BaseEntityNoID    `bson:",inline" json:"-"`
	ContractAddress   string  `bson:"contractAddress" json:"contractAddress"`
	TokenID           string  `bson:"tokenid" json:"tokenID"`
	TokenIDInt        int64   `bson:"tokenIDInt" json:"tokenIDInt"`
	Name              string  `bson:"name" json:"name"`
	CreatorAddress    string  `bson:"creatorAddress" json:"creatorAddrr"`
	CreatorAddressBTC string  `bson:"creatorAddrrBTC" json:"creatorAddrrBTC"`
	Thumbnail         string  `bson:"thumbnail" json:"thumbnail"`
	MaxSupply         int64   `bson:"maxSupply" json:"maxSupply"`
	MintVolume        float64 `bson:"mint_volume" json:"mint_volume"`
	Volume            float64 `bson:"volume" json:"volume"`
	FloorPrice        float64 `bson:"floor_price" json:"floor_price"`
	CEXVolume         float64 `bson:"cex_volume" json:"cex_volume"`
	Listed            int     `bson:"listed" json:"listed"`
	IsBuyable         bool    `bson:"is_buyable" json:"is_buyable"`
}

type ProjectsProtabAPI struct {
	BaseEntityNoID `bson:",inline" json:"-"`
	ProjectsProtab `bson:",inline" json:"-"`
	Project        ProjectsProtabAPIProject `bson:"project" json:"project"`
	Owner          ProjectsProtabAPIOwner   `bson:"owner" json:"owner"`
}

type ProjectsProtabAPIProject struct {
	Name            string                              `json:"name" bson:"name"`
	TokenId         string                              `json:"tokenid" bson:"tokenid"`
	Thumbnail       string                              `json:"thumbnail" bson:"thumbnail"`
	ContractAddress string                              `json:"contractAddress" bson:"contractAddress"`
	CreatorAddress  string                              `json:"creatorAddress" bson:"creatorAddress"`
	MaxSupply       int                                 `json:"maxSupply" bson:"maxSupply"`
	IsMintedOut     bool                                `json:"isMintedOut" bson:"isMintedOut"`
	MintingInfo     ProjectsProtabAPIProjectMintingInfo `json:"mintingInfo" bson:"mintingInfo"`
}

type ProjectsProtabAPIProjectMintingInfo struct {
	Index        int `json:"index" bson:"index"`
	IndexReverse int `json:"index_reverse" bson:"indexReverse"`
}

type ProjectsProtabAPIOwner struct {
	WalletAddress           string `json:"wallet_address" bson:"wallet_address"`
	WalletAddressPayment    string `json:"wallet_address_payment" bson:"wallet_address_payment"`
	WalletAddressBtc        string `json:"wallet_address_btc" bson:"wallet_address_btc"`
	WalletAddressBtcTaproot string `json:"wallet_address_btc_taproot" bson:"wallet_address_btc_taproot"`
	DisplayName             string `json:"display_name" bson:"display_name"`
	Avatar                  string `json:"avatar" bson:"avatar"`
}

func (u ProjectsProtab) TableName() string {
	return utils.COLLECTION_PROJECT_PROTAB
}

func (u ProjectsProtab) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}

func (u UpdateProjectsProtab) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
