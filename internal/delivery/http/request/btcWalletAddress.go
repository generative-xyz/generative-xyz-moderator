package request

type CreateBtcWalletAddressReq struct {
	WalletAddress string `json:"walletAddress"`
	ProjectID string `json:"projectID"`
}

type CreateMintReq struct {
	Address string `json:"address"` //ord_walletaddress
	
}