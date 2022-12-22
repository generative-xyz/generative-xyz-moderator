package usecase

import (
	"rederinghub.io/external/nfts"
	"rederinghub.io/internal/repository"
	"rederinghub.io/utils/config"
	"rederinghub.io/utils/global"
	"rederinghub.io/utils/googlecloud"
	"rederinghub.io/utils/logger"
	"rederinghub.io/utils/mqttClient"
	"rederinghub.io/utils/oauth2service"
	"rederinghub.io/utils/redis"
	"rederinghub.io/utils/tracer"

	"github.com/opentracing/opentracing-go"
)

type Usecase struct {
	Repo          repository.Repository
	Logger        logger.Ilogger
	Config        *config.Config
	Span          opentracing.Span
	Tracer        tracer.ITracer
	PubSub        redis.IPubSubClient
	Cache       redis.IRedisCache
	MqttClient mqttClient.IDeviceMqtt
	Auth2 oauth2service.Auth2
	GCS           googlecloud.IGcstorage
	MoralisNft nfts.MoralisNfts
}

func NewUsecase(global *global.Global, r repository.Repository) (*Usecase, error) {
	u := new(Usecase)
	u.Logger = global.Logger
	u.Config = global.Conf
	u.Tracer = global.Tracer
	u.Repo = r
	u.PubSub = global.Pubsub
	u.Cache = global.Cache
	u.Auth2 = global.Auth2
	u.GCS = global.GCS
	u.MoralisNft = global.MoralisNFT
	return u, nil
}

func (uc *Usecase) Version() string {
	return "Generateve-API Server - version 1"
}

func (uc *Usecase) SetSpan(span opentracing.Span) {
	uc.Span = span
}

func (uc *Usecase) StartSpan(name string,  rootSpan opentracing.Span) (opentracing.Span, *tracer.TraceLog) {
	span := uc.Tracer.StartSpanFromRoot(rootSpan, name)
	log := tracer.NewTraceLog()
	return span, log
}