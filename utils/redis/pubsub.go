package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
	"rederinghub.io/utils/config"
)

type PubSubPayload struct {
	Data             interface{}
	InjectionTracing map[string]string
}

type IPubSubClient interface {
	Producer(channel string, payload PubSubPayload) error
	GetClient() *redis.Client
	GetChannelName(name string) string
	GetChannelNames(channels ...string)
	Parsepayload(payload string) (interface{}, map[string]string, error)
}

type pubsub struct {
	cfg    config.RedisConfig
	client *redis.Client
	ctx    context.Context
}

func NewPubsubClient(cfg config.RedisConfig) *pubsub {
	r := new(pubsub)
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password, // no password set
		DB:       0,            // use default DB
	})
	r.cfg = cfg
	r.client = rdb
	r.ctx = ctx
	return r
}

func (r *pubsub) GetClient() *redis.Client {
	return r.client
}

func (r *pubsub) GetChannelName(channel string) string {
	return fmt.Sprintf("%s:%s", r.cfg.ENV, channel)
}

func (r *pubsub) GetChannelNames(channels ...string) {
	if len(channels) > 0 {
		for i, name := range channels {
			channels[i] = r.GetChannelName(name)
		}
	}
}

func (r *pubsub) Producer(channel string, payload PubSubPayload) error {
	channel = r.GetChannelName(channel)
	bytesData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	err = r.client.Publish(channel, string(bytesData)).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *pubsub) Parsepayload(payload string) (interface{}, map[string]string, error) {
	dataBytes := []byte(payload)
	p := PubSubPayload{}
	err := json.Unmarshal(dataBytes, &p)
	if err != nil {
		return nil, nil, err
	}
	return p.Data, p.InjectionTracing, nil
}
