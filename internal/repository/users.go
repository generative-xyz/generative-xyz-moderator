package repository

import (
	"context"

	"rederinghub.io/internal/model"
)

type IUserRepository interface {
	IRepository
	FindUserByWalletAddress(ctx context.Context, walletAddress string) (*model.Users, error)
}

type userRepository struct {
	repository
}

func NewUserRepository(db model.Database) IUserRepository {
	repository := repository{}
	repository.CollectionName = model.Users{}.CollectionName()
	repository.DB =           db.DB()
	return &userRepository{repository: repository}
}

func (u userRepository) FindUserByWalletAddress(ctx context.Context,walletAddress string) (*model.Users, error) {
	// find in mongo
	user := &model.Users{}
	err := u.FindOne(ctx, map[string]interface{}{
		"wallet_address":         walletAddress,
	}, &user)

	if err != nil {
		return nil, err
	}

	return user, nil
}