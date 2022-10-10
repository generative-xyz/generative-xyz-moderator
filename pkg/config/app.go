package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/gommon/log"
)

const (
	DevelopEnvironment    = "develop"
	ProductionEnvironment = "production"
)

var (
	server    ServerCfg
	appConfig AppCfg
)

type ServerCfg struct {
	SERVERUrl string `envconfig:"SERVER_URL" default:"0.0.0.0"`
	GRPCPort  int    `envconfig:"CORE_GRPC_PORT" default:"10000"`
	HTTPPort  int    `envconfig:"CORE_HTTP_PORT" default:"8000"`
}

type AppCfg struct {
	Environment   string `envconfig:"ENVIRONMENT" default:"develop"`
	MoralisURL    string `envconfig:"MORALIS_URL" default:"https://deep-index.moralis.io/api/v2/nft/"`
	MoralisAPIKey string `envconfig:"MORALIS_API_KEY" default:"6pELUXoEuCjQO1S92nEEQW6c1wNk1Qv4YdPNHJZPzkYeb3EOWlxF0pVPcWxd6J9u"`
}

func InitConfig() {
	configs := []interface{}{
		&server,
		&appConfig,
	}
	for _, instance := range configs {
		err := envconfig.Process("", instance)
		if err != nil {
			log.Fatalf("unable to init config: %v, err: %v", instance, err)
		}
	}
}

func ServerConfig() ServerCfg {
	return server
}

func AppConfig() AppCfg {
	return appConfig
}
