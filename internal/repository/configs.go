package repository

import (
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/helpers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
func (r Repository) FindConfigCustom(key string, result interface{}) error {
	usr, err := r.FilterOne(entity.Configs{}.TableName(), bson.D{{"key", key}})
	if err != nil {
		return err
	}

	err = helpers.Transform(usr, result)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) DeleteConfig(uuid string) (*mongo.DeleteResult, error) {
	filter := bson.D{{"uuid", uuid}}
	result, err := r.DeleteOne(entity.Configs{}.TableName(), filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) InsertConfig(data *entity.Configs) error {
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) ListConfigs(filter entity.FilterConfigs) (*entity.Pagination, error) {
	confs := []entity.Configs{}
	resp := &entity.Pagination{}
	f := bson.M{}

	p, err := r.Paginate(entity.Configs{}.TableName(), filter.Page, filter.Limit, f, bson.D{}, []Sort{}, &confs)
	if err != nil {
		return nil, err
	}
	resp.Result = confs
	resp.Page = p.Pagination.Page
	resp.Total = p.Pagination.Total
	resp.PageSize = filter.Limit
	return resp, nil
}

func (r Repository) UpdateConfig(uuid string, conf *entity.Configs) (*mongo.UpdateResult, error) {
	filter := bson.D{{"uuid", uuid}}
	result, err := r.UpdateOne(conf.TableName(), filter, conf)
	if err != nil {
		return nil, err
	}

	return result, nil
}
