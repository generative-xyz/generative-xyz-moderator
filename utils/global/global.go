package global

import (
<<<<<<< HEAD
	"rederinghub.io/external/nfts"
=======
	"rederinghub.io/utils/blockchain"
>>>>>>> add tx-consumer
	"rederinghub.io/utils/config"
	_pConnection "rederinghub.io/utils/connections"
	"rederinghub.io/utils/googlecloud"
	_logger "rederinghub.io/utils/logger"
	"rederinghub.io/utils/oauth2service"
	"rederinghub.io/utils/redis"
	"rederinghub.io/utils/tracer"

	"github.com/gorilla/mux"
)

type Global struct {
	Conf         *config.Config
	Logger       _logger.Ilogger
	MuxRouter    *mux.Router
	DBConnection _pConnection.IConnection
	Cache        redis.IRedisCache
	CacheAuthService        redis.IRedisCache
	Pubsub       redis.IPubSubClient
	Tracer       tracer.ITracer
	Auth2 oauth2service.Auth2
	GCS           googlecloud.IGcstorage
	MoralisNFT nfts.MoralisNfts
	Blockchain blockchain.Blockchain
}
