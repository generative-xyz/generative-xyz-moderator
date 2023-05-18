package gm_crontab_sever

import (
	"rederinghub.io/internal/usecase"
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
		tx.Usecase.JobGetChartDataForGMCollection()
		time.Sleep(time.Minute * 5)
	}
}
