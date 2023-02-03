package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type FilterMarketplaceOffers struct {
	BaseFilters
	CollectionContract *string
	TokenId *string
	Erc20Token *string
	BuyerAddress *string
	Closed             *bool  
	Finished           *bool  
	OwnerAddress       *string
}

type MarketplaceOffers struct {
	BaseEntity         `bson:",inline"`
	ID                 string `bson:"id"`
	OfferingId         string `bson:"offering_id"`
	CollectionContract string `bson:"collection_contract"`
	TokenId            string `bson:"token_id"`
	Buyer              string `bson:"buyer"`
	Erc20Token         string `bson:"erc_20_token"`
	Price              string `bson:"price"`
	Closed             bool   `bson:"closed"`
	Finished           bool   `bson:"finished"`
	DurationTime       string `bson:"duration_time"`
	BlockNumber   		 uint64 `bson:"block_number"`
	OwnerAddress 			 *string `bson:"owner_address"`
	Token TokenUri `bson:"-"`
	BuyerInfo Users `bson:"-"`
}

func (u MarketplaceOffers) TableName() string { 
	return utils.COLLECTION_MARKETPLACE_OFFERS
}

func (u MarketplaceOffers) ToBson()  (*bson.D, error) { 
	return helpers.ToDoc(u)
}
