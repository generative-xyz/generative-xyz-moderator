package usecase

import (
	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
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
