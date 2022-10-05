package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/gommon/log"
)

var (
	server ServerCfg
)

type ServerCfg struct {
	SERVERUrl string `envconfig:"SERVER_URL" default:"0.0.0.0"`
	GRPCPort  int    `envconfig:"CORE_GRPC_PORT" default:"10000"`
	HTTPPort  int    `envconfig:"CORE_HTTP_PORT" default:"8000"`
	LogLevel  string `envconfig:"LOG_LEVEL" default:"debug"`
}

func InitConfig() {
	configs := []interface{}{
		&server,
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
