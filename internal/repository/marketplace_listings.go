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
