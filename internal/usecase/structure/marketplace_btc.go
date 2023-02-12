package structure

type MarketplaceBTC_ListingInfo struct {
	InscriptionID  string `bson:"inscriptionID"` // tokenID in btc
	Name           string `bson:"name"`
	Description    string `bson:"description"`
	SellOrdAddress string `bson:"seller_ord_address"` //user's wallet address from FE
	SellerAddress  string `json:"seller_address"`
	Price          string `bson:"amount"`
	float64        int    `bson:"service_fee"`
	Min            int    `bson:"service_fee"`
}

type MarketplaceBTC_BuyOrderInfo struct {
	InscriptionID string `bson:"inscriptionID"`   // tokenID in btc
	OrderID       string `bson:"order_id"`        //
	BuyOrdAddress string `bson:"buy_ord_address"` //user's wallet address from FE
}
