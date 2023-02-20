package mint_nft_btc

import (
	"sync"
	"time"

	"github.com/opentracing/opentracing-go"
	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils/global"
	"rederinghub.io/utils/logger"
	"rederinghub.io/utils/redis"
	"rederinghub.io/utils/tracer"
)

type CronMintNftBtcHandler struct {
	Logger  logger.Ilogger
	Tracer  tracer.ITracer
	Cache   redis.IRedisCache
	Usecase usecase.Usecase
}

func NewCronMintNftBtcHandler(global *global.Global, uc usecase.Usecase) *CronMintNftBtcHandler {
	return &CronMintNftBtcHandler{
		Logger:  global.Logger,
		Tracer:  global.Tracer,
		Cache:   global.Cache,
		Usecase: uc,
	}
}

func (h CronMintNftBtcHandler) StartServer() {

	var wg sync.WaitGroup

	for {
		wg.Add(5)

		span := h.Tracer.StartSpan("CronMintNftBtcHandler.DispatchCron")
		defer span.Finish()

		log := tracer.NewTraceLog()
		defer log.ToSpan(span)

		// job check balance:
		go func(rootSpan opentracing.Span, wg *sync.WaitGroup) {

			span := h.Tracer.StartSpanFromRoot(rootSpan, "CronMintNftBtcHandler.JobMint_CheckBalance")
			defer wg.Done()
			defer span.Finish()

			h.Usecase.JobMint_CheckBalance(span)

		}(span, &wg)

		// job check tx:
		go func(rootSpan opentracing.Span, wg *sync.WaitGroup) {

			span := h.Tracer.StartSpanFromRoot(rootSpan, "CronMintNftBtcHandler.JobMint_CheckTxMintSend")
			defer wg.Done()
			defer span.Finish()

			h.Usecase.JobMint_CheckTxMintSend(span)

		}(span, &wg)

		// job send btc to ord address:
		go func(rootSpan opentracing.Span, wg *sync.WaitGroup) {

			span := h.Tracer.StartSpanFromRoot(rootSpan, "CronMintNftBtcHandler.JobInscribeSendBTCToOrdWallet")
			defer wg.Done()
			defer span.Finish()

			h.Usecase.JobInscribeSendBTCToOrdWallet(span)

		}(span, &wg)

		// job mint nft:
		go func(rootSpan opentracing.Span, wg *sync.WaitGroup) {

			span := h.Tracer.StartSpanFromRoot(rootSpan, "CronMintNftBtcHandler.JobInscribeMintNft")
			defer wg.Done()
			defer span.Finish()

			h.Usecase.JobMint_MintNftBtc(span)

		}(span, &wg)

		// job send nft to user:
		go func(rootSpan opentracing.Span, wg *sync.WaitGroup) {

			span := h.Tracer.StartSpanFromRoot(rootSpan, "Inscribe.JobInscribeSendNft")
			defer wg.Done()
			defer span.Finish()

			h.Usecase.JobMin_SendNftToUser(span)

		}(span, &wg)

		log.SetData("wait", "wait")
		wg.Wait()
		time.Sleep(5 * time.Minute)
	}
}
