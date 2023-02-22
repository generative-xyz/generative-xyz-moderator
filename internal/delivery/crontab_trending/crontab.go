package crontab_trending

import (
	"gopkg.in/robfig/cron.v2"
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
	c := cron.New()
	// cronjob to sync projects trending
	c.AddFunc("*/15 * * * *", func() {
		
		err := h.Usecase.SyncProjectTrending()
		if err != nil {
			h.Logger.Error("DispatchCron.OneMinute.GetTheCurrentBlockNumber", err.Error(), err)
		}
	})
	c.Start()
}
