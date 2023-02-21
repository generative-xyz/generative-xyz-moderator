package response

import "time"

type TokenRes struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type GeneratedMessage struct {
	Message string `json:"message"`
}

type VerifyResponse struct {
	IsVerified   bool   `json:"isVerified"`
	Token        string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type ProfileResponse struct {
	BaseResponse
	WalletAddress           string        `json:"walletAddress"`
	WalletAddressBTC        string        `json:"walletAddressBtc,omitempty"`
	WalletAddressBTCTaproot string        `json:"walletAddressBtcTaproot,omitempty"`
	DisplayName             string        `json:"displayName"`
	Bio                     string        `json:"bio"`
	Avatar                  string        `json:"avatar"`
	CreatedAt               *time.Time    `json:"createdAt"`
	ProfileSocial           ProfileSocial `json:"profileSocial"`
}

type ArtistResponse struct {
	ProfileResponse `json:",inline"`
	Projects        []*ProjectBasicInfo `json:"projects"`
}

type ProjectBasicInfo struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	WalletAddress string `json:"walletAddress"`
}

type ProfileSocial struct {
	Web       string `json:"web"`
	Twitter   string `json:"twitter"`
	Discord   string `json:"discord"`
	Medium    string `json:"medium"`
	Instagram string `json:"instagram"`
	EtherScan string `json:"etherScan"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}
