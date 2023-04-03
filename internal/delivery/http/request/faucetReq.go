package request

type FaucetReq struct {
	Url               string `json:"url"`
	RecaptchaResponse string `json:"g-recaptcha-response"`
}
