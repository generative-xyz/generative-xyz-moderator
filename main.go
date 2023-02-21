package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"rederinghub.io/utils/delegate"

	"github.com/gorilla/mux"
	"rederinghub.io/external/nfts"
	"rederinghub.io/external/ord_service"
	"rederinghub.io/internal/delivery"
	"rederinghub.io/internal/delivery/crontab"
	"rederinghub.io/internal/delivery/crontab_btc"
	"rederinghub.io/internal/delivery/crontab_marketplace"
	"rederinghub.io/internal/delivery/crontab_trending"
	httpHandler "rederinghub.io/internal/delivery/http"
	"rederinghub.io/internal/delivery/incribe_btc"
	"rederinghub.io/internal/delivery/mint_nft_btc"
	"rederinghub.io/internal/delivery/pubsub"
	"rederinghub.io/internal/delivery/txserver"
	"rederinghub.io/internal/repository"
	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils/blockchain"
	"rederinghub.io/utils/config"
	"rederinghub.io/utils/connections"
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

	l := _logger.NewLogger()

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

// @securityDefinitions.apikey Authorization
// @in header
// @name Authorization

// @securityDefinitions.apikey Api-Key
// @in header
// @name Api-Key

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
	ord := ord_service.NewBtcOrd(conf, cache)
	covalent := nfts.NewCovalentNfts(conf)
	slack := slack.NewSlack(conf.Slack)
	rPubsub := redis.NewPubsubClient(conf.Redis)
	delegateService, err := delegate.NewService(ethClient.GetClient())
	if err != nil {
		logger.Error("error initializing delegate service", err)
		return
	}
	// hybrid auth
	auth2Service := oauth2service.NewAuth2()
	g := global.Global{
		Logger:          logger,
		MuxRouter:       r,
		Conf:            conf,
		DBConnection:    mongoConnection,
		Cache:           cache,
		Auth2:           *auth2Service,
		GCS:             gcs,
		S3Adapter:       s3Adapter,
		MoralisNFT:      *moralis,
		CovalentNFT:     *covalent,
		Blockchain:      *ethClient,
		Slack:           *slack,
		Pubsub:          rPubsub,
		OrdService:      ord,
		DelegateService: delegateService,
	}

	repo, err := repository.NewRepository(&g)
	if err != nil {
		logger.Error("Cannot init repository", err)
		return
	}

	err = repo.CreateCollectionIndexes()
	if err != nil {
		logger.Error("CreateCollectionIndexes - Cannot created index ", err)
		return
	}

	uc, err := usecase.NewUsecase(&g, *repo)
	if err != nil {
		logger.Error("LoadUsecases - Cannot init usecase", err)
		return
	}

	h, _ := httpHandler.NewHandler(&g, *uc)
	txConsumer, _ := txserver.NewTxServer(&g, *uc, *conf)
	cron := crontab.NewScronHandler(&g, *uc)
	btcCron := crontab_btc.NewScronBTCHandler(&g, *uc)
	mkCron := crontab_marketplace.NewScronMarketPlace(&g, *uc)
	inscribeCron := incribe_btc.NewScronBTCHandler(&g, *uc)

	mintNftBtcCron := mint_nft_btc.NewCronMintNftBtcHandler(&g, *uc)

	trendingCron := crontab_trending.NewScronTrendingHandler(&g, *uc)
	ph := pubsub.NewPubsubHandler(*uc, rPubsub, logger)

	servers := make(map[string]delivery.AddedServer)
	servers["http"] = delivery.AddedServer{
		Server:  h,
		Enabled: conf.StartHTTP,
	}

	servers["txconsumer"] = delivery.AddedServer{
		Server:  txConsumer,
		Enabled: conf.TxConsumerConfig.Enabled,
	}

	servers["crontab"] = delivery.AddedServer{
		Server:  cron,
		Enabled: conf.Crontab.Enabled,
	}

	servers["btc_crontab"] = delivery.AddedServer{
		Server:  btcCron,
		Enabled: conf.Crontab.BTCEnabled,
	}

	servers["btc_crontab_v2"] = delivery.AddedServer{
		Server:  inscribeCron,
		Enabled: conf.Crontab.BTCV2Enabled,
	}

	servers["mint_nft_btc"] = delivery.AddedServer{
		Server:  mintNftBtcCron,
		Enabled: conf.Crontab.MintNftBtcEnabled,
	}

	servers["marketplace_crontab"] = delivery.AddedServer{
		Server:  mkCron,
		Enabled: conf.Crontab.MarketPlaceEnabled,
	}

	servers["trending_crontab"] = delivery.AddedServer{
		Server:  trendingCron,
		Enabled: conf.Crontab.TrendingEnabled,
	}

	servers["pubsub"] = delivery.AddedServer{
		Server:  ph,
		Enabled: conf.StartPubsub,
	}

	//var wait time.Duration
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)
	// Run our server in a goroutine so that it doesn't block.

	for name, server := range servers {
		if server.Enabled {
			if server.Server != nil {
				go server.Server.StartServer()
			}
			h.Logger.Info(fmt.Sprintf("%s is enabled", name))
		} else {
			h.Logger.Info(fmt.Sprintf("%s is disabled", name))
		}
	}

	// Block until we receive our signal.
	<-c
	wait := time.Second
	// // Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// // Doesn't block if no connections, but will otherwise wait
	// // until the timeout deadline.
	// err := srv.Shutdown(ctx)
	// if err != nil {
	// 	h.Logger.Error("httpDelivery.StartServer - Server can not shutdown", err)
	// 	return
	// }
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	<-ctx.Done() //if your application should wait for other services
	// to finalize based on context cancellation.
	h.Logger.Warning("httpDelivery.StartServer - server is shutting down")
	os.Exit(0)

}
