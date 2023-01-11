package slack

import (
	"fmt"

	"github.com/slack-go/slack"
)

type Slack struct {
	client    *slack.Client
	token     string
	channelId string
	env       string
}

func NewSlack(cnf Config) *Slack {
	client := slack.New(cnf.Token)

	return &Slack{
		token:     cnf.Token,
		client:    client,
		env:       cnf.Env,
		channelId: cnf.ChannelId,
	}
}

func (s Slack) postMessage(channelId string, messages ...slack.MsgOption) (string, string, error) {
	return s.client.PostMessage(channelId, messages...)
}

func (s Slack) SendMessageToSlack(pretext string, title string, text string) (string, string, error) {
	attachment := slack.Attachment{
		Color:   "#1766ff",
		Pretext: fmt.Sprintf("[%s] %s", s.env, pretext),
		Title:   title,
		Text:    text,
	}

	return s.postMessage(s.channelId, slack.MsgOptionAttachments(attachment))
}
