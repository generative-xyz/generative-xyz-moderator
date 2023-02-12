package repository

import (
	"rederinghub.io/internal/entity"
)

func (r Repository) CreateTokenUriHistory(data *entity.TokenUriHistories) error {
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}
