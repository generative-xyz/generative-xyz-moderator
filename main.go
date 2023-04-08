package main

import (
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"
	"rederinghub.io/utils/btc"
	"rederinghub.io/utils/delegate"
	"rederinghub.io/utils/eth"
	"rederinghub.io/utils/redisv9"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
	migrate "github.com/xakep666/mongo-migrate"
	"rederinghub.io/external/nfts"
	"rederinghub.io/external/ord_service"
	"rederinghub.io/internal/repository"
	"rederinghub.io/internal/usecase"
	_ "rederinghub.io/mongo/migrate"
	"rederinghub.io/utils/blockchain"
	"rederinghub.io/utils/config"
	"rederinghub.io/utils/connections"
	discordclient "rederinghub.io/utils/discord"
	"rederinghub.io/utils/global"
	"rederinghub.io/utils/googlecloud"
	_logger "rederinghub.io/utils/logger"
	"rederinghub.io/utils/oauth2service"
	"rederinghub.io/utils/redis"
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

// @title Generative.xyz APIs
// @version 1.0.0
// @description This is a sample server Generative.

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @BasePath /rederinghub.io/v1
func main() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered. Error:\n", r)
		}
	}()

	// log.Println("init sentry ...")
	// sentry.InitSentry(conf)
	startServer()
}

func startServer() {
	logger.AtLog().Logger.Info("starting server ...")
	cache, redisClient := redis.NewRedisCache(conf.Redis)
	redisV9 := redisv9.NewClient(conf.Redis)
	r := mux.NewRouter()

	gcs, err := googlecloud.NewDataGCStorage(*conf)
	if err != nil {
		_logger.AtLog.Logger.Error("Cannot init gcs", zap.Error(err))
		return
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
		return
	}

	// init tc client
	tcClientWrap, err := ethclient.Dial(conf.BlockchainConfig.TCEndpoint)
	if err != nil {
		_logger.AtLog.Logger.Error("error initializing tcClient service", zap.Error(err))
		return
	}
	tcClients := eth.NewClient(tcClientWrap)

	// init eth client
	ethClientWrap, err := ethclient.Dial(conf.BlockchainConfig.ETHEndpoint)
	if err != nil {
		_logger.AtLog.Logger.Error("error initializing ethClients service", zap.Error(err))
		return
	}
	ethClients := eth.NewClient(ethClientWrap)

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

		EthClient: ethClients, // for eth chain
		TcClient:  tcClients,  // for tc chain
		BsClient:  bsClient,   // for btc/blockcypher service

	}

	repo, err := repository.NewRepository(&g)
	if err != nil {
		_logger.AtLog.Logger.Error("Cannot init repository", zap.Error(err))
		return
	}

	err = repo.CreateCollectionIndexes()
	if err != nil {
		_logger.AtLog.Logger.Error("CreateCollectionIndexes - Cannot created index ", zap.Error(err))
		// return
	}

	// migration
	migrate.SetDatabase(repo.DB)
	if migrateErr := migrate.Up(-1); migrateErr != nil {
		_logger.AtLog.Errorf("migrate failed", zap.Error(err))
	}

	uc, err := usecase.NewUsecase(&g, *repo)
	if err != nil {
		_logger.AtLog.Errorf("LoadUsecases - Cannot init usecase", zap.Error(err))
		return
	}

	uc.GetTokenArtworkName()

}
