package crontab_btc

import (
	"sync"
	"time"

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

	
	//waiting for blancing
	go func() {
		var wg sync.WaitGroup
		//Waiting for balance + Mint
		for {
			wg.Add(3)
			go func(wg *sync.WaitGroup) {
				span := h.Tracer.StartSpan("BTC.WaitingForBalancing")
				defer wg.Done()
				defer span.Finish()

				h.Usecase.WaitingForETHBalancing(span) // ETH
			}(&wg)

			go func(wg *sync.WaitGroup) {
				span := h.Tracer.StartSpan("ETH.SendBTCToMaster")
				defer wg.Done()
				defer span.Finish()

				h.Usecase.WaitingForBalancing(span) // BTC
			}(&wg)

			go func(wg *sync.WaitGroup) {
				span := h.Tracer.StartSpan("ETH.SendBTCToMaster")
				defer wg.Done()
				defer span.Finish()

				h.Usecase.JobBtcSendBtcToMaster(span) // BTC
			}(&wg)
			time.Sleep(5 * time.Minute)
		}
	}()

	//waiting for minting
	go func() {
		var wg sync.WaitGroup
		//Waiting for Send
		for {
			wg.Add(2)
			go func(wg *sync.WaitGroup) {

				span := h.Tracer.StartSpan("BTC.WaitingForMinting")
				defer wg.Done()
				defer span.Finish()

				h.Usecase.WaitingForMinting(span)
				
			}(&wg)

			//TODO mint with ETH payment?
			go func(wg *sync.WaitGroup) {
				span := h.Tracer.StartSpan("ETH.WaitingForETHMinting")
				defer wg.Done()
				defer span.Finish()

				h.Usecase.WaitingForETHMinting(span)

			}(&wg)

			wg.Wait()
			time.Sleep(1 * time.Minute)
		}
	}()


	//Waiting for minted and send
	go func() {
		var wg sync.WaitGroup
		//Waiting for Send
		for {
			wg.Add(2)
			go func(wg *sync.WaitGroup) {

				span := h.Tracer.StartSpan("BTC.WaitingForMinted")
				defer wg.Done()
				defer span.Finish()

				h.Usecase.WaitingForMinted(span)

			}(&wg)

			//TODO mint with ETH payment?
			go func(wg *sync.WaitGroup) {
				span := h.Tracer.StartSpan("ETH.WaitingForETHMinted")
				defer wg.Done()
				defer span.Finish()

				h.Usecase.WaitingForETHMinted(span)

			}(&wg)

			wg.Wait()
			time.Sleep(1 * time.Minute)
		}
	}()

	c := cron.New()
	// cronjob to sync inscription index
	c.AddFunc("*/30 * * * *", func() {
		span := h.Tracer.StartSpan("DispatchCron.EveryTwentyMinutes")
		defer span.Finish()

		log := tracer.NewTraceLog()
		defer log.ToSpan(span)

		err := h.Usecase.SyncTokenInscribeIndex(span)
		if err != nil {
			log.Error("DispatchCron.OneMinute.GetTheCurrentBlockNumber", err.Error(), err)
		}
	})
	c.Start()
}
