package dex_btc_cron

import (
	"sync"
	"time"

	"go.uber.org/zap"
	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils/global"
	"rederinghub.io/utils/logger"
	"rederinghub.io/utils/redis"
)

type ScronBTCHandler struct {
	Logger logger.Ilogger

	Cache   redis.IRedisCache
	Usecase usecase.Usecase
}

func NewScronDexBTCHandler(global *global.Global, uc usecase.Usecase) *ScronBTCHandler {
	return &ScronBTCHandler{
		Logger:  global.Logger,
		Cache:   global.Cache,
		Usecase: uc,
	}
}

func (h ScronBTCHandler) StartServer() {

	var wg sync.WaitGroup

	for {
		wg.Add(1)

		logger.AtLog.Logger.Info("h.Usecase.JobWatchPendingDexBTCListing", zap.Any("start", "start"))
		// job check tx:
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			h.Usecase.JobWatchPendingDexBTCListing()
		}(&wg)

		logger.AtLog.Logger.Info("h.Usecase.JobWatchPendingDexBTCListing", zap.Any("wait", "wait"))
		wg.Wait()
		time.Sleep(3 * time.Minute)
	}
}
