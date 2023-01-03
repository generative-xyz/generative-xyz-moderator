package crontab

import (
	"fmt"
	"os"
	"time"

	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils/global"
	"rederinghub.io/utils/logger"
	"rederinghub.io/utils/redis"
	"rederinghub.io/utils/tracer"

	"gopkg.in/robfig/cron.v2"
)


type ScronHandler struct {
	Logger   logger.Ilogger
	Tracer   tracer.ITracer
	Cache    redis.IRedisCache
	Usecase    usecase.Usecase
}

func NewScronHandler(global *global.Global, uc usecase.Usecase) *ScronHandler {
	return &ScronHandler{
		Logger: global.Logger,
		Tracer: global.Tracer,
		Cache: global.Cache,
		Usecase: uc,
	}
}


func (h ScronHandler) StartServer() {
	c := cron.New()

	disPatchOn := os.Getenv("CRONTAB_SCHEDULE")
	h.Logger.Info(fmt.Sprintf("Cron is listerning: %s", disPatchOn))
	//check device's statues each 1 hours
	c.AddFunc(disPatchOn, func() {

		span := h.Tracer.StartSpan("DispatchCron.CRYPTO_PING")
		defer span.Finish()

		log := tracer.NewTraceLog()
		defer log.ToSpan(span)

		log.SetTag("cron", true)
		log.SetData("dispatch", disPatchOn)
		log.SetData("time", time.Now().UTC())

		err := h.Usecase.GetProjectsFromChain(span)
		if err != nil {
			log.Error("h.Usecase.UpdateProductPrice", err.Error(), err)
		}

	})

	c.Start()
}
