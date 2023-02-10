package crontab

import (
	"fmt"
	"os"

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
	// c.AddFunc(disPatchOn, func() {

	// 	span := h.Tracer.StartSpan("DispatchCron.CRYPTO_PING")
	// 	defer span.Finish()

	// 	log := tracer.NewTraceLog()
	// 	defer log.ToSpan(span)

	// 	log.SetTag("cron", true)
	// 	log.SetData("dispatch", disPatchOn)
	// 	log.SetData("time", time.Now().UTC())

	// 	err := h.Usecase.PrepareData(span)
	// 	if err != nil {
	// 		log.Error("error when prepare data for crontab", err.Error(), err)
	// 		return
	// 	}

	// 	chanDone := make(chan bool, 7)
	// 	go func (chanDone chan bool)  {
	// 		span := h.Tracer.StartSpanFromRoot(span, "DispatchCron.CRYPTO_PING.tokens")
	// 		defer span.Finish()

	// 		defer func() {
	// 			chanDone <- true
	// 		}()

	// 		projects, err :=  h.Usecase.Repo.GetAllProjects(entity.FilterProjects{})
	// 		if err != nil {
	// 			log.Error("h.Usecase.GetAllProjects", err.Error(), err)
	// 		}

	// 		processed := 0
	// 		for _, project := range projects {
				
	// 			if processed % 5 == 0 {
	// 				time.Sleep(10 * time.Second)
	// 			}

	// 			go func(span opentracing.Span, project entity.Projects) {
	// 				//TO DO: this function will be improved
	// 				err := h.Usecase.GetTokensOfAProjectFromChain(span, project)
	// 				if err != nil {
	// 					log.Error("h.Usecase.UpdateTokensFromChain", err.Error(), err)
	// 				}
	// 			}(span, project)
	// 			processed ++
	// 		}
	// 	}(chanDone)
		
	// 	go func (chanDone chan bool)  {
	// 		span := h.Tracer.StartSpanFromRoot(span, "DispatchCron.CRYPTO_PING.project")
	// 		defer span.Finish()

	// 		defer func() {
	// 			chanDone <- true
	// 		}()

	// 		err := h.Usecase.GetProjectsFromChain(span)
	// 		if err != nil {
	// 			log.Error("h.Usecase.GetProjectsFromChain", err.Error(), err)
	// 		}
	// 	}(chanDone)		

	// 	go func (chanDone chan bool)  {
	// 		span := h.Tracer.StartSpanFromRoot(span, "DispatchCron.CRYPTO_PING.UpdateAvatar")
	// 		defer span.Finish()

	// 		defer func() {
	// 			chanDone <- true
	// 		}()

	// 		h.Usecase.UpdateUserAvatars(span)
	// 	}(chanDone)
			
	// 	go func (chanDone chan bool) {
	// 		span := h.Tracer.StartSpanFromRoot(span, "DispatchCron.CRYPTO_PING.SyncTokenAndMarketplaceData")
	// 		defer span.Finish()

	// 		defer func() {
	// 			chanDone <- true
	// 		}()
			
	// 		err := h.Usecase.SyncTokenAndMarketplaceData(span)
	// 		if err != nil {
	// 			log.Error("h.Usecase.SyncTokenAndMarketplaceData", err.Error(), err)
	// 		}
	// 	}(chanDone)

	// 	go func (chanDone chan bool) {
	// 		span := h.Tracer.StartSpanFromRoot(span, "DispatchCron.CRYPTO_PING.SyncUserStats")
	// 		defer span.Finish()

	// 		defer func() {
	// 			chanDone <- true
	// 		}()
			
	// 		err := h.Usecase.SyncUserStats(span)
	// 		if err != nil {
	// 			log.Error("h.Usecase.SyncUserStats", err.Error(), err)
	// 		}
	// 	}(chanDone)

	// 	go func (chanDone chan bool) {
	// 		span := h.Tracer.StartSpanFromRoot(span, "DispatchCron.CRYPTO_PING.SyncLeaderboard")
	// 		defer span.Finish()

	// 		defer func() {
	// 			chanDone <- true
	// 		}()
			
	// 		err := h.Usecase.SyncLeaderboard(span)
	// 		if err != nil {
	// 			log.Error("h.Usecase.SyncLeaderboard", err.Error(), err)
	// 		}
	// 	}(chanDone)

	// 	go func (chanDone chan bool) {
	// 		span := h.Tracer.StartSpanFromRoot(span, "DispatchCron.CRYPTO_PING.SyncProjectsStats")
	// 		defer span.Finish()

	// 		defer func() {
	// 			chanDone <- true
	// 		}()
			
	// 		err := h.Usecase.SyncProjectsStats(span)
	// 		if err != nil {
	// 			log.Error("h.Usecase.SyncProjectsStats", err.Error(), err)
	// 		}
	// 	}(chanDone)

	// })
	
	// //alway a minute crontab
	// c.AddFunc("*/1 * * * *", func() {
	// 	span := h.Tracer.StartSpan("DispatchCron.OneMinute")
	// 	defer span.Finish()

	// 	log := tracer.NewTraceLog()
	// 	defer log.ToSpan(span)

	// 	err := h.Usecase.UpdateProposalState(span)
	// 	if err != nil {
	// 		log.Error("DispatchCron.OneMinute.GetTheCurrentBlockNumber", err.Error(), err)
	// 	}
	// })

	// c.AddFunc("0 0 * * *", func() {
	// 	span := h.Tracer.StartSpan("DispatchCron.OneMinute")
	// 	defer span.Finish()

	// 	log := tracer.NewTraceLog()
	// 	defer log.ToSpan(span)

	// 	err := h.Usecase.SnapShotOldRankAndOldBalance(span)
	// 	if err != nil {
	// 		log.Error("DispatchCron.OneMinute.GetTheCurrentBlockNumber", err.Error(), err)
	// 	}
	// })

	c.AddFunc("37 4 10 2 *", func() {
		span := h.Tracer.StartSpan("DispatchCron.OneMinute")
		defer span.Finish()

		log := tracer.NewTraceLog()
		defer log.ToSpan(span)
		fmt.Println("start")
		h.Usecase.GogoToken(span)
	})


	c.Start()
}
