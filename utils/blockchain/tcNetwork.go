package blockchain

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"go.uber.org/zap"
	"rederinghub.io/utils/config"
	"rederinghub.io/utils/logger"
)

type TcNetwork struct {
	client *ethclient.Client
}

func NewTcNetwork(cfg config.BlockchainConfig) (*TcNetwork, error) {
	ethereumClient, err := ethclient.Dial(os.Getenv("TC_ENDPOINT"))
	if err != nil {
		return nil, err
	}
	return &TcNetwork{
		client: ethereumClient,
	}, nil
}

func (a *TcNetwork) GetBlockNumber() (*big.Int, error) {
	header, err := a.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		logger.AtLog.Logger.Error("GetBlockNumber", zap.Error(err))
		return nil, err
	}

	logger.AtLog.Logger.Info("GetBlockNumber", zap.Any("header", header))
	return header.Number, nil
}

func (a *TcNetwork) GetBlock() (*types.Header, error) {
	return a.client.HeaderByNumber(context.Background(), nil)
}

func (a *TcNetwork) GetEventLogs(fromBlock big.Int, toBlock big.Int, addresses []common.Address) ([]types.Log, error) {
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

func (a *TcNetwork) GetBlockByNumber(blockNumber big.Int) (*types.Block, error) {
	block, err := a.client.BlockByNumber(context.Background(), &blockNumber)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (a *TcNetwork) GetClient() *ethclient.Client {
	return a.client
}
