package structure

import (
	"time"

	"rederinghub.io/external/nfts"
	"rederinghub.io/internal/entity"
)

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
