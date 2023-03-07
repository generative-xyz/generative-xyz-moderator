package entity

import (
	"time"

	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"

	"go.mongodb.org/mongo-driver/bson"
)

type UserType int

const (
	User  UserType = 0
	Admin UserType = 1
)

type FilterUsers struct {
	BaseFilters
	Search           *string
	Email            *string
	WalletAddress    *string
	WalletAddressBTC *string
	UserType         *UserType
	IsUpdatedAvatar  *bool
}

type FilteArtist struct {
	BaseFilters
}

type UserStats struct {
	CollectionCreated int32   `bson:"collection_created" json:"collection_created"`
	NftMinted         int32   `bson:"nft_minted" json:"nft_minted"`
	OutputMinted      int32   `bson:"output_minted" json:"output_minted"`
	VolumeMinted      float64 `bson:"volume_minted" json:"volume_minted"`
}

type Users struct {
	BaseEntity              `bson:",inline" json:"-"`
	IsVerified              bool          `bson:"is_verified"`
	VerifiedAt              *time.Time    `bson:"verified_at"`
	Message                 string        `bson:"message"`
	ID                      string        `bson:"id" json:"id,omitempty"`
	WalletAddress           string        `bson:"wallet_address" json:"wallet_address,omitempty"`                         // eth wallet define user in platform by connect wallet and sign
	WalletAddressPayment    string        `bson:"wallet_address_payment" json:"wallet_address_payment,omitempty"`         // eth wallet artist receive royalty
	WalletAddressBTC        string        `bson:"wallet_address_btc" json:"wallet_address_btc,omitempty"`                 // btc wallet artist receive royalty
	WalletAddressBTCTaproot string        `bson:"wallet_address_btc_taproot" json:"wallet_address_btc_taproot,omitempty"` // btc wallet receive minted nft
	DisplayName             string        `bson:"display_name" json:"display_name,omitempty"`
	Bio                     string        `bson:"bio" json:"bio,omitempty"`
	Avatar                  string        `bson:"avatar" json:"avatar"`
	IsUpdatedAvatar         *bool         `bson:"is_updated_avatar" json:"is_updated_avatar,omitempty"`
	CreatedAt               *time.Time    `bson:"created_at" json:"created_at,omitempty"`
	ProfileSocial           ProfileSocial `json:"profile_social,omitempty" bson:"profile_social"`
	Stats                   UserStats     `bson:"stats" json:"stats"`
	IsAdmin                 bool          `bson:"isAdmin" json:"isAdmin"`
}

type ProfileSocial struct {
	Web             string `bson:"web" json:"web,omitempty"`
	Twitter         string `bson:"twitter" json:"twitter,omitempty"`
	Discord         string `bson:"discord" json:"discord,omitempty"`
	Medium          string `bson:"medium" json:"medium,omitempty"`
	Instagram       string `bson:"instagram" json:"instagram,omitempty"`
	EtherScan       string `bson:"etherScan" json:"ether_scan,omitempty"`
	TwitterVerified bool   `bson:"twitter_verified" json:"twitterVerified,omitempty"`
}

func (u Users) TableName() string {
	return utils.COLLECTION_USERS
}

func (u Users) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
