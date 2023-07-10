package usecase

import (
	"context"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/contracts/erc20"
	"rederinghub.io/utils/logger"
)

func (u Usecase) ListArtist(req entity.FilteArtist) (*entity.Pagination, error) {
	artists, err := u.Repo.ListArtist(req)
	if err != nil {
		logger.AtLog.Logger.Error("u.Repo.ListArtist", zap.Error(err))
		return nil, err
	}

	logger.AtLog.Logger.Info("artists", zap.Any("artists.Total", artists.Total))
	return artists, nil
}

func (u Usecase) ListUsers(req structure.FilterUsers) (*entity.Pagination, error) {
	users, err := u.Repo.ListUsers(req)
	if err != nil {
		logger.AtLog.Logger.Error("u.Repo.ListUsers", zap.Error(err))
		return nil, err
	}

	logger.AtLog.Logger.Info("ListUsers", zap.Any("users.Total", users.Total))
	return users, nil
}

func (u Usecase) GetUsersMap(addresses []string) (map[string]*entity.Users, error) {
	users, err := u.Repo.FindUserByAddresses(addresses)
	if err != nil {
		return nil, err
	}
	addressToUser := map[string]*entity.Users{}
	for id, user := range users {
		if user.WalletAddress != "" {
			addressToUser[user.WalletAddress] = &users[id]
		}
		if user.WalletAddressBTC != "" {
			addressToUser[user.WalletAddressBTC] = &users[id]
		}
		if user.WalletAddressBTCTaproot != "" {
			addressToUser[user.WalletAddressBTCTaproot] = &users[id]
		}
		if user.WalletAddressPayment != "" {
			addressToUser[user.WalletAddressPayment] = &users[id]
		}
	}
	return addressToUser, nil
}

func (u Usecase) GetGMBalance(walletAddress string) (*big.Int, error) {
	erc20Contract, err := erc20.NewErc20(common.HexToAddress(strings.ToLower(os.Getenv("GM_CONTRACT_ADDRESS"))), u.TcClientPublicNode.GetClient())
	if err != nil {
		return nil, err
	}
	balance, err := erc20Contract.BalanceOf(&bind.CallOpts{Context: context.Background()}, common.HexToAddress(walletAddress))
	if err != nil {
		return nil, err
	}

	return balance, nil
}
