package request

type CreateBtcWalletAddressReq struct {
	WalletAddress string `json:"walletAddress"`
	ProjectID string `json:"projectID"`
}

type CheckBalanceAddressReq struct {
	Address string `json:"address"`
}

type CreateMintReq struct {
	Address string `json:"address"` //ord_walletaddress
	
}

type CreateBtcWalletAddressReqV2 struct {
	WalletAddress string `json:"walletAddress"`
	Name string `json:"name"`
	File string `json:"file"`
}
