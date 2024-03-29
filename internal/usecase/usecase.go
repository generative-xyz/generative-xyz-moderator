package usecase

import (
	"rederinghub.io/external/coin_market_cap"
	"rederinghub.io/external/etherscan"
	"rederinghub.io/external/mempool_space"
	"rederinghub.io/external/nfts"
	"rederinghub.io/external/ord_service"
	"rederinghub.io/external/token_explorer"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/repository"
	"rederinghub.io/utils/blockchain"
	"rederinghub.io/utils/btc"
	"rederinghub.io/utils/config"
	"rederinghub.io/utils/delegate"
	discordclient "rederinghub.io/utils/discord"
	"rederinghub.io/utils/eth"
	"rederinghub.io/utils/global"
	"rederinghub.io/utils/googlecloud"
	"rederinghub.io/utils/logger"
	"rederinghub.io/utils/oauth2service"
	"rederinghub.io/utils/redis"
	"rederinghub.io/utils/redisv9"
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
	Repo                repository.Repository
	Logger              logger.Ilogger
	Config              *config.Config
	PubSub              redis.IPubSubClient
	Cache               redis.IRedisCache
	Auth2               oauth2service.Auth2
	GCS                 googlecloud.IGcstorage
	S3Adapter           googlecloud.S3Adapter
	MoralisNft          nfts.MoralisNfts
	CovalentNft         nfts.CovalentNfts
	Blockchain          blockchain.Blockchain
	Slack               slack.Slack
	DiscordClient       *discordclient.Client
	OrdService          *ord_service.BtcOrd
	OrdServiceDeveloper *ord_service.BtcOrd
	gData               gData
	DelegateService     *delegate.Service
	RedisV9             redisv9.Client

	TcClient, TcClientPublicNode, EthClient, EthClientDex *eth.Client
	BsClient                                              *btc.BlockcypherService
	EtherscanService                                      *etherscan.EtherscanService
	MempoolService                                        *mempool_space.MempoolService
	CoinMarketCap                                         *coin_market_cap.CoinMarketCap
	TokenExplorer                                         *token_explorer.TokenExplorer
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
	u.DiscordClient = global.DiscordClient
	u.OrdService = global.OrdService
	u.OrdServiceDeveloper = global.OrdServiceDeveloper
	u.DelegateService = global.DelegateService
	u.RedisV9 = global.RedisV9

	u.TcClient = global.TcClient
	u.TcClientPublicNode = global.TcClientPublicNode
	u.EthClientDex = global.EthClientDex
	u.EthClient = global.EthClient
	u.BsClient = global.BsClient
	u.EtherscanService = global.EtherscanService
	u.MempoolService = global.MempoolService
	u.CoinMarketCap = global.CoinMarketCap
	u.TokenExplorer = global.TokenExplorer

	return u, nil
}

func (uc *Usecase) Version() string {
	return "Generateve-API Server - version 1"
}
