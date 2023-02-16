package repository

import (
	"rederinghub.io/internal/entity"
)

func (r Repository) InsertActitvy(data *entity.Activity) error {
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}
