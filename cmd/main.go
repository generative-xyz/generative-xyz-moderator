package main

import (
	"os"
	"strings"

	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"rederinghub.io/internal/api/grpc"
	"rederinghub.io/pkg/config"
)

var (
	rootCmd = &cobra.Command{
		Use:     "renderinghub",
		Short:   "Rendering Hub application",
		Long:    `Rendering Hub application`,
		Version: "1.0.0",
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(grpc.ServiceCmd)
}

func initConfig() {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("error while reading config file: %s", err)
	}

	for _, env := range viper.AllKeys() {
		if viper.GetString(env) != "" {
			_ = os.Setenv(env, viper.GetString(env))
			_ = os.Setenv(strings.ToUpper(env), viper.GetString(env))
		}
	}

	config.InitConfig()
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		log.Errorf("error while execute: %s", err)
	}
}
