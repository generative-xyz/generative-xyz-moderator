package ordinals

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/logger"
)

type Service struct {
	fromAddress     common.Address
	ethClient       *ethclient.Client
	auth            *bind.TransactOpts
	ordinalsService *Ordinals
}

func NewService(ordinalsContractStr, privateKeyStr string, chainId int64) (*Service, error) {
	ethClient, err := helpers.EthDialer()
	if err != nil {
		logger.AtLog.Logger.Error("EthDialer", zap.Error(err))
		return nil, err
	}
	ordinalsContract := common.HexToAddress(ordinalsContractStr)
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		logger.AtLog.Logger.Error("HexToECDSA failed", zap.Error(err))
		return nil, err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		logger.AtLog.Logger.Error("publicKeyECDSA failed")
		return nil, errors.New("Get PublicKeyECDSA failed")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	ords, err := NewOrdinals(ordinalsContract, ethClient)
	if err != nil {
		logger.AtLog.Logger.Error("NewOrdinals failed", zap.Error(err))
		return nil, err
	}
	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		logger.AtLog.Logger.Error("SuggestGasPrice failed", zap.Error(err))
		return nil, err
	}
	// set caller
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(chainId))
	if err != nil {
		logger.AtLog.Logger.Error("NewKeyedTransactorWithChainID failed", zap.Error(err))
		return nil, err
	}
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(100000)
	auth.GasPrice = gasPrice
	return &Service{
		fromAddress:     fromAddress,
		ethClient:       ethClient,
		auth:            auth,
		ordinalsService: ords,
	}, nil
}

func (s *Service) AddContractToOrdinalsContract(ctx context.Context, tokenAddr, tokenId, inscriptionID string) (string, uint64, error) {
	tokenAddress := common.HexToAddress(tokenAddr)
	tokenID := new(big.Int)
	tokenID, ok := tokenID.SetString(tokenId, 10)
	if !ok {
		return "", 0, errors.New("tokenId is wrong")
	}
	nonce, err := s.ethClient.PendingNonceAt(ctx, s.fromAddress)
	if err != nil {
		return "", 0, err
	}
	s.auth.Nonce = big.NewInt(int64(nonce))
	tx, err := s.ordinalsService.SetInscription(s.auth, tokenAddress, tokenID, inscriptionID)
	if err != nil {
		return "", 0, err
	}
	var status uint64
	receipt, err := s.ethClient.TransactionReceipt(ctx, tx.Hash())
	if err == nil {
		status = receipt.Status
	}
	return tx.Hash().Hex(), status, nil
}
