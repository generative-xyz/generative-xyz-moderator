package structure

type InscribeBtcReceiveAddrRespReq struct {
	WalletAddress string `json:"walletAddress"`
	Name          string `json:"name"`
	FileName      string `json:"fileName"`
	File          string `json:"file"`
	FeeRate       int32  `json:"fee_rate"`
}
