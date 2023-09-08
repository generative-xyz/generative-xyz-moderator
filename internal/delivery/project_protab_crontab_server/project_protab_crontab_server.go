package project_protab_crontab_server

import (
	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils/logger"
	"time"
)

type ProjectProtabCrontabServer struct {
	Usecase *usecase.Usecase
}

func NewProjectProtabCrontabServer(uc *usecase.Usecase) (*ProjectProtabCrontabServer, error) {
	t := &ProjectProtabCrontabServer{}
	t.Usecase = uc
	return t, nil
}

func (tx ProjectProtabCrontabServer) StartServer() {
	logger.AtLog.Logger.Info("ProjectProtabCrontabServer")
	for {

		tx.Usecase.JobProjectProtab()
		time.Sleep(time.Minute * 5)

		tx.Usecase.JobProjectProtabUniqueOwner()
		time.Sleep(time.Minute * 10)
	}
}
