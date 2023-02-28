package request

type CreateInscribeBtcReq struct {
	WalletAddress string `json:"walletAddress"`
	Name          string `json:"name"`
	File          string `json:"file"`
	FeeRate       int32  `json:"fee_rate"`
	FileName      string `json:"fileName"`
	TokenAddress  string `json:"tokenAddress"`
	TokenId       string `json:"tokenId"`
}
