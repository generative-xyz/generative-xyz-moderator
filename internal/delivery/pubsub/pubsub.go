package pubsub

import (
	"fmt"
	redis2 "github.com/go-redis/redis"
	"go.uber.org/zap"
	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils"
	"rederinghub.io/utils/logger"
	"rederinghub.io/utils/redis"
	"time"
)

type Handler struct {
	useCase usecase.Usecase
	pubsub  redis.IPubSubClient
	log     logger.Ilogger
}

func NewPubsubHandler(useCase usecase.Usecase, pubsub redis.IPubSubClient, log logger.Ilogger) *Handler {
	return &Handler{
		useCase: useCase,
		pubsub:  pubsub,
		log:     log,
	}
}

func (h Handler) StartServer() {
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
			if errCount > 0 {
				time.Sleep(1 * time.Second)
			}
			errCount++
			continue
		}
		errCount = 0
		if m, ok := msg.(*redis2.Message); ok {
			processCount++
			go func(message *redis2.Message) {
				h.handlerMessage(message)
				processCount--
			}(m)

			for processCount >= 5 {
				time.Sleep(1 * time.Second)
			}
		}
	}
	return
}
func (h Handler) handlerMessage(msg *redis2.Message) error {
	defer func() {
		if rcv := recover(); rcv != nil {
			logger.AtLog.Error("panic error", zap.Any("recover", rcv))
		}
	}()

	chanName := msg.Channel
	logger.AtLog.Info("pubsubHandler.handlerMessage", zap.String("channel", chanName), zap.Any("payload", msg.Payload))

	payload, tracingInjection, err := h.pubsub.Parsepayload(msg.Payload)
	if err != nil {
		return err
	}
	switch chanName {
	case h.pubsub.GetChannelName(utils.PUBSUB_TOKEN_THUMBNAIL):
		h.useCase.PubSubCreateTokenThumbnail(tracingInjection, chanName, payload)
	case h.pubsub.GetChannelName(utils.PUBSUB_PROJECT_UNZIP):
		h.useCase.PubSubProjectUnzip(tracingInjection, chanName, payload)
	case h.pubsub.GetChannelName(utils.PUBSUB_CAPTURE_THUMBNAIL):
		h.useCase.PubSubCaptureThumbnail(tracingInjection, chanName, payload)
	}
	return nil
}
