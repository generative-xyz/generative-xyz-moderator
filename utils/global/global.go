package global

import (
	"rederinghub.io/external/nfts"
	"rederinghub.io/external/ord_service"
	"rederinghub.io/utils/blockchain"
	"rederinghub.io/utils/config"
	_pConnection "rederinghub.io/utils/connections"
	"rederinghub.io/utils/delegate"
	discordclient "rederinghub.io/utils/discord"
	"rederinghub.io/utils/googlecloud"
	_logger "rederinghub.io/utils/logger"
	"rederinghub.io/utils/oauth2service"
	"rederinghub.io/utils/redis"
	"rederinghub.io/utils/redisv9"
	"rederinghub.io/utils/slack"

	"github.com/gorilla/mux"
)

type Global struct {
	Conf                *config.Config
	Logger              _logger.Ilogger
	MuxRouter           *mux.Router
	DBConnection        _pConnection.IConnection
	Cache               redis.IRedisCache
	CacheAuthService    redis.IRedisCache
	RedisV9             redisv9.Client
	Pubsub              redis.IPubSubClient
	Auth2               oauth2service.Auth2
	GCS                 googlecloud.IGcstorage
	S3Adapter           googlecloud.S3Adapter
	MoralisNFT          nfts.MoralisNfts
	CovalentNFT         nfts.CovalentNfts
	OrdService          *ord_service.BtcOrd
	OrdServiceDeveloper *ord_service.BtcOrd
	Blockchain          blockchain.Blockchain
	Slack               slack.Slack
	DiscordClient       *discordclient.Client
	DelegateService     *delegate.Service
}
