package response

type FaucetStatusRes struct {
	Txhash    string `json:"txhash"`
	Status    string `json:"status"`
	CreatedAt int64  `json:"created_at"`
}

type FaucetConfigRes struct {
	FaucetAmounts map[string]string `json:"faucet_amounts"`
}
