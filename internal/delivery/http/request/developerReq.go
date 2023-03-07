package request

type GetApiKeyReq struct {
	RecaptchaResponse string `json:"g-recaptcha-response"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	Email             string `json:"email"`
	Company           string `json:"company"`
}
