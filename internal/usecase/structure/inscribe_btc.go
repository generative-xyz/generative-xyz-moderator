package structure

type InscribeBtcReceiveAddrRespReq struct {
	WalletAddress string `json:"walletAddress"`
	Name          string `json:"name"`
	FileName      string `json:"fileName"`
	File          string `json:"file"`
	FeeRate       int32  `json:"fee_rate"`
	UserUuid      string `json:"-"`
	TokenAddress  string `json:"tokenAddress"`
	TokenId       string `json:"tokenId"`
}

func (s InscribeBtcReceiveAddrRespReq) NeedVerifyAuthentic() bool {
	return s.TokenAddress != "" && s.TokenId != ""
}

func (s *InscribeBtcReceiveAddrRespReq) SetFields(fns ...func(*InscribeBtcReceiveAddrRespReq)) {
	for _, fn := range fns {
		fn(s)
	}
}
func (s InscribeBtcReceiveAddrRespReq) WithUserUuid(userUuid string) func(*InscribeBtcReceiveAddrRespReq) {
	return func(ibrarr *InscribeBtcReceiveAddrRespReq) {
		ibrarr.UserUuid = userUuid
	}
}
