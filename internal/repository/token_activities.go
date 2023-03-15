package repository

import (
	"rederinghub.io/internal/entity"
)

func (r Repository) InsertTokenActivity(data *entity.TokenActivity) error {
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}
