package crontab_btc

import (
	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils/global"
	"rederinghub.io/utils/logger"
	"rederinghub.io/utils/redis"
	"rederinghub.io/utils/tracer"

	"gopkg.in/robfig/cron.v2"
)


type ScronBTCHandler struct {
	Logger   logger.Ilogger
	Tracer   tracer.ITracer
	Cache    redis.IRedisCache
	Usecase    usecase.Usecase
}

func NewScronBTCHandler(global *global.Global, uc usecase.Usecase) *ScronBTCHandler {
	return &ScronBTCHandler{
		Logger: global.Logger,
		Tracer: global.Tracer,
		Cache: global.Cache,
		Usecase: uc,
	}
}


func (h ScronBTCHandler) StartServer() {
	span := h.Tracer.StartSpan("ScronBTCHandler.DispatchCron.OneMinute")
	defer span.Finish()

	c := cron.New()
	c.AddFunc("*/5 * * * *", func() {
		span := h.Tracer.StartSpan("ScronBTCHandler.DispatchCron.OneMinute")
		defer span.Finish()

		log := tracer.NewTraceLog()
		defer log.ToSpan(span)

		go func(){
			h.Usecase.WaitingForBalancing(span)
		}()
		
		go func(){
			h.Usecase.WaitingForMinted(span)
			
		}()
	
	})

	c.Start()
}
