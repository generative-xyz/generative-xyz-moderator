package crontab_btc

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
		wg.Add(10)

		span := h.Tracer.StartSpan("ScronBTCHandler.DispatchCron")
		defer span.Finish()

		log := tracer.NewTraceLog()
		defer log.ToSpan(span)

		go func(rootSpan opentracing.Span, wg *sync.WaitGroup) {
			span := h.Tracer.StartSpanFromRoot(rootSpan, "BTC.WaitingForBalancing")
			defer wg.Done()
			defer span.Finish()

			h.Usecase.WaitingForBalancing(span) // BTC
		}(span, &wg)

		go func(rootSpan opentracing.Span, wg *sync.WaitGroup) {
			span := h.Tracer.StartSpanFromRoot(rootSpan, "ETH.WaitingForETHBalancing")
			defer wg.Done()
			defer span.Finish()

			h.Usecase.WaitingForETHBalancing(span) // ETH
		}(span, &wg)

		go func(rootSpan opentracing.Span, wg *sync.WaitGroup) {

			span := h.Tracer.StartSpanFromRoot(rootSpan, "BTC.WaitingForMinted")
			defer wg.Done()
			defer span.Finish()

			h.Usecase.WaitingForMinted(span)

		}(span, &wg)

		//TODO mint with ETH payment?
		go func(rootSpan opentracing.Span, wg *sync.WaitGroup) {
			span := h.Tracer.StartSpanFromRoot(rootSpan, "ETH.WaitingForETHMinted")
			defer wg.Done()
			defer span.Finish()

			h.Usecase.WaitingForETHMinted(span)

		}(span, &wg)

		go func(rootSpan opentracing.Span, wg *sync.WaitGroup) {
			span := h.Tracer.StartSpanFromRoot(rootSpan, "ScronBTCHandler.BtcChecktListNft")
			defer wg.Done()
			defer span.Finish()

			h.Usecase.BtcChecktListNft(span)

		}(span, &wg)

		go func(rootSpan opentracing.Span, wg *sync.WaitGroup) {
			span := h.Tracer.StartSpanFromRoot(rootSpan, "ScronBTCHandler.BtcCheckReceivedBuyingNft")
			defer wg.Done()
			defer span.Finish()

			h.Usecase.BtcCheckReceivedBuyingNft(span)

		}(span, &wg)

		go func(rootSpan opentracing.Span, wg *sync.WaitGroup) {
			span := h.Tracer.StartSpanFromRoot(rootSpan, "ScronBTCHandler.BtcSendBTCForBuyOrder")
			defer wg.Done()
			defer span.Finish()

			h.Usecase.BtcSendBTCForBuyOrder(span)

		}(span, &wg)

		go func(rootSpan opentracing.Span, wg *sync.WaitGroup) {
			span := h.Tracer.StartSpanFromRoot(rootSpan, "ScronBTCHandler.BtcCheckSendBTCForBuyOrder")
			defer wg.Done()
			defer span.Finish()
			h.Usecase.BtcCheckSendBTCForBuyOrder(span)

		}(span, &wg)

		go func(rootSpan opentracing.Span, wg *sync.WaitGroup) {
			span := h.Tracer.StartSpanFromRoot(rootSpan, "ScronBTCHandler.BtcSendNFTForBuyOrder")
			defer wg.Done()
			defer span.Finish()
			h.Usecase.BtcSendNFTForBuyOrder(span)

		}(span, &wg)

		go func(rootSpan opentracing.Span, wg *sync.WaitGroup) {
			span := h.Tracer.StartSpanFromRoot(rootSpan, "ScronBTCHandler.BtcCheckSendNFTForBuyOrder")
			defer wg.Done()
			defer span.Finish()
			h.Usecase.BtcCheckSendNFTForBuyOrder(span)

		}(span, &wg)

		log.SetData("wait", "wait")
		wg.Wait()
		time.Sleep(1 * time.Minute)
	}
}
