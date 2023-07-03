package request

type FaucetReq struct {
	Url               string `json:"url"`
	Address           string `json:"address"`
	RecaptchaResponse string `json:"g-recaptcha-response"`
	Type              string `json:"type"`
	Txhash            string `json:"txhash"`
	Source            string `json:"source"`
}

type FaucetAdminReq struct {
	Url         string             `json:"url"`
	Address     string             `json:"address"`
	Type        string             `json:"type"`
	Txhash      string             `json:"txhash"`
	Source      string             `json:"source"`
	ListAddress []string           `json:"listAddress"`
	MapAddress  map[string]float64 `json:"mapAddress"`
	Amount      float64            `json:"amount"`
}
