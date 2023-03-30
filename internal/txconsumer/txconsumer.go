package txconsumer

import (
	"fmt"
	"math"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils"
	"rederinghub.io/utils/blockchain"
	"rederinghub.io/utils/config"
	"rederinghub.io/utils/global"
	"rederinghub.io/utils/logger"
	"rederinghub.io/utils/redis"
)

type HttpTxConsumer struct {
	Blockchain                blockchain.TcNetwork
	DefaultLastProcessedBlock int64
	CronJobPeriod             int32
	BatchLogSize              int32
	Addresses                 []common.Address
	Cache                     redis.IRedisCache
	Logger                    logger.Ilogger
	RedisKey                  string
	Usecase                   usecase.Usecase
	Config                    *config.Config
}

func NewHttpTxConsumer(global *global.Global, uc usecase.Usecase, cfg config.Config) (*HttpTxConsumer, error) {
	txConsumer := new(HttpTxConsumer)
	txConsumer.DefaultLastProcessedBlock = cfg.TxConsumerConfig.StartBlock
	txConsumer.CronJobPeriod = cfg.TxConsumerConfig.CronJobPeriod
	txConsumer.BatchLogSize = cfg.TxConsumerConfig.BatchLogSize
	txConsumer.Addresses = make([]common.Address, 0)
	for _, address := range cfg.TxConsumerConfig.Addresses {
		txConsumer.Addresses = append(txConsumer.Addresses, common.HexToAddress(address))
	}
	txConsumer.Cache = global.Cache
	txConsumer.Logger = global.Logger
	txConsumer.Blockchain = global.TcNetwotkchain
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
		logger.AtLog.Logger.Error("c.Cache.Exists", zap.String("redisKey", redisKey), zap.Error(err))
		return 0, err
	}
	if *exists {
		processed, err := c.Cache.GetData(redisKey)
		if err != nil {
			logger.AtLog.Logger.Error("error get from redis", zap.Error(err))
			return 0, err
		}
		if processed == nil {
			return (c.DefaultLastProcessedBlock), nil
		}
		lastProcessedSavedOnRedis, err := strconv.ParseInt(*processed, 10, 64)
		if err != nil {
			logger.AtLog.Logger.Error("err.getLastProcessedBlock", zap.Error(err))
			return 0, err
		}
		lastProcessed = int64(math.Max(float64(lastProcessed), float64(lastProcessedSavedOnRedis)))
	}
	return lastProcessed, nil
}

