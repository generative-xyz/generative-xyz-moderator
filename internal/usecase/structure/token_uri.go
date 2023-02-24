package structure

import (
	"errors"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"rederinghub.io/external/nfts"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/contracts/generative_project_contract"
)

type GetTokenMessageReq struct {
	ContractAddress string
	TokenID         string
}

type GetProjectDetailMessageReq struct {
	ContractAddress string
	ProjectID       string
}

type GetTokenMessageResp struct {
	ContractAddress string
	TokenID         string
}

type ProjectDetail struct {
	ProjectDetail    *generative_project_contract.NFTProjectProject
	Status           bool
	NftTokenUri      string
	Royalty          ProjectRoyalty
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
	Err           error
}

type ProjectStatusChan struct {
	Status *bool
	Err    error
}

type ProjectNftTokenUriChan struct {
	TokenURI *string
	Err      error
}

type RoyaltyChan struct {
	Data *big.Int
	Err  error
}

type NftProjectDetailChan struct {
	Data *NftProjectDetail
	Err  error
}

type NftMoralisChan struct {
	Data *nfts.MoralisToken
	Err  error
}

type TokenAnimationURI struct {
	Thumbnail   string
	ParsedImage string
	Traits      []entity.TokenUriAttr
	TraitsStr   []entity.TokenUriAttrStr
	CapturedAt  *time.Time
	IsUpdated   bool
}

type TokenAnimationURIChan struct {
	Data *TokenAnimationURI
	Err  error
}

type GetNftTransactionsReq struct {
	Chain           *string
	ContractAddress string
	TokenID         string
}

type TokenUriAttrReq struct {
	TraitType string
	Values    []string
}

type FilterTokens struct {
	BaseFilters
	Keyword         *string
	ContractAddress *string
	Search          *string
	OwnerAddr       *string
	CreatorAddr     *string
	GenNFTAddr      *string
	CollectionIDs   []string
	TokenIDs        []string
	Attributes      []TokenUriAttrReq
	HasPrice        *bool
	FromPrice       *int64
	ToPrice         *int64
}

func (f *FilterTokens) CreateFilter(r *http.Request) error {
	contractAddress := r.URL.Query().Get("contract_address")
	geNftAddr := r.URL.Query().Get("gen_nft_address")
	ownerAddress := r.URL.Query().Get("owner_address")
	creatorAddress := r.URL.Query().Get("creator_address")
	keyword := r.URL.Query().Get("keyword")
	search := r.URL.Query().Get("search")
	hasPrice := r.URL.Query().Get("has_price")
	fromPrice := r.URL.Query().Get("from_price")
	toPrice := r.URL.Query().Get("to_price")

	attributesRaw := r.URL.Query().Get("attributes")
	if len(attributesRaw) > 0 {
		attributesStrs := strings.Split(attributesRaw, ",")
		attrs := []TokenUriAttrReq{}
		for _, attributesStr := range attributesStrs {
			parts := strings.Split(attributesStr, ":")
			if len(parts) != 2 {
				return errors.New("errors when parse attribute query")
			}

			added := false
			for i := 0; i < len(attrs); i++ {
				if attrs[i].TraitType == parts[0] {
					attrs[i].Values = append(attrs[i].Values, parts[1])
					added = true
				}
			}
			if !added {
				attrs = append(attrs, TokenUriAttrReq{
					TraitType: parts[0],
					Values:    []string{parts[1]},
				})
			}
		}
		f.Attributes = attrs
	}

	tokenID := r.URL.Query().Get("tokenID")
	if tokenID != "" {
		f.TokenIDs = append(f.TokenIDs, tokenID)
	}
	if contractAddress != "" {
		f.ContractAddress = &contractAddress
	}

	if search != "" {
		f.Search = &search
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
	if keyword != "" {
		f.Keyword = &keyword
	}

	if hasPrice != "" {
		hPrice, err := strconv.ParseBool(hasPrice)
		if err != nil {
			return err
		}
		f.HasPrice = &hPrice
	}

	if fromPrice != "" {
		fPrice, err := strconv.ParseInt(fromPrice, 10, 64)
		if err != nil {
			return err
		}
		f.FromPrice = &fPrice
	}

	if toPrice != "" {
		tPrice, err := strconv.ParseInt(toPrice, 10, 64)
		if err != nil {
			return err
		}
		f.ToPrice = &tPrice
	}

	return nil
}

type FilterMkListing struct {
	BaseFilters
	CollectionContract *string
	TokenId            *string
	Erc20Token         *string
	SellerAddress      *string
	Closed             *bool
	Finished           *bool
}

type FilterMkOffers struct {
	BaseFilters
	CollectionContract *string
	TokenId            *string
	Erc20Token         *string
	BuyerAddress       *string
	Closed             *bool
	Finished           *bool
	OwnerAddress       *string
}

type UpdateTokenReq struct {
	TokenID        string `json:"tokenID"`
	Priority       *int   `json:"priority"`
	ContracAddress string `json:"contractAddress"`
}

type UpdateTokenThumbnailReq struct {
	TokenID   string `json:"tokenID"`
	Thumbnail string `json:"thumbnail"`
}
