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

func (r Repository) FindEthWalletAddress(key string) (*entity.ETHWalletAddress, error) {
	resp := &entity.ETHWalletAddress{}
	usr, err := r.FilterOne(entity.ETHWalletAddress{}.TableName(), bson.D{{utils.KEY_UUID, key}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) FindEthWalletAddressByUserAddress(userAddress string) (*entity.ETHWalletAddress, error) {
	resp := &entity.ETHWalletAddress{}
	usr, err := r.FilterOne(entity.ETHWalletAddress{}.TableName(), bson.D{{"user_address", userAddress}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) FindDelegateEthWalletAddressByUserAddress(userAddress string) (*entity.ETHWalletAddress, error) {
	resp := &entity.ETHWalletAddress{}
	usr, err := r.FilterOne(entity.ETHWalletAddress{}.TableName(), bson.D{{"user_address", userAddress}, {"isUseDelegate", true}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) FindEthWalletAddressByOrd(ordAddress string) (*entity.ETHWalletAddress, error) {
	resp := &entity.ETHWalletAddress{}
	usr, err := r.FilterOne(entity.ETHWalletAddress{}.TableName(), bson.D{{"ordAddress", ordAddress}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) DeleteEthWalletAddress(id string) (*mongo.DeleteResult, error) {
	filter := bson.D{{utils.KEY_UUID, id}}
	result, err := r.DeleteOne(entity.ETHWalletAddress{}.TableName(), filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) InsertEthWalletAddress(data *entity.ETHWalletAddress) error {
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) ListEthWalletAddress(filter entity.FilterBTCWalletAddress) (*entity.Pagination, error) {
	confs := []entity.ETHWalletAddress{}
	resp := &entity.Pagination{}
	f := bson.M{}

	p, err := r.Paginate(entity.ETHWalletAddress{}.TableName(), filter.Page, filter.Limit, f, bson.D{}, []Sort{}, &confs)
	if err != nil {
		return nil, err
	}

	resp.Result = confs
	resp.Page = p.Pagination.Page
	resp.Total = p.Pagination.Total
	resp.PageSize = filter.Limit
	return resp, nil
}

func (r Repository) UpdateEthWalletAddressByOrdAddr(ordAddress string, conf *entity.ETHWalletAddress) (*mongo.UpdateResult, error) {
	filter := bson.D{{"ordAddress", ordAddress}}
	result, err := r.UpdateOne(conf.TableName(), filter, conf)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) ListProcessingETHWalletAddress() ([]entity.ETHWalletAddress, error) {
	confs := []entity.ETHWalletAddress{}
	f := bson.M{}
	f["$or"] = []interface{}{
		bson.M{"isMinted": false},
		bson.M{"isConfirm": false},
	}
	f["balanceCheckTime"] = bson.M{"$lt": utils.MAX_CHECK_BALANCE}

	opts := options.Find()
	cursor, err := r.DB.Collection(utils.COLLECTION_ETH_WALLET_ADDRESS).Find(context.TODO(), f, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &confs); err != nil {
		return nil, err
	}

	return confs, nil
}

func (r Repository) ListETHAddress() ([]entity.ETHWalletAddress, error) {
	confs := []entity.ETHWalletAddress{}

	f := bson.M{}
	f["mintResponse"] = bson.M{"$not": bson.M{"$eq": nil}}
	f["mintResponse.issent"] = false
	f["mintResponse.inscription"] = bson.M{"$not": bson.M{"$eq": ""}}

	opts := options.Find()
	cursor, err := r.DB.Collection(utils.COLLECTION_ETH_WALLET_ADDRESS).Find(context.TODO(), f, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &confs); err != nil {
		return nil, err
	}

	return confs, nil
}
