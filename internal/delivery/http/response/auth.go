package response

type TokenRes struct{
	AccessToken string `json:"access_token"`
	RefreshToken  string `json:"refresh_token"`
}

type GeneratedMessage struct {
	Message string `json:"message"`
}



type VerifyResponse struct {
	IsVerified bool `json:"is_verified"`
	Token string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}