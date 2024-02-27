package blockchain

import (
	"context"
	"rederinghub.io/utils/logger"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
	"math/big"
	"os"
	"time"
)

type EthNetwork struct {
	client *ethclient.Client
}

func NewEthNetwork() (*EthNetwork, error) {
	ethereumClient, err := ethclient.Dial(os.Getenv("ETH_ENDPOINT"))
	if err != nil {
		return nil, err
	}

	return &EthNetwork{
		client: ethereumClient,
	}, nil
}

func (a *EthNetwork) GetBlockNumber() (*big.Int, error) {
	header, err := a.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		logger.AtLog.Logger.Error("GetBlockNumber", zap.Error(err))
		return nil, err
	}

	return header.Number, nil
}

func (a *EthNetwork) GetBlock() (*types.Header, error) {
	return a.client.HeaderByNumber(context.Background(), nil)
}

func (a *EthNetwork) GetEventLogs(fromBlock big.Int, toBlock big.Int, addresses []common.Address) ([]types.Log, error) {
	query := ethereum.FilterQuery{
		FromBlock: &fromBlock,
		ToBlock:   &toBlock,
		Addresses: addresses,
	}
	logs, err := a.client.FilterLogs(context.Background(), query)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (a *EthNetwork) GetBlockByNumber(blockNumber big.Int) (*types.Block, error) {
	block, err := a.client.BlockByNumber(context.Background(), &blockNumber)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (a *EthNetwork) GetBlockTimeByNumber(blockNumber big.Int) (*time.Time, error) {
	block, err := a.client.BlockByNumber(context.Background(), &blockNumber)
	if err != nil {
		return nil, err
	}
	blockTime := time.Unix(int64(block.Time()), 0)
	return &blockTime, nil
}

func (a *EthNetwork) TransactionByHash(hash common.Hash) (*types.Transaction, bool, error) {
	block, isPending, err := a.client.TransactionByHash(context.Background(), hash)
	if err != nil {
		return nil, false, err
	}
	return block, isPending, nil
}

func (a *EthNetwork) HeaderByHash(hash common.Hash) (*types.Header, error) {
	block, err := a.client.HeaderByHash(context.Background(), hash)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (a *EthNetwork) GetClient() *ethclient.Client {
	return a.client
}
