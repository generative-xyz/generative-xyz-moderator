package crontabManager

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/robfig/cron/v3"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils/global"
	"rederinghub.io/utils/logger"
	"rederinghub.io/utils/redis"
)

type CrontabManager struct {
	JobKey  string
	Logger  logger.Ilogger
	Cache   redis.IRedisCache
	Usecase usecase.Usecase
}

func NewCrontabManager(jobKey string, global *global.Global, uc usecase.Usecase) *CrontabManager {
	return &CrontabManager{
		JobKey:  jobKey,
		Logger:  global.Logger,
		Cache:   global.Cache,
		Usecase: uc,
	}
}

var secondParser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)

func (h *CrontabManager) StartServer() {
	start := time.Now()
	h.Logger.Info(fmt.Sprintf("\nTime now (UTC): %v \n", start.UTC().Format(time.RFC3339)))

	c := cron.New(cron.WithLocation(time.UTC))
	h.Logger.Info(fmt.Sprintf("\nCron location: %v \n", c.Location()))

	eventList, err := h.Usecase.Repo.FindCronJobManagerByJobKey(h.JobKey)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("\nGet all items has error %v\n", err))
		panic(err)
	}

	if eventList == nil {
		h.Logger.Info("list event empty")
		return
		// panic(errors.New("List event is empty"))
	}

	for _, item := range eventList {
		eventName := fmt.Sprintf("[%v-%v]", item.Group, item.JobName)
		domain := item.FunctionName

		h.Logger.Info(fmt.Sprintf("\n=======Add task %v, webhook %v ======= \n", eventName, domain))

		err := h.Validation(item)
		if err != nil {
			// go e.slack.PostMsg(err.Error())
			// TODO: notify
			h.Logger.Info(fmt.Sprintf("\nHas error--> %v\n", err))
			go h.trackCronJobManagerLogs(fmt.Sprintf("[%v-%v]", item.Group, item.JobName), "Cron job manager -> StartServer", item.TableName(), item.Enabled, "h.Validation(item) - with err", err.Error())
			continue
		}

		entryId, err := h.AddTask(c, item)
		if err != nil {
			h.Logger.Info(fmt.Sprintf("\nAdding taskName %v, webhook %v, err %v\n", eventName, domain, err))
			continue
		}

		h.Logger.Info(fmt.Sprintf("\nAdd task completed with id %v of %v\n", entryId, eventName))
	}

	c.Start()
}

func (h *CrontabManager) Validation(cronJobManager entity.CronJobManager) error {
	eventName := fmt.Sprintf("[%v-%v-%v]", cronJobManager.JobKey, cronJobManager.Group, cronJobManager.JobName)

	if len(cronJobManager.Schedule) < 0 {
		msg := fmt.Sprintf("\nTask Name %v, Schedule %v invalid", eventName, cronJobManager.Schedule)
		return errors.New(msg)
	}

	if len(cronJobManager.FunctionName) < 0 {
		msg := fmt.Sprintf("\nTask Name %v, FunctionName %v invalid", eventName, cronJobManager.FunctionName)
		return errors.New(msg)
	}

	_, err := secondParser.Parse(cronJobManager.Schedule)
	if err != nil {
		msg := fmt.Sprintf("\nTask Name %v, Parse time error %v", eventName, err)
		return errors.New(msg)
	}

	return nil
}

func (h *CrontabManager) AddTask(c *cron.Cron, cronJobManager entity.CronJobManager) (cron.EntryID, error) {
	eventName := fmt.Sprintf("[%v-%v-%v]", cronJobManager.JobKey, cronJobManager.Group, cronJobManager.JobName)
	isRunning := false

	idEntries, err := c.AddFunc(cronJobManager.Schedule, func() {
		start := time.Now()
		h.Logger.Info(fmt.Sprintf("Time now (UTC): %v \n", start.UTC().Format(time.RFC3339)))

		h.Logger.Info(fmt.Sprintf("=======Start call task %v ======= \n", cronJobManager.FunctionName))

		if isRunning {
			h.Logger.Info(fmt.Sprintf("\nTaskName %v, webhook %v is running\n", eventName, cronJobManager.FunctionName))
			h.Usecase.Repo.UpdateCronJobManagerLastSatus(cronJobManager.UUID, "Job is waiting...Skip now!")
			return
		}

		isRunning = true
		defer func() {
			isRunning = false
		}()

		// check
		eventItem, err := h.Usecase.Repo.FindCronJobManagerByUUID(cronJobManager.UUID)
		if err != nil {
			h.Logger.Info(fmt.Sprintf("\nTaskName %v, webhook %v error: %v", eventName, cronJobManager.FunctionName, err.Error()))
			return
		}

		if eventItem == nil {
			h.Logger.Info(fmt.Sprintf("\nTaskName %v, webhook %v is disable\n", eventName, cronJobManager.FunctionName))
			return
		}

		if eventItem.Enabled == false {
			h.Usecase.Repo.UpdateCronJobManagerLastSatus(cronJobManager.UUID, "Job is paused!")
			return
		} else {
			h.Usecase.Repo.UpdateCronJobManagerLastSatus(cronJobManager.UUID, "Job is running...")
		}

		defer func() {
			if err := recover(); err != nil {
				h.Logger.Error(fmt.Sprintf("\nError calling %s->%s: %v", eventName, cronJobManager.FunctionName, err))
				go h.trackCronJobManagerLogs(fmt.Sprintf("[%v-%v]", eventItem.Group, eventItem.JobName), "Cron job manager -> AddTask", eventItem.TableName(), cronJobManager.FunctionName, "call the function - with err", err)
				return
			}
		}()

		// important: execute the job:
		meth := reflect.ValueOf(h.Usecase).MethodByName(cronJobManager.FunctionName)

		result := meth.Call(nil)
		errResult := result[0].Interface()
		if errResult == nil {
			h.Logger.Info(fmt.Sprintf("\nCalling %s->%s: is OK!!!!", eventName, cronJobManager.FunctionName))
		} else {
			h.Logger.Error(fmt.Sprintf("\nError calling %s->%s: %v", eventName, cronJobManager.FunctionName, errResult))
		}

		duration := time.Since(start)
		h.Logger.Info(fmt.Sprintf("\nEnd %v \n", duration.String()))

		if err != nil {
			h.Logger.Info(fmt.Sprintf("TaskName %v, webhook %v, err %v\n", eventName, cronJobManager.FunctionName, err))

			return
		}

		h.Logger.Info(fmt.Sprintf("=======TaskName %v, webhook %v completefully=======\n", eventName, cronJobManager.FunctionName))
	})

	return idEntries, err
}

func (h *CrontabManager) trackCronJobManagerLogs(id, name, table string, status interface{}, requestMsg interface{}, responseMsg interface{}) {
	trackData := &entity.CronJobManagerLogs{
		RecordID:    id,
		Name:        name,
		Table:       table,
		Status:      status,
		RequestMsg:  requestMsg,
		ResponseMsg: responseMsg,
	}
	err := h.Usecase.Repo.InsertCronJobManagerLogs(trackData)
	if err != nil {
		fmt.Printf("trackDeveloperInscribeHistory.%s.Error:%s", name, err.Error())
	}

}
