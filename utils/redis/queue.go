package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hibiken/asynq"
	"rederinghub.io/utils/config"
)

type QueuePayload struct {
	Data             interface{}
	InjectionTracing map[string]string
}

type QueueClient struct {
	cfg    config.RedisConfig
	client *asynq.Client
	ctx    context.Context
}

func NewQueueClient(cfg config.RedisConfig) *QueueClient {
	r := new(QueueClient)
	ctx := context.Background()

	dbInt, err := strconv.Atoi(cfg.DB)
	if err != nil {
		dbInt = 1
	}

	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     cfg.Address,
		DB:       dbInt,
		Password: cfg.Password,
	})
	//defer client.Close()

	r.cfg = cfg
	r.client = client
	r.ctx = ctx
	return r
}

func (r *QueueClient) GetClient() *asynq.Client {
	return r.client
}

func (r *QueueClient) GetChannelName(channel string) string {
	return fmt.Sprintf("%s:%s", r.cfg.ENV, channel)
}

func (r *QueueClient) GetChannelNames(channels ...string) {
	if len(channels) > 0 {
		for i, name := range channels {
			channels[i] = r.GetChannelName(name)
		}
	}
}

func (r *QueueClient) CreateTask(channel string, payload interface{}, retries int, priority string) (*asynq.Task, error) {
	payLoadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	// task options can be passed to NewTask, which can be overridden at enqueue time.
	return asynq.NewTask(channel, 
		payLoadBytes,
		asynq.MaxRetry(retries), 
		asynq.Queue(priority)), 
	nil
}

func (r *QueueClient) CreateSchedule(channel string, payload interface{}, priority string) (*asynq.TaskInfo, error) {
	task, err := r.CreateTask(channel, payload, 20, priority)
	if err != nil {
		return nil, err
	}

	info, err := r.client.Enqueue(task)
    if err != nil {
        return nil, err
    }

	return info, nil
}

