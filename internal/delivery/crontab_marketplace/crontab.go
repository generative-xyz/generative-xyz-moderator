package crontab_marketplace

import (
	"sync"
	"time"

	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils/global"
	"rederinghub.io/utils/logger"
	"rederinghub.io/utils/redis"
)

type ScronMarketplaceHandler struct {
	Logger logger.Ilogger

	Cache   redis.IRedisCache
	Usecase usecase.Usecase
}

func NewScronMarketPlace(global *global.Global, uc usecase.Usecase) *ScronMarketplaceHandler {
	return &ScronMarketplaceHandler{
		Logger:  global.Logger,
		Cache:   global.Cache,
		Usecase: uc,
	}
}

func (h ScronMarketplaceHandler) StartServer() {

	var wg sync.WaitGroup

	for {
		wg.Add(7)

		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			h.Usecase.JobMKP_CheckReceivedNftFromSeller()

		}(&wg)

		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			h.Usecase.JobMKP_WaitForBalanceFromBuyer()

		}(&wg)

		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			h.Usecase.JobMKP_Payment()

		}(&wg)

		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			h.Usecase.JobMKP_CheckTxSendPayment()

		}(&wg)

		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			h.Usecase.JobMKP_SendNftToBuyer()

		}(&wg)

		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			h.Usecase.JobMKP_CheckTxSendNftToBuyer()

		}(&wg)

		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			h.Usecase.JobMKP_CrawlToUpdateNftInfo()
		}(&wg)

		h.Logger.Info("MaketPlace.wait", "wait")
		wg.Wait()
		time.Sleep(5 * time.Minute)
	}
}
