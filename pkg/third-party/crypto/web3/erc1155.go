package web3

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/laziercoder/go-web3"
	"github.com/laziercoder/go-web3/eth"
	"go.uber.org/zap"
	"rederinghub.io/pkg/logger"
	"rederinghub.io/pkg/third-party/crypto/constants/abi"
	"rederinghub.io/pkg/third-party/crypto/web3/nftdata"

	"math/big"
)

type clientERC1155Impl struct {
	contract *eth.Contract
}

func (c clientERC1155Impl) NFTFeeds(req *BalanceRequest) (interface{}, error) {
	return nil, nil
}

func NewClientERC1155(providerURL string, contractAddress string) Strategy {
	if providerURL == "" {
		return nil
	}

	web3Client, err := web3.NewWeb3(providerURL)
	if err != nil {
		logger.AtLog.Logger.Error("failed to create web3 client", zap.Error(err))
		return nil
	}

	contract, err := web3Client.Eth.NewContract(abi.EthereumERC1155, contractAddress)
	if err != nil {
		logger.AtLog.Logger.Error("failed to create web3 contract", zap.Error(err))
		return nil
	}

	return &clientERC1155Impl{
		contract: contract,
	}
}

func (c clientERC1155Impl) BalanceOf(req *BalanceRequest) (*nftdata.NFTCustomerInfo, error) {
	ownerAddress := common.HexToAddress(req.Address)
	b, err := c.contract.Call("balanceOf", ownerAddress, big.NewInt(req.TokenID))
	if err != nil {
		return nil, err
	}
	balance, ok := b.(*big.Int)
	if !ok {
		return nil, errors.New("invalid type")
	}

	return getCustomerInfoFromBalance(balance, req.Address)
}
