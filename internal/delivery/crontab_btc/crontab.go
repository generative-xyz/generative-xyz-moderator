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

	

	go func() {
		var wg sync.WaitGroup
		span := h.Tracer.StartSpan("ScronBTCHandler.CheckBalance")
		defer span.Finish()

		log := tracer.NewTraceLog()
		defer log.ToSpan(span)

		//Waiting for balance + Mint
		for {
			wg.Add(3)
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
				span := h.Tracer.StartSpanFromRoot(rootSpan, "ETH.SendBTCToMaster")
				defer wg.Done()
				defer span.Finish()

				h.Usecase.JobBtcSendBtcToMaster(span) // BTC
			}(span, &wg)

			log.SetData("wait.CheckBlance", "wait")
			wg.Wait()
			time.Sleep(5 * time.Minute)
		}

	}()

	go func() {
		var wg sync.WaitGroup
		span := h.Tracer.StartSpan("ScronBTCHandler.SendNft")
		defer span.Finish()

		log := tracer.NewTraceLog()
		defer log.ToSpan(span)

		//Waiting for Send
		for {
			wg.Add(2)
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

			log.SetData("wait.SendNft", "wait")
			wg.Wait()
			time.Sleep(1 * time.Minute)
		}
	}()
}
