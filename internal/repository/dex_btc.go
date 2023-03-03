package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

func (r Repository) CreateDexBTCListing(listing *entity.DexBTCListing) error {
	err := r.InsertOne(listing.TableName(), listing)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) GetDexBTCListingOrderByID(id string) (*entity.DexBTCListing, error) {
	resp := &entity.DexBTCListing{}

	f := bson.D{
		{Key: "uuid", Value: id},
	}

	orderInfo, err := r.FilterOne(utils.COLLECTION_DEX_BTC_LISTING, f)
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(orderInfo, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) UpdateDexBTCListingOrderCancelTx(model *entity.DexBTCListing) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: "uuid", Value: model.UUID}}

	update := bson.M{
		"$set": bson.M{
			"cancel_tx": model.CancelTx,
			"cancel_at": model.CancelAt,
			"cancelled": model.Cancelled,
		},
	}

	result, err := r.DB.Collection(model.TableName()).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (r Repository) UpdateDexBTCListingOrderMatchTx(model *entity.DexBTCListing) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: "uuid", Value: model.UUID}}

	update := bson.M{
		"$set": bson.M{
			"matched":    model.Matched,
			"matched_tx": model.MatchedTx,
			"matched_at": model.MatchAt,
		},
	}

	result, err := r.DB.Collection(model.TableName()).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (r Repository) GetDexBTCListingOrderUserPending(user_address string) ([]entity.DexBTCListing, error) {
	listings := []entity.DexBTCListing{}
	f := bson.D{{
		Key:   "seller_address",
		Value: user_address,
	},
		{Key: "matched", Value: false},
		{Key: "cancelled", Value: false}}

	cursor, err := r.DB.Collection(utils.COLLECTION_DEX_BTC_LISTING).Find(context.TODO(), f, &options.FindOptions{
		Sort: bson.D{{Key: "created_at", Value: -1}},
	})
	if err != nil {
		return nil, err
	}

	if err = cursor.All((context.TODO()), &listings); err != nil {
		return nil, err
	}

	return listings, nil
}

func (r Repository) GetDexBTCListingOrderUser(user_address string, limit, offset int64) ([]entity.DexBTCListing, error) {
	listings := []entity.DexBTCListing{}
	f := bson.D{{
		Key:   "seller_address",
		Value: user_address,
	}}

	cursor, err := r.DB.Collection(utils.COLLECTION_DEX_BTC_LISTING).Find(context.TODO(), f, &options.FindOptions{
		Sort:  bson.D{{Key: "created_at", Value: -1}},
		Limit: &limit,
		Skip:  &offset,
	})
	if err != nil {
		return nil, err
	}

	if err = cursor.All((context.TODO()), &listings); err != nil {
		return nil, err
	}

	return listings, nil
}

func (r Repository) GetDexBTCListingOrderPendingByInscriptionID(id string) (*entity.DexBTCListing, error) {
	resp := &entity.DexBTCListing{}

	f := bson.D{
		{Key: "inscription_id", Value: id},
		{Key: "matched", Value: false},
		{Key: "cancelled", Value: false},
	}

	orderInfo, err := r.FilterOne(utils.COLLECTION_DEX_BTC_LISTING, f, &options.FindOneOptions{
		Sort: bson.D{{Key: "created_at", Value: -1}},
	})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(orderInfo, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) GetDexBTCListingOrderPending() ([]entity.DexBTCListing, error) {
	listings := []entity.DexBTCListing{}
	f := bson.D{
		{Key: "matched", Value: false},
		{Key: "cancelled", Value: false}}

	cursor, err := r.DB.Collection(utils.COLLECTION_DEX_BTC_LISTING).Find(context.TODO(), f, &options.FindOptions{
		Sort: bson.D{{Key: "created_at", Value: -1}},
	})
	if err != nil {
		return nil, err
	}

	if err = cursor.All((context.TODO()), &listings); err != nil {
		return nil, err
	}

	return listings, nil
}
