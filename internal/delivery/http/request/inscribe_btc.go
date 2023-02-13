package request

type CreateInscribeBtcReq struct {
	WalletAddress string `json:"walletAddress"`
	Name          string `json:"name"`
	File          string `json:"file"`
	FeeRate       int32  `json:"fee_rate"`
}
