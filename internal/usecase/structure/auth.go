package structure

import "time"


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
	FullName string `json:"full_name"`
	FirstName string `json:"-"`
	LastName string `json:"-"`
	Email string `json:"-"`
	VerifiedAt *time.Time `json:"verified_at"`
	IsVerified bool `json:"is_verified"`
	NickName string `json:"nickname"`
	Bio string `json:"bio"`
	Avatar string `json:"avatar"`
	CoverPhoto string `json:"cover_photo"`
	TotalItems int `json:"total_items"`
	EstimateValue float64 `json:"estimate_value"`
	TotalOwned int `json:"total_owned"`
	TotalCreated int `json:"total_created"`
	TotalForging int `json:"total_forging"`
	LinkOpensea string `json:"link_opensea"`
	LinkSocial string `json:"link_social"`
	Address string `json:"address"`
	Apartment string `json:"apartment"`
	City string `json:"city"`
	State string `json:"state"`
	ZipCode string `json:"zip_code"`
	CreatedAt *time.Time `json:"created_at"`
}