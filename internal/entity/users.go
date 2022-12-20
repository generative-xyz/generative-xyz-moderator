package entity

import (
	"time"

	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"

	"go.mongodb.org/mongo-driver/bson"
)
type UserType int 
const (
	User UserType = 0
	Admin UserType = 1	
)

type FilterUsers struct {
	BaseFilters
	Email *string
	WalletAddress *string
	UserType *UserType
}

type Users struct {
	BaseEntity`bson:",inline"`
	WalletAddress string `bson:"wallet_address"`
	IsVerified bool `bson:"is_verified"`
	Message string `bson:"message"`
	Nickname string `bson:"nickname"`
	Bio string `bson:"bio"`
	FullName string `bson:"full_name"`
	FirstName string `bson:"first_name"`
	LastName string `bson:"last_name"`
	Email string `bson:"email"`
	Avatar string `bson:"avatar"`
	Address string `bson:"address"`
	Apartment string `bson:"apartment"`
	City string `bson:"city"`
	State string `bson:"state"`
	ZipCode string `bson:"zip_code"`
	CoverPhoto string `bson:"cover_photo"`
	UserType UserType `bson:"user_type"`
	LinkOpensea string `bson:"link_opensea"`
	LinkSocial string `bson:"link_social"`
	UserID int32 `bson:"user_id"` // only admin can have this ID
	VerifiedAt *time.Time `bson:"verified_at"`
}

func (u Users) TableName() string { 
	return utils.COLLECTION_USERS
}

func (u Users) ToBson()  (*bson.D, error) { 
	return helpers.ToDoc(u)
}