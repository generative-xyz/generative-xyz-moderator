package blockchain

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"rederinghub.io/utils/config"
)

type Blockchain struct {
	client *ethclient.Client
}

func NewBlockchain(cfg config.BlockchainConfig) (*Blockchain, error) {
	ethereumClient, err := ethclient.Dial(cfg.ETHEndpoint)
	if err != nil {
		return nil, err
	}
	return &Blockchain{
		client: ethereumClient,
	}, nil
}

func (a *Blockchain) GetBlockNumber() (*big.Int, error) {
	header, err := a.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return header.Number, nil
}

func (a *Blockchain) GetEventLogs(fromBlock big.Int, toBlock big.Int, addresses []common.Address) ([]types.Log, error) {
	query := ethereum.FilterQuery{
		FromBlock: &fromBlock,
		ToBlock: &toBlock, 
		Addresses: addresses,
	}
	logs, err := a.client.FilterLogs(context.Background(), query)
	if err != nil {
		return nil, err
	}
	return logs, nil
}
