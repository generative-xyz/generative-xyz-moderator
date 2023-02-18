package entity

import (
	"time"

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

type MarketplaceBTCListing struct {
	BaseEntity     `bson:",inline"`
	SellOrdAddress string    `bson:"seller_ord_address"`
	SellerAddress  string    `bson:"seller_address"`
	HoldOrdAddress string    `bson:"hold_ord_address"`
	Price          string    `bson:"amount"`
	ServiceFee     string    `bson:"service_fee"`
	IsConfirm      bool      `bson:"isConfirm"`
	IsSold         bool      `bson:"isSold"`
	TxNFT          string    `bson:"tx_nft"`
	InscriptionID  string    `bson:"inscriptionID"` // tokenID in btc
	Name           string    `bson:"name"`
	Description    string    `bson:"description"`
	ExpiredAt      time.Time `bson:"expired_at"`

	// for filter
	CollectionID     string    `bson:"collection_id"`
	CollectionName   string    `bson:"collection_name"`
	InscriptionName  string    `bson:"inscription_name"`
	InscriptionIndex string    `bson:"inscription_index"`
	Inscription      *TokenUri `bson:"inscription"`
}

type MarketplaceBTCListingFilterPipeline struct {
	ID             string    `bson:"_id"`
	UUID           string    `bson:"uuid"`
	SellOrdAddress string    `bson:"seller_ord_address"`
	SellerAddress  string    `bson:"seller_address"`
	HoldOrdAddress string    `bson:"hold_ord_address"`
	Price          string    `bson:"amount"`
	ServiceFee     string    `bson:"service_fee"`
	IsConfirm      bool      `bson:"isConfirm"`
	IsSold         bool      `bson:"isSold"`
	TxNFT          string    `bson:"tx_nft"`
	InscriptionID  string    `bson:"inscriptionID"` // tokenID in btc
	Name           string    `bson:"name"`
	Description    string    `bson:"description"`
	ExpiredAt      time.Time `bson:"expired_at"`
	CreatedAt      time.Time `bson:"created_at"`

	// for filter
	CollectionID     string    `bson:"collection_id"`
	CollectionName   string    `bson:"collection_name"`
	InscriptionName  string    `bson:"inscription_name"`
	InscriptionIndex string    `bson:"inscription_index"`
	Inscription      *TokenUri `bson:"inscription"`
}

func (u MarketplaceBTCListing) TableName() string {
	return utils.COLLECTION_MARKETPLACE_BTC_LISTING
}

func (u MarketplaceBTCListing) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}

type MarketplaceBTCLogs struct {
	BaseEntity  `bson:",inline"`
	RecordID    string      `bson:"record_id"`
	Table       string      `bson:"table"`
	Name        string      `bson:"name"`
	Status      interface{} `bson:"status"`
	RequestMsg  interface{} `bson:"request_msg"`
	ResponseMsg interface{} `bson:"response_msg"`
}

func (u MarketplaceBTCLogs) TableName() string {
	return utils.COLLECTION_MARKETPLACE_BTC_LOGS
}

func (u MarketplaceBTCLogs) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
