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
	ContractDecimals     int           
	ContractName         string        
	ContractTickerSymbol string        
	ContractAddress      string        
	SupportsErc          interface{}   
	LogoURL              string        
	Address              string        
	Balance              string        
	TotalSupply          string        
	BlockHeight          int           
	Profile              *entity.Users 
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
