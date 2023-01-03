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
		c.Logger.Error("Error when get last processed block", err)
		return err
	}
	fromBlock := lastProcessedBlock + 1
	blockNumber, err := c.Blockchain.GetBlockNumber()
	if err != nil {
		return err
	}

	toBlock := int64(math.Min(float64(blockNumber.Int64()), float64(fromBlock + int64(c.BatchLogSize))))
	c.Logger.Info(fmt.Sprintf("Searching log from %v to %v", fromBlock, toBlock))
	logs, err := c.Blockchain.GetEventLogs(*big.NewInt(fromBlock), *big.NewInt(toBlock), c.Addresses)
	if err != nil {
		return err
	}

	for _, log := range logs {
		// marketplace logs
		if strings.ToLower(log.Address.String()) == c.Config.MarketplaceEvents.Contract {
			topic :=  strings.ToLower(log.Topics[0].String())
			switch topic {
			case c.Config.MarketplaceEvents.ListToken:
				c.Usecase.ResolveMarketplaceListTokenEvent(log)
			case c.Config.MarketplaceEvents.PurchaseToken:
				c.Usecase.ResolveMarketplacePurchaseTokenEvent(log)
			case c.Config.MarketplaceEvents.MakeOffer:
				c.Usecase.ResolveMarketplaceMakeOffer(log)
			case c.Config.MarketplaceEvents.AcceptMakeOffer:
				c.Usecase.ResolveMarketplaceAcceptOfferEvent(log)
			case c.Config.MarketplaceEvents.CancelListing:
				c.Usecase.ResolveMarketplaceCancelListing(log)
			case c.Config.MarketplaceEvents.CancelMakeOffer:
				c.Usecase.ResolveMarketplaceCancelOffer(log)
			}
		}
		// do switch case with log.Address and log.Topics
		if log.Topics[0].String() == os.Getenv("TRANSFER_NFT_SIGNATURE") {
			c.Usecase.UpdateProjectWithListener(log)
		}
	}

	// if no error occured, save toBlock as lastProcessedBlock
	err = c.Cache.SetStringData(c.getRedisKey(), strconv.FormatInt(toBlock, 10))
	if err != nil {
		return err
	}

	return nil
}

func (c *HttpTxConsumer) StartListen() {
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


