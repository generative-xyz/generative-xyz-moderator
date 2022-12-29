package nfts

import (
	"rederinghub.io/utils/config"
)

type CovalentNfts struct {
	conf      *config.Config
	serverURL string
	apiKey    string
}

func NewCovalentNfts(conf *config.Config)
