package repository

import (
	"context"
	"log"
	"time"

	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r Repository) FindDeveloperInscribeBTC(key string) (*entity.DeveloperInscribe, error) {
	resp := &entity.DeveloperInscribe{}
	usr, err := r.FilterOne(entity.DeveloperInscribe{}.TableName(), bson.D{{utils.KEY_UUID, key}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) FindDeveloperInscribeBTCByNftID(uuid string) (*entity.DeveloperInscribeBTCResp, error) {

	log.Println("uuid:", uuid)

	resp := &entity.DeveloperInscribeBTCResp{}
	usr, err := r.FilterOne(entity.DeveloperInscribe{}.TableName(), bson.D{{"uuid", uuid}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) InsertDeveloperInscribeBTC(data *entity.DeveloperInscribe) error {
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) ListDeveloperInscribeBTC(filter *entity.FilterDeveloperInscribeBT) (*entity.Pagination, error) {
	confs := []entity.DeveloperInscribeBTCResp{}
	resp := &entity.Pagination{}
	f := bson.M{}
	if filter.UserUuid != nil {
		f["user_uuid"] = *filter.UserUuid
	}
	if filter.TokenAddress != nil {
		f["token_address"] = *filter.TokenAddress
	}
	if filter.TokenId != nil {
		f["token_id"] = *filter.TokenId
	}
	if len(filter.NeStatuses) > 0 {
		f["status"] = bson.M{"$nin": filter.NeStatuses}
	}
	p, err := r.Paginate(entity.DeveloperInscribe{}.TableName(), filter.Page, filter.Limit, f, bson.D{}, []Sort{}, &confs)
	if err != nil {
		return nil, err
	}

	resp.Result = confs
	resp.Page = p.Pagination.Page
	resp.Total = p.Pagination.Total
	resp.PageSize = filter.Limit
	return resp, nil
}

func (r Repository) DeveloperUpdateBtcWalletAddressByOrdAddrV2(ordAddress string, conf *entity.DeveloperInscribe) (*mongo.UpdateResult, error) {
	filter := bson.D{{"ordAddress", ordAddress}}
	result, err := r.UpdateOne(conf.TableName(), filter, conf)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) DeveloperListProcessingWalletAddressV2() ([]entity.DeveloperInscribe, error) {
	confs := []entity.DeveloperInscribe{}
	f := bson.M{}
	f["$or"] = []interface{}{
		bson.M{"isMinted": bson.M{"$not": bson.M{"$eq": true}}},
		bson.M{"isConfirm": bson.M{"$not": bson.M{"$eq": true}}},
	}

	opts := options.Find()
	cursor, err := r.DB.Collection(entity.DeveloperInscribe{}.TableName()).Find(context.TODO(), f, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &confs); err != nil {
		return nil, err
	}

	return confs, nil
}

func (r Repository) DeveloperListBTCAddressV2() ([]entity.DeveloperInscribe, error) {
	confs := []entity.DeveloperInscribe{}

	f := bson.M{}
	f["mintResponse"] = bson.M{"$not": bson.M{"$eq": nil}}
	f["mintResponse.issent"] = false
	f["mintResponse.inscription"] = bson.M{"$not": bson.M{"$eq": ""}}

	opts := options.Find()
	cursor, err := r.DB.Collection(entity.DeveloperInscribe{}.TableName()).Find(context.TODO(), f, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &confs); err != nil {
		return nil, err
	}

	return confs, nil
}

// new:
func (r Repository) DeveloperListBTCInscribePending() ([]entity.DeveloperInscribe, error) {
	resp := []entity.DeveloperInscribe{}
	filter := bson.M{
		"status":     entity.StatusDeveloperInscribe_Pending,
		"expired_at": bson.M{"$gte": primitive.NewDateTimeFromTime(time.Now().UTC())},
	}

	cursor, err := r.DB.Collection(entity.DeveloperInscribe{}.TableName()).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
func (r Repository) DeveloperListBTCInscribeByStatus(statuses []entity.StatusDeveloperInscribe) ([]entity.DeveloperInscribe, error) {
	resp := []entity.DeveloperInscribe{}
	filter := bson.M{
		"status": bson.M{"$in": statuses},
	}

	cursor, err := r.DB.Collection(entity.DeveloperInscribe{}.TableName()).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (r Repository) UpdateDeveloperInscribe(model *entity.DeveloperInscribe) (*mongo.UpdateResult, error) {

	filter := bson.D{{Key: "uuid", Value: model.UUID}}
	result, err := r.UpdateOne(model.TableName(), filter, model)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) DeveloperCreateDeveloperInscribeBTCLog(logs *entity.DeveloperInscribeBTCLogs) error {
	err := r.InsertOne(logs.TableName(), logs)
	if err != nil {
		return err
	}
	return nil
}
