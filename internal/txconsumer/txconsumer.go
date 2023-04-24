package txconsumer

import (
	"fmt"
	"go.uber.org/zap"
	"math/big"
	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils/blockchain"
	"rederinghub.io/utils/config"
	"rederinghub.io/utils/global"
	"rederinghub.io/utils/logger"
	"rederinghub.io/utils/redis"
	"strconv"
	"strings"
	"time"
)

type HttpTxConsumer struct {
	Blockchain                blockchain.TcNetwork
	DefaultLastProcessedBlock int64
	CronJobPeriod             int32
	BatchLogSize              int32
	//Addresses                 []common.Address
	Cache    redis.IRedisCache
	Logger   logger.Ilogger
	RedisKey string
	Usecase  usecase.Usecase
	Config   *config.Config
}

func NewHttpTxConsumer(global *global.Global, uc usecase.Usecase, cfg config.Config) (*HttpTxConsumer, error) {
	txConsumer := new(HttpTxConsumer)
	txConsumer.DefaultLastProcessedBlock = cfg.TxConsumerConfig.StartBlock
	txConsumer.CronJobPeriod = cfg.TxConsumerConfig.CronJobPeriod
	txConsumer.BatchLogSize = cfg.TxConsumerConfig.BatchLogSize
	//txConsumer.Addresses = make([]common.Address, 0)
	//for _, address := range cfg.TxConsumerConfig.Addresses {
	//	txConsumer.Addresses = append(txConsumer.Addresses, common.HexToAddress(address))
	//}
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
	defaultStartBlock := c.DefaultLastProcessedBlock
	redisKey := c.getRedisKey()
	value, _ := c.Cache.GetData(redisKey)
	var processingBlock int64 = 1
	if value != nil {
		processingBlock, _ = strconv.ParseInt(*value, 10, 64)
	}

	if processingBlock > defaultStartBlock {
		defaultStartBlock = processingBlock
	}
	return defaultStartBlock, nil
}

func (c *HttpTxConsumer) setLastProcessedBlock(block int64) error {
	redisKey := c.getRedisKey()
	return c.Cache.SetData(redisKey, block)
}

func (c *HttpTxConsumer) resolveTransaction() error {
	// get last processed block from redis
	ProcessingBlock, err := c.getLastProcessedBlock()
	if err != nil {
		logger.AtLog.Logger.Error("resolveTransaction", zap.Any("err", err))
		return err
	}

	// get new block from db
	lastBlockOnChain, err := c.Blockchain.GetBlockNumber()
	if err != nil {
		logger.AtLog.Logger.Error("resolveTransaction", zap.Any("err", err))
		return err
	}

	for ProcessingBlock <= lastBlockOnChain.Int64() {

		logs, err := c.Blockchain.GetEventLogs(*big.NewInt(ProcessingBlock), *big.NewInt(ProcessingBlock + 1))
		if err != nil {
			logger.AtLog.Logger.Error("err.GetEventLogs", zap.String("err", err.Error()))
			return err
		}

		logger.AtLog.Logger.Info("resolveTransaction", zap.Int64("processing block number", ProcessingBlock))
		if len(logs) > 0 {
			fmt.Println("logs", logs)
		}
		for _, _log := range logs {

			topic := strings.ToLower(_log.Topics[0].String())
			logger.AtLog.Logger.Info("topic", zap.Any("topic", topic), zap.Any("_log.TxHash", _log.TxHash), zap.Any("_log.BlockNumber", _log.BlockNumber))

			switch topic {
			//case c.Config.MarketplaceEvents.PurchaseToken:
			//	err = c.Usecase.ResolveMarketplacePurchaseTokenEvent(_log)
			//case c.Config.MarketplaceEvents.MakeOffer:
			//	err = c.Usecase.ResolveMarketplaceMakeOffer(_log)
			//case c.Config.MarketplaceEvents.AcceptMakeOffer:
			//	err = c.Usecase.ResolveMarketplaceAcceptOfferEvent(_log)
			//case c.Config.MarketplaceEvents.CancelListing:
			//	err = c.Usecase.ResolveMarketplaceCancelListing(_log)
			//case c.Config.MarketplaceEvents.CancelMakeOffer:
			//	err = c.Usecase.ResolveMarketplaceCancelOffer(_log)
			//case c.Config.DAOEvents.ProposalCreated:
			//	err = c.Usecase.DAOProposalCreated(_log)
			//case c.Config.DAOEvents.CastVote:
			//	err = c.Usecase.DAOCastVote(_log)
			//case c.Config.DAOEvents.TransferNFTSignature:
			//	c.Usecase.UpdateProjectWithListener(_log)
			//case c.Config.MarketplaceEvents.ListToken:
			//	err = c.Usecase.ResolveMarketplaceListTokenEvent(_log)
			}
		}
		if ProcessingBlock < lastBlockOnChain.Int64() {
			ProcessingBlock++
			c.setLastProcessedBlock(ProcessingBlock)
		} else {
			break
		}
	}
	return nil
}

func (c *HttpTxConsumer) StartServer() {
	logger.AtLog.Logger.Info("HttpTxConsumer start listening")
	for {
		err := c.resolveTransaction()
		if err != nil {
			logger.AtLog.Logger.Error("Error when resolve transactions", zap.String("err", err.Error()))
		}
		time.Sleep(1 * time.Minute)
	}
}
