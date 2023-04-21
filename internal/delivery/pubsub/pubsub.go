package pubsub

import (
	"fmt"
	redis2 "github.com/redis/go-redis/v9"
	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils"
	"rederinghub.io/utils/logger"
	"rederinghub.io/utils/redis"
	"time"
)

type PubsubHandler struct {
	usecase usecase.Usecase
	pubsub  redis.IPubSubClient
	log     logger.Ilogger
}

func NewPubsubHandler(usecase usecase.Usecase, pubsub redis.IPubSubClient, log logger.Ilogger) *PubsubHandler {
	return &PubsubHandler{
		usecase: usecase,
		pubsub:  pubsub,
		log:     log,
	}
}

func (h PubsubHandler) StartServer() {
	names := []string{
		utils.PUBSUB_TOKEN_THUMBNAIL,
		utils.PUBSUB_PROJECT_UNZIP,
		utils.PUBSUB_CAPTURE_THUMBNAIL,
	}

	h.pubsub.GetChannelNames(names...)
	pubsub := h.pubsub.GetClient().Subscribe(names...)

	// Wait for confirmation that subscription is created before publishing anything.
	_, err := pubsub.Receive()
	if err != nil {
		panic(err)
	}

	errCount := 0
	processCount := 0
	logger.AtLog.Info(fmt.Sprintf("pubsubHandler.SubscribeMessageRoute - Listen on channel name: %s ", names))
	for {
		msg, err := pubsub.Receive()
		if err != nil {
			if err == redis2.ErrClosed {
				panic("redis: pubsub connection closed")
				return
			}
			if errCount > 0 {
				time.Sleep(1 * time.Second)
			}
			errCount++
			continue
		}
		errCount = 0
		switch msg := msg.(type) {
		case *redis2.Subscription, *redis2.Pong:
		case *redis2.Message:
			processCount++
			go func(message *redis2.Message) {
				h.handlerMessage(message)
				processCount--
			}(msg)

			for processCount >= 5 {
				time.Sleep(1 * time.Second)
			}
		default:
			logger.AtLog.Info(fmt.Sprintf("pubsubHandler.SubscribeMessageRoute - unknown message type: %s ", msg))
		}
	}
	return
}
func (h PubsubHandler) handlerMessage(msg *redis2.Message) error {
	chanName := msg.Channel
	payload, tracingInjection, err := h.pubsub.Parsepayload(msg.Payload)
	if err != nil {
		return err
	}
	switch chanName {
	case h.pubsub.GetChannelName(utils.PUBSUB_TOKEN_THUMBNAIL):
		h.usecase.PubSubCreateTokenThumbnail(tracingInjection, chanName, payload)
	case h.pubsub.GetChannelName(utils.PUBSUB_PROJECT_UNZIP):
		h.usecase.PubSubProjectUnzip(tracingInjection, chanName, payload)
	case h.pubsub.GetChannelName(utils.PUBSUB_CAPTURE_THUMBNAIL):
		h.usecase.PubSubCaptureThumbnail(tracingInjection, chanName, payload)
	}
	return nil
}
