package discordclient

type Config struct {
	WebhookURL string `yaml:"webhookUrl" validate:"required"`
}
