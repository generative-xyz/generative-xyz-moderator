package repository

import (
	"context"

	"rederinghub.io/internal/model"
)

type ITokenUriRepository interface {
	IRepository
	FindTokenBy(ctx context.Context,contractAddr string, tokenID string) (*model.TokenUri, error)
}

type tokenURIRepository struct {
	repository
}

func NewTokenURIRepository(db model.Database) ITokenUriRepository {
	repository := repository{}
	repository.CollectionName = model.TokenUri{}.CollectionName()
	repository.DB =           db.DB()
	return &tokenURIRepository{repository: repository}
}

func (u tokenURIRepository) FindTokenBy(ctx context.Context,contractAddr string, tokenID string) (*model.TokenUri, error) {
	// find in mongo
	tokenUri := &model.TokenUri{}
	err := u.FindOne(ctx, map[string]interface{}{
		"contract_address":         contractAddr,
		"token_id":         tokenID,
	}, &tokenUri)

	if err != nil {
		return nil, err
	}

	return tokenUri, nil
}