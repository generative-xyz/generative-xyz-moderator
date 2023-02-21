package request

type UpsertRedisRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ListNftIdsReq struct {
	InscriptionID []string `json:"inscriptionID"`

	SellOrdAddress string `json:"seller_ord_address"`
	SellerAddress  string `json:"seller_address"`
	Price          string `json:"amount"`
}
