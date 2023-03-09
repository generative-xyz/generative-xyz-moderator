package main

import (
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"
	"rederinghub.io/utils/delegate"

	"github.com/gorilla/mux"
	migrate "github.com/xakep666/mongo-migrate"
	"rederinghub.io/external/dev5service"
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
	l.LogAny("config", zap.Any("config.NewConfig", c))

	mongoCnn := fmt.Sprintf("%s://%s:%s@%s/?retryWrites=true&w=majority", c.Databases.Mongo.Scheme, c.Databases.Mongo.User, c.Databases.Mongo.Pass, c.Databases.Mongo.Host)
	mongoDbConnection, err := connections.NewMongo(mongoCnn)
	if err != nil {
		log.Println("Cannot connect mongoDB ", err)
		panic(err)
	}

	ethClient, err = blockchain.NewBlockchain(c.BlockchainConfig)
	if err != nil {
		log.Println("Cannot connect to eth client ", err)
		panic(err)
	}

	logger = l
	conf = c
	mongoConnection = mongoDbConnection
}

// @title Generative.xyz APIs
// @version 1.0.0
// @description This is a sample server Autonomous devices management server.

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
	log.Println("starting server ...")
	cache, redisClient := redis.NewRedisCache(conf.Redis)
	r := mux.NewRouter()

	gcs, err := googlecloud.NewDataGCStorage(*conf)
	if err != nil {
		logger.Error("Cannot init gcs", err)
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
		logger.Error("error initializing delegate service", err)
		return
	}

	devSer := dev5service.NewDev5Service(conf, cache)
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
		Slack:               *slack,
		DiscordClient:       discordclient.NewCLient(),
		Pubsub:              rPubsub,
		OrdService:          ord,
		OrdServiceDeveloper: ordForDeveloper,
		DelegateService:     delegateService,
		Dev5Services: devSer,
	}

	repo, err := repository.NewRepository(&g)
	if err != nil {
		logger.Error("Cannot init repository", err)
		return
	}

	err = repo.CreateCollectionIndexes()
	if err != nil {
		logger.Error("CreateCollectionIndexes - Cannot created index ", err)
		// return
	}

	// migration
	migrate.SetDatabase(repo.DB)
	if migrateErr := migrate.Up(-1); migrateErr != nil {
		_logger.AtLog.Errorf("migrate failed", zap.Error(err))
	}

	uc, err := usecase.NewUsecase(&g, *repo)
	if err != nil {
		logger.Error("LoadUsecases - Cannot init usecase", err)
		return
	}

	//uc.FindInscriptions()
	uc.CreateInscriptionFiles()
}
