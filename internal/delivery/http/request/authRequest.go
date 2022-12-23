package request

type RefreshTokenData struct {
	RefreshToken  string `json:"refreshToken"`
	RedirectUri string `json:"redirectUri"`
}

type GenerateMessageRequest struct {
	Address string `json:"address"`
}

type VerifyMessageRequest struct {
	Sinature string `json:"signature"`
	Address string `json:"address"`
}

type UpdateProfileRequest struct {
	DisplayName *string `json:"displayName"`
	Bio *string `json:"bio"`
	Avatar *string `json:"avatar"`
}

