package crontab_trending

import (
	"time"

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
	for {
		err := h.Usecase.SyncProjectTrending()
		if err != nil {
			h.Logger.Error("SyncProjectTrendingError", err.Error(), err)
		}
		time.Sleep(10 * time.Minute)
	}
}
