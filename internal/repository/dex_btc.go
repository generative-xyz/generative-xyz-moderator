package repository

import (
	"context"

	"github.com/pkg/errors"
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
func (r Repository) CreateDexBTCBuyWithETH(listing *entity.DexBTCBuyWithETH) error {
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
			"buyer":      model.Buyer,
		},
	}

	result, err := r.DB.Collection(model.TableName()).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (r Repository) UpdateDexBTCListingOrderConfirm(model *entity.DexBTCListing) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: "uuid", Value: model.UUID}}

	update := bson.M{
		"$set": bson.M{
			"verified": model.Verified,
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
		{Key: "cancelled", Value: false},
	}

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

func (r Repository) GetDexBtcsAlongWithProjectInfo(req entity.GetDexBtcListingWithProjectInfoReq) ([]entity.DexBtcListingWithProjectInfo, error) {
	aggregates := bson.A{
		bson.D{{"$sort", bson.D{{"created_at", -1}}}},
		bson.D{{"$skip", (req.Page - 1) * req.Limit}},
		bson.D{{"$limit", req.Limit}},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "token_uri"},
					{"localField", "inscription_id"},
					{"foreignField", "token_id"},
					{"pipeline",
						bson.A{
							bson.D{{"$project", bson.D{{"project_id", 1}}}},
						},
					},
					{"as", "project_info"},
				},
			},
		},
	}
	cursor, err := r.DB.Collection(entity.DexBTCListing{}.TableName()).Aggregate(context.TODO(), aggregates)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	listings := []entity.DexBtcListingWithProjectInfo{}

	if err = cursor.All((context.TODO()), &listings); err != nil {
		return nil, errors.WithStack(err)
	}

	return listings, nil
}

func (r Repository) GetDexBTCBuyETHOrderByStatus(statuses []entity.DexBTCETHBuyStatus) ([]entity.DexBTCBuyWithETH, error) {
	listings := []entity.DexBTCBuyWithETH{}
	f := bson.M{
		"status": bson.M{"$in": statuses},
	}
	cursor, err := r.DB.Collection(utils.COLLECTION_DEX_BTC_BUY_ETH).Find(context.TODO(), f, &options.FindOptions{
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

func (r Repository) GetDexBTCBuyETHOrderByUserID(userID string, limit, offset int64) ([]entity.DexBTCBuyWithETH, error) {
	listings := []entity.DexBTCBuyWithETH{}
	f := bson.M{
		"user_id": bson.M{"$eq": userID},
	}
	cursor, err := r.DB.Collection(utils.COLLECTION_DEX_BTC_BUY_ETH).Find(context.TODO(), f, &options.FindOptions{
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