func (c *HttpTxConsumer) resolveTransaction() error {
	lastProcessedBlock, err := c.getLastProcessedBlock()
	if err != nil {
		logger.AtLog.Logger.Error("resolveTransaction", zap.Any("err", err))
		return err
	}

	fromBlock := lastProcessedBlock + 1
	blockNumber, err := c.Blockchain.GetBlockNumber()
	if err != nil {
		logger.AtLog.Logger.Error("resolveTransaction", zap.Any("err", err))
		return err
	}

	toBlock := int64(math.Min(float64(blockNumber.Int64()), float64(fromBlock+int64(c.BatchLogSize))))
	if toBlock < fromBlock {
		fromBlock = toBlock
	}
	
	logger.AtLog.Logger.Info("resolveTransaction", zap.Int64("currentBlockNumber", blockNumber.Int64()), zap.Int64("fromBlock", fromBlock), zap.Int64("toBlock", toBlock), zap.Int64("lastProcessedBlock",lastProcessedBlock))
	logs, err := c.Blockchain.GetEventLogs(*big.NewInt(fromBlock), *big.NewInt(toBlock), c.Addresses)
	if err != nil {
		logger.AtLog.Logger.Error("err.GetEventLogs", zap.String("err", err.Error()))
		return err
	}

	logger.AtLog.Logger.Info("resolveTransaction", zap.Int64("currentBlockNumber", blockNumber.Int64()), zap.Int64("fromBlock", fromBlock), zap.Int64("toBlock", toBlock), zap.Int64("lastProcessedBlock",lastProcessedBlock), zap.Int("log.number", len(logs)))
	for _, _log := range logs {
		// marketplace logs
		logger.AtLog.Logger.Info("resolveTransaction", zap.Any("_log.Address", _log.Address.String()))
		//MAKET PLACE
		// if strings.ToLower(_log.Address.String()) == c.Config.MarketplaceEvents.Contract {
		// 	topic := strings.ToLower(_log.Topics[0].String())
		// 	logger.AtLog.Logger.Info("topic", zap.Any("topic", topic), zap.Any("_log.TxHash", _log.TxHash), zap.Any("_log.BlockNumber", _log.BlockNumber))

		// 	switch topic {
		// 	case c.Config.MarketplaceEvents.ListToken:
		// 		err = c.Usecase.ResolveMarketplaceListTokenEvent(_log)
		// 	case c.Config.MarketplaceEvents.PurchaseToken:
		// 		err = c.Usecase.ResolveMarketplacePurchaseTokenEvent(_log)
		// 	case c.Config.MarketplaceEvents.MakeOffer:
		// 		err = c.Usecase.ResolveMarketplaceMakeOffer(_log)
		// 	case c.Config.MarketplaceEvents.AcceptMakeOffer:
		// 		err = c.Usecase.ResolveMarketplaceAcceptOfferEvent(_log)
		// 	case c.Config.MarketplaceEvents.CancelListing:
		// 		err = c.Usecase.ResolveMarketplaceCancelListing(_log)
		// 	case c.Config.MarketplaceEvents.CancelMakeOffer:
		// 		err = c.Usecase.ResolveMarketplaceCancelOffer(_log)
		// 	}
		// }

		//DAO
		// if strings.ToLower(_log.Address.String()) == c.Config.DAOEvents.Contract {
		// 	topic := strings.ToLower(_log.Topics[0].String())
		// 	logger.AtLog.Logger.Info("topic", zap.Any("topic", topic), zap.Any("_log.TxHash", _log.TxHash),  zap.Any("_log.BlockNumber", _log.BlockNumber))

		// 	switch topic {
		// 	case c.Config.DAOEvents.ProposalCreated:
		// 		err = c.Usecase.DAOProposalCreated(_log)

		// 	case c.Config.DAOEvents.CastVote:
		// 		err = c.Usecase.DAOCastVote(_log)
		// 	}
		// }

		// do switch case with log.Address and log.Topics
		if _log.Topics[0].String() == os.Getenv("TRANSFER_NFT_SIGNATURE") {
			c.Usecase.UpdateProjectWithListener(_log)
		}

		if err != nil {
			logger.AtLog.Logger.Error("err", zap.Error(err))
			//return err
			continue
		}
	}

	lock, err := c.Cache.GetData(utils.REDIS_KEY_LOCK_TX_CONSUMER_CONSUMER_BLOCK)
	if err == nil && *lock == "true" {
		logger.AtLog.Logger.Info("resolveTransaction", zap.Any("true", true))
	} else {
		// if no error occured, save toBlock as lastProcessedBlock
		err = c.Cache.SetStringData(c.getRedisKey(), strconv.FormatInt(toBlock, 10))
		if err != nil {
			logger.AtLog.Logger.Error("resolveTransaction", zap.Error(err))
			return err
		}
	}

	return nil
}

func (c *HttpTxConsumer) StartServer() {
	logger.AtLog.Logger.Info("HttpTxConsumer start listening")
	for {
		previousTime := time.Now()
		err := c.resolveTransaction()
		if err != nil {
			logger.AtLog.Logger.Error("Error when resolve transactions", zap.String("err", err.Error()))
		}
		processedTime := time.Now().Unix() - previousTime.Unix()
		if processedTime < int64(c.CronJobPeriod) {
			time.Sleep(time.Duration(c.CronJobPeriod-int32(processedTime)) * time.Second)
		}
	}
}
