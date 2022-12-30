package response

import "time"

type TokenRes struct{
	AccessToken string `json:"accessToken"`
	RefreshToken  string `json:"refreshToken"`
}

type GeneratedMessage struct {
	Message string `json:"message"`
}

type VerifyResponse struct {
	IsVerified bool `json:"isVerified"`
	Token string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}


type ProfileResponse struct {
	BaseResponse
	WalletAddress string `json:"walletAddress"`
	DisplayName string `json:"displayName"`
	Bio string `json:"bio"`
	Avatar string `json:"avatar"`
	CreatedAt *time.Time `json:"createdAt"`
	ProfileSocial ProfileSocial `json:"profileSocial"`
}

type ProfileSocial  struct{
    Web string `bson:"web"`;
    Twitter string `bson:"web"`;
    Discord string `bson:"discord"`;
    Medium string `bson:"medium"`;
	Instagram string `bson:"instagram"`;
}


type LogoutResponse struct {
	Message string `json:"message"`
}