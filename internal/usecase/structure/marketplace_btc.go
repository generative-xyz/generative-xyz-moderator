package structure

import "time"

type MarketplaceBTC_ListingInfo struct {
	InscriptionID  string `json:"inscriptionID"` // tokenID in btc
	Name           string `json:"name"`
	Description    string `json:"description"`
	SellOrdAddress string `json:"seller_ord_address"` //user's wallet address from FE
	SellerAddress  string `json:"seller_address"`
	Price          string `json:"amount"`
	ServiceFee     string `json:"service_fee"`
	Min            string `json:"min"`
}

type MarketplaceBTC_BuyOrderInfo struct {
	InscriptionID string `json:"inscriptionID"`   // tokenID in btc
	OrderID       string `json:"order_id"`        //
	BuyOrdAddress string `json:"buy_ord_address"` //user's wallet address from FE
}

type MarketplaceNFTDetail struct {
	InscriptionID string    `json:"inscriptionID"`
	Price         string    `json:"price"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	OrderID       string    `json:"orderID"`
	IsConfirmed   bool      `json:"isConfirmed"`
	Buyable       bool      `json:"buyable"`
	IsCompleted   bool      `json:"isCompleted"`
	CreatedAt     time.Time `json:"createdAt"`
	CollectionID  string    `json:"collectionID"`
}

type MarketplaceCollectionDetail struct {
	CollectionID  string `json:"collectionID"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Total         string `json:"total"`
	InscriptionID string `json:"inscriptionID"`
	Creator       []struct {
		UUID     string `json:"uuid"`
		Name     string `json:"display_name"`
		Avatar   string `json:"avatar"`
		Verified bool   `json:"is_verified"`
	} `json:"creator"`
	CreatedAt time.Time `json:"createdAt"`
}
