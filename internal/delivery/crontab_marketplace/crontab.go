package crontab_marketplace

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

type ScronMarketplaceHandler struct {
	Logger  logger.Ilogger
	Tracer  tracer.ITracer
	Cache   redis.IRedisCache
	Usecase usecase.Usecase
}

func NewScronMarketPlace(global *global.Global, uc usecase.Usecase) *ScronMarketplaceHandler {
	return &ScronMarketplaceHandler{
		Logger:  global.Logger,
		Tracer:  global.Tracer,
		Cache:   global.Cache,
		Usecase: uc,
	}
}

func (h ScronMarketplaceHandler) StartServer() {

	var wg sync.WaitGroup

	for {
		wg.Add(6)

		span := h.Tracer.StartSpan("ScronMarketPlace.DispatchCron")
		defer span.Finish()

		log := tracer.NewTraceLog()
		defer log.ToSpan(span)
	

		go func(rootSpan opentracing.Span, wg *sync.WaitGroup) {
			span := h.Tracer.StartSpanFromRoot(rootSpan, "ScronMarketPlace.BtcChecktListNft")
			defer wg.Done()
			defer span.Finish()

			h.Usecase.BtcChecktListNft(span)

		}(span, &wg)

		go func(rootSpan opentracing.Span, wg *sync.WaitGroup) {
			span := h.Tracer.StartSpanFromRoot(rootSpan, "ScronMarketPlace.BtcCheckReceivedBuyingNft")
			defer wg.Done()
			defer span.Finish()

			h.Usecase.BtcCheckReceivedBuyingNft(span)

		}(span, &wg)

		go func(rootSpan opentracing.Span, wg *sync.WaitGroup) {
			span := h.Tracer.StartSpanFromRoot(rootSpan, "ScronMarketPlace.BtcSendBTCForBuyOrder")
			defer wg.Done()
			defer span.Finish()

			h.Usecase.BtcSendBTCForBuyOrder(span)

		}(span, &wg)

		go func(rootSpan opentracing.Span, wg *sync.WaitGroup) {
			span := h.Tracer.StartSpanFromRoot(rootSpan, "ScronMarketPlace.BtcCheckSendBTCForBuyOrder")
			defer wg.Done()
			defer span.Finish()
			h.Usecase.BtcCheckSendBTCForBuyOrder(span)

		}(span, &wg)

		go func(rootSpan opentracing.Span, wg *sync.WaitGroup) {
			span := h.Tracer.StartSpanFromRoot(rootSpan, "ScronMarketPlace.BtcSendNFTForBuyOrder")
			defer wg.Done()
			defer span.Finish()
			h.Usecase.BtcSendNFTForBuyOrder(span)

		}(span, &wg)

		go func(rootSpan opentracing.Span, wg *sync.WaitGroup) {
			span := h.Tracer.StartSpanFromRoot(rootSpan, "ScronMarketPlace.BtcCheckSendNFTForBuyOrder")
			defer wg.Done()
			defer span.Finish()
			h.Usecase.BtcCheckSendNFTForBuyOrder(span)

		}(span, &wg)

		log.SetData("MaketPlace.wait", "wait")
		wg.Wait()
		time.Sleep(1 * time.Minute)
	}
}
