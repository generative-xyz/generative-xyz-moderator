package crontab

import (
	"context"
	"fmt"
	"os"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils/global"
	"rederinghub.io/utils/logger"
	"rederinghub.io/utils/redis"
)

type ScronHandler struct {
	Logger  logger.Ilogger
	Cache   redis.IRedisCache
	Usecase usecase.Usecase
}

func NewScronHandler(global *global.Global, uc usecase.Usecase) *ScronHandler {
	return &ScronHandler{
		Logger:  global.Logger,
		Cache:   global.Cache,
		Usecase: uc,
	}
}

func (h ScronHandler) StartServer() {
	c := cron.New()

	disPatchOn := os.Getenv("CRONTAB_SCHEDULE")
	logger.AtLog.Logger.Info(fmt.Sprintf("Cron is listerning: %s", disPatchOn))

	//check device's statues each 1 hours
	c.AddFunc(disPatchOn, func() {

		// logger.AtLog.Logger.Info("dispatch", zap.Any("disPatchOn", disPatchOn))
		// logger.AtLog.Logger.Info("time", zap.Any("time.Now().UTC()", time.Now().UTC()))

		// err := h.Usecase.PrepareData()
		// if err != nil {
		// 	logger.AtLog.Logger.Error("err", zap.Error(err))
		// 	return
		// }

		// err = h.Usecase.SyncUserStats()
		// if err != nil {
		// 	logger.AtLog.Logger.Error("err", zap.Error(err))
		// }

		// chanDone := make(chan bool, 1)
		// go func (chanDone chan bool)  {

		// 	defer func() {
		// 		chanDone <- true
		// 	}()

		// 	projects, err :=  h.Usecase.Repo.GetAllProjects(entity.FilterProjects{})
		// 	if err != nil {
		// 		logger.AtLog.Logger.Error("err", zap.Error(err))
		// 	}

		// 	processed := 0
		// 	for _, project := range projects {
		// 				if processed % 5 == 0 {
		// 			time.Sleep(10 * time.Second)
		// 		}

		// 		go func( project entity.Projects) {
		// 			//TO DO: this function will be improved
		// 			err := h.Usecase.GetTokensOfAProjectFromChain(project)
		// 			if err != nil {
		// 				logger.AtLog.Logger.Error("err", zap.Error(err))
		// 			}
		// 		}(project)
		// 		processed ++
		// 	}
		// }(chanDone)
		// go func (chanDone chan bool)  {
		// 	defer func() {
		// 		chanDone <- true
		// 	}()

		// 	err := h.Usecase.GetProjectsFromChain()
		// 	if err != nil {
		// 		logger.AtLog.Logger.Error("err", zap.Error(err))
		// 	}
		// }(chanDone)
		// go func (chanDone chan bool)  {
		// 	defer func() {
		// 		chanDone <- true
		// 	}()

		// 	h.Usecase.UpdateUserAvatars()
		// }(chanDone)
		// 	go func (chanDone chan bool) {
		// 	defer func() {
		// 		chanDone <- true
		// 	}()
		// 		err := h.Usecase.SyncTokenAndMarketplaceData()
		// 	if err != nil {
		// 		logger.AtLog.Logger.Error("err", zap.Error(err))
		// 	}
		// }(chanDone)

		// go func (chanDone chan bool) {
		// 	defer func() {
		// 		chanDone <- true
		// 	}()

		// }(chanDone)

		// go func (chanDone chan bool) {
		// 	defer func() {
		// 		chanDone <- true
		// 	}()
		// 		err := h.Usecase.SyncLeaderboard()
		// 	if err != nil {
		// 		logger.AtLog.Logger.Error("err", zap.Error(err))
		// 	}
		// }(chanDone)

		// go func (chanDone chan bool) {
		// 		defer func() {
		// 		chanDone <- true
		// 	}()
		// 		err := h.Usecase.SyncProjectsStats()
		// 	if err != nil {
		// 		logger.AtLog.Logger.Error("err", zap.Error(err))
		// 	}
		// }(chanDone)

	})

	//alway 10 minutes crontab
	// c.AddFunc("*/1 * * * *", func() {
	// 	err := h.Usecase.UpdateProposalState()
	// 	if err != nil {
	// 		logger.AtLog.Logger.Error("err", zap.Error(err))
	// 	}
	// })

	//At minute 0.
	c.AddFunc("@hourly", func() {
		h.Usecase.JobAggregateVolumns()
	})

	c.AddFunc("*/10 * * * *", func() {
		h.Usecase.JobAggregateReferral()
	})

	c.AddFunc("*/10 * * * *", func() {
		err := h.Usecase.JobSyncTraitStats()
		if err != nil {
			logger.AtLog.Logger.Error("error when sync trait stats", zap.Error(err))
		}
	})

	// At 05:00. UTC
	c.AddFunc("0 5 * * *", func() {
		err := h.Usecase.CalUserStats(context.Background())
		if err != nil {
			logger.AtLog.Logger.Error("CalUserStats failed", zap.Error(err))
		}
	})

	c.Start()
}
