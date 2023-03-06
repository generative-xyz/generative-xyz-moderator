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
	go func () {
		h.Logger.Info("StartCrontabSyncTrending")
		for {
			err := h.Usecase.SyncProjectTrending()
			if err != nil {
				h.Logger.ErrorAny("SyncProjectTrendingError", zap.Any("err", err.Error()))
			}
			time.Sleep(10 * time.Minute)
		}
	}()

	go func () {
		h.Logger.Info("StartCrontabDeleteActivities")
		for {
			err := h.Usecase.DeleteOldActivities()
			if err != nil {
				h.Logger.ErrorAny("h.Usecase.DeleteOldActivities", zap.Any("err", err.Error()))
			}
			time.Sleep(24 * time.Hour)	
		}
	}()
}
