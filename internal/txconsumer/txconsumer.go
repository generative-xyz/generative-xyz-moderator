package txconsumer

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"rederinghub.io/utils/blockchain"
	"rederinghub.io/utils/config"
	"rederinghub.io/utils/global"
	"rederinghub.io/utils/logger"
	"rederinghub.io/utils/redis"
)

const (
	TRANSFER_NFT_SIGNATURE = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
)

type HttpTxConsumer struct {
	Blockchain blockchain.Blockchain
	DefaultLastProcessedBlock int64
	CronJobPeriod int32
	Addresses []common.Address
	Cache redis.IRedisCache
	Logger logger.Ilogger
	RedisKey string
}

func NewHttpTxConsumer(global *global.Global, cfg config.TxConsumerConfig) (*HttpTxConsumer, error) {
	txConsumer := new(HttpTxConsumer)
	txConsumer.DefaultLastProcessedBlock = cfg.StartBlock
	txConsumer.CronJobPeriod = cfg.CronJobPeriod
	txConsumer.Addresses = make([]common.Address, 0)
	fmt.Println(cfg.Addresses)
	for _, address := range cfg.Addresses {
		fmt.Println(address)
		txConsumer.Addresses = append(txConsumer.Addresses, common.HexToAddress(address))
	}
	fmt.Println(txConsumer.Addresses)
	txConsumer.Cache = global.Cache
	txConsumer.Logger = global.Logger
	txConsumer.Blockchain = global.Blockchain
	txConsumer.RedisKey = "tx-consumer"
	return txConsumer, nil
}

func (c *HttpTxConsumer) getRedisKey() string {
	return fmt.Sprintf("%s:lastest_processed", c.RedisKey)
}

func (c *HttpTxConsumer) getLastProcessedBlock() (int64, error) {
	redisKey := c.getRedisKey()

	exists, err := c.Cache.Exists(redisKey)
	if err != nil {
		return 0, err
	}
	lastProcessed := c.DefaultLastProcessedBlock
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
	toBlock, err := c.Blockchain.GetBlockNumber()
	if err != nil {
		return err
	}

	c.Logger.Info(fmt.Sprintf("Searching log from %v to %v", fromBlock, toBlock))

	logs, err := c.Blockchain.GetEventLogs(*big.NewInt(fromBlock), *toBlock, c.Addresses)
	if err != nil {
		return err
	}

	for _, log := range logs {
		fmt.Println(log.Address)
		fmt.Println(log.Topics[0])
		// do switch case with log.Address and log.Topics
		if log.Topics[0].String() == TRANSFER_NFT_SIGNATURE {
			fmt.Println(log.Topics)
		}
	}

	// if no error occured, save toBlock as lastProcessedBlock
	err = c.Cache.SetStringData(c.getRedisKey(), toBlock.String())
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


