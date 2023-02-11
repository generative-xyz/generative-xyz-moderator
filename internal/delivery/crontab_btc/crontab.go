package crontab_btc

import (
	"gopkg.in/robfig/cron.v2"
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

	if true {
		return
	}

	span := h.Tracer.StartSpan("ScronBTCHandler.DispatchCron.OneMinute")
	defer span.Finish()

	c := cron.New()
	c.AddFunc("*/1 * * * *", func() {
		span := h.Tracer.StartSpan("ScronBTCHandler.DispatchCron.OneMinute")
		defer span.Finish()

		h.Usecase.WaitingForETHBalancing(span) // ETH

		log := tracer.NewTraceLog()
		defer log.ToSpan(span)

		go func() {
			span := h.Tracer.StartSpan("ScronBTCHandler.WaitingForBalancing")
			defer span.Finish()

			h.Usecase.WaitingForBalancing(span) // BTC
		}()
		go func() {
			span := h.Tracer.StartSpan("ScronBTCHandler.WaitingForETHBalancing")
			defer span.Finish()

			h.Usecase.WaitingForETHBalancing(span) // ETH
		}()

		go func() {

			span := h.Tracer.StartSpan("ScronBTCHandler.WaitingForMinted")
			defer span.Finish()

			h.Usecase.WaitingForMinted(span)

		}()

		//TODO mint with ETH payment?
		go func() {

			span := h.Tracer.StartSpan("ScronBTCHandler.WaitingForETHMinted")
			defer span.Finish()

			h.Usecase.WaitingForETHMinted(span)

		}()
		go func() {

			span := h.Tracer.StartSpan("ScronBTCHandler.BtcChecktListNft")
			defer span.Finish()

			h.Usecase.BtcChecktListNft(span)

		}()
		go func() {
			span := h.Tracer.StartSpan("ScronBTCHandler.BtcCheckReceivedBuyingNft")
			defer span.Finish()

			h.Usecase.BtcCheckReceivedBuyingNft(span)

		}()

		go func() {
			span := h.Tracer.StartSpan("ScronBTCHandler.BtcSendBTCForBuyOrder")
			defer span.Finish()

			h.Usecase.BtcSendBTCForBuyOrder(span)

		}()

		go func() {
			span := h.Tracer.StartSpan("ScronBTCHandler.BtcCheckSendBTCForBuyOrder")
			defer span.Finish()
			h.Usecase.BtcCheckSendBTCForBuyOrder(span)

		}()

		go func() {
			span := h.Tracer.StartSpan("ScronBTCHandler.BtcSendNFTForBuyOrder")
			defer span.Finish()
			h.Usecase.BtcSendNFTForBuyOrder(span)

		}()

		go func() {
			span := h.Tracer.StartSpan("ScronBTCHandler.BtcCheckSendNFTForBuyOrder")
			defer span.Finish()
			h.Usecase.BtcCheckSendNFTForBuyOrder(span)

		}()

	})

	c.Start()
}
