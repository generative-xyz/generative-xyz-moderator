package dex_btc_cron

import (
	"sync"
	"time"

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

		h.Logger.Info("h.Usecase.JobWatchPendingDexBTCListing", "start")
		// job check tx:
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			h.Usecase.JobWatchPendingDexBTCListing()
		}(&wg)

		h.Logger.Info("h.Usecase.JobWatchPendingDexBTCListing", "wait")
		wg.Wait()
		time.Sleep(6 * time.Minute)
	}
}
