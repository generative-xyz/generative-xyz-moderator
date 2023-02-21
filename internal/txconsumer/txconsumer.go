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
	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils"
	"rederinghub.io/utils/blockchain"
	"rederinghub.io/utils/config"
	"rederinghub.io/utils/global"
	"rederinghub.io/utils/logger"
	"rederinghub.io/utils/redis"
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
	return txConsumer, nil
}

func (c *HttpTxConsumer) getRedisKey() string {
	return fmt.Sprintf("%s:lastest_processed", c.RedisKey)
}

func (c *HttpTxConsumer) getLastProcessedBlock() (int64, error) {

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
		c.Logger.Info("processed", processed)
		c.Logger.Info("lastProcessedSavedOnRedis", lastProcessedSavedOnRedis)
		if err != nil {
			return 0, err
		}
		lastProcessed = int64(math.Max(float64(lastProcessed), float64(lastProcessedSavedOnRedis)))
	}
return lastProcessed, nil
}

func (c *HttpTxConsumer) resolveTransaction() error {

	lastProcessedBlock, err := c.getLastProcessedBlock()
	if err != nil {
		c.Logger.Error(err)
		return err
	}
	
	fromBlock := lastProcessedBlock + 1
	
	blockNumber, err := c.Blockchain.GetBlockNumber()
	if err != nil {
		return err
	}
	
	toBlock := int64(math.Min(float64(blockNumber.Int64()), float64(fromBlock + int64(c.BatchLogSize))))
	
	c.Logger.Info(fmt.Sprintf("Searching log from %v to %v", fromBlock, toBlock))
	c.Logger.Info("from block", fromBlock)
	c.Logger.Info("to block", toBlock)
	c.Logger.Info("block number", blockNumber.Int64())

	logs, err := c.Blockchain.GetEventLogs(*big.NewInt(fromBlock), *big.NewInt(toBlock), c.Addresses)
	if err != nil {
		return err
	}

	c.Logger.Info("logs", logs)
	for _, _log := range logs {
		// marketplace logs
		c.Logger.Info("_log.Address.String()", _log.Address.String())
		c.Logger.Info("c.Config.MarketplaceEvents.Contract", c.Config.MarketplaceEvents.Contract)
		//MAKET PLACE
		if strings.ToLower(_log.Address.String()) == c.Config.MarketplaceEvents.Contract {
			topic :=  strings.ToLower(_log.Topics[0].String())
			c.Logger.Info("topic", topic)
			c.Logger.Info("topic tx hash", _log.TxHash)
			c.Logger.Info("topic block number", _log.BlockNumber)
	
			switch topic {
			case c.Config.MarketplaceEvents.ListToken:
				err = c.Usecase.ResolveMarketplaceListTokenEvent(_log)
			case c.Config.MarketplaceEvents.PurchaseToken:
				c.Logger.Info("event", "PurchaseToken")
				err = c.Usecase.ResolveMarketplacePurchaseTokenEvent(_log)
			case c.Config.MarketplaceEvents.MakeOffer:
				c.Logger.Info("event", "MakeOffer")
				err = c.Usecase.ResolveMarketplaceMakeOffer(_log)
			case c.Config.MarketplaceEvents.AcceptMakeOffer:
				c.Logger.Info("event", "AcceptMakeOffer")
				err = c.Usecase.ResolveMarketplaceAcceptOfferEvent(_log)
			case c.Config.MarketplaceEvents.CancelListing:
				c.Logger.Info("event", "CancelListing")
				err = c.Usecase.ResolveMarketplaceCancelListing(_log)
			case c.Config.MarketplaceEvents.CancelMakeOffer:
				c.Logger.Info("event", "CancelMakeOffer")
				err = c.Usecase.ResolveMarketplaceCancelOffer(_log)
			}
		}
		//DAO
		if strings.ToLower(_log.Address.String()) == c.Config.DAOEvents.Contract {
			topic :=  strings.ToLower(_log.Topics[0].String())
			c.Logger.Info("topic", topic)
			c.Logger.Info("topic tx hash", _log.TxHash)
			c.Logger.Info("topic block number", _log.BlockNumber)
					switch topic {
			case c.Config.DAOEvents.ProposalCreated:
						err = c.Usecase.DAOProposalCreated(_log)

			case c.Config.DAOEvents.CastVote:
						err = c.Usecase.DAOCastVote(_log)
			}
			}

		// do switch case with log.Address and log.Topics
		if _log.Topics[0].String() == os.Getenv("TRANSFER_NFT_SIGNATURE") {
			c.Usecase.UpdateProjectWithListener(_log)
		}


		if err != nil {
			c.Logger.Error(err)
			//return err
		}
	}

	lock, err := c.Cache.GetData(utils.REDIS_KEY_LOCK_TX_CONSUMER_CONSUMER_BLOCK)
	if err == nil && *lock == "true" {
		c.Logger.Info("lock-tx-consumer-update-last-processed-block", true)
	} else {
		// if no error occured, save toBlock as lastProcessedBlock
		err = c.Cache.SetStringData(c.getRedisKey(), strconv.FormatInt(toBlock, 10))
		c.Logger.Info("set-last-processed", strconv.FormatInt(toBlock, 10))
		if err != nil {
			c.Logger.Error(err)
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


