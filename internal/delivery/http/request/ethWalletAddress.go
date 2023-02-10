package request

type CreateEthWalletAddressReq struct {
	WalletAddress string `json:"walletAddress"`
	ProjectID     string `json:"projectID"`
}
