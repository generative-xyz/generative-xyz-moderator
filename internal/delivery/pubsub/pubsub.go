package pubsub

import (
	"fmt"
	redis2 "github.com/go-redis/redis"
	"go.uber.org/zap"
	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils"
	"rederinghub.io/utils/logger"
	"rederinghub.io/utils/redis"
	"sync"
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
	var wg sync.WaitGroup
	processing := 0
	maxProcessing := 5

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
		m, ok := msg.(*redis2.Message)
		if ok {
			wg.Add(1)
			processing++
			go h.worker(&wg, m)

			if processing > 0 && processing%maxProcessing == 0 {
				wg.Wait()
				processing = 0
			}

		}
	}
	return
}

func (h Handler) worker(wg *sync.WaitGroup, message *redis2.Message) {
	defer wg.Done()
	h.handlerMessage(message)

}

func (h Handler) handlerMessage(msg *redis2.Message) error {
	defer func() {
		if rcv := recover(); rcv != nil {
			logger.AtLog.Logger.Error("handlerMessage", zap.Any("recover", rcv))
		}
	}()

	chanName := msg.Channel

	payload, tracingInjection, err := h.pubsub.Parsepayload(msg.Payload)
	if err != nil {
		return err
	}

	logger.AtLog.Logger.Info(fmt.Sprintf("handlerMessage - %s", chanName), zap.String("chanName", chanName), zap.Any("payload", payload))

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
