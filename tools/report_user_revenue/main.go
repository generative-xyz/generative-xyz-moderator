package main

import (
	"fmt"
	_ "rederinghub.io/mongo/migrate"
	"rederinghub.io/tools"
	"rederinghub.io/utils/helpers"
)

// @title Generative.xyz APIs
// @version 1.0.0
// @description This is a sample server Generative.

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @BasePath /rederinghub.io/v1
func main() {
	uc := tools.StartFactory()
	if uc != nil {
		btcRate, err := helpers.GetExternalPrice("BTC")
		if err != nil {
			return
		}

		ethRate, err := helpers.GetExternalPrice("ETH")
		if err != nil {
			return
		}

		uc.ReportUserRevenue(btcRate, ethRate)
	}

	fmt.Println("Reports have been created")
}
