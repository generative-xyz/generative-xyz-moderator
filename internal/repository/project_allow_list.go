package repository

import (
	"rederinghub.io/internal/entity"
)

func (r Repository) CreateProjectAllowList(data *entity.ProjectAllowList) error {
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}