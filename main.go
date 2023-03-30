package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"rederinghub.io/utils/delegate"
	"rederinghub.io/utils/eth"
	"rederinghub.io/utils/redisv9"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
	migrate "github.com/xakep666/mongo-migrate"
	"rederinghub.io/external/nfts"
	"rederinghub.io/external/ord_service"
	"rederinghub.io/internal/delivery"
	"rederinghub.io/internal/delivery/crontabManager"
	"rederinghub.io/internal/delivery/crontab_ordinal_collections"
	httpHandler "rederinghub.io/internal/delivery/http"
	"rederinghub.io/internal/delivery/pubsub"
	"rederinghub.io/internal/delivery/txserver"
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

	tcClientWrap, err := ethclient.Dial(conf.BlockchainConfig.TCEndpoint)
	if err != nil {
		_logger.AtLog.Logger.Error("error initializing tcClient service", zap.Error(err))
		return
	}
	tcClientw := eth.NewClient(tcClientWrap)

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

		TcClient: tcClientw,
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

	servers := make(map[string]delivery.AddedServer)
	// api fixed run:
	h, _ := httpHandler.NewHandler(&g, *uc)
	servers["http"] = delivery.AddedServer{
		Server:  h,
		Enabled: conf.StartHTTP,
	}

	ph := pubsub.NewPubsubHandler(*uc, rPubsub, logger)
	servers["pubsub"] = delivery.AddedServer{
		Server:  ph,
		Enabled: conf.StartPubsub,
	}

	// job ORDINAL_COLLECTION_CRONTAB_START: @Dac TODO move all function to Usercase.
	ordinalCron := crontab_ordinal_collections.NewScronOrdinalCollectionHandler(&g, *uc)
	servers["ordinal_collections_crontab"] = delivery.AddedServer{
		Server:  ordinalCron,
		Enabled: conf.Crontab.OrdinalCollectionEnabled,
	}

	txConsumer, _ := txserver.NewTxServer(&g, *uc, *conf)
	servers["txconsumer"] = delivery.AddedServer{
		Server:  txConsumer,
		Enabled: conf.TxConsumerConfig.Enabled,
	}

	// job init:
	/*
			txConsumer, _ := txserver.NewTxServer(&g, *uc, *conf)
			cron := crontab.NewScronHandler(&g, *uc)
			btcCron := crontab_btc.NewScronBTCHandler(&g, *uc)
			mkCron := crontab_marketplace.NewScronMarketPlace(&g, *uc)
			inscribeCron := incribe_btc.NewScronBTCHandler(&g, *uc)
			mintNftBtcCron := mint_nft_btc.NewCronMintNftBtcHandler(&g, *uc)
			trendingCron := crontab_trending.NewScronTrendingHandler(&g, *uc)
			ordinalCron := crontab_ordinal_collections.NewScronOrdinalCollectionHandler(&g, *uc)
			inscriptionIndexCron := crontab_inscription_info.NewScronInscriptionInfoHandler(&g, *uc)
			dexBTCCron := dex_btc_cron.NewScronDexBTCHandler(&g, *uc)

		ph := pubsub.NewPubsubHandler(*uc, rPubsub, logger)

			servers["developer_inscribe"] = delivery.AddedServer{
				Server:  developerInscribeCron,
				Enabled: conf.Crontab.CrontabDeveloperInscribeEnabled,
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

			servers["ordinal_collections_crontab"] = delivery.AddedServer{
				Server:  ordinalCron,
				Enabled: conf.Crontab.OrdinalCollectionEnabled,
			}
			servers["inscription_index_crontab"] = delivery.AddedServer{
				Server:  inscriptionIndexCron,
				Enabled: conf.Crontab.InscriptionIndexEnabled,
			}

			servers["dex_btc_cron"] = delivery.AddedServer{
				Server:  dexBTCCron,
				Enabled: conf.Crontab.DexBTCEnabled,
			}
			servers["pubsub"] = delivery.AddedServer{
			Server:  ph,
			Enabled: conf.StartPubsub,
		}
	*/

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
			_logger.AtLog.Logger.Info(fmt.Sprintf("%s is enabled", name))
		} else {
			_logger.AtLog.Logger.Info(fmt.Sprintf("%s is disabled", name))
		}
	}

	// start a group cron:
	if len(conf.CronTabList) > 0 {
		for _, cronKey := range conf.CronTabList {
			_logger.AtLog.Logger.Info(fmt.Sprintf("%s is running...", cronKey))
			crontabManager.NewCrontabManager(cronKey, &g, *uc).StartServer()
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
	// 	logger.AtLog.Logger.Error("httpDelivery.StartServer - Server can not shutdown", err)
	// 	return
	// }
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	<-ctx.Done() //if your application should wait for other services
	// to finalize based on context cancellation.
	_logger.AtLog.Logger.Warn("httpDelivery.StartServer - server is shutting down")
	tracer.Stop()
	os.Exit(0)

}
