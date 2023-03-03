package incribe_btc

import (
	"sync"
	"time"

	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils/global"
	"rederinghub.io/utils/logger"
	"rederinghub.io/utils/redis"
)

type ScronDeveloperInscribeHandler struct {
	Logger logger.Ilogger

	Cache   redis.IRedisCache
	Usecase usecase.Usecase
}

func NewScronDeveloperInscribeHandler(global *global.Global, uc usecase.Usecase) *ScronDeveloperInscribeHandler {
	return &ScronDeveloperInscribeHandler{
		Logger:  global.Logger,
		Cache:   global.Cache,
		Usecase: uc,
	}
}

func (h ScronDeveloperInscribeHandler) StartServer() {

	var wg sync.WaitGroup

	for {
		wg.Add(4)

		// job check tx:
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			h.Usecase.JobDeveloperInscribe_CheckTxSend()

		}(&wg)

		// job check balance:
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			h.Usecase.JobDeveloperInscribe_WaitingBalance()

		}(&wg)

		// job send btc to ord address:
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			h.Usecase.JobDeveloperInscribe_SendBTCToOrdWallet()

		}(&wg)

		// job mint nft:
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			h.Usecase.JobDeveloperInscribe_MintNft()

		}(&wg)

		h.Logger.Info("wait", "wait")
		wg.Wait()
		time.Sleep(5 * time.Minute)
	}
}
