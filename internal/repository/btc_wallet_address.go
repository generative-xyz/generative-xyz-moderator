package repository

import (
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r Repository) FindBtcWalletAddress(key string) (*entity.BTCWalletAddress, error) {
	resp := &entity.BTCWalletAddress{}
	usr, err := r.FilterOne(entity.BTCWalletAddress{}.TableName(), bson.D{{utils.KEY_UUID, key}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) DeleteBtcWalletAddress(id string) (*mongo.DeleteResult, error) {
	filter := bson.D{{utils.KEY_UUID, id}}
	result, err := r.DeleteOne(entity.BTCWalletAddress{}.TableName(), filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) InsertBtcWalletAddress(data *entity.BTCWalletAddress) error {
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) ListBtcWalletAddress(filter entity.FilterBTCWalletAddress) (*entity.Pagination, error)  {
	confs := []entity.BTCWalletAddress{}
	resp := &entity.Pagination{}
	f := bson.M{}

	p, err := r.Paginate(entity.BTCWalletAddress{}.TableName(), filter.Page, filter.Limit, f, bson.D{},[]Sort{} , &confs)
	if err != nil {
		return nil, err
	}
	
	resp.Result = confs
	resp.Page = p.Pagination.Page
	resp.Total = p.Pagination.Total
	resp.PageSize = filter.Limit
	return resp, nil
}

func (r Repository) UpdateBtcWalletAddress(key string, conf *entity.BTCWalletAddress) (*mongo.UpdateResult, error) {
	filter := bson.D{{"key", key}}
	result, err := r.UpdateOne(conf.TableName(), filter, conf)
	if err != nil {
		return nil, err
	}

	return result, nil
}
