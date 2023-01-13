package structure

import (
	"math/big"
	"net/http"

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

type GetNftTransactionsReq struct {
	Chain *string
	ContractAddress string
	TokenID string
}

type FilterTokens struct {
	BaseFilters
	ContractAddress *string
	OwnerAddr *string
	CreatorAddr *string
	GenNFTAddr *string
	CollectionIDs []string
	TokenIDs []string
}

func (f *FilterTokens) CreateFilter(r *http.Request) {
	contractAddress := r.URL.Query().Get("contract_address")
	geNftAddr := r.URL.Query().Get("gen_nft_address")
	ownerAddress := r.URL.Query().Get("owner_address")
	creatorAddress := r.URL.Query().Get("creator_address")


	tokenID := r.URL.Query().Get("tokenID")
	if tokenID != "" {
		f.TokenIDs = append(f.TokenIDs, tokenID)
	}
	
	if contractAddress != "" {
		f.ContractAddress = &contractAddress
	}
	
	if geNftAddr != "" {
		f.GenNFTAddr = &geNftAddr
	}
	
	if ownerAddress != "" {
		f.OwnerAddr = &ownerAddress
	}
	
	if creatorAddress != "" {
		f.CreatorAddr = &creatorAddress
	}
}

type FilterMkListing struct {
	BaseFilters
	CollectionContract *string
	TokenId *string
	Erc20Token *string
	SellerAddress *string
	Closed             *bool  
	Finished           *bool  
}

type FilterMkOffers struct {
	BaseFilters
	CollectionContract *string
	TokenId *string
	Erc20Token *string
	BuyerAddress *string
	Closed             *bool  
	Finished           *bool  
}

type UpdateTokenReq struct {
	TokenID string `json:"tokenID"`
	Priority *int `json:"priority"`
	ContracAddress string `json:"contractAddress"`
}