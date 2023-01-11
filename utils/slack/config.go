package slack

import "github.com/slack-go/slack"

type Config struct {
	Token     string `yaml:"token" json:"token"`
	ChannelId string `yaml:"channelId" json:"channel_id"`
	Env       string `yaml:"env" json:"env"`
}

type SlackData struct {
	ChannelName string           `json:"channel_name"`
	Data        slack.Attachment `json:"data"`
}
