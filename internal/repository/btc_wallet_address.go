package repository

import (
	"context"

	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (r Repository) FindBtcWalletAddressByOrd(ordAddress string) (*entity.BTCWalletAddress, error) {
	resp := &entity.BTCWalletAddress{}
	usr, err := r.FilterOne(entity.BTCWalletAddress{}.TableName(), bson.D{{"ordAddress", ordAddress}})
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

func (r Repository) ListBtcWalletAddress(filter entity.FilterBTCWalletAddress) (*entity.Pagination, error) {
	confs := []entity.BTCWalletAddress{}
	resp := &entity.Pagination{}
	f := bson.M{}

	p, err := r.Paginate(entity.BTCWalletAddress{}.TableName(), filter.Page, filter.Limit, f, bson.D{}, []Sort{}, &confs)
	if err != nil {
		return nil, err
	}

	resp.Result = confs
	resp.Page = p.Pagination.Page
	resp.Total = p.Pagination.Total
	resp.PageSize = filter.Limit
	return resp, nil
}

func (r Repository) UpdateBtcWalletAddressByOrdAddr(ordAddress string, conf *entity.BTCWalletAddress) (*mongo.UpdateResult, error) {
	filter := bson.D{{"ordAddress", ordAddress}}
	result, err := r.UpdateOne(conf.TableName(), filter, conf)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) UpdateBtcWalletAddress(model *entity.BTCWalletAddress) (*mongo.UpdateResult, error) {

	filter := bson.D{{"uuid", model.UUID}}
	result, err := r.UpdateOne(model.TableName(), filter, model)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) ListProcessingWalletAddress() ([]entity.BTCWalletAddress, error) {
	confs := []entity.BTCWalletAddress{}
	f := bson.M{}
	f["$or"] = []interface{}{
		bson.M{"isMinted": false},
		bson.M{"isConfirm": false},
	}
	f["balanceCheckTime"] = bson.M{"$lt": utils.MAX_CHECK_BALANCE}

	opts := options.Find()
	cursor, err := r.DB.Collection(utils.COLLECTION_BTC_WALLET_ADDRESS).Find(context.TODO(), f, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &confs); err != nil {
		return nil, err
	}

	return confs, nil
}

func (r Repository) ListBTCAddress() ([]entity.BTCWalletAddress, error) {
	confs := []entity.BTCWalletAddress{}

	f := bson.M{}
	f["mintResponse"] = bson.M{"$not": bson.M{"$eq": nil}}
	f["mintResponse.issent"] = false
	f["mintResponse.inscription"] = bson.M{"$not": bson.M{"$eq": ""}}

	opts := options.Find()
	cursor, err := r.DB.Collection(utils.COLLECTION_BTC_WALLET_ADDRESS).Find(context.TODO(), f, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &confs); err != nil {
		return nil, err
	}

	return confs, nil
}

//list to claim btc:
func (r Repository) ListWalletAddressToClaimBTC() ([]entity.BTCWalletAddress, error) {
	resp := []entity.BTCWalletAddress{}

	filter := bson.M{
		"isConfirm":      true,
		"uuid":           bson.M{"$gt": "63ea272eb020796632eb8811"},
		"is_sent_master": false,
	}

	cursor, err := r.DB.Collection(utils.COLLECTION_BTC_WALLET_ADDRESS).Find(context.TODO(), filter)

	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
