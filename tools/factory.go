package tools

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"log"
	"os"
	"rederinghub.io/external/coin_market_cap"
	"rederinghub.io/external/etherscan"
	"rederinghub.io/external/mempool_space"
	"rederinghub.io/external/nfts"
	"rederinghub.io/external/ord_service"
	"rederinghub.io/internal/repository"
	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils/blockchain"
	"rederinghub.io/utils/btc"
	"rederinghub.io/utils/config"
	"rederinghub.io/utils/connections"
	"rederinghub.io/utils/delegate"
	discordclient "rederinghub.io/utils/discord"
	"rederinghub.io/utils/eth"
	"rederinghub.io/utils/global"
	"rederinghub.io/utils/googlecloud"
	_logger "rederinghub.io/utils/logger"
	"rederinghub.io/utils/oauth2service"
	"rederinghub.io/utils/redis"
	"rederinghub.io/utils/redisv9"
	"rederinghub.io/utils/slack"
)

var conf *config.Config
var logger _logger.Ilogger
var mongoConnection connections.IConnection
var ethClient *blockchain.Blockchain
var tcClient *blockchain.TcNetwork

func init() {
	c, err := config.NewConfig()
	if err != nil {
		log.Println("Error while reading config from env", err)
		panic(err)
	}

	if c.Debug {
		log.Println("Service RUN on DEBUG mode")
	}

	l := _logger.NewLogger(c.Debug)
	//l.AtLog().Logger.Info("config", zap.Any("config.NewConfig", c))

	mongoCnn := fmt.Sprintf("%s://%s:%s@%s/?retryWrites=true&w=majority", c.Databases.Mongo.Scheme, c.Databases.Mongo.User, c.Databases.Mongo.Pass, c.Databases.Mongo.Host)
	mongoDbConnection, err := connections.NewMongo(mongoCnn)
	if err != nil {
		l.AtLog().Logger.Error("Cannot connect mongoDB ", zap.Error(err))
		panic(err)
	}

	ethClient, err = blockchain.NewBlockchain(c.BlockchainConfig)
	if err != nil {
		l.AtLog().Logger.Error("Cannot connect ethClient ", zap.Error(err))
		panic(err)
	}

	tcClient, err = blockchain.NewTcNetwork(c.BlockchainConfig)
	if err != nil {
		log.Println("Cannot connect to tc client ", err)
		panic(err)
	}

	logger = l
	conf = c
	mongoConnection = mongoDbConnection
}

func startServer() *usecase.Usecase {
	fmt.Println("starting server ...")
	cache, redisClient := redis.NewRedisCache(conf.Redis)
	redisV9 := redisv9.NewClient(conf.Redis)
	r := mux.NewRouter()

	gcs, err := googlecloud.NewDataGCStorage(*conf)
	if err != nil {
		_logger.AtLog.Logger.Error("Cannot init gcs", zap.Error(err))
		return nil
	}
	s3Adapter := googlecloud.NewS3Adapter(googlecloud.S3AdapterConfig{
		BucketName: conf.Gcs.Bucket,
		Endpoint:   conf.Gcs.Endpoint,
		Region:     conf.Gcs.Region,
		AccessKey:  conf.Gcs.AccessKey,
		SecretKey:  conf.Gcs.SecretKey,
	}, redisClient)

	moralis := nfts.NewMoralisNfts(conf, cache)
	ord := ord_service.NewBtcOrd(conf, cache, "")

	ordForDeveloper := ord_service.NewBtcOrd(conf, cache, os.Getenv("ORD_SERVER_FOR_DEVELOPER"))

	covalent := nfts.NewCovalentNfts(conf)
	slack := slack.NewSlack(conf.Slack)
	rPubsub := redis.NewPubsubClient(conf.Redis)
	delegateService, err := delegate.NewService(ethClient.GetClient())
	if err != nil {
		_logger.AtLog.Logger.Error("error initializing delegate service", zap.Error(err))
		return nil
	}

	// init tc client public
	tcClientPublicWrap, err := ethclient.Dial(conf.BlockchainConfig.TCPublicEndpoint)
	if err != nil {
		_logger.AtLog.Logger.Error("error initializing tcClient public service", zap.Error(err))
		return nil
	}
	tcClientPublic := eth.NewClient(tcClientPublicWrap)

	// init tc client
	tcClientWrap, err := ethclient.Dial(conf.BlockchainConfig.TCEndpoint)
	if err != nil {
		_logger.AtLog.Logger.Error("error initializing tcClient service", zap.Error(err))
		return nil
	}
	tcClients := eth.NewClient(tcClientWrap)

	// init eth client
	ethClientWrap, err := ethclient.Dial(conf.BlockchainConfig.ETHEndpoint)
	if err != nil {
		_logger.AtLog.Logger.Error("error initializing ethClients service", zap.Error(err))
		return nil
	}
	ethClients := eth.NewClient(ethClientWrap)

	// init eth client for dex:
	ethDexClientWrap, err := ethclient.Dial(conf.BlockchainConfig.ETHEndpointDex)
	if err != nil {
		_logger.AtLog.Logger.Error("error initializing ethDexClientWrap service", zap.Error(err))
		return nil
	}
	ethClientDex := eth.NewClient(ethDexClientWrap)
	eScan := etherscan.NewEtherscanService(conf, cache)
	mpService := mempool_space.NewMempoolService(conf, cache)
	coinMKC := coin_market_cap.NewCoinMarketCap(conf, cache)

	// init blockcypher service:
	bsClient := btc.NewBlockcypherService(conf.BlockcypherAPI, "", conf.BlockcypherToken, &chaincfg.MainNetParams)

	// hybrid auth
	auth2Service := oauth2service.NewAuth2()
	g := global.Global{
		Logger:              logger,
		MuxRouter:           r,
		Conf:                conf,
		DBConnection:        mongoConnection,
		Cache:               cache,
		Auth2:               *auth2Service,
		GCS:                 gcs,
		S3Adapter:           s3Adapter,
		MoralisNFT:          *moralis,
		CovalentNFT:         *covalent,
		Blockchain:          *ethClient,
		TcNetwotkchain:      *tcClient,
		Slack:               *slack,
		DiscordClient:       discordclient.NewCLient(),
		Pubsub:              rPubsub,
		OrdService:          ord,
		OrdServiceDeveloper: ordForDeveloper,
		DelegateService:     delegateService,
		RedisV9:             redisV9,

		EthClient:          ethClients,     // for eth chain (for mint...)
		EthClientDex:       ethClientDex,   // for eth chain (dex)
		TcClient:           tcClients,      // for tc chain
		TcClientPublicNode: tcClientPublic, // for tc chain
		BsClient:           bsClient,       // for btc/blockcypher service
		EtherscanService:   eScan,
		MempoolService:     mpService,
		CoinMarketCap:      coinMKC,
	}

	repo, err := repository.NewRepository(&g)
	if err != nil {
		_logger.AtLog.Logger.Error("Cannot init repository", zap.Error(err))
		return nil
	}

	err = repo.CreateCollectionIndexes()
	if err != nil {
		_logger.AtLog.Logger.Error("CreateCollectionIndexes - Cannot created index ", zap.Error(err))
		// return
	}

	uc, err := usecase.NewUsecase(&g, *repo)
	if err != nil {
		_logger.AtLog.Errorf("LoadUsecases - Cannot init usecase", zap.Error(err))
		return nil
	}

	return uc
}

func StartFactory() *usecase.Usecase {
	return startServer()
}
