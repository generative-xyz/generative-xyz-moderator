package gm_crontab_sever

import (
	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils"
	"time"
)

type GmCrontabServer struct {
	Usecase *usecase.Usecase
}

func NewGmCrontabServer(uc *usecase.Usecase) (*GmCrontabServer, error) {
	t := &GmCrontabServer{}
	t.Usecase = uc
	return t, nil
}

func (tx GmCrontabServer) StartServer() {
	for {
		processing := "0"
		cached, err := tx.Usecase.Cache.GetData(utils.GM_CRONTAB_PROCESSING_KEY)
		if cached == nil && err != nil {
			cached = &processing
		}

		if *cached == "1" {
			continue
		}

		tx.Usecase.JobGetChartDataForGMCollection()
		time.Sleep(time.Minute * 5)
	}
}
