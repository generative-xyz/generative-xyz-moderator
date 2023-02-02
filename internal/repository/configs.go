package repository

import (
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
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

func (r Repository) DeleteConfig(key string) (*mongo.DeleteResult, error) {
	filter := bson.D{{"key", key}}
	result, err := r.DeleteOne(utils.COLLECTION_CONFIGS, filter)
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

func (r Repository) ListConfigs(filter entity.FilterConfigs) (*entity.Pagination, error)  {
	confs := []entity.Configs{}
	resp := &entity.Pagination{}
	f := bson.M{}

	p, err := r.Paginate(utils.COLLECTION_CONFIGS, filter.Page, filter.Limit, f,bson.D{}, []Sort{}, &confs)
	if err != nil {
		return nil, err
	}
	
	resp.Result = confs
	resp.Page = p.Pagination.Page
	resp.Total = p.Pagination.Total
	resp.PageSize = filter.Limit
	return resp, nil
}

func (r Repository) UpdateConfig(key string, conf *entity.Configs) (*mongo.UpdateResult, error) {
	filter := bson.D{{"key", key}}
	result, err := r.UpdateOne(conf.TableName(), filter, conf)
	if err != nil {
		return nil, err
	}

	return result, nil
}
