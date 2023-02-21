package request

type CreateMintReceiveAddressReq struct {
	WalletAddress string `json:"walletAddress"`
	ProjectID     string `json:"projectID"`
	PayType       string `json:"payType"`
}
