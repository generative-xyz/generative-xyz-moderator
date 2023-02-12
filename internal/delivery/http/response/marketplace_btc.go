package response

type MarketplaceNFTDetail struct {
	InscriptionID string `json:"inscriptionID"`
	Price         string `json:"price"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	OrderID       string `json:"orderID"`
	IsConfirmed   bool   `json:"isConfirmed"`
}

type CreateMarketplaceBTCBuyOrder struct {
	ReceiveAddress string `json:"receiveAddress"`
	TimeoutAt      string `json:"timeoutAt"`
}

type CreateMarketplaceBTCListing struct {
	ReceiveAddress string `json:"receiveAddress"`
	TimeoutAt      string `json:"timeoutAt"`
}
