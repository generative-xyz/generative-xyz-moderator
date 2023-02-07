package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
)

func (r Repository) CreateTokenHolders(tokenHolders []entity.TokenHolder) error {
	_tokenHolders := make([]entity.IEntity, 0, len(tokenHolders))
	for _, tokenHolder := range tokenHolders {
		_tokenHolder := tokenHolder
		_tokenHolders = append(_tokenHolders, &_tokenHolder)
	}
	err := r.InsertMany(utils.COLLECTION_LEADERBOARD_TOKEN_HOLDER, _tokenHolders)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) GetAllTokenHolders() ([]entity.TokenHolder, error) {
	tokenHolders := []entity.TokenHolder{}
	f := bson.M{}

	cursor, err := r.DB.Collection(utils.COLLECTION_LEADERBOARD_TOKEN_HOLDER).Find(context.TODO(), f)
	if err != nil {
		return nil, err
	}

	if err = cursor.All((context.TODO()), &tokenHolders); err != nil {
		return nil, err
	}

	return tokenHolders, nil
}

func (r Repository) DeleteAllTokenHolders() error {
	_, err := r.DeleteMany(utils.COLLECTION_LEADERBOARD_TOKEN_HOLDER, bson.D{})
	return err
}

func (r Repository) FilterTokenHolders(filter entity.FilterTokenHolders) (*entity.Pagination, error) {
	confs := []entity.TokenHolder{}
	resp := &entity.Pagination{}
	p, err := r.Paginate(utils.COLLECTION_LEADERBOARD_TOKEN_HOLDER, filter.Page, filter.Limit, bson.M{}, bson.M{}, []Sort{}, &confs)
	if err != nil {
		return nil, err
	}

	resp.Result = confs
	resp.Page = p.Pagination.Page
	resp.Total = p.Pagination.Total
	resp.PageSize = filter.Limit
	return resp, nil
} 
