package usecase

import (
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
)

func (u Usecase) ListArtist(req entity.FilteArtist) (*entity.Pagination, error) {
	artists, err := u.Repo.ListArtist(req)
	if err != nil {
		u.Logger.Error("u.Repo.ListArtist", err.Error(), err)
		return nil, err
	}

	u.Logger.Info("artists", artists.Total)
	return artists, nil
}

func (u Usecase) ListUsers(req structure.FilterUsers) (*entity.Pagination, error) {
	users, err := u.Repo.ListUsers(req)
	if err != nil {
		u.Logger.Error("u.Repo.ListUsers", err.Error(), err)
		return nil, err
	}

	u.Logger.Info("ListUsers", users.Total)
	return users, nil
}

func (u Usecase) GetUsersMap(addresses []string) (map[string]*entity.Users, error) {
	users, err := u.Repo.FindUserByAddresses(addresses)
	if err != nil {
		return nil, err
	}
	addressToUser := map[string]*entity.Users{}
	for _, user := range users {
		if user.WalletAddress != "" {
			addressToUser[user.WalletAddress] = &user
		}
		if user.WalletAddressBTC != "" {
			addressToUser[user.WalletAddressBTC] = &user
		}
		if user.WalletAddressBTCTaproot != "" {
			addressToUser[user.WalletAddressBTCTaproot] = &user
		}
		if user.WalletAddressPayment != "" {
			addressToUser[user.WalletAddressPayment] = &user
		}
	}
	return addressToUser, nil
}
