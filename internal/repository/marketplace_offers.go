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

func (r Repository) CancelOfferByOfferingID(offeringID string) error {
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

	obj.Closed = true
	_, err = r.UpdateOne(obj.TableName(), f, obj)
	
	return err
}

func (r Repository) filterOffers(filter entity.FilterMarketplaceOffers) bson.M {
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
	
	if filter.Closed != nil {
		f["closed"] = *filter.Closed
	}
	
	if filter.Finished != nil {
		f["finished"] = *filter.Finished
	}
	
	return f
}

func (r Repository) FilterMarketplaceOffers(filter entity.FilterMarketplaceOffers) (*entity.Pagination, error)  {
	confs := []entity.MarketplaceOffers{}
	resp := &entity.Pagination{}
	f := r.filterOffers(filter)

	p, err := r.Paginate(utils.COLLECTION_MARKETPLACE_OFFERS, filter.Page, filter.Limit, f, filter.SortBy, filter.Sort, &confs)
	if err != nil {
		return nil, err
	}
	
	resp.Result = confs
	resp.Page = p.Pagination.Page
	resp.Total = p.Pagination.Total
	return resp, nil
}
