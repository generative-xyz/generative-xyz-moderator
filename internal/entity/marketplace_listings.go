package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type FilterMarketplaceListings struct {
	BaseFilters
	CollectionContract *string
	TokenId *string
	Erc20Token *string
	SellerAddress *string
	Closed             *bool  
	Finished           *bool  
}

type MarketplaceListings struct {
	BaseEntity         `bson:",inline"`
	ID                 string `bson:"id"`
	OfferingId         string `bson:"offering_id"`
	CollectionContract string `bson:"collection_contract"`
	TokenId            string `bson:"token_id"`
	Seller             string `bson:"seller"`
	Erc20Token         string `bson:"erc_20_token"`
	Price              string `bson:"price"`
	Closed             bool   `bson:"closed"`
	Finished           bool   `bson:"finished"`
	DurationTime       string `bson:"duration_time"`
	Token TokenUri `bson:"-"`
}

func (u MarketplaceListings) TableName() string { 
	return utils.COLLECTION_MARKETPLACE_LISTINGS
}

func (u MarketplaceListings) ToBson()  (*bson.D, error) { 
	return helpers.ToDoc(u)
}
