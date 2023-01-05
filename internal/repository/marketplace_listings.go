package repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

func (r Repository) CreateMarketplaceListing(listing *entity.MarketplaceListings) error {
	err := r.InsertOne(listing.TableName(), listing)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) FindListingByOfferingID(offeringID string) (*entity.MarketplaceListings, error) {
	resp := &entity.MarketplaceListings{}

	f := bson.D{
		{Key: "offering_id", Value: offeringID,},
	}

	listing, err := r.FilterOne(utils.COLLECTION_MARKETPLACE_LISTINGS, f)
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(listing, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) PurchaseTokenByOfferingID(offeringID string) error {
	obj := &entity.MarketplaceListings{}
	
	f := bson.D{
		{Key: "offering_id", Value: offeringID,},
	}

	listing, err := r.FilterOne(utils.COLLECTION_MARKETPLACE_LISTINGS, f)
	if err != nil {
		return err
	}

	err = helpers.Transform(listing, obj)
	if err != nil {
		return err
	}

	obj.Finished = true
	obj.Closed = true
	_, err = r.UpdateOne(obj.TableName(), f, obj)
	
	return err
}

func (r Repository) CancelListingByOfferingID(offeringID string) error {
	obj := &entity.MarketplaceListings{}
	
	f := bson.D{
		{Key: "offering_id", Value: offeringID,},
	}

	listing, err := r.FilterOne(utils.COLLECTION_MARKETPLACE_LISTINGS, f)
	if err != nil {
		return err
	}

	err = helpers.Transform(listing, obj)
	if err != nil {
		return err
	}

	obj.Closed = true
	_, err = r.UpdateOne(obj.TableName(), f, obj)
	
	return err
}

func (r Repository) filterListings(filter entity.FilterMarketplaceListings) bson.M {
	f := bson.M{}
	f[utils.KEY_DELETED_AT] = nil

	if filter.CollectionContract != nil {
		if *filter.CollectionContract != "" {
			f["collection_contract"] = *filter.CollectionContract
		}
	}
	
	if filter.TokenId != nil {
		if *filter.TokenId != "" {
			f["token_id"] = *filter.TokenId
		}
	}
	
	if filter.Erc20Token != nil {
		if *filter.Erc20Token != "" {
			f["erc_20_token"] = *filter.Erc20Token
		}
	}
	
	if filter.SellerAddress != nil {
		if *filter.SellerAddress != "" {
			f["seller"] = *filter.SellerAddress
		}
	}

	return f
}

func (r Repository) FilterMarketplaceListings(filter entity.FilterMarketplaceListings) (*entity.Pagination, error)  {
	confs := []entity.MarketplaceListings{}
	resp := &entity.Pagination{}
	f := r.filterListings(filter)

	p, err := r.Paginate(utils.COLLECTION_MARKETPLACE_LISTINGS, filter.Page, filter.Limit, f, filter.SortBy, filter.Sort, &confs)
	if err != nil {
		return nil, err
	}
	
	resp.Result = confs
	resp.Page = p.Pagination.Page
	resp.Total = p.Pagination.Total
	return resp, nil
}
