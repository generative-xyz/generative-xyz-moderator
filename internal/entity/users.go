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
	ID string `bson:"id"`
	WalletAddress string `bson:"wallet_address"`
	DisplayName string `bson:"display_name"`
	Bio string `bson:"bio"`
	Avatar string `bson:"avatar"`
	CreatedAt *time.Time `bson:"created_at"`
	ProfileSocial ProfileSocial
}


type ProfileSocial  struct{
    Web string `bson:"web"`;
    Twitter string `bson:"twitter"`;
    Discord string `bson:"discord"`;
    Medium string `bson:"medium"`;
	Instagram string `bson:"instagram"`;
}


func (u Users) TableName() string { 
	return utils.COLLECTION_USERS
}

func (u Users) ToBson()  (*bson.D, error) { 
	return helpers.ToDoc(u)
}