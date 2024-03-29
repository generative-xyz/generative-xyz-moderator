package crontab_trending

import (
	"time"

	"go.uber.org/zap"
	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils/global"
	"rederinghub.io/utils/logger"
)

type ScronTrendingHandler struct {
	Logger  logger.Ilogger
	Usecase usecase.Usecase
}

func NewScronTrendingHandler(global *global.Global, uc usecase.Usecase) *ScronTrendingHandler {
	return &ScronTrendingHandler{
		Logger:  global.Logger,
		Usecase: uc,
	}
}

func (h ScronTrendingHandler) StartServer() {
	go func() {
		logger.AtLog.Logger.Info("StartCrontabSyncTrending")
		for {
			err := h.Usecase.JobSyncProjectTrending()
			if err != nil {
				logger.AtLog.Logger.Error("JobSyncProjectTrendingError", zap.Any("err", err.Error()))
			}
			time.Sleep(10 * time.Minute)
		}
	}()

	go func() {
		logger.AtLog.Logger.Info("StartCrontabDeleteActivities")
		for {
			err := h.Usecase.JobDeleteOldActivities()
			if err != nil {
				logger.AtLog.Logger.Error("h.Usecase.JobDeleteOldActivities", zap.Any("err", err.Error()))
			}
			time.Sleep(24 * time.Hour)
		}
	}()
}
