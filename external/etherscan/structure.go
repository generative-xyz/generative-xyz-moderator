package etherscan

type WalletAddressResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Result  string `json:"result"`
}

type WalletAddressTxResponse struct {
	Message string                   `json:"message"`
	Status  string                   `json:"status"`
	Result  []*AddressTxItemResponse `json:"result"`
}

type AddressTxItemResponse struct {
	ENS             string  `json:"ens"`
	Avatar          string  `json:"avatar"`
	From            string  `json:"from"`
	To              string  `json:"to"`
	Value           string  `json:"value"`
	UsdtValue       float64 `json:"usdt_value"`
	UsdtValueExtra  float64 `json:"usdt_value_extra"`
	ExtraPercent    float64 `json:"extra_percent"`
	Percent         float64 `json:"percent"`
	GMReceive       float64 `json:"gm_receive"`
	GMReceiveString string  `json:"gm_receive_string"`
	Currency        string  `json:"currency"`
}
