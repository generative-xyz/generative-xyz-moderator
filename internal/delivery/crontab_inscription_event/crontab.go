package crontab_scription_events

import (
	"gopkg.in/robfig/cron.v2"
	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils/global"
	"rederinghub.io/utils/logger"
	"rederinghub.io/utils/tracer"
)

type ScronTrendingHandler struct {
	Logger  logger.Ilogger
	Tracer  tracer.ITracer
	Usecase usecase.Usecase
}

func NewScronTrendingHandler(global *global.Global, uc usecase.Usecase) *ScronTrendingHandler {
	return &ScronTrendingHandler{
		Logger:  global.Logger,
		Tracer:  global.Tracer,
		Usecase: uc,
	}
}

func (h ScronTrendingHandler) StartServer() {
	c := cron.New()
	// cronjob to sync projects trending
	c.AddFunc("*/10 * * * *", func() {
		span := h.Tracer.StartSpan("DispatchCron.EveryTenMinutes.InscriptionEvents")
		defer span.Finish()

		log := tracer.NewTraceLog()
		defer log.ToSpan(span)

		err := h.Usecase.SyncInscriptionEvents(span)
		if err != nil {
			log.Error("DispatchCron.EveryTenMinutes.InscriptionEvents", err.Error(), err)
		}
	})
	c.Start()
}
