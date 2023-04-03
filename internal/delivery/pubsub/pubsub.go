package pubsub

import (
	"fmt"

	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils"
	"rederinghub.io/utils/logger"
	"rederinghub.io/utils/redis"
)

type PubsubHandler struct {
	usecase usecase.Usecase
	pubsub     redis.IPubSubClient
	log logger.Ilogger
}

func NewPubsubHandler(usecase usecase.Usecase, pubsub redis.IPubSubClient, log logger.Ilogger) *PubsubHandler {
	return &PubsubHandler{
		usecase: usecase,
		pubsub:     pubsub,
		log: log,
	}
}

func (h PubsubHandler) StartServer() {
	names := []string{
		utils.PUBSUB_TOKEN_THUMBNAIL,
		utils.PUBSUB_PROJECT_UNZIP,
	}

	h.pubsub.GetChannelNames(names...)
	pubsub := h.pubsub.GetClient().Subscribe(names...)

	// Wait for confirmation that subscription is created before publishing anything.
	_, err := pubsub.Receive()
	if err != nil {
		panic(err)
	}

	logger.AtLog.Info(fmt.Sprintf("pubsubHandler.SubscribeMessageRoute - Listen on channel name: %s ", names))
	// Go channel which receives messages.
	ch := pubsub.Channel()
	for msg := range ch {

		chanName := msg.Channel
		payload, _, err := h.pubsub.Parsepayload(msg.Payload)
		if err != nil {
			continue
		}

		switch chanName {
		case h.pubsub.GetChannelName(utils.PUBSUB_TOKEN_THUMBNAIL):
			h.usecase.PubSubCreateTokenThumbnail(chanName, payload)
			break
		case h.pubsub.GetChannelName(utils.PUBSUB_PROJECT_UNZIP):
			h.usecase.PubSubProjectUnzip(chanName, payload)
			break
		}}
	<-ch
	return
}

