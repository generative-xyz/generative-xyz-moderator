package crontab_inscription_info

import (
	"sync"

	"go.uber.org/zap"
	"gopkg.in/robfig/cron.v2"
	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils/global"
	"rederinghub.io/utils/logger"
)

type ScronInscriptionInfoHandler struct {
	Logger  logger.Ilogger
	Usecase usecase.Usecase
}

func NewScronInscriptionInfoHandler(global *global.Global, uc usecase.Usecase) *ScronInscriptionInfoHandler {
	return &ScronInscriptionInfoHandler{
		Logger:  global.Logger,
		Usecase: uc,
	}
}

func (h ScronInscriptionInfoHandler) StartServer() {
	c := cron.New()
	// cronjob to sync inscription index

	// mutex to make sure 2 cronjob do not overlap
	var mu sync.Mutex
	c.AddFunc("*/10 * * * *", func() {
		mu.Lock()
		defer func() {
			mu.Unlock()
		}()
		err := h.Usecase.JobSyncTokenInscribeIndex()
		if err != nil {
			logger.AtLog.Logger.Error("DispatchCron.OneMinute.JobSyncTokenInscribeIndex", zap.Error(err))
		}
	})
	c.Start()
}
