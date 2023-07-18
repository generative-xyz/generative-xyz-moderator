package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"rederinghub.io/external/coin_market_cap"
	"rederinghub.io/external/etherscan"
	"rederinghub.io/external/mempool_space"
	"rederinghub.io/external/token_explorer"
	gm_crontab_sever "rederinghub.io/internal/delivery/gm_crontab_server"
	"rederinghub.io/internal/delivery/project_protab_crontab_server"
	"strconv"
	"time"

	"go.uber.org/zap"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
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
// @description This is a sample server Generative.

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @BasePath /rederinghub.io/v1
func main() {

	//defer func() {
	//	if r := recover(); r != nil {
	//		fmt.Println("Recovered. Error:\n", r)
	//	}
	//}()

	// log.Println("init sentry ...")
	// sentry.InitSentry(conf)
	startServer()
}

func startServer() {
	fmt.Println("starting server ...")
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

	// init tc client public
	tcClientPublicWrap, err := ethclient.Dial(conf.BlockchainConfig.TCPublicEndpoint)
	if err != nil {
		_logger.AtLog.Logger.Error("error initializing tcClient public service", zap.Error(err))
		return
	}
	tcClientPublic := eth.NewClient(tcClientPublicWrap)

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

	// init eth client for dex:
	ethDexClientWrap, err := ethclient.Dial(conf.BlockchainConfig.ETHEndpointDex)
	if err != nil {
		_logger.AtLog.Logger.Error("error initializing ethDexClientWrap service", zap.Error(err))
		return
	}
	ethClientDex := eth.NewClient(ethDexClientWrap)
	eScan := etherscan.NewEtherscanService(conf, cache)
	mpService := mempool_space.NewMempoolService(conf, cache)
	coinMKC := coin_market_cap.NewCoinMarketCap(conf, cache)
	te := token_explorer.NewTokenExplorer(cache)

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
		TokenExplorer:      te,
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

	// job ORDINAL_COLLECTION_CRONTAB_START: @Dac TODO move all function to Usercase and crontab db config.
	ordinalCron := crontab_ordinal_collections.NewScronOrdinalCollectionHandler(&g, *uc)
	servers["ordinal_collections_crontab"] = delivery.AddedServer{
		Server:  ordinalCron,
		Enabled: conf.Crontab.OrdinalCollectionEnabled,
	}

	// TODO move all function to Usercase and crontab db config.
	txConsumer, _ := txserver.NewTxServer(&g, *uc, *conf)
	servers["txconsumer"] = delivery.AddedServer{
		Server:  txConsumer,
		Enabled: conf.TxConsumerConfig.Enabled,
	}

	pProtabBool := false
	pProtab := os.Getenv("PROJECT_PROTAB_ENABLED")
	if pProtab != "" {
		if pProtab == "true" {
			pProtabBool = true
		}
	}
	protab, _ := project_protab_crontab_server.NewProjectProtabCrontabServer(uc)
	servers["project_protab"] = delivery.AddedServer{
		Server:  protab,
		Enabled: pProtabBool,
	}

	isGmEnabled := false
	gmEnabled := os.Getenv("START_GM_CRONTAB")
	if gmEnabled != "" {
		isGmEnabled, err = strconv.ParseBool(gmEnabled)
		if err != nil {
			isGmEnabled = false
		}
	}

	gmCrontab, _ := gm_crontab_sever.NewGmCrontabServer(uc)
	servers["gm_crontab"] = delivery.AddedServer{
		Server:  gmCrontab,
		Enabled: isGmEnabled,
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
			fmt.Printf("%s is enabled \n", name)
		} else {
			fmt.Printf("%s is disabled \n", name)
		}
	}

	// start a group cron:
	if len(conf.CronTabList) > 0 {
		for _, cronKey := range conf.CronTabList {
			fmt.Printf("%s is running... \n", cronKey)
			crontabManager.NewCrontabManager(cronKey, &g, *uc).StartServer()
		}
	}
	// testing purpose
	//uc.TestSendNoti()

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
