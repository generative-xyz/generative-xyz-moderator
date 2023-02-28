package structure

import (
	"time"

	"rederinghub.io/internal/entity"
)

type MarketplaceBTC_ListingInfo struct {
	InscriptionID string `json:"inscriptionID"` // tokenID in btc
	Name          string `json:"name"`
	Description   string `json:"description"`

	Price string `json:"price"`

	SellOrdAddress string `json:"sellOrdAddress"`

	PayType map[string]string `bson:"payType"`

	ServiceFee string `json:"serviceFee"`
}

type MarketplaceBTC_BuyOrderInfo struct {
	InscriptionID string `json:"inscriptionID"`   // tokenID in btc
	OrderID       string `json:"order_id"`        //
	BuyOrdAddress string `json:"buy_ord_address"` //user's wallet address from FE
	PayType       string `json:"pay_type"`
}

type PaymentInfoForBuyOrder struct {
	PaymentAddress string `json:"paymentAddress"`
	Price          string `json:"price"`
	// PriceFloat     float64 `json:"priceFloat"`
}

type MarketplaceNFTDetail struct {
	InscriptionID     string            `json:"inscriptionID"`
	Price             string            `json:"price"`
	Name              string            `json:"name"`
	Description       string            `json:"description"`
	OrderID           string            `json:"orderID"`
	IsConfirmed       bool              `json:"isConfirmed"`
	Buyable           bool              `json:"buyable"`
	PayType           map[string]string `bson:"pay_type"`
	IsCompleted       bool              `json:"isCompleted"`
	CreatedAt         time.Time         `json:"createdAt"`
	InscriptionNumber string            `json:"inscriptionNumber"`
	ContentType       string            `json:"contentType"`
	ContentLength     string            `json:"contentLength"`

	// for filter
	CollectionID     string           `json:"collection_id"`
	CollectionName   string           `json:"collection_name"`
	InscriptionName  string           `json:"inscription_name"`
	InscriptionIndex string           `json:"inscription_index"`
	Inscription      *entity.TokenUri `json:"inscription"`

	// paytype:
	PaymentListingInfo map[string]PaymentInfoForBuyOrder `json:"paymentListingInfo"`
}

type MarketplaceCollectionStats struct {
	FloorPrice uint64 `json:"floor_price"`
}
