package crontab_btc_v2

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

type ScronBTCHandler struct {
	Logger  logger.Ilogger
	Tracer  tracer.ITracer
	Cache   redis.IRedisCache
	Usecase usecase.Usecase
}

func NewScronBTCHandler(global *global.Global, uc usecase.Usecase) *ScronBTCHandler {
	return &ScronBTCHandler{
		Logger:  global.Logger,
		Tracer:  global.Tracer,
		Cache:   global.Cache,
		Usecase: uc,
	}
}

func (h ScronBTCHandler) StartServer() {

	var wg sync.WaitGroup

	for {
		wg.Add(4)

		span := h.Tracer.StartSpan("ScronBTCHandler.DispatchCron")
		defer span.Finish()

		log := tracer.NewTraceLog()
		defer log.ToSpan(span)

		// job check tx:
		go func(rootSpan opentracing.Span, wg *sync.WaitGroup) {

			span := h.Tracer.StartSpanFromRoot(rootSpan, "Inscribe.JobInscribeCheckTxSend")
			defer wg.Done()
			defer span.Finish()

			h.Usecase.JobInscribeCheckTxSend(span)

		}(span, &wg)

		// job send btc to ord address:
		go func(rootSpan opentracing.Span, wg *sync.WaitGroup) {

			span := h.Tracer.StartSpanFromRoot(rootSpan, "Inscribe.JobInscribeSendBTCToOrdWallet")
			defer wg.Done()
			defer span.Finish()

			h.Usecase.JobInscribeSendBTCToOrdWallet(span)

		}(span, &wg)

		// job mint nft:
		go func(rootSpan opentracing.Span, wg *sync.WaitGroup) {

			span := h.Tracer.StartSpanFromRoot(rootSpan, "Inscribe.JobInscribeMintNft")
			defer wg.Done()
			defer span.Finish()

			h.Usecase.JobInscribeMintNft(span)

		}(span, &wg)

		// job send nft to user:
		go func(rootSpan opentracing.Span, wg *sync.WaitGroup) {

			span := h.Tracer.StartSpanFromRoot(rootSpan, "Inscribe.JobInscribeSendNft")
			defer wg.Done()
			defer span.Finish()

			h.Usecase.JobInscribeSendNft(span)

		}(span, &wg)

		log.SetData("wait", "wait")
		wg.Wait()
		time.Sleep(1 * time.Minute)
	}
}
