package repository

import (
	"rederinghub.io/internal/entity"
)

func (r Repository) CreateTokenUriMetadata(data *entity.TokenUriMetadata) error {
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}
