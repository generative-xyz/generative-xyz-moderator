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

	

	go func() {
		var wg sync.WaitGroup
		//Waiting for balance + Mint
		for {
			wg.Add(3)
			go func(wg *sync.WaitGroup) {
				span := h.Tracer.StartSpan("BTC.WaitingForBalancing")
				defer wg.Done()
				defer span.Finish()

				h.Usecase.WaitingForBalancing(span) // BTC
			}(&wg)

			go func(wg *sync.WaitGroup) {
				span := h.Tracer.StartSpan("ETH.WaitingForETHBalancing")
				defer wg.Done()
				defer span.Finish()

				h.Usecase.WaitingForETHBalancing(span) // ETH
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
}
