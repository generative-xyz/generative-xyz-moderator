package structure

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"rederinghub.io/utils/contracts/generative_project_contract"
)


type GetTokenMessageReq struct {
	ContractAddress string
	TokenID string
}

type GetProjectDetailMessageReq struct {
	ContractAddress string
	ProjectID string
}

type GetTokenMessageResp struct {
	ContractAddress string
	TokenID string
}

type ProjectDetail struct {
	ProjectDetail *generative_project_contract.NFTProjectProject
	Status bool
	NftTokenUri string
	Royalty ProjectRoyalty
	NftProjectDetail NftProjectDetail
}

type ProjectRoyalty struct {
	Data big.Int
}

type NftProjectDetail struct {
		ProjectAddr     common.Address
		ProjectId       *big.Int
		MaxSupply       *big.Int
		Limit           *big.Int
		Index           *big.Int
		IndexReserve    *big.Int
		Creator         string
		MintPrice       *big.Int
		MintPriceAddr   common.Address
		Name            string
		MintingSchedule interface{}
}


type ProjectDetailChan struct {
	ProjectDetail *generative_project_contract.NFTProjectProject
	Err error
}

type ProjectStatusChan struct {
	Status *bool
	Err error
}

type ProjectNftTokenUriChan struct {
	TokenURI *string
	Err error
}

type RoyaltyChan struct {
	Data *big.Int
	Err error
}


type NftProjectDetailChan struct {
	Data *NftProjectDetail
	Err error
}

