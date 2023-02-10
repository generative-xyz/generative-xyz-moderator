package response

type MarketplaceNFTDetail struct {
	InscriptionID string `json:"inscriptionID"`
	Price         string `json:"price"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	OrderID       string `json:"orderID"`
}

type CreateMarketplaceBTCBuyOrder struct {
	ReceiveAddress string `json:"receiveAddress"`
}

type CreateMarketplaceBTCListing struct {
	ReceiveAddress string `json:"receiveAddress"`
}
