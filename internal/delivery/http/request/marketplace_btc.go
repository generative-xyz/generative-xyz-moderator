package request

type CreateMarketplaceBTCBuyOrder struct {
	WalletAddress string `json:"walletAddress"`
	InscriptionID string `json:"inscriptionID"`
	OrderID       string `json:"orderID"`
}

type CreateMarketplaceBTCListing struct {
	ReceiveAddress    string `json:"receiveAddress"`
	ReceiveOrdAddress string `json:"receiveOrdAddress"`
	InscriptionID     string `json:"inscriptionID"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	Price             string `json:"price"`
}

type SendNFT struct {
	WalletName        string `json:"WalletName"`
	ReceiveOrdAddress string `json:"ReceiveOrdAddress"`
	InscriptionID     string `json:"InscriptionID"`
}
