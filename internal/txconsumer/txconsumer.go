package txconsumer

import (
	"fmt"
	"math"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/opentracing/opentracing-go"
	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils"
	"rederinghub.io/utils/blockchain"
	"rederinghub.io/utils/config"
	"rederinghub.io/utils/global"
	"rederinghub.io/utils/logger"
	"rederinghub.io/utils/redis"
	"rederinghub.io/utils/tracer"
)

type HttpTxConsumer struct {
	Blockchain blockchain.Blockchain
	DefaultLastProcessedBlock int64
	CronJobPeriod int32
	BatchLogSize int32
	Addresses []common.Address
	Cache redis.IRedisCache
	Logger logger.Ilogger
	RedisKey string
	Usecase    usecase.Usecase
	Config        *config.Config
	Tracer  tracer.ITracer
}

func NewHttpTxConsumer(global *global.Global,uc usecase.Usecase, cfg config.Config) (*HttpTxConsumer, error) {
	txConsumer := new(HttpTxConsumer)
	txConsumer.DefaultLastProcessedBlock = cfg.TxConsumerConfig.StartBlock
	txConsumer.CronJobPeriod = cfg.TxConsumerConfig.CronJobPeriod
	txConsumer.BatchLogSize = cfg.TxConsumerConfig.BatchLogSize
	txConsumer.Addresses = make([]common.Address, 0)
	for _, address := range cfg.TxConsumerConfig.Addresses {
		fmt.Println(address)
		txConsumer.Addresses = append(txConsumer.Addresses, common.HexToAddress(address))
	}
	txConsumer.Cache = global.Cache
	txConsumer.Logger = global.Logger
	txConsumer.Blockchain = global.Blockchain
	txConsumer.RedisKey = "tx-consumer"
	txConsumer.Usecase = uc
	txConsumer.Config = &cfg
	txConsumer.Tracer = global.Tracer
	return txConsumer, nil
}

func (h *HttpTxConsumer) StartSpan(name string,  rootSpan opentracing.Span) (opentracing.Span, *tracer.TraceLog) {
	span := h.Tracer.StartSpanFromRoot(rootSpan, name)
	log := tracer.NewTraceLog()
	return span, log
}

func (h *HttpTxConsumer) StartSpanWithoutRoot(name string) (opentracing.Span, *tracer.TraceLog) {
	span := h.Tracer.StartSpan(name)
	log := tracer.NewTraceLog()
	return span, log
}

func (c *HttpTxConsumer) getRedisKey() string {
	return fmt.Sprintf("%s:lastest_processed", c.RedisKey)
}

func (c *HttpTxConsumer) getLastProcessedBlock(rootSpan opentracing.Span) (int64, error) {
	span, log := c.StartSpan("getLastProcessedBlock", rootSpan)
	defer c.Tracer.FinishSpan(span, log)
	lastProcessed := c.DefaultLastProcessedBlock
	
	redisKey := c.getRedisKey()
	exists, err := c.Cache.Exists(redisKey)
	if err != nil {
		return 0, err
	}
	
	if *exists {
		processed, err := c.Cache.GetData(redisKey)
		if err != nil {
			fmt.Println(err)
			c.Logger.Error("error get from redis", err)
			return 0, err
		}
		if processed == nil {
			return (c.DefaultLastProcessedBlock), nil
		}
		lastProcessedSavedOnRedis, err := strconv.ParseInt(*processed, 10, 64)
		log.SetData("processed", processed)
		log.SetData("lastProcessedSavedOnRedis", lastProcessedSavedOnRedis)
		
		if err != nil {
			return 0, err
		}
		lastProcessed = int64(math.Max(float64(lastProcessed), float64(lastProcessedSavedOnRedis)))
	}
	
	return lastProcessed, nil
}

