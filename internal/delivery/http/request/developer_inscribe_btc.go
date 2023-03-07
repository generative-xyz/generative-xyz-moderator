package request

type DeveloperCreateInscribeBtcReq struct {
	WalletAddress string `json:"walletAddress"`
	Name          string `json:"name"`
	File          string `json:"file"`
	FeeRate       int32  `json:"feeRate"`
	FileName      string `json:"fileName"`
	TokenAddress  string `json:"tokenAddress"`
	TokenId       string `json:"tokenId"`
}
