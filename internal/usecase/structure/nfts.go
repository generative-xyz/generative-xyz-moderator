package structure

import "time"

type GetNftMintedTimeReq struct {
	ContractAddress string
	TokenID string
}

type NftMintedTime struct {
	BlockNumberMinted *string
	MintedTime *time.Time
}

type NftMintedTimeChan struct {
	NftMintedTime *NftMintedTime
	Err error
}
