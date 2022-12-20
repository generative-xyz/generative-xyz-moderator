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

type UpdateProfileRequest struct {
	DisplayName *string `json:"display_name"`
	Bio *string `json:"bio"`
}

