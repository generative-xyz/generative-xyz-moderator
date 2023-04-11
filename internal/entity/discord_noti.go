package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type DiscordNotiStatus int

const (
	PENDING DiscordNotiStatus = 0
	DONE    DiscordNotiStatus = 1
	FAILED  DiscordNotiStatus = 2
)

type DiscordNotiType string

const (
	UNREGCONIZE                DiscordNotiType = ""
	NEW_AIRDROP                DiscordNotiType = "new_airdrop"
	DISCORD_SALE_CHANNEL       DiscordNotiType = "new_sale"
	DISCORD_PERCEPTRON_CHANNEL DiscordNotiType = "new_sale_perceptron"
	DISCORD_ART_CHANNEL        DiscordNotiType = "new_sale_art"
	NEW_LISTING                DiscordNotiType = "new_listing"
	NEW_MINT                   DiscordNotiType = "new_mint"
	DISCORD_PFPS_CHANNEL       DiscordNotiType = "new_mint_pfps"
	DISCORD_DAO_CHANNEL        DiscordNotiType = "project_report"
	NEW_PROJECT                DiscordNotiType = "new_project"
	NEW_PROJECT_PROPOSED       DiscordNotiType = "new_project_proposed"
	NEW_PROJECT_APPROVED       DiscordNotiType = "new_project_approved"
	NEW_BID                    DiscordNotiType = "new_bid"
)

type GetDiscordNotiReq struct {
	Page   int64
	Limit  int64
	Status *DiscordNotiStatus
}

type DiscordNotiMeta struct {
	InscriptionID string `bson:"inscription_id"`
	ProjectID     string `bson:"project_id"`
	SentTo        string `bson:"sent_to"`
}

type DiscordNoti struct {
	BaseEntity `bson:",inline"`
	Message    DiscordMessage    `bson:"message"`
	Status     DiscordNotiStatus `bson:"status"`
	NumRetried int               `bson:"num_retried"`
	Webhook    string            `bson:"webhook"`
	Type       DiscordNotiType   `bson:"type"`
	Meta       DiscordNotiMeta   `bson:"meta"`
}

type DiscordMessage struct {
	Username        string          `bson:"username"`
	AvatarUrl       string          `bson:"avatar_url"`
	Content         string          `bson:"content"`
	Embeds          []Embed         `bson:"embeds"`
	AllowedMentions AllowedMentions `bson:"allowed_mentions"`
}

type Embed struct {
	Title       string    `bson:"title"`
	Url         string    `bson:"url"`
	Description string    `bson:"description"`
	Color       string    `bson:"color"`
	Author      Author    `bson:"author"`
	Fields      []Field   `bson:"fields"`
	Thumbnail   Thumbnail `bson:"thumbnail"`
	Image       Image     `bson:"image"`
	Footer      Footer    `bson:"footer"`
}

type Author struct {
	Name    string `bson:"name"`
	Url     string `bson:"url"`
	IconUrl string `bson:"icon_url"`
}

type Field struct {
	Name   string `bson:"name"`
	Value  string `bson:"value"`
	Inline bool   `bson:"_inline"`
}

type Thumbnail struct {
	Url string `bson:"url"`
}

type Image struct {
	Url string `bson:"url"`
}

type Footer struct {
	Text    string `bson:"text"`
	IconUrl string `bson:"icon_url"`
}

type AllowedMentions struct {
	Parse []string `bson:"parse"`
	Users []string `bson:"users"`
	Roles []string `bson:"roles"`
}

func (u DiscordNoti) TableName() string {
	return utils.COLLECTION_DISCORD_NOTI
}

func (u DiscordNoti) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
