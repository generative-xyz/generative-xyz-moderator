package repository

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/helpers"
)

func (r Repository) FindFaucetByStatus(status int) ([]*entity.Faucet, error) {
	var resp []*entity.Faucet
	filter := bson.D{{"status", status}}
	cursor, err := r.DB.Collection(entity.Faucet{}.TableName()).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (r Repository) FindFaucetByTx(tx string) (*entity.Faucet, error) {
	resp := &entity.Faucet{}
	tx = strings.ToLower(tx)
	usr, err := r.FilterOne(entity.Faucet{}.TableName(), bson.D{{"tx", tx}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (r Repository) FindFaucetByAddress(address string) (*entity.Faucet, error) {
	address = strings.ToLower(address)
	resp := &entity.Faucet{}
	usr, err := r.FilterOne(entity.Faucet{}.TableName(), bson.D{{"address", address}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) InsertFaucet(data *entity.Faucet) error {
	data.Address = strings.ToLower(data.Address)
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) UpdateFaucetByUUid(uuid, tx, txBtc string, status int) (*mongo.UpdateResult, error) {
	filter := bson.D{{"uuid", uuid}}
	update := bson.M{"$set": bson.M{"status": status, "tx": tx, "btc_tx": txBtc}}
	result, err := r.DB.Collection(entity.Faucet{}.TableName()).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) UpdateFaucetStatusByTx(tx string, status int) (*mongo.UpdateResult, error) {
	filter := bson.D{{"tx", tx}}
	update := bson.M{"$set": bson.M{"status": status}}
	result, err := r.DB.Collection(entity.Faucet{}.TableName()).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}
