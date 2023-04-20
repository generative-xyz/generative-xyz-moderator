package response

type FaucetStatusRes struct {
	Txhash string `json:"txhash"`
	Status string `json:"status"`
}
