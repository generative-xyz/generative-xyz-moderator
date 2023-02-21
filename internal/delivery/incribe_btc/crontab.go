package incribe_btc

import (
	"sync"
	"time"

	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils/global"
	"rederinghub.io/utils/logger"
	"rederinghub.io/utils/redis"
)

type ScronBTCHandler struct {
	Logger  logger.Ilogger
	
	Cache   redis.IRedisCache
	Usecase usecase.Usecase
}

func NewScronBTCHandler(global *global.Global, uc usecase.Usecase) *ScronBTCHandler {
	return &ScronBTCHandler{
		Logger:  global.Logger,
		Cache:   global.Cache,
		Usecase: uc,
	}
}

func (h ScronBTCHandler) StartServer() {

	var wg sync.WaitGroup

	for {
		wg.Add(5)


		// job check tx:
		go func( wg *sync.WaitGroup) {
			defer wg.Done()
			h.Usecase.JobInscribeCheckTxSend()

		}(&wg)

		// job check balance:
		go func( wg *sync.WaitGroup) {
			defer wg.Done()
			h.Usecase.JobInscribeWaitingBalance()

		}(&wg)

		// job send btc to ord address:
		go func( wg *sync.WaitGroup) {
			defer wg.Done()
			h.Usecase.JobInscribeSendBTCToOrdWallet()

		}(&wg)

		// job mint nft:
		go func( wg *sync.WaitGroup) {
			defer wg.Done()
			h.Usecase.JobInscribeMintNft()

		}(&wg)

		// job send nft to user:
		go func( wg *sync.WaitGroup) {

				defer wg.Done()
			h.Usecase.JobInscribeSendNft()

		}(&wg)

		h.Logger.Info("wait", "wait")
		wg.Wait()
		time.Sleep(5 * time.Minute)
	}
}
