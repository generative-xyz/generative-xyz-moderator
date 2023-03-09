package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type DiscordNotiStatus int
const (
	PENDING DiscordNotiStatus = 0
	DONE		DiscordNotiStatus	= 1
	FAILED  DiscordNotiStatus = 2
)

type GetDiscordNotiReq struct {
	Page int64
	Limit int64
	Status *DiscordNotiStatus
}

type DiscordNoti struct {
	BaseEntity    `bson:",inline"`
	Message    DiscordMessage           `bson:"message"`
	Status     DiscordNotiStatus `bson:"status"`
	NumRetried int               `bson:"num_retried"`
	Webhook    string `bson:"webhook"`
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
	Inline bool   `bson:"inline"`
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
