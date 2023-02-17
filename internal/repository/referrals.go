package repository

import (
	"rederinghub.io/internal/entity"
)

func (r Repository) InsertReferral(data *entity.Referral) error {
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}
