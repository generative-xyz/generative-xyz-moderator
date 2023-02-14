package request

import "errors"

type RefreshTokenData struct {
	RefreshToken  string `json:"refreshToken"`
	RedirectUri string `json:"redirectUri"`
}

type GenerateMessageRequest struct {
	Address *string `json:"address"`
}

type VerifyMessageRequest struct {
	Sinature *string `json:"signature"`
	Address *string `json:"address"`
}

type UpdateProfileRequest struct {
	DisplayName *string `json:"displayName"`
	Bio *string `json:"bio"`
	Avatar *string `json:"avatar"`
	ProfileSocial ProfileSocial `json:"profileSocial"`
	WalletAddressBTC   string        `json:"wallet_address_btc"`
}

type ProfileSocial  struct{
    Web string `json:"web"`;
    Twitter string `json:"twitter"`;
    Discord string `json:"discord"`;
    Medium string `json:"medium"`;
	Instagram string `json:"instagram"`;
	EtherScan string `json:"etherScan"`;
}

func (g GenerateMessageRequest) SelfValidate() error {
	if g.Address == nil {
		return errors.New("Address is required")
	}

	if *g.Address == "" {
		return errors.New("Address is not empty")
	}

	return nil
}


func (g VerifyMessageRequest) SelfValidate() error {
	if g.Address == nil {
		return errors.New("Address is required")
	}

	if *g.Address == "" {
		return errors.New("Address is not empty")
	}
	
	if g.Sinature == nil {
		return errors.New("Sinature is required")
	}

	if *g.Sinature == "" {
		return errors.New("Sinature is not empty")
	}

	return nil
}
