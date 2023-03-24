package repository

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"rederinghub.io/internal/entity"
)

func (r Repository) InsertTokenTx(data *entity.TokenTx) error {
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) UpsertTokenTx(inscriptionID string, tx string, data *entity.TokenTx) (*mongo.UpdateResult, error) {
	data.SetUpdatedAt()
	data.SetCreatedAt()
	data.SetTokenTxID()
	filter := bson.D{
		{Key: "inscription_id", Value: inscriptionID},
		{Key: "tx", Value: tx},
	}
	bData, _ := data.ToBson()
	update := bson.D{{Key: "$set", Value: bData}}
	updateOpts := options.Update().SetUpsert(true)
	result, err := r.DB.Collection(entity.TokenTx{}.TableName()).UpdateOne(context.TODO(), filter, update, updateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (r Repository) GetTailTokenTxs(page int64, limit int64) (*entity.Pagination, error) {
	confs := []entity.TokenTx{}
	resp := &entity.Pagination{}
	f := bson.M{"next_tx": "", "num_failed": bson.M{"$lte": 5}}
	s := []Sort{
		{SortBy: "last_time_check", Sort: entity.SORT_ASC},
		{SortBy: "priority", Sort: entity.SORT_DESC},
	}
	p, err := r.Paginate(entity.TokenTx{}.TableName(), page, limit, f, bson.D{}, s, &confs)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp.Result = confs
	resp.Page = p.Pagination.Page
	resp.Total = p.Pagination.Total
	resp.PageSize = limit
	return resp, nil
}

func (r Repository) GetUnresolvedTokenTx(page int64, limit int64) (*entity.Pagination, error) {
	confs := []entity.TokenTx{}
	resp := &entity.Pagination{}
	f := bson.M{"resolved": bson.M{"$ne": true}, "retried_resolve": bson.M{"$lte": 5}}
	s := []Sort{
		{SortBy: "created_at", Sort: entity.SORT_ASC},
	}
	p, err := r.Paginate(entity.TokenTx{}.TableName(), page, limit, f, bson.D{}, s, &confs)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp.Result = confs
	resp.Page = p.Pagination.Page
	resp.Total = p.Pagination.Total
	resp.PageSize = limit
	return resp, nil
}

func (r Repository) UpdateResolvedTx(inscriptionID string, tx string) (*mongo.UpdateResult, error) {
	filter := bson.D{
		{Key: "inscription_id", Value: inscriptionID},
		{Key: "tx", Value: tx},
	}
	update := bson.M{
		"$set": bson.M{"resolved": true},
	}

	result, err := r.DB.Collection(entity.TokenTx{}.TableName()).UpdateMany(context.TODO(), filter, update)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, errors.WithStack(err)
}

func (r Repository) AddTokenTxRetryResolve(inscriptionID string, tx string) (*mongo.UpdateResult, error) {
	filter := bson.D{
		{Key: "inscription_id", Value: inscriptionID},
		{Key: "tx", Value: tx},
	}
	update := bson.M{
		"$inc": bson.M{"retried_resolve": 1},
	}

	result, err := r.DB.Collection(entity.TokenTx{}.TableName()).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, errors.WithStack(err)
}
