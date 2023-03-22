package structure

type CreateProjectAllowListReq struct {
	ProjectID        *string 
	UserWalletAddress *string 
}

type Erc20Config struct{
	Value int64 `json:"value"`
	Decimal int64 `json:"decimal"`
}