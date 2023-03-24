package thumbor

type Config struct {
	SecretKey string `yaml:"secretKey,omitempty" json:"secret_key,omitempty" validate:"required"`
	ServerUrl string `yaml:"serverUrl,omitempty" json:"server_url,omitempty" validate:"required"`
}
