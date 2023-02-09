package usecase

import (
	"rederinghub.io/external/nfts"
	"rederinghub.io/external/ord_service"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/repository"
	"rederinghub.io/utils/blockchain"
	"rederinghub.io/utils/config"
	"rederinghub.io/utils/global"
	"rederinghub.io/utils/googlecloud"
	"rederinghub.io/utils/logger"
	"rederinghub.io/utils/mqttClient"
	"rederinghub.io/utils/oauth2service"
	"rederinghub.io/utils/redis"
	"rederinghub.io/utils/slack"
	"rederinghub.io/utils/tracer"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// global data to handle cronjob
type gData struct {
	AllListings []entity.MarketplaceListings
	AllOffers []entity.MarketplaceOffers
	AllTokens []entity.TokenUri
	AllProfile []entity.Users
	AllProjects []entity.Projects
}

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
	CovalentNft nfts.CovalentNfts
	Blockchain blockchain.Blockchain
	Slack slack.Slack
	OrdService *ord_service.BtcOrd
	gData gData
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
	u.CovalentNft = global.CovalentNFT
	u.Blockchain = global.Blockchain
	u.Slack = global.Slack
	u.OrdService = global.OrdService
	u.MqttClient = global.MqttClient
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

func (uc *Usecase) StartSpanWithoutRoot(name string) (opentracing.Span, *tracer.TraceLog) {
	span := uc.Tracer.StartSpan(name)
	log := tracer.NewTraceLog()
	return span, log
}

func (uc *Usecase) StartSpanFromInjecttion(tracingInjection map[string]string, name string) (opentracing.Span, *tracer.TraceLog) {
	spanCtx, _ := uc.Tracer.GetTrace().Extract(opentracing.TextMap, opentracing.TextMapCarrier(tracingInjection))
	span := uc.Tracer.GetTrace().StartSpan(name, ext.RPCServerOption(spanCtx))

	defer span.Finish()

	log := tracer.NewTraceLog()
	return span, log
}
