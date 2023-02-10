package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

// type FilterMarketplaceListings struct {
// 	BaseFilters
// 	CollectionContract *string
// 	TokenId            *string
// 	Erc20Token         *string
// 	SellerAddress      *string
// 	Closed             *bool
// 	Finished           *bool
// }

type MarketplaceBTCBuyOrder struct {
	BaseEntity    `bson:",inline"`
	OrdAddress    string `bson:"ord_address"`
	ItemID        string `bson:"item_id"`
	InscriptionID string `bson:"inscriptionID"` // tokenID in btc
}

type MarketplaceBTCListing struct {
	BaseEntity     `bson:",inline"`
	SellOrdAddress string `bson:"seller_ord_address"` //user's wallet address from FE
	HoldOrdAddress string `bson:"hold_ord_address"`
	Price          string `bson:"amount"`
	ServiceFee     string `bson:"service_fee"`
	IsConfirm      bool   `bson:"isConfirm"`
	IsSold         bool   `bson:"isSold"`
	InscriptionID  string `bson:"inscriptionID"` // tokenID in btc
	Name           string `bson:"name"`
	Description    string `bson:"description"`
}

func (u MarketplaceBTCListing) TableName() string {
	return utils.COLLECTION_MARKETPLACE_BTC_LISTING
}

func (u MarketplaceBTCListing) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