func (c *HttpTxConsumer) resolveTransaction() error {
	span, log := c.StartSpanWithoutRoot("resolveTransaction")
	defer c.Tracer.FinishSpan(span, log)
	lastProcessedBlock, err := c.getLastProcessedBlock(span)
	if err != nil {
		log.Error("Error when get last processed block", err.Error(), err)
		return err
	}
	log.SetTag("lastProcessedBlock", lastProcessedBlock)
	fromBlock := lastProcessedBlock + 1
	log.SetTag("fromBlock", fromBlock)
	blockNumber, err := c.Blockchain.GetBlockNumber()
	if err != nil {
		return err
	}
	log.SetTag("blockNumber", blockNumber)
	toBlock := int64(math.Min(float64(blockNumber.Int64()), float64(fromBlock + int64(c.BatchLogSize))))
	log.SetTag("toBlock", toBlock)
	c.Logger.Info(fmt.Sprintf("Searching log from %v to %v", fromBlock, toBlock))
	log.SetData("from block", fromBlock)
	log.SetData("to block", toBlock)
	log.SetData("block number", blockNumber.Int64())

	logs, err := c.Blockchain.GetEventLogs(*big.NewInt(fromBlock), *big.NewInt(toBlock), c.Addresses)
	if err != nil {
		return err
	}

	log.SetData("logs", logs)
	for _, _log := range logs {
		// marketplace logs
		log.SetData("_log.Address.String()", _log.Address.String())
		log.SetData("c.Config.MarketplaceEvents.Contract", c.Config.MarketplaceEvents.Contract)
		
		//MAKET PLACE
		if strings.ToLower(_log.Address.String()) == c.Config.MarketplaceEvents.Contract {
			topic :=  strings.ToLower(_log.Topics[0].String())
			log.SetData("topic", topic)
			log.SetData("topic tx hash", _log.TxHash)
			log.SetData("topic block number", _log.BlockNumber)
			log.SetTag("blockNumber", _log.BlockNumber)

			switch topic {
			case c.Config.MarketplaceEvents.ListToken:
				log.SetTag("event", "ListToken")
				err = c.Usecase.ResolveMarketplaceListTokenEvent(span, _log)
			case c.Config.MarketplaceEvents.PurchaseToken:
				log.SetData("event", "PurchaseToken")
				err = c.Usecase.ResolveMarketplacePurchaseTokenEvent(span, _log)
			case c.Config.MarketplaceEvents.MakeOffer:
				log.SetData("event", "MakeOffer")
				err = c.Usecase.ResolveMarketplaceMakeOffer(span, _log)
			case c.Config.MarketplaceEvents.AcceptMakeOffer:
				log.SetData("event", "AcceptMakeOffer")
				err = c.Usecase.ResolveMarketplaceAcceptOfferEvent(span, _log)
			case c.Config.MarketplaceEvents.CancelListing:
				log.SetData("event", "CancelListing")
				err = c.Usecase.ResolveMarketplaceCancelListing(span, _log)
			case c.Config.MarketplaceEvents.CancelMakeOffer:
				log.SetData("event", "CancelMakeOffer")
				err = c.Usecase.ResolveMarketplaceCancelOffer(span, _log)
			}
		}
		
		//DAO
		if strings.ToLower(_log.Address.String()) == c.Config.DAOEvents.Contract {
			topic :=  strings.ToLower(_log.Topics[0].String())
			log.SetData("topic", topic)
			log.SetData("topic tx hash", _log.TxHash)
			log.SetData("topic block number", _log.BlockNumber)
			log.SetTag("blockNumber", _log.BlockNumber)
			
			switch topic {
			case c.Config.DAOEvents.ProposalCreated:
				log.SetTag("event", "ProposalCreated")
				err = c.Usecase.DAOProposalCreated(span, _log)
			}
		}

		// do switch case with log.Address and log.Topics
		if _log.Topics[0].String() == os.Getenv("TRANSFER_NFT_SIGNATURE") {
			c.Usecase.UpdateProjectWithListener(_log)
		}
		if err != nil {
			log.Error("error resolve event", err.Error(), err)
			return err
		}
	}

	lock, err := c.Cache.GetData(utils.REDIS_KEY_LOCK_TX_CONSUMER_CONSUMER_BLOCK)
	if err == nil && *lock == "true" {
		log.SetData("lock-tx-consumer-update-last-processed-block", true)
	} else {
		// if no error occured, save toBlock as lastProcessedBlock
		err = c.Cache.SetStringData(c.getRedisKey(), strconv.FormatInt(toBlock, 10))
		log.SetData("set-last-processed", strconv.FormatInt(toBlock, 10))
		if err != nil {
			log.Error("error set redis", err.Error(), err)
			return err
		}
	}

	return nil
}

func (c *HttpTxConsumer) StartServer() {
	c.Logger.Info("Start listening")
	for {
		previousTime := time.Now()
		err := c.resolveTransaction()
		if err != nil {
			c.Logger.Error("Error when resolve transactions", err)
		}
		processedTime := time.Now().Unix() - previousTime.Unix()
		if processedTime < int64(c.CronJobPeriod) {
			time.Sleep(time.Duration(c.CronJobPeriod - int32(processedTime)) * time.Second)
		}
	}
}


