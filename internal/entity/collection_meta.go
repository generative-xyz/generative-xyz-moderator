package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type CollectionMeta struct {
	BaseEntity`bson:",inline"`
	Name            string `bson:"name" json:"name"`
	InscriptionIcon string `bson:"inscription_icon" json:"inscription_icon"`
	Supply          string `bson:"supply" json:"supply"`
	Slug            string `bson:"slug" json:"slug"`
	Description     string `bson:"description" json:"description"`
	TwitterLink     string `bson:"twitter_link" json:"twitter_link"`
	DiscordLink     string `bson:"discord_link" json:"discord_link"`
	WebsiteLink     string `bson:"website_link" json:"website_link"`
	ProjectCreated  bool `bson:"project_created"`
	Royalty					string  `bson:"royalty" json:"royalty"`
	Source          string `bson:"source"`
	WalletAddress   string `bson:"wallet_address" json:"wallet_address"`
	ProjectID       string `bson:"project_id" json:"project_id"`
	ProjectExisted  bool   `bson:"project_existed" json:"project_existed"`
	From            string `bson:"from" json:"source"`
}

func (u CollectionMeta) TableName() string { 
	return utils.COLLECTION_COLLECTION_META
}

func (u CollectionMeta) ToBson()  (*bson.D, error) { 
	return helpers.ToDoc(u)
}
