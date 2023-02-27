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

type StatusListing int

const (
	StatusListing_Pending        StatusListing = iota // 0: pending: waiting for fund.
	StatusListing_ReceivedNft                         // 1: received nft from seller // isConfirm true.
	StatusListing_RequestdCancel                      // 3: user requested a cancel.
	StatusListing_Canceling                           // 4: submit send nft back to the seller.
	StatusListing_Canceled                            // 5: send nft success.
	StatusListing_TxCancelFailed                      // 6 tx cancel failed, need to retry.
	StatusListing_Invalid                             // 7 the listing invalid // not show on the mkp.
)

type MarketplaceBTCListing struct {
	BaseEntity `bson:",inline"`

	HoldOrdAddress string `bson:"hold_ord_address"` // address temp to receive nft from user.

	SellOrdAddress string `bson:"seller_ord_address"` // the address to refund nft, it is same the address in PayType if pay btc.

	SellerAddress string `bson:"seller_address"` // old flow, it's segwit address, new flow dont need

	ServiceFee string `bson:"service_fee"`

	Status StatusListing `bson:"status"` // TODO: need to migrate old data.

	Price string `bson:"amount"` // amount by btc

	PayType map[string]string `bson:"pay_type"` // {"eth": <seller_receive_payment_address>, "btc": <seller_receive_payment_address>, ...}

	IsConfirm bool `bson:"isConfirm"`
	IsSold    bool `bson:"isSold"`
	IsCancel  bool `bson:"isCancel"`

	TxNFT         string    `bson:"tx_nft"`
	TxCancel      string    `bson:"tx_cancel"`
	InscriptionID string    `bson:"inscriptionID"` // tokenID in btc
	Name          string    `bson:"name"`
	Description   string    `bson:"description"`
	ExpiredAt     time.Time `bson:"expired_at"`

	// for filter
	CollectionID     string    `bson:"collection_id"`
	CollectionName   string    `bson:"collection_name"`
	InscriptionName  string    `bson:"inscription_name"`
	InscriptionIndex string    `bson:"inscription_index"`
	Inscription      *TokenUri `bson:"inscription"`
}

type MarketplaceBTCListingFilterPipeline struct {
	ID   string `bson:"_id"`
	UUID string `bson:"uuid"`

	SellOrdAddress string `bson:"seller_ord_address"` // refund nft when cancel
	HoldOrdAddress string `bson:"hold_ord_address"`   // address temp that user send nft

	SellerAddress string `bson:"seller_address"` // address to receive btc

	Price         string    `bson:"amount"`
	ServiceFee    string    `bson:"service_fee"`
	IsConfirm     bool      `bson:"isConfirm"`
	IsSold        bool      `bson:"isSold"`
	TxNFT         string    `bson:"tx_nft"`
	InscriptionID string    `bson:"inscriptionID"` // tokenID in btc
	Name          string    `bson:"name"`
	Description   string    `bson:"description"`
	ExpiredAt     time.Time `bson:"expired_at"`
	CreatedAt     time.Time `bson:"created_at"`

	//listing payment info:
	PayType map[string]string `bson:"pay_type"`

	// for filter
	CollectionID     string    `bson:"collection_id"`
	CollectionName   string    `bson:"collection_name"`
	InscriptionName  string    `bson:"inscription_name"`
	InscriptionIndex string    `bson:"inscription_index"`
	Inscription      *TokenUri `bson:"inscription"`
}

type MarketplaceBTCListingFloorPrice struct {
	ID    string `bson:"_id"`
	Price uint64 `bson:"price"`
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
