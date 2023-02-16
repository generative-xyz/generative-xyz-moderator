package crontab_btc

import (
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
	go func() {
		//it does not call our ORD server
		for {
			h.Usecase.JobBtcSendBtcToMaster() // BTC
			
			time.Sleep(5 * time.Minute)
		}
	}()

	//waiting for minting - CALL our ord server
	go func() {
		//All process will be >= 30 minutes
		for {
			h.Usecase.WaitingForBalancing() // BTC

			h.Usecase.WaitingForETHBalancing() //ETH

			//Sleetp 5 minutes after check balancing
			time.Sleep(5 * time.Minute)

			h.Usecase.WaitingForMinting() // BTC
				
			h.Usecase.WaitingForETHMinting() //ETH


			//Sleep 15 minutes after mint
			time.Sleep(15 * time.Minute)
			
			h.Usecase.WaitingForMinted() // BTC

			h.Usecase.WaitingForETHMinted() //ETH

			//Sleep 5 minutes after mint
			time.Sleep(5 * time.Minute)
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
