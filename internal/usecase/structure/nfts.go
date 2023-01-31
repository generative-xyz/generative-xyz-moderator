package structure

import (
	"time"

	"rederinghub.io/external/nfts"
	"rederinghub.io/internal/entity"
)

type GetTokenHolderRequest struct {
	Chain *string
	ContractAddress string
	Page int32
	Limit int32
}

type TokenHolder struct {
	ContractDecimals     int           `json:"contract_decimals,omitempty"`
	ContractName         string        `json:"contract_name,omitempty"`
	ContractTickerSymbol string        `json:"contract_ticker_symbol,omitempty"`
	ContractAddress      string        `json:"contract_address,omitempty"`
	SupportsErc          interface{}   `json:"supports_erc,omitempty"`
	LogoURL              string        `json:"logo_url,omitempty"`
	Address              string        `json:"address,omitempty"`
	Balance              string        `json:"balance,omitempty"`
	TotalSupply          string        `json:"total_supply,omitempty"`
	BlockHeight          int           `json:"block_height,omitempty"`
	Profile              *entity.Users `json:"profile,omitempty"`
}

type GetNftMintedTimeReq struct {
	ContractAddress string
	TokenID string
}

type NftMintedTime struct {
	BlockNumberMinted *string
	MintedTime *time.Time
	Nft *nfts.MoralisToken
}

type TokenDataChan struct {
	Data *entity.TokenUri
	Err error
}

type NftMintedTimeChan struct {
	NftMintedTime *NftMintedTime
	Err error
}
