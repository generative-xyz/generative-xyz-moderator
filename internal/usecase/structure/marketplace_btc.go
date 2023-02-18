package structure

import (
	"time"

	"rederinghub.io/internal/entity"
)

type MarketplaceBTC_ListingInfo struct {
	InscriptionID  string `json:"inscriptionID"` // tokenID in btc
	Name           string `json:"name"`
	Description    string `json:"description"`
	SellOrdAddress string `json:"seller_ord_address"` //user's wallet address from FE
	SellerAddress  string `json:"seller_address"`
	Price          string `json:"amount"`
	ServiceFee     string `json:"service_fee"`
}

type MarketplaceBTC_BuyOrderInfo struct {
	InscriptionID string `json:"inscriptionID"`   // tokenID in btc
	OrderID       string `json:"order_id"`        //
	BuyOrdAddress string `json:"buy_ord_address"` //user's wallet address from FE
}

type MarketplaceNFTDetail struct {
	InscriptionID     string    `json:"inscriptionID"`
	Price             string    `json:"price"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	OrderID           string    `json:"orderID"`
	IsConfirmed       bool      `json:"isConfirmed"`
	Buyable           bool      `json:"buyable"`
	IsCompleted       bool      `json:"isCompleted"`
	CreatedAt         time.Time `json:"createdAt"`
	InscriptionNumber string    `json:"inscriptionNumber"`
	ContentType       string    `json:"contentType"`
	ContentLength     string    `json:"contentLength"`

	// for filter
	CollectionID    string           `json:"collection_id"`
	CollectionName  string           `json:"collection_name"`
	InscriptionName string           `json:"inscription_name"`
	Inscription     *entity.TokenUri `json:"inscription"`
}
