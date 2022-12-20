package response

import "time"

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


type ProfileResponse struct {
	ID string `json:"id"`
	WalletAddress string `json:"wallet_address"`
	DisplayName string `json:"display_name"`
	Bio string `json:"bio"`
	Avatar string `json:"avatar"`
	CreatedAt *time.Time `json:"created_at"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}