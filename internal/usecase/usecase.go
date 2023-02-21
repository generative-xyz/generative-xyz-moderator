package usecase

import (
	"rederinghub.io/external/nfts"
	"rederinghub.io/external/ord_service"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/repository"
	"rederinghub.io/utils/blockchain"
	"rederinghub.io/utils/config"
	"rederinghub.io/utils/delegate"
	"rederinghub.io/utils/global"
	"rederinghub.io/utils/googlecloud"
	"rederinghub.io/utils/logger"
	"rederinghub.io/utils/oauth2service"
	"rederinghub.io/utils/redis"
	"rederinghub.io/utils/slack"
)

// global data to handle cronjob
type gData struct {
	AllListings []entity.MarketplaceListings
	AllOffers   []entity.MarketplaceOffers
	AllTokens   []entity.TokenUri
	AllProfile  []entity.Users
	AllProjects []entity.Projects
}

type Usecase struct {
	Repo            repository.Repository
	Logger          logger.Ilogger
	Config          *config.Config
	PubSub          redis.IPubSubClient
	Cache           redis.IRedisCache
	Auth2           oauth2service.Auth2
	GCS             googlecloud.IGcstorage
	S3Adapter  googlecloud.S3Adapter
	MoralisNft      nfts.MoralisNfts
	CovalentNft     nfts.CovalentNfts
	Blockchain      blockchain.Blockchain
	Slack           slack.Slack
	OrdService      *ord_service.BtcOrd
	gData           gData
	DelegateService *delegate.Service
}

func NewUsecase(global *global.Global, r repository.Repository) (*Usecase, error) {
	u := new(Usecase)
	u.Logger = global.Logger
	u.Config = global.Conf
	u.Repo = r
	u.PubSub = global.Pubsub
	u.Cache = global.Cache
	u.Auth2 = global.Auth2
	u.GCS = global.GCS
	u.S3Adapter = global.S3Adapter
	u.MoralisNft = global.MoralisNFT
	u.CovalentNft = global.CovalentNFT
	u.Blockchain = global.Blockchain
	u.Slack = global.Slack
	u.OrdService = global.OrdService
	u.DelegateService = global.DelegateService
	return u, nil
}

func (uc *Usecase) Version() string {
	return "Generateve-API Server - version 1"
}
