package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"rederinghub.io/utils/config"

	"github.com/go-redis/redis"
	"github.com/opentracing/opentracing-go"
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
	ProducerWithTrace(rootSpan opentracing.Span, channel string, payload PubSubPayload) error
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

func (r *pubsub) ProducerWithTrace(rootSpan opentracing.Span, channel string, payload PubSubPayload) error {
	span := rootSpan.Tracer().StartSpan("ProduceAMessage", opentracing.ChildOf(rootSpan.Context()))
	defer span.Finish()

	//span.Tracer().Inject(span.Context(), opentracing.Binary, bytesData)
	textInjection := map[string]string{}

	channel = r.GetChannelName(channel)

	span.Tracer().Inject(span.Context(), opentracing.TextMap, opentracing.TextMapCarrier(textInjection))

	payload.InjectionTracing = textInjection
	bytesData, err := json.Marshal(payload)
	if err != nil {
		//Logger.log.Error(fmt.Sprintf("pubsub.Producer - Error - Can not marshal data for channel %s", channel), err)
		return err
	}
	stringData := string(bytesData)
	err = r.client.Publish(channel, stringData).Err()
	if err != nil {
		//Logger.log.Error(fmt.Sprintf("pubsub.Producer - Error - Published message into channel %s failure", channel), err)
		return err
	}

	//Logger.log.Info(fmt.Sprintf("pubsub.Producer - Success - Published message into channel %s successfully", channel))
	return nil
}
