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
	UNREGCONIZE          DiscordNotiType = ""
	NEW_AIRDROP          DiscordNotiType = "new_airdrop"
	NEW_SALE             DiscordNotiType = "new_sale"
	NEW_SALE_PERCEPTRON  DiscordNotiType = "new_sale_perceptron"
	NEW_SALE_ART         DiscordNotiType = "new_sale_art"
	NEW_SALE_PFPS        DiscordNotiType = "new_sale_pfps"
	NEW_SALE_RECALL      DiscordNotiType = "new_sale_recall"
	NEW_LISTING          DiscordNotiType = "new_listing"
	NEW_MINT             DiscordNotiType = "new_mint"
	NEW_MINT_PERCEPTRON  DiscordNotiType = "new_mint_perceptron"
	NEW_MINT_ART         DiscordNotiType = "new_mint_art"
	NEW_MINT_PFPS        DiscordNotiType = "new_mint_pfps"
	NEW_PROJECT          DiscordNotiType = "new_project"
	NEW_PROJECT_PROPOSED DiscordNotiType = "new_project_proposed"
	NEW_ART_WEBHOOK                      = "new_art_webhook"
	NEW_PROJECT_APPROVED DiscordNotiType = "new_project_approved"
	NEW_BID              DiscordNotiType = "new_bid"
	NEW_PROJECT_REPORT   DiscordNotiType = "new_project_report"
	NEW_PROJECT_REMOVE   DiscordNotiType = "new_project_remove"
	NEW_PROJECT_VOTE     DiscordNotiType = "new_project_vote"

	WaitingMintNotification = "waiting_mint_notification"
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
	Category      string `bson:"category"`
	Amount        uint64 `bson:"amount"`
}

type ImageSourceType int
type ImagePosition int

const (
	ImageFromInscriptionID ImageSourceType = 1
	ThumbNailPosition      ImagePosition   = 1
	FullImagePosition      ImagePosition   = 2
)

type DiscordNoti struct {
	BaseEntity      `bson:",inline"`
	Message         DiscordMessage    `bson:"message"`
	Status          DiscordNotiStatus `bson:"status"`
	NumRetried      int               `bson:"num_retried"`
	Webhook         string            `bson:"webhook"`
	Type            DiscordNotiType   `bson:"type"`
	Meta            DiscordNotiMeta   `bson:"meta"`
	RequireImage    bool              `bson:"require_image"`
	ImageSourceType ImageSourceType   `bson:"image_source_type"`
	ImageSourceID   string            `bson:"image_source_id"`
	Note            string            `bson:"note"`
	ImagePosition   ImagePosition     `bson:"image_position"`
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
