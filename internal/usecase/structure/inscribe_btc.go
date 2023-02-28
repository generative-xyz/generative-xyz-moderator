package structure

type InscribeBtcReceiveAddrRespReq struct {
	WalletAddress string `json:"walletAddress"`
	Name          string `json:"name"`
	FileName      string `json:"fileName"`
	File          string `json:"file"`
	FeeRate       int32  `json:"fee_rate"`
	UserUuid      string `json:"-"`
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
