package request

type CreateMarketplaceBTCBuyOrder struct {
	WalletAddress string `json:"walletAddress"`
	InscriptionID string `json:"inscriptionID"`
	OrderID       string `json:"orderID"`
	PayType       string `json:"payType"`
}

type CreateMarketplaceBTCListing struct {
	OrdWalletAddress string `json:"ordWalletAddress"`
	InscriptionID    string `json:"inscriptionID"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Price            string `json:"price"`

	PayType map[string]string `bson:"payType"`
}

type SendNFT struct {
	WalletName        string `json:"WalletName"`
	ReceiveOrdAddress string `json:"ReceiveOrdAddress"`
	InscriptionID     string `json:"InscriptionID"`
}

type ListingFee struct {
	InscriptionID string `json:"InscriptionID"`
}
