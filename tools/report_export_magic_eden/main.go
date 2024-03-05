package main

import (
	"fmt"
	_ "rederinghub.io/mongo/migrate"
	"rederinghub.io/tools"
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
		uc.ExportMagicEdend("1000001")
	}

	fmt.Println("Reports have been created")
}
