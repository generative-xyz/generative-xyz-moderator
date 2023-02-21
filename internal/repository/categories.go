package repository

import (
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r Repository) FindCategory(key string) (*entity.Categories, error) {
	resp := &entity.Categories{}
	usr, err := r.FilterOne(entity.Categories{}.TableName(), bson.D{{utils.KEY_UUID, key}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) DeleteCategory(id string) (*mongo.DeleteResult, error) {
	filter := bson.D{{utils.KEY_UUID, id}}
	result, err := r.DeleteOne(entity.Categories{}.TableName(), filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) InsertCategory(data *entity.Categories) error {
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) ListCategories(filter entity.FilterCategories) (*entity.Pagination, error)  {
	confs := []entity.Categories{}
	resp := &entity.Pagination{}
	f := bson.M{}

	p, err := r.Paginate(entity.Categories{}.TableName(), filter.Page, filter.Limit, f, bson.D{},[]Sort{{SortBy: filter.SortBy, Sort: filter.Sort}} , &confs)
	if err != nil {
		return nil, err
	}
resp.Result = confs
	resp.Page = p.Pagination.Page
	resp.Total = p.Pagination.Total
	resp.PageSize = filter.Limit
	return resp, nil
}

func (r Repository) UpdateCategory(key string, conf *entity.Categories) (*mongo.UpdateResult, error) {
	filter := bson.D{{"key", key}}
	result, err := r.UpdateOne(conf.TableName(), filter, conf)
	if err != nil {
		return nil, err
	}

	return result, nil
}
