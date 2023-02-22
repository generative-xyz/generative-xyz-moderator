package mint_nft_btc

import (
	"sync"
	"time"

	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils/global"
	"rederinghub.io/utils/logger"
	"rederinghub.io/utils/redis"
)

type CronMintNftBtcHandler struct {
	Logger logger.Ilogger

	Cache   redis.IRedisCache
	Usecase usecase.Usecase
}

func NewCronMintNftBtcHandler(global *global.Global, uc usecase.Usecase) *CronMintNftBtcHandler {
	return &CronMintNftBtcHandler{
		Logger:  global.Logger,
		Cache:   global.Cache,
		Usecase: uc,
	}
}

func (h CronMintNftBtcHandler) StartServer() {

	var wg sync.WaitGroup

	for {
		wg.Add(7)

		// job check balance:
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			h.Usecase.JobMint_CheckBalance()

		}(&wg)

		// job mint nft:
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			h.Usecase.JobMint_MintNftBtc()

		}(&wg)

		// job check tx:
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			h.Usecase.JobMint_CheckTxMintSend()

		}(&wg)

		// job send nft to user:
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			h.Usecase.JobMint_SendNftToUser()

		}(&wg)

		// job send btc to master wallet:
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			h.Usecase.JobMint_SendFundToMaster()

		}(&wg)

		// job refund btc/eth
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			h.Usecase.JobMint_RefundBtc()

		}(&wg)

		// job check tx refund btc/eth
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			h.Usecase.JobMint_CheckTxMasterAndRefund()

		}(&wg)

		h.Logger.Info("wait", "wait")
		wg.Wait()
		time.Sleep(5 * time.Minute)
	}
}
