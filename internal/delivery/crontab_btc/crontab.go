package crontab_btc

import (
	"sync"
	"time"

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

	wg.Add(10)
	for {
		span := h.Tracer.StartSpan("ScronBTCHandler.DispatchCron.OneMinute")
		defer span.Finish()
		log := tracer.NewTraceLog()
		defer log.ToSpan(span)

		go func(wg *sync.WaitGroup) {
			span := h.Tracer.StartSpan("ScronBTCHandler.WaitingForBalancing")
			defer wg.Done()
			defer span.Finish()

			h.Usecase.WaitingForBalancing(span) // BTC
		}(&wg)
		
		go func(wg *sync.WaitGroup) {
			span := h.Tracer.StartSpan("ScronBTCHandler.WaitingForETHBalancing")
			defer wg.Done()
			defer span.Finish()

			h.Usecase.WaitingForETHBalancing(span) // ETH
		}(&wg)

		go func(wg *sync.WaitGroup) {

			span := h.Tracer.StartSpan("ScronBTCHandler.WaitingForMinted")
			defer wg.Done()
			defer span.Finish()

			h.Usecase.WaitingForMinted(span)

		}(&wg)

		//TODO mint with ETH payment?
		go func(wg *sync.WaitGroup) {
			span := h.Tracer.StartSpan("ScronBTCHandler.WaitingForETHMinted")
			defer wg.Done()
			defer span.Finish()

			h.Usecase.WaitingForETHMinted(span)

		}(&wg)

		go func(wg *sync.WaitGroup) {
			span := h.Tracer.StartSpan("ScronBTCHandler.BtcChecktListNft")
			defer wg.Done()
			defer span.Finish()

			h.Usecase.BtcChecktListNft(span)

		}(&wg)

		go func(wg *sync.WaitGroup) {
			span := h.Tracer.StartSpan("ScronBTCHandler.BtcCheckReceivedBuyingNft")
			defer wg.Done()
			defer span.Finish()

			h.Usecase.BtcCheckReceivedBuyingNft(span)

		}(&wg)

		go func(wg *sync.WaitGroup) {
			span := h.Tracer.StartSpan("ScronBTCHandler.BtcSendBTCForBuyOrder")
			defer wg.Done()
			defer span.Finish()

			h.Usecase.BtcSendBTCForBuyOrder(span)

		}(&wg)

		go func(wg *sync.WaitGroup) {
			span := h.Tracer.StartSpan("ScronBTCHandler.BtcCheckSendBTCForBuyOrder")
			defer wg.Done()
			defer span.Finish()
			h.Usecase.BtcCheckSendBTCForBuyOrder(span)

		}(&wg)

		go func(wg *sync.WaitGroup) {
			span := h.Tracer.StartSpan("ScronBTCHandler.BtcSendNFTForBuyOrder")
			defer wg.Done()
			defer span.Finish()
			h.Usecase.BtcSendNFTForBuyOrder(span)

		}(&wg)

		go func(wg *sync.WaitGroup) {
			span := h.Tracer.StartSpan("ScronBTCHandler.BtcCheckSendNFTForBuyOrder")
			defer wg.Done()
			defer span.Finish()
			h.Usecase.BtcCheckSendNFTForBuyOrder(span)

		}(&wg)

		wg.Wait()
		time.Sleep(1 * time.Minute)
	}
}
