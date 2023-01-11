package structure

import (
	"time"

	"rederinghub.io/internal/entity"
)


type GenerateMessage struct {
	Address string
}

type VerifyMessage struct {
	Signature string
	Address string
}

type VerifyResponse struct {
	IsVerified bool
	Token string
	RefreshToken string
}

type ProfileResponse struct {
	ID string `json:"id"`
	WalletAddress string `json:"wallet_address"`
	DisplayName string `json:"display_name"`
	Bio string `json:"bio"`
	Avatar string `json:"avatar"`
	CreatedAt *time.Time `json:"created_at"`
}

type UpdateProfile struct {
	DisplayName *string 
	Bio *string 
	ProfileSocial ProfileSocial
	Avatar *string
}

type ProfileSocial  struct{
    Web *string 
    Twitter *string
    Discord *string 
    Medium *string 
	Instagram *string 
}

type ProfileChan struct {
	Data *entity.Users
	Err error
}