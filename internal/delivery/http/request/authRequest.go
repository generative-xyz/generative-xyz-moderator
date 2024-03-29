package request

import (
	"errors"

	"rederinghub.io/internal/entity"
)

type RefreshTokenData struct {
	RefreshToken string `json:"refreshToken"`
	RedirectUri  string `json:"redirectUri"`
}

type GenerateMessageRequest struct {
	Address    *string `json:"address"`
	WalletType string  `json:"walletType"`
}

type VerifyMessageRequest struct {
	ETHSinature      *string `json:"ethSignature"`
	Sinature         *string `json:"signature"`
	Address          *string `json:"address"`
	AddressBTC       *string `json:"addressBtc"`
	AddressBTCSegwit *string `json:"addressBtcSegwit"`
	MessagePrefix    *string `json:"messagePrefix"`
	AddressPayment   string  `json:"addressPayment"`
}

type UpdateProfileRequest struct {
	DisplayName          *string       `json:"displayName"`
	Bio                  *string       `json:"bio"`
	Avatar               *string       `json:"avatar"`
	Banner               *string       `json:"banner"`
	ProfileSocial        ProfileSocial `json:"profileSocial"`
	WalletAddressBTC     string        `json:"walletAddressBtc"`
	WalletAddressPayment string        `json:"walletAddressPayment"`
	EnableNotification   *bool         `json:"enableNotification"`
}

type ProfileSocial struct {
	Web       *string `json:"web"`
	Twitter   *string `json:"twitter"`
	Discord   *string `json:"discord"`
	Medium    *string `json:"medium"`
	Instagram *string `json:"instagram"`
	EtherScan *string `json:"etherScan"`
}

func (g GenerateMessageRequest) SelfValidate() error {
	if g.Address == nil {
		return errors.New("Address is required")
	}

	if *g.Address == "" {
		return errors.New("Address is not empty")
	}

	switch g.WalletType {
	case "", entity.WalletType_BTC_PRVKEY:
		break
	default:
		return errors.New("invalid wallet type")
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
