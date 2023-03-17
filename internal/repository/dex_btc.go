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

func (r Repository) GetDexBTListingOrderByListID(ids []string) ([]entity.DexBTCListing, error) {
	resp := []entity.DexBTCListing{}
	filter := bson.M{
		"uuid": bson.M{"$in": ids},
	}

	cursor, err := r.DB.Collection(entity.DexBTCListing{}.TableName()).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &resp); err != nil {
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

func (r Repository) UpdateDexBTCListingOrderInvalidMatch(model *entity.DexBTCListing) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: "uuid", Value: model.UUID}}

	update := bson.M{
		"$set": bson.M{
			"invalid_match":    model.InvalidMatch,
			"invalid_match_tx": model.InvalidMatchTx,
			"cancelled":        model.Cancelled,
			"cancel_at":        model.CancelAt,
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

func (r Repository) UpdateDexBTCListingTimeseriesData(model *entity.DexBTCListing) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: "uuid", Value: model.UUID}}
	update := bson.M{
		"$set": bson.M{
			"is_time_series_data": model.IsTimeSeriesData,
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

func (r Repository) GetAllDexBTCListingByInscriptionID(inscriptionID string) ([]entity.DexBTCListing, error) {
	listings := []entity.DexBTCListing{}
	f := bson.D{
		{Key: "inscription_id", Value: inscriptionID},
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

func (r Repository) UpdateDexBTCBuyETHOrderStatus(model *entity.DexBTCBuyWithETH) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: "uuid", Value: model.UUID}}

	update := bson.M{
		"$set": bson.M{
			"status": model.Status,
		},
	}

	result, err := r.DB.Collection(model.TableName()).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (r Repository) UpdateDexBTCBuyETHOrderConfirmation(model *entity.DexBTCBuyWithETH) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: "uuid", Value: model.UUID}}

	update := bson.M{
		"$set": bson.M{
			"confirmation": model.Confirmation,
		},
	}

	result, err := r.DB.Collection(model.TableName()).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (r Repository) UpdateDexBTCBuyETHOrderSendMaster(model *entity.DexBTCBuyWithETH) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: "uuid", Value: model.UUID}}

	update := bson.M{
		"$set": bson.M{
			"status":     model.Status,
			"master_tx":  model.MasterTx,
			"updated_at": model.UpdatedAt,
		},
	}

	result, err := r.DB.Collection(model.TableName()).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (r Repository) UpdateDexBTCBuyETHOrderRefund(model *entity.DexBTCBuyWithETH) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: "uuid", Value: model.UUID}}

	update := bson.M{
		"$set": bson.M{
			"status":     model.Status,
			"refund_tx":  model.RefundTx,
			"updated_at": model.UpdatedAt,
		},
	}

	result, err := r.DB.Collection(model.TableName()).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (r Repository) UpdateDexBTCBuyETHOrderBuy(model *entity.DexBTCBuyWithETH) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: "uuid", Value: model.UUID}}

	update := bson.M{
		"$set": bson.M{
			"status":     model.Status,
			"buy_tx":     model.BuyTx,
			"updated_at": model.UpdatedAt,
			"split_tx":   model.SplitTx,
		},
	}

	result, err := r.DB.Collection(model.TableName()).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	return result, err
}

// func (r Repository) UpdateDexBTCBuyETHOrder(model *entity.DexBTCBuyWithETH) (*mongo.UpdateResult, error) {
// 	filter := bson.D{{Key: "uuid", Value: model.UUID}}

// 	result, err := r.DB.Collection(model.TableName()).UpdateOne(context.TODO(), filter, model)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return result, err
// }

// func (r Repository) UpdateDexBTCBuyETHOrderTx(model *entity.DexBTCBuyWithETH) (*mongo.UpdateResult, error) {
// 	filter := bson.D{{Key: "uuid", Value: model.UUID}}

// 	update := bson.M{
// 		"$set": bson.M{
// 			"eth_tx": model.ETHTx,
// 		},
// 	}

// 	result, err := r.DB.Collection(model.TableName()).UpdateOne(context.TODO(), filter, update)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return result, err
// }

func (r Repository) GetDexBTCBuyETHOrderByID(buyOrderID string) (*entity.DexBTCBuyWithETH, error) {
	f := bson.D{{Key: "uuid", Value: buyOrderID}}

	resp := &entity.DexBTCBuyWithETH{}
	usr, err := r.FilterOne(utils.COLLECTION_DEX_BTC_BUY_ETH, f)
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) GetNotCreatedActivitiesListing(page int64, limit int64) (*entity.Pagination, error) {
	confs := []entity.DexBTCListing{}
	resp := &entity.Pagination{}
	f := bson.M{"$or": bson.A{
		bson.M{"verified": true, "careated_verified_activity": bson.M{"$ne": true}},
		bson.M{"cancelled": true, "careated_cancelled_activity": bson.M{"$ne": true}},
		bson.M{"matched": true, "careated_matched_activity": bson.M{"$ne": true}},
	}}
	s := []Sort{{SortBy: "created_at", Sort: entity.SORT_ASC}}
	p, err := r.Paginate(entity.DexBTCListing{}.TableName(), page, limit, f, bson.D{}, s, &confs)
	if err != nil {
		return nil, err
	}

	resp.Result = confs
	resp.Page = p.Pagination.Page
	resp.Total = p.Pagination.Total
	resp.PageSize = limit
	return resp, nil
}

func (r Repository) UpdateListingCreatedVerifiedActivity(id string) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: "uuid", Value: id}}
	update := bson.M{
		"$set": bson.M{"created_verified_activity": true},
	}

	result, err := r.DB.Collection(entity.DexBTCListing{}.TableName()).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (r Repository) UpdateListingCreatedCancelledActivity(id string) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: "uuid", Value: id}}
	update := bson.M{
		"$set": bson.M{"created_cancelled_activity": true},
	}

	result, err := r.DB.Collection(entity.DexBTCListing{}.TableName()).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (r Repository) UpdateListingCreatedMatchedActivity(id string) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: "uuid", Value: id}}
	update := bson.M{
		"$set": bson.M{"created_matched_activity": true},
	}

	result, err := r.DB.Collection(entity.DexBTCListing{}.TableName()).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (r Repository) ListAllDexBTCBuyETH() ([]entity.DexBTCBuyWithETH, error) {
	listings := []entity.DexBTCBuyWithETH{}

	f := bson.M{
		"status": bson.M{"$in": []entity.DexBTCETHBuyStatus{entity.StatusDEXBuy_Expired, entity.StatusDEXBuy_Pending}},
	}

	cursor, err := r.DB.Collection(utils.COLLECTION_DEX_BTC_BUY_ETH).Find(context.TODO(), f, &options.FindOptions{})
	if err != nil {
		return nil, err
	}

	if err = cursor.All((context.TODO()), &listings); err != nil {
		return nil, err
	}

	return listings, nil
}
