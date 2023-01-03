package repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

func (r Repository) CreateMarketplaceOffer(offer *entity.MarketplaceOffers) error {
	err := r.InsertOne(offer.TableName(), offer)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) FindOfferByOfferingID(offeringID string) (*entity.MarketplaceOffers, error) {
	resp := &entity.MarketplaceOffers{}

	f := bson.D{
		{Key: "offering_id", Value: offeringID,},
	}

	Offer, err := r.FilterOne(utils.COLLECTION_MARKETPLACE_OFFERS, f)
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(Offer, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) AcceptOfferByOfferingID(offeringID string) error {
	obj := &entity.MarketplaceOffers{}
	
	f := bson.D{
		{Key: "offering_id", Value: offeringID,},
	}

	Offer, err := r.FilterOne(utils.COLLECTION_MARKETPLACE_OFFERS, f)
	if err != nil {
		return err
	}

	err = helpers.Transform(Offer, obj)
	if err != nil {
		return err
	}

	obj.Finished = true
	obj.Closed = true
	_, err = r.UpdateOne(obj.TableName(), f, obj)
	
	return err
}
