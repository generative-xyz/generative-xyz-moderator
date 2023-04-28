package request

type FaucetReq struct {
	Url               string `json:"url"`
	Address           string `json:"address"`
	RecaptchaResponse string `json:"g-recaptcha-response"`
	Type              string `json:"type"`
	Txhash            string `json:"txhash"`
	Source            string `json:"source"`
}
