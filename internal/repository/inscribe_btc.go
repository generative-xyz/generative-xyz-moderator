package repository

import (
	"context"
	"time"

	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r Repository) FindInscribeBTC(key string) (*entity.InscribeBTC, error) {
	resp := &entity.InscribeBTC{}
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

func (r Repository) InsertInscribeBTC(data *entity.InscribeBTC) error {
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) ListInscribeBTC(filter entity.FilterBTCWalletAddress) (*entity.Pagination, error) {
	confs := []entity.InscribeBTC{}
	resp := &entity.Pagination{}
	f := bson.M{}

	p, err := r.Paginate(entity.InscribeBTC{}.TableName(), filter.Page, filter.Limit, f, bson.D{}, []Sort{}, &confs)
	if err != nil {
		return nil, err
	}

	resp.Result = confs
	resp.Page = p.Pagination.Page
	resp.Total = p.Pagination.Total
	resp.PageSize = filter.Limit
	return resp, nil
}

func (r Repository) UpdateBtcWalletAddressByOrdAddrV2(ordAddress string, conf *entity.InscribeBTC) (*mongo.UpdateResult, error) {
	filter := bson.D{{"ordAddress", ordAddress}}
	result, err := r.UpdateOne(conf.TableName(), filter, conf)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) ListProcessingWalletAddressV2() ([]entity.InscribeBTC, error) {
	confs := []entity.InscribeBTC{}
	f := bson.M{}
	f["$or"] = []interface{}{
		bson.M{"isMinted": bson.M{"$not": bson.M{"$eq": true}}},
		bson.M{"isConfirm": bson.M{"$not": bson.M{"$eq": true}}},
	}

	opts := options.Find()
	cursor, err := r.DB.Collection(utils.INSCRIBE_BTC).Find(context.TODO(), f, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &confs); err != nil {
		return nil, err
	}

	return confs, nil
}

func (r Repository) ListBTCAddressV2() ([]entity.InscribeBTC, error) {
	confs := []entity.InscribeBTC{}

	f := bson.M{}
	f["mintResponse"] = bson.M{"$not": bson.M{"$eq": nil}}
	f["mintResponse.issent"] = false
	f["mintResponse.inscription"] = bson.M{"$not": bson.M{"$eq": ""}}

	opts := options.Find()
	cursor, err := r.DB.Collection(utils.INSCRIBE_BTC).Find(context.TODO(), f, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &confs); err != nil {
		return nil, err
	}

	return confs, nil
}

// new:
func (r Repository) ListBTCInscribePending() ([]entity.InscribeBTC, error) {
	resp := []entity.InscribeBTC{}
	filter := bson.M{
		"status":     entity.StatusInscribe_Pending,
		"expired_at": bson.M{"$gte": primitive.NewDateTimeFromTime(time.Now().UTC())},
	}

	cursor, err := r.DB.Collection(utils.INSCRIBE_BTC).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
func (r Repository) ListBTCInscribeByStatus(statuses []entity.StatusInscribe) ([]entity.InscribeBTC, error) {
	resp := []entity.InscribeBTC{}
	filter := bson.M{
		"status": bson.M{"$in": statuses},
	}

	cursor, err := r.DB.Collection(utils.INSCRIBE_BTC).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (r Repository) UpdateBtcInscribe(model *entity.InscribeBTC) (*mongo.UpdateResult, error) {

	filter := bson.D{{Key: "uuid", Value: model.UUID}}
	result, err := r.UpdateOne(model.TableName(), filter, model)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) CreateInscribeBTCLog(logs *entity.InscribeBTCLogs) error {
	err := r.InsertOne(logs.TableName(), logs)
	if err != nil {
		return err
	}
	return nil
}
