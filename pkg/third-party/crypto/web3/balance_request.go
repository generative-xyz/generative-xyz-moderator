package web3

type BalanceRequest struct {
	Address string `json:"address"`
	ChainID int    `json:"chainID"`
	TokenID int64  `json:"tokenID"`
}
