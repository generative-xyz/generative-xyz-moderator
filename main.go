package main

import (
	"fmt"
	"log"

	"rederinghub.io/internal/delivery/crontab"

	"rederinghub.io/external/nfts"
	httpHandler "rederinghub.io/internal/delivery/http"
	"rederinghub.io/internal/repository"
	"rederinghub.io/internal/txconsumer"
	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils/blockchain"
	"rederinghub.io/utils/config"
	"rederinghub.io/utils/connections"
	"rederinghub.io/utils/global"
	"rederinghub.io/utils/googlecloud"
	_logger "rederinghub.io/utils/logger"
	"rederinghub.io/utils/oauth2service"
	"rederinghub.io/utils/redis"
	"rederinghub.io/utils/tracer"

	"github.com/gorilla/mux"
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

	l := _logger.NewLogger()


	mongoCnn := fmt.Sprintf("mongodb://%s:%s@%s:%s/", c.Databases.Mongo.User,c.Databases.Mongo.Pass, c.Databases.Mongo.Host, c.Databases.Mongo.Port )
	mongoDbConnection, err := connections.NewMongo(mongoCnn)
	if err != nil {
		log.Println("Can not connect mongoDB ", err)
		panic(err)
	}

	ethClient, err = blockchain.NewBlockchain(c.BlockchainConfig)
	if err != nil {
		log.Println("Can not connect to eth client ", err)
		panic(err)
	}

	logger = l
	conf = c
	mongoConnection = mongoDbConnection
}

// @title Generative.xyz APIs
// @version 1.0.0
// @description This is a sample server Autonomous devices management server.

// @securityDefinitions.apikey Authorization
// @in header
// @name Authorization

// @securityDefinitions.apikey Api-Key
// @in header
// @name Api-Key

// @BasePath /rederinghub.io/v1
func main() {
	// log.Println("init sentry ...")
	// sentry.InitSentry(conf)
	startServer()
}

func startServer() {
	log.Println("starting server ...")
	cache := redis.NewRedisCache(conf.Redis)
	t := tracer.NewTracing(logger)
	r := mux.NewRouter()

	gcs, err := googlecloud.NewDataGCStorage(*conf)
	if err != nil {
		logger.Error("Can not init gcs", err)
		return
	}

	moralis := nfts.NewMoralisNfts(conf, t, cache)
	covalent := nfts.NewCovalentNfts(conf);

	// hybrid auth
	auth2Service := oauth2service.NewAuth2()
	g := global.Global{
		Tracer: t,
		Logger:       logger,
		MuxRouter:    r,
		Conf:         conf,
		DBConnection: mongoConnection,
		Cache:        cache,
		Auth2: *auth2Service,
		GCS: gcs,
		MoralisNFT: *moralis,
		CovalentNFT: *covalent,
		Blockchain: *ethClient,
	}

	repo, err := repository.NewRepository(&g)
	if err != nil {
		logger.Error("Can not init repository", err)
		return
	}
	
	uc, err := usecase.NewUsecase(&g, *repo)
	if err != nil {
		logger.Error("LoadUsecases - Can not init usecase", err)
		return
	}

	h, err := httpHandler.NewHandler(&g, *uc)
	if err != nil {
		logger.Error("Init handler failure", err)
		return
	}

	if (conf.TxConsumerConfig.Enabled) {
		txConsumer, err := txconsumer.NewHttpTxConsumer(&g, *uc, conf.TxConsumerConfig)
		if err != nil {
			logger.Error("Failed to init tx consumer")
			return
		}
		go func (txConsumer *txconsumer.HttpTxConsumer)  {
			txConsumer.StartListen()
		}(txConsumer)
		
	}

	cron := crontab.NewScronHandler(&g, *uc)
	go func (cron *crontab.ScronHandler)  {
		logger.Info("Cron is listening")
		cron.StartServer()
	}(cron)

	log.Println("started server and listening")
	h.StartServer()
}
