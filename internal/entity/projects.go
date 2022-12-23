package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type Projects struct {
	BaseEntity`bson:",inline"`
	ContractAddress string `bson:"contractAddress"`
	TokenID string `bsonbson:"tokenID"`
	MaxSupply int `bson:"maxSupply"`
	LimitSupply int `bson:"limitSupply"`
	MintPrice string `bson:"mintPrice"`
	Name string `bson:"name"`
	CreatorName string `bson:"creatorName"`
	Description string `bson:"description"`
	Thumbnail string `bson:"thumbnail"`
	ThirdPartyScripts []string `bson:"thirdPartyScripts"`
	Scripts []string `bson:"scripts"`
	ReservationList []string `bson:"reservationList"`
	MintFee int `bson:"mintFee"`
	OpenMintUnixTimestamp int `bson:"openMintUnixTimestamp"`
	TokenDescription string `bson:"tokenDescription"`
	Styles string `bson:"styles"`
	Royalty int `bson:"royalty"`
	SocialWeb string `bson:"socialWeb"`
	SocialTwitter string `bson:"socialTwitter"`
	SocialDiscord string `bson:"socialDiscord"`
	SocialMedium string `bson:"socialMedium"`
	SocialInstagram string `bson:"socialInstagram"`
	License string `bson:"license"`
	MintTokenAddress string `bson:"mintTokenAddress"`
	Hash string `bson:"hash"`
	Tags []string `bson:"tags"`
	Categories []string `bson:"categories"`
}

type FilterProjects struct {
	BaseFilters
}

func (u Projects) TableName() string { 
	return utils.COLLECTION_PROJECTS
}

func (u Projects) ToBson()  (*bson.D, error) { 
	return helpers.ToDoc(u)
}