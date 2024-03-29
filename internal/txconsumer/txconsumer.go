package txconsumer

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils/blockchain"
	"rederinghub.io/utils/config"
	"rederinghub.io/utils/global"
	"rederinghub.io/utils/logger"
	"rederinghub.io/utils/redis"
)

type HttpTxConsumer struct {
	Blockchain                *blockchain.TcNetwork
	DefaultLastProcessedBlock int64
	CronJobPeriod             int32
	BatchLogSize              int32
	Addresses                 []common.Address
	FetchedAddress            []string
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
		txConsumer.FetchedAddress = append(txConsumer.FetchedAddress, strings.ToLower(address))
		txConsumer.Addresses = append(txConsumer.Addresses, common.HexToAddress(address))
	}

	tcClient, err := blockchain.NewTcPublicNodeNetwork()
	if err != nil {
		return nil, err
	}

	txConsumer.Cache = global.Cache
	txConsumer.Logger = global.Logger
	txConsumer.Blockchain = tcClient
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
		ProcessingBlockTo := ProcessingBlock + int64(c.BatchLogSize)
		logs, err := c.Blockchain.GetEventLogs(*big.NewInt(ProcessingBlock), *big.NewInt(ProcessingBlockTo), c.Addresses)
		if err != nil {
			logger.AtLog.Logger.Error("err.GetEventLogs", zap.String("err", err.Error()))
			return err
		}

		logger.AtLog.Logger.Info("resolveTransaction",
			zap.Int64("fromFlock", ProcessingBlock),
			zap.Int64("toBlock", ProcessingBlockTo),
			zap.Int64("lastBlockOnChain", lastBlockOnChain.Int64()),
			zap.Int64("logs", int64(len(logs))))

		for _, _log := range logs {

			address := strings.ToLower(_log.Address.String())
			topic := strings.ToLower(_log.Topics[0].String())

			logger.AtLog.Logger.Info("resolveTransaction",
				zap.Int64("fromFlock", ProcessingBlock),
				zap.Int64("toBlock", ProcessingBlockTo),
				zap.Int64("lastBlockOnChain", lastBlockOnChain.Int64()),
				zap.String("address", address),
				zap.String("topic", topic),
				zap.String("GENERATIVE_PROJECT", strings.ToLower(os.Getenv("GENERATIVE_PROJECT"))),
				zap.String("TransferNFT", strings.ToLower(c.Config.BlockChainEvent.TransferNFT)),
			)

			switch address {
			case strings.ToLower(c.Config.MarketplaceEvents.Contract):
				switch topic {
				case strings.ToLower(c.Config.MarketplaceEvents.PurchaseToken):
					c.Usecase.ResolveMarketplacePurchaseTokenEvent(_log)
				case strings.ToLower(c.Config.MarketplaceEvents.MakeOffer):
					c.Usecase.ResolveMarketplaceMakeOffer(_log)
				case strings.ToLower(c.Config.MarketplaceEvents.AcceptMakeOffer):
					c.Usecase.ResolveMarketplaceAcceptOfferEvent(_log)
				case strings.ToLower(c.Config.MarketplaceEvents.CancelListing):
					c.Usecase.ResolveMarketplaceCancelListing(_log)
				case strings.ToLower(c.Config.MarketplaceEvents.CancelMakeOffer):
					c.Usecase.ResolveMarketplaceCancelOffer(_log)
				case strings.ToLower(c.Config.MarketplaceEvents.ListToken):
					c.Usecase.ResolveMarketplaceListTokenEvent(_log)
				}
			case strings.ToLower(c.Config.DAOEvents.Contract):
				//switch topic {
				//case c.Config.DAOEvents.ProposalCreated:
				//	err = c.Usecase.DAOProposalCreated(_log)
				//case c.Config.DAOEvents.CastVote:
				//	err = c.Usecase.DAOCastVote(_log)
				//}
			case strings.ToLower(os.Getenv("GENERATIVE_PROJECT")): // master project
				switch topic {
				case strings.ToLower(c.Config.BlockChainEvent.TransferNFT):
					c.Usecase.UpdateProjectWithListener(_log)
				}
			default:
				switch topic {
				case strings.ToLower(c.Config.BlockChainEvent.TransferNFT):
					// handle transfer
					err := c.Usecase.UpdateTokenOwner(_log, c.Blockchain)
					if err != nil {
						logger.AtLog.Error("err.UpdateTokenOwner", zap.String("err", err.Error()))
					}
				}
			}
		}

		if ProcessingBlockTo < lastBlockOnChain.Int64() {
			ProcessingBlock = ProcessingBlockTo + 1
			c.setLastProcessedBlock(ProcessingBlock)
		} else {
			ProcessingBlock = lastBlockOnChain.Int64() + 1
			c.setLastProcessedBlock(ProcessingBlock)
			break
		}
	}
	return nil
}

func (c *HttpTxConsumer) getTcAddress() error {
	savedMap := make(map[string]bool)

	for _, address := range c.FetchedAddress {
		savedMap[address] = true
	}

	projects, err := c.Usecase.Repo.GetTCProject(c.FetchedAddress)
	if err != nil {
		return err
	}

	for _, project := range projects {
		address := common.HexToAddress(project.GenNFTAddr)
		if !savedMap[project.GenNFTAddr] {
			savedMap[project.GenNFTAddr] = true
			c.Addresses = append(c.Addresses, address)
			c.FetchedAddress = append(c.FetchedAddress, project.GenNFTAddr)
		}
	}
	return nil
}

func (c *HttpTxConsumer) StartServer() {
	logger.AtLog.Logger.Info("HttpTxConsumer start listening")

	for {
		err := c.getTcAddress()
		if err != nil {
			logger.AtLog.Logger.Error("Error when resolve getTcAddress", zap.String("err", err.Error()))
		}

		err = c.resolveTransaction()
		if err != nil {
			logger.AtLog.Logger.Error("Error when resolve transactions", zap.String("err", err.Error()))
		}
		//time.Sleep(1 * time.Minute)
	}
}
