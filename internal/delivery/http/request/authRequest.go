package request

type RefreshTokenData struct {
	RefreshToken  string `json:"refresh_token"`
	RedirectUri string `json:"redirect_uri"`
}

type GenerateMessageRequest struct {
	Address string `json:"address"`
}

type VerifyMessageRequest struct {
	Sinature string `json:"signature"`
	Address string `json:"address"`
}
