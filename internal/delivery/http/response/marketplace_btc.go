package response

type MarketplaceNFTDetail struct {
	InscriptionID string `json:"inscriptionID"`
	Price         string `json:"price"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	OrderID       string `json:"orderID"`
	IsConfirmed   bool   `json:"isConfirmed"`
	Buyable       bool   `json:"buyable"`
	IsCompleted   bool   `json:"isCompleted"` //last order is completed
	// LastPrice     int64  `json:"lastPrice"`
}

type CreateMarketplaceBTCBuyOrder struct {
	ReceiveAddress string `json:"receiveAddress"`
	TimeoutAt      string `json:"timeoutAt"`
}

type CreateMarketplaceBTCListing struct {
	ReceiveAddress string `json:"receiveAddress"`
	TimeoutAt      string `json:"timeoutAt"`
}
