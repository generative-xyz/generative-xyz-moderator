package crontab

import (
	"fmt"
	"os"
	"time"

	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils/global"
	"rederinghub.io/utils/logger"
	"rederinghub.io/utils/redis"

	"gopkg.in/robfig/cron.v2"
)


type ScronHandler struct {
	Logger   logger.Ilogger
	Cache    redis.IRedisCache
	Usecase    usecase.Usecase
}

func NewScronHandler(global *global.Global, uc usecase.Usecase) *ScronHandler {
	return &ScronHandler{
		Logger: global.Logger,
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

		h.Logger.Info("dispatch", disPatchOn)
		h.Logger.Info("time", time.Now().UTC())

		err := h.Usecase.PrepareData()
		if err != nil {
			h.Logger.Error(err)
			return
		}

		// chanDone := make(chan bool, 1)
		// go func (chanDone chan bool)  {

		// 	defer func() {
		// 		chanDone <- true
		// 	}()

		// 	projects, err :=  h.Usecase.Repo.GetAllProjects(entity.FilterProjects{})
		// 	if err != nil {
		// 		h.Logger.Error(err)
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
		// 				h.Logger.Error(err)
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
		// 		h.Logger.Error(err)
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
		// 		h.Logger.Error(err)
		// 	}
		// }(chanDone)

		// go func (chanDone chan bool) {
		// 	defer func() {
		// 		chanDone <- true
		// 	}()
		err = h.Usecase.SyncUserStats()
		if err != nil {
			h.Logger.Error(err)
		}
		// }(chanDone)

		// go func (chanDone chan bool) {
		// 	defer func() {
		// 		chanDone <- true
		// 	}()
		// 		err := h.Usecase.SyncLeaderboard()
		// 	if err != nil {
		// 		h.Logger.Error(err)
		// 	}
		// }(chanDone)

		// go func (chanDone chan bool) {
		// 		defer func() {
		// 		chanDone <- true
		// 	}()
		// 		err := h.Usecase.SyncProjectsStats()
		// 	if err != nil {
		// 		h.Logger.Error(err)
		// 	}
		// }(chanDone)

	})
//alway a minute crontab
	// c.AddFunc("*/1 * * * *", func() {
	// 	err := h.Usecase.UpdateProposalState()
	// 	if err != nil {
	// 		h.Logger.Error(err)
	// 	}
	// })

	// c.AddFunc("0 0 * * *", func() {
	// 	err := h.Usecase.SnapShotOldRankAndOldBalance()
	// 	if err != nil {
	// 		h.Logger.Error(err)
	// 	}
	// })

	c.Start()
}
