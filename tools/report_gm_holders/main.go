package main

import (
	"fmt"
	_ "rederinghub.io/mongo/migrate"
	"rederinghub.io/tools"
	"sync"
)

// @title Generative.xyz APIs
// @version 1.0.0
// @description This is a sample server Generative.

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @BasePath /rederinghub.io/v1
func main() {
	level1Domain := "https://explorer.trustless.computer"
	level2Domain := "https://explorer.l2.trustless.computer"
	wg := sync.WaitGroup{}
	wg.Add(2)

	uc := tools.StartFactory()
	if uc != nil {
		go uc.ReportGMHolders(&wg, level1Domain, "0x2fe8d5a64affc1d703aeca8a566f5e9faee0c003", "level-1")
		go uc.ReportGMHolders(&wg, level2Domain, "0x0170435186a9a2Af5881C6236CF47211D046cAE6", "level-2")
	}

	wg.Wait()
	fmt.Println("Reports have been created")
}
