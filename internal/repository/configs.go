package repository

import (
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/helpers"

	"go.mongodb.org/mongo-driver/bson"
)

func (r Repository) FindConfig(key string) (*entity.Configs, error) {
	resp := &entity.Configs{}
	usr, err := r.FilterOne(entity.Configs{}.TableName(), bson.D{{"key", key}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) InsertConfig(data *entity.Configs) error {
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}

