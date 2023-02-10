package request

type CreateMarketplaceBTCBuyOrder struct {
	WalletAddress string `json:"walletAddress`
	InscriptionID string `json:"inscriptionID`
	OrderID       string `json:"orderID"`
}

type CreateMarketplaceBTCListing struct {
	ReceiveAddress string `json:"receiveAddress`
	InscriptionID  string `json:"inscriptionID`
	Name           string `json:"name`
	Description    string `json:"description`
	Price          string `json:"price`
}
