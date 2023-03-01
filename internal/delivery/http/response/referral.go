package response

type ReferralResp struct {
	BaseResponse
	ReferreeID string          `json:"referreeID"`
	ReferrerID string          `json:"referrerID"`
	Referree   ProfileResponse `json:"referree"`
	Referrer   ProfileResponse `json:"referrer"`
	ReferreeVolumn   ReferralVolumnResp	 `json:"referreeVolumn"`
}

type ReferralVolumnResp struct {
	Amount string `json:"amount"`
	AmountType string `json:"amountType"`
	ProjectID string  `json:"projectID"`
	Percent int `json:"percent"`
	Earn string  `json:"earn"`
	GenEarn string `json:"GenEarn"`
}
