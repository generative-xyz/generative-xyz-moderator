package repository

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/helpers"
)

func (r Repository) ListFaucetByStatus(statuses []int) ([]*entity.Faucet, error) {
	resp := []*entity.Faucet{}
	filter := bson.M{
		"status": bson.M{"$in": statuses},
	}

	cursor, err := r.DB.Collection(entity.Faucet{}.TableName()).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (r Repository) FindFaucetByStatus(status int) ([]*entity.Faucet, error) {
	var resp []*entity.Faucet
	filter := bson.D{{"status", status}, {"faucet_type", bson.D{{"$ne", "dev"}}}}
	cursor, err := r.DB.Collection(entity.Faucet{}.TableName()).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (r Repository) FindFaucetBySharedID(sharedID string) ([]*entity.Faucet, error) {
	var resp []*entity.Faucet
	filter := bson.D{{"twitter_share_id", sharedID}}
	cursor, err := r.DB.Collection(entity.Faucet{}.TableName()).Find(context.TODO(), filter)
	if err != nil {
		return resp, err
	}

	if err = cursor.All(context.TODO(), &resp); err != nil {
		return resp, err
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

func (r Repository) InsertFaucet(data *entity.Faucet) error {
	data.Address = strings.ToLower(data.Address)
	data.TwitterName = strings.ToLower(data.TwitterName)
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

func (r Repository) FindFaucetByAddress(address string) ([]*entity.Faucet, error) {
	address = strings.ToLower(address)
	var resp []*entity.Faucet
	filter := bson.D{{"address", address}}
	cursor, err := r.DB.Collection(entity.Faucet{}.TableName()).Find(context.TODO(), filter)
	if err != nil {
		return resp, err
	}

	if err = cursor.All(context.TODO(), &resp); err != nil {
		return resp, err
	}

	return resp, nil
}

func (r Repository) FindFaucetByTwitterName(twitterName string) ([]*entity.Faucet, error) {
	twitterName = strings.ToLower(twitterName)
	var resp []*entity.Faucet
	filter := bson.D{{"twitter_name", twitterName}}
	cursor, err := r.DB.Collection(entity.Faucet{}.TableName()).Find(context.TODO(), filter)
	if err != nil {
		return resp, err
	}

	if err = cursor.All(context.TODO(), &resp); err != nil {
		return resp, err
	}

	return resp, nil
}

func (r Repository) FindFaucetByTwitterNameOrAddress(twitterName, address string) ([]*entity.Faucet, error) {

	twitterName = strings.ToLower(twitterName)
	address = strings.ToLower(address)

	var resp []*entity.Faucet

	filters := make(bson.M)

	filters["$or"] = bson.A{
		bson.M{"twitter_name": twitterName},
		bson.M{"address": address},
	}

	cursor, err := r.DB.Collection(entity.Faucet{}.TableName()).Find(context.TODO(), filters, &options.FindOptions{
		Sort: bson.D{{Key: "created_at", Value: -1}},
	})

	if err != nil {
		return resp, err
	}

	if err = cursor.All(context.TODO(), &resp); err != nil {
		return resp, err
	}

	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (r Repository) UpdateFaucet(faucet *entity.Faucet) (*mongo.UpdateResult, error) {
	filter := bson.D{{"uuid", faucet.UUID}}
	result, err := r.UpdateOne(entity.Faucet{}.TableName(), filter, faucet)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) UpdateFaucetByTxTc(tx, txBtc string, status int) (*mongo.UpdateResult, error) {
	filter := bson.D{{"tx", tx}}
	update := bson.M{"$set": bson.M{"status": status, "btc_tx": txBtc}}
	result, err := r.DB.Collection(entity.Faucet{}.TableName()).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) UpdateStatusFaucetByTxTc(tx string, status int) (*mongo.UpdateResult, error) {
	filter := bson.D{{"tx", tx}}
	update := bson.M{"$set": bson.M{"status": status}}
	result, err := r.DB.Collection(entity.Faucet{}.TableName()).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) ListFaucetPending() ([]entity.Faucet, error) {
	faucets := []entity.Faucet{}
	f := bson.M{}

	// var limit int64 = 1
	// var skip int64 = 1

	cursor, err := r.DB.Collection(entity.AuctionCollectionBidder{}.TableName()).Find(context.TODO(), f, &options.FindOptions{
		Sort: bson.D{{"status", 1}},
		// Limit: &limit,
		// Skip:  &skip,
	})
	if err != nil {
		return faucets, err
	}

	if err = cursor.All(context.TODO(), &faucets); err != nil {
		return faucets, err
	}

	return faucets, nil
}
