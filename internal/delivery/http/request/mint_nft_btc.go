package request

type CreateMintReceiveAddressReq struct {
	WalletAddress     string `json:"walletAddress"`
	ProjectID         string `json:"projectID"`
	PayType           string `json:"payType"`
	RefundUserAddress string `json:"refundUserAddress"`
	Quantity          int    `json:"quantity"`
	FeeRate           int32  `json:"feeRate"`

	IsCustomFeeRate bool        `json:"isCustomFeeRate"`
	EstMintFeeInfo  interface{} `json:"estMintFeeInfo"`
}
