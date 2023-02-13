package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

func (r Repository) CreateMarketplaceListingBTC(listing *entity.MarketplaceBTCListing) error {
	err := r.InsertOne(listing.TableName(), listing)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) CreateMarketplaceBuyOrder(order *entity.MarketplaceBTCBuyOrder) error {
	err := r.InsertOne(order.TableName(), order)
	if err != nil {
		return err
	}
	return nil
}

// get item valid:
func (r Repository) FindBtcNFTListingByNFTID(inscriptionID string) (*entity.MarketplaceBTCListing, error) {
	resp := &entity.MarketplaceBTCListing{}

	f := bson.D{
		{Key: "inscriptionID", Value: inscriptionID},
	}

	listing, err := r.FilterOne(utils.COLLECTION_MARKETPLACE_BTC_LISTING, f, &options.FindOneOptions{
		Sort: bson.D{{Key: "created_at", Value: -1}},
	})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(listing, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// get item valid:
func (r Repository) FindBtcNFTListingLastSoldByNFTID(inscriptionID string) (*entity.MarketplaceBTCListing, error) {
	resp := &entity.MarketplaceBTCListing{}

	f := bson.D{
		{Key: "inscriptionID", Value: inscriptionID},
		{Key: "isConfirm", Value: true},
		{Key: "isSold", Value: true},
	}

	listing, err := r.FilterOne(utils.COLLECTION_MARKETPLACE_BTC_LISTING, f, &options.FindOneOptions{
		Sort: bson.D{{Key: "created_at", Value: -1}},
	})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(listing, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// get item valid + unsold:
func (r Repository) FindBtcNFTListingUnsoldByNFTID(inscriptionID string) (*entity.MarketplaceBTCListing, error) {
	resp := &entity.MarketplaceBTCListing{}

	f := bson.D{
		{Key: "inscriptionID", Value: inscriptionID},
		{Key: "isConfirm", Value: true},
		{Key: "isSold", Value: false},
	}

	listing, err := r.FilterOne(utils.COLLECTION_MARKETPLACE_BTC_LISTING, f)
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(listing, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) FindBtcNFTListingByOrderID(id string) (*entity.MarketplaceBTCListing, error) {
	resp := &entity.MarketplaceBTCListing{}

	f := bson.D{
		{Key: "uuid", Value: id},
	}

	listing, err := r.FilterOne(utils.COLLECTION_MARKETPLACE_BTC_LISTING, f)
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(listing, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) GetBTCListingHaveOngoingOrder(id string) ([]entity.MarketplaceBTCBuyOrder, error) {
	resp := []entity.MarketplaceBTCBuyOrder{}
	filter := bson.M{
		"item_id": id,
	}

	cursor, err := r.DB.Collection(utils.COLLECTION_MARKETPLACE_BTC_BUY).Find(context.TODO(), filter)
	if err != nil {
		if cursor.Err() == mongo.ErrNoDocuments || cursor.Err() == mongo.ErrNilDocument {
			return nil, nil
		}
		return nil, cursor.Err()
	}

	if err = cursor.All(context.TODO(), &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// get item valid to get info:
func (r Repository) FindBtcNFTListingByOrderIDValid(uuid string) (*entity.MarketplaceBTCListing, error) {
	resp := &entity.MarketplaceBTCListing{}

	f := bson.D{
		{Key: "uuid", Value: uuid},
		{Key: "isConfirm", Value: true},
	}

	listing, err := r.FilterOne(utils.COLLECTION_MARKETPLACE_BTC_LISTING, f)
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(listing, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) RetrieveBTCNFTPendingListings() ([]entity.MarketplaceBTCListing, error) {
	resp := []entity.MarketplaceBTCListing{}
	filter := bson.M{
		"isConfirm":  false,
		"expired_at": bson.M{"$gte": primitive.NewDateTimeFromTime(time.Now().UTC())},
	}

	cursor, err := r.DB.Collection(utils.COLLECTION_MARKETPLACE_BTC_LISTING).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (r Repository) UpdateBTCNFTConfirmListings(model *entity.MarketplaceBTCListing) (*mongo.UpdateResult, error) {

	filter := bson.D{{"uuid", model.UUID}}
	result, err := r.UpdateOne(model.TableName(), filter, model)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) retrieveBTCNFTListingsByFilter(filter bson.D, limit, offset int64) ([]entity.MarketplaceBTCListingFilterPipeline, error) {
	resp := []entity.MarketplaceBTCListingFilterPipeline{}

	cursor, err := r.DB.Collection(utils.COLLECTION_MARKETPLACE_BTC_LISTING).Aggregate(context.TODO(), bson.A{
		bson.D{
			{"$project",
				bson.D{
					{"uuid", 1},
					{"inscriptionID", 1},
					{"isConfirm", 1},
					{"isSold", 1},
					{"created_at", 1},
					{"expired_at", 1},
					{"name", 1},
					{"description", 1},
					{"seller_address", 1},
					{"seller_ord_address", 1},
					{"hold_ord_address", 1},
					{"amount", 1},
					{"service_fee", 1},
				},
			},
		},
		bson.D{{"$match", filter}},
		bson.D{{"$sort", bson.D{{"created_at", -1}}}},
		bson.D{
			{"$group",
				bson.D{
					{"_id", "$inscriptionID"},
					{"uuid", bson.D{{"$first", "$uuid"}}},
					{"inscriptionID", bson.D{{"$first", "$inscriptionID"}}},
					{"isConfirm", bson.D{{"$first", "$isConfirm"}}},
					{"isSold", bson.D{{"$first", "$isSold"}}},
					{"created_at", bson.D{{"$first", "$created_at"}}},
					{"expired_at", bson.D{{"$first", "$expired_at"}}},
					{"name", bson.D{{"$first", "$name"}}},
					{"description", bson.D{{"$first", "$description"}}},
					{"seller_address", bson.D{{"$first", "$seller_address"}}},
					{"seller_ord_address", bson.D{{"$first", "$seller_ord_address"}}},
					{"hold_ord_address", bson.D{{"$first", "$hold_ord_address"}}},
					{"amount", bson.D{{"$first", "$amount"}}},
					{"service_fee", bson.D{{"$first", "$service_fee"}}},
				},
			},
		},
		bson.D{{"$limit", limit}},
		bson.D{{"$skip", offset}},
	})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (r Repository) RetrieveBTCNFTListings(limit, offset int64) ([]entity.MarketplaceBTCListingFilterPipeline, error) {
	filter := bson.D{
		{Key: "isConfirm", Value: true},
	}

	resp, err := r.retrieveBTCNFTListingsByFilter(filter, limit, offset)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r Repository) RetrieveBTCNFTListingsUnsold(limit, offset int64) ([]entity.MarketplaceBTCListingFilterPipeline, error) {
	filter := bson.D{
		{Key: "isConfirm", Value: true},
		{Key: "isSold", Value: false},
	}

	resp, err := r.retrieveBTCNFTListingsByFilter(filter, limit, offset)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r Repository) RetrieveBTCNFTPendingBuyOrders() ([]entity.MarketplaceBTCBuyOrder, error) {
	resp := []entity.MarketplaceBTCBuyOrder{}
	filter := bson.M{
		"status":     entity.StatusBuy_Pending,
		"expired_at": bson.M{"$gte": primitive.NewDateTimeFromTime(time.Now().UTC())},
	}

	cursor, err := r.DB.Collection(utils.COLLECTION_MARKETPLACE_BTC_BUY).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (r Repository) RetrieveBTCNFTBuyOrdersByStatus(status entity.BuyStatus) ([]entity.MarketplaceBTCBuyOrder, error) {
	resp := []entity.MarketplaceBTCBuyOrder{}
	filter := bson.M{
		"status": status,
	}

	cursor, err := r.DB.Collection(utils.COLLECTION_MARKETPLACE_BTC_BUY).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (r Repository) UpdateBTCNFTBuyOrder(model *entity.MarketplaceBTCBuyOrder) (*mongo.UpdateResult, error) {

	filter := bson.D{{Key: "uuid", Value: model.UUID}}
	result, err := r.UpdateOne(model.TableName(), filter, model)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// func (r Repository) UpdateBTCNFTListingSoldStatus(id string) error {
// 	filter := bson.D{{Key: "uuid", Value: id}}
// 	update := bson.M{
// 		"$set": bson.M{
// 			"isSold": true,
// 		},
// 	}
// 	result, err := r.DB.Collection(utils.COLLECTION_MARKETPLACE_BTC_LISTING).UpdateOne(context.TODO(), filter, update)
// 	if err != nil {
// 		return err
// 	}
// 	if result.MatchedCount == 0 {
// 		return fmt.Errorf("not found listing %v", id)
// 	}
// 	return err
// }

func (r Repository) CreateMarketplaceBTCLog(listing *entity.MarketplaceBTCLogs) error {
	err := r.InsertOne(listing.TableName(), listing)
	if err != nil {
		return err
	}
	return nil
}

// func (r Repository) FindListingByOfferingID(offeringID string) (*entity.MarketplaceListings, error) {
// 	resp := &entity.MarketplaceListings{}

// 	f := bson.D{
// 		{Key: "offering_id", Value: offeringID},
// 	}

// 	listing, err := r.FilterOne(utils.COLLECTION_MARKETPLACE_LISTINGS, f)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = helpers.Transform(listing, resp)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return resp, nil
// }

// func (r Repository) PurchaseTokenByOfferingID(offeringID string) error {
// 	obj := &entity.MarketplaceListings{}

// 	f := bson.D{
// 		{Key: "offering_id", Value: offeringID},
// 	}

// 	listing, err := r.FilterOne(utils.COLLECTION_MARKETPLACE_LISTINGS, f)
// 	if err != nil {
// 		return err
// 	}

// 	err = helpers.Transform(listing, obj)
// 	if err != nil {
// 		return err
// 	}

// 	obj.Finished = true
// 	obj.Closed = true
// 	_, err = r.UpdateOne(obj.TableName(), f, obj)

// 	return err
// }

// func (r Repository) CancelListingByOfferingID(offeringID string) error {
// 	obj := &entity.MarketplaceListings{}

// 	f := bson.D{
// 		{Key: "offering_id", Value: offeringID},
// 	}

// 	listing, err := r.FilterOne(utils.COLLECTION_MARKETPLACE_LISTINGS, f)
// 	if err != nil {
// 		return err
// 	}

// 	err = helpers.Transform(listing, obj)
// 	if err != nil {
// 		return err
// 	}

// 	obj.Closed = true
// 	_, err = r.UpdateOne(obj.TableName(), f, obj)

// 	return err
// }

// func (r Repository) filterListings(filter entity.FilterMarketplaceListings) bson.M {
// 	f := bson.M{}
// 	f[utils.KEY_DELETED_AT] = nil

// 	if filter.CollectionContract != nil {
// 		if *filter.CollectionContract != "" {
// 			f["collection_contract"] = *filter.CollectionContract
// 		}
// 	}

// 	if filter.TokenId != nil {
// 		if *filter.TokenId != "" {
// 			f["token_id"] = *filter.TokenId
// 		}
// 	}

// 	if filter.Erc20Token != nil {
// 		if *filter.Erc20Token != "" {
// 			f["erc_20_token"] = *filter.Erc20Token
// 		}
// 	}

// 	if filter.SellerAddress != nil {
// 		if *filter.SellerAddress != "" {
// 			f["seller"] = *filter.SellerAddress
// 		}
// 	}

// 	if filter.Closed != nil {
// 		f["closed"] = *filter.Closed
// 	}

// 	if filter.Finished != nil {
// 		f["finished"] = *filter.Finished
// 	}

// 	return f
// }

// func (r Repository) FilterMarketplaceListings(filter entity.FilterMarketplaceListings) (*entity.Pagination, error) {
// 	confs := []entity.MarketplaceListings{}
// 	resp := &entity.Pagination{}
// 	f := r.filterListings(filter)

// 	p, err := r.Paginate(utils.COLLECTION_MARKETPLACE_LISTINGS, filter.Page, filter.Limit, f, bson.D{}, []Sort{}, &confs)
// 	if err != nil {
// 		return nil, err
// 	}

// 	resp.Result = confs
// 	resp.Page = p.Pagination.Page
// 	resp.Total = p.Pagination.Total
// 	resp.PageSize = filter.Limit
// 	return resp, nil
// }

// func (r Repository) GetListingBySeller(sellerAddress string) ([]entity.MarketplaceListings, error) {
// 	resp := []entity.MarketplaceListings{}
// 	filter := entity.FilterMarketplaceListings{
// 		SellerAddress: &sellerAddress,
// 	}

// 	f := r.filterListings(filter)

// 	cursor, err := r.DB.Collection(utils.COLLECTION_MARKETPLACE_LISTINGS).Find(context.TODO(), f)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err = cursor.All(context.TODO(), &resp); err != nil {
// 		return nil, err
// 	}

// 	return resp, nil
// }

// func (r Repository) GetAllListingByCollectionContract(contract string) ([]entity.MarketplaceListings, error) {
// 	listings := []entity.MarketplaceListings{}
// 	f := bson.D{{
// 		Key:   utils.KEY_LISTING_CONTRACT,
// 		Value: contract,
// 	}}

// 	cursor, err := r.DB.Collection(utils.COLLECTION_MARKETPLACE_LISTINGS).Find(context.TODO(), f)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err = cursor.All((context.TODO()), &listings); err != nil {
// 		return nil, err
// 	}

// 	return listings, nil
// }

// func (r Repository) GetAllListings() ([]entity.MarketplaceListings, error) {
// 	listings := []entity.MarketplaceListings{}
// 	f := bson.D{}

// 	cursor, err := r.DB.Collection(utils.COLLECTION_MARKETPLACE_LISTINGS).Find(context.TODO(), f)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err = cursor.All((context.TODO()), &listings); err != nil {
// 		return nil, err
// 	}

// 	return listings, nil
// }

// func (r Repository) GetAllActiveListings() ([]entity.MarketplaceListings, error) {
// 	listings := []entity.MarketplaceListings{}
// 	f := bson.D{{
// 		Key:   "closed",
// 		Value: false,
// 	}}

// 	cursor, err := r.DB.Collection(utils.COLLECTION_MARKETPLACE_LISTINGS).Find(context.TODO(), f)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err = cursor.All((context.TODO()), &listings); err != nil {
// 		return nil, err
// 	}

// 	return listings, nil
// }

// func (r Repository) UpdateListingOwnerAddress(offeringID string, ownerAddress string) (*mongo.UpdateResult, error) {
// 	f := bson.D{
// 		{Key: "offering_id", Value: offeringID},
// 	}

// 	update := bson.M{
// 		"$set": bson.M{
// 			"owner_address": ownerAddress,
// 		},
// 	}

// 	result, err := r.DB.Collection(utils.COLLECTION_MARKETPLACE_LISTINGS).UpdateOne(context.TODO(), f, update)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return result, err
// }
