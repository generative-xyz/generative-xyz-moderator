package structure

import (
	"time"

	"rederinghub.io/external/etherscan"
)

type FilterProjects struct {
	BaseFilters
	WalletAddress   *string
	ContractAddress *string
	TxHash          *string
	TxHex           *string
	Name            *string
	CategoryIds     []string
	IsHidden        *bool
	Status          *bool
	IsSynced        *bool
	Ids             []string
	CommitTxHash    *string
	RevealTxHash    *string
}

type FilterProposal struct {
	BaseFilters
	Proposer   *string
	State      *int
	ProposalID *string
}

type FilterProposalVote struct {
	BaseFilters
	Voter      *string
	Support    *int
	ProposalID *string
}

type CreateProjectReq struct {
	CreateBtcProjectReq
	TxHash       string  `json:"txHash"`
	TxHex        string  `json:"txHex"`
	ZipLink      *string `json:"ziplink"`
	RevealTxHash *string `json:"revealTxHash"`
}

type CreateBtcProjectReq struct {
	MaxSupply              int64    `json:"maxSupply"`
	LimitSupply            int64    `json:"limitSupply"`
	MintPrice              string   `json:"mintPrice"`
	Name                   string   `json:"name"`
	CreatorName            string   `json:"creatorName"`
	CreatorAddrr           string   `json:"creatorAddrr"`
	CreatorAddrrBTC        string   `json:"creatorAddrrBTC"`
	Description            string   `json:"description"`
	OpenMintUnixTimestamp  int      `json:"openMintUnixTimestamp"`
	CloseMintUnixTimestamp int      `json:"closeMintUnixTimestamp"`
	Thumbnail              string   `json:"thumbnail"`
	ThirdPartyScripts      []string `json:"thirdPartyScripts"`
	Scripts                []string `json:"scripts"`
	TokenDescription       string   `json:"tokenDescription"`
	Styles                 string   `json:"styles"`
	SocialWeb              string   `json:"socialWeb"`
	SocialTwitter          string   `json:"socialTwitter"`
	SocialDiscord          string   `json:"socialDiscord"`
	SocialMedium           string   `json:"socialMedium"`
	SocialInstagram        string   `json:"socialInstagram"`
	License                string   `json:"license"`
	Tags                   []string `json:"tags"`
	Categories             []string `json:"categories"`
	ZipLink                *string  `json:"zipLink"`
	AnimationURL           *string  `json:"animationURL"`
	Royalty                int      `json:"royalty"`
	IsFullChain            bool     `json:"isFullChain"`
	CaptureImageTime       *int     `json:"captureImageTime"`
	Reservers              []string `json:"reservers"`
	ReserveMintPrice       string   `json:"reserveMintPrice"`
	ReserveMintLimit       int      `json:"reserveMintLimit"`
	FromAuthentic          bool     `json:"fromAuthentic"`
	TokenAddress           string   `json:"tokenAddress"`
	TokenId                string   `json:"tokenId"`
	OwnerOf                string   `json:"ownerOf"`
	OrdinalsTx             string   `json:"ordinalsTx"`
	InscribedBy            string   `json:"inscribedBy"`
	LimitMintPerProcess    int      `json:"limitMintPerProcess"`
	Index                  int64    `json:"-"`
	IsHidden               *bool    `json:"-"`

	// GM whitelist holder
	IsSupportGMHolder bool   `json:"isSupportGMHolder"`
	MinimumGMSupport  string `json:"minimumGMSupport"`
}

type UpdateBTCProjectReq struct {
	ProjectID           *string  `json:"projectID"`
	Name                *string  `json:"name"`
	Description         *string  `json:"description"`
	Thumbnail           *string  `json:"thumbnail"`
	IsHidden            *bool    `json:"isHidden"`
	Royalty             *int     `json:"royalty"`
	MintPrice           *string  `json:"mintPrice"`
	MaxSupply           *int64   `json:"maxSupply"`
	CreatetorAddress    *string  `json:"createtorAddress"`
	Categories          []string `json:"categories"`
	CaptureImageTime    *int     `json:"captureImageTime"`
	LimitMintPerProcess *int     `json:"limitMintPerProcess"`
	Reservers           []string `json:"reservers"`
	ReserveMintPrice    *string  `json:"reserveMintPrice"`
	ReserveMintLimit    *int     `json:"reserveMintLimit"`
	Index               *int64   `json:"index"`

	// GM whitelist holder
	IsSupportGMHolder *bool   `json:"isSupportGMHolder"`
	MinimumGMSupport  *string `json:"minimumGMSupport"`
}

type UpdateBTCProjectCategoriesReq struct {
	ProjectID  *string  `json:"projectID"`
	Categories []string `json:"categories"`
}

type CreateProposaltReq struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	TokenType       string `json:"tokenType"`
	Amount          string `json:"amount"`
	ReceiverAddress string `json:"receiverAddress"`
}

type UpdateProjectReq struct {
	TokenID  string `json:"tokenID"`
	Priority *int   `json:"priority"`
	// GM whitelist holder
	IsSupportGMHolder *bool   `json:"isSupportGMHolder"`
	MinimumGMSupport  *string `json:"minimumGMSupport"`

	ContracAddress string `json:"contractAddress"`
}

type UpdateProjectHash struct {
	ContractAddress *string `json:"contractAddress"`
	TxHash          *string `json:"txHash"`
	CommitTxHash    *string `json:"commitTxHash"`
	RevealTxHash    *string `json:"revealTxHash"`
}

type GetProjectReq struct {
	ContractAddr string
	TokenID      string
}

type ProjectAnimationUrl struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Image        string `json:"image"`
	AnimationUrl string `json:"animation_url"`
}

type VolumnLogs struct {
	ProjectID   string `json:"projectID"`
	Paytype     string `json:"payType"`
	OldMinted   int    `json:"oldMinted"`
	NewMinted   int    `json:"newMinted"`
	TotalMinted int    `json:"totaMinted"`

	OldAmount        string `json:"oldEthAmount"`
	NewAmount        string `json:"newAmount"`
	TotalAmount      string `json:"totalAmount"`
	TotalEarnings    string `json:"totalEarnings"`
	ApprovedWithdraw string `json:"approvedWithdraw"`
	Available        string `json:"available"`
	SeparateRate     string `json:"separateRate"`
	GenEarnings      string `json:"genEarnings"`
	MintPrice        int    `json:"mintPrice"`
}

type ProjectFloorPrice struct {
	ID    string `bson:"_id" json:"id"`
	Floor uint64 `bson:"floor" json:"floor"`
}

type AnalyticsProjectDeposit struct {
	Value            string                             `json:"value"`
	Currency         string                             `json:"currency"`
	CurrencyRate     float64                            `json:"currencyRate"`
	UsdtValue        float64                            `json:"usdtValue"`
	UsdtExtra        float64                            `json:"usdtExtra"`
	TotalGMReceive   float64                            `json:"totalGMReceive"`
	Items            []*etherscan.AddressTxItemResponse `json:"items"`
	MapItems         map[string]*etherscan.AddressTxItemResponse
	MapTokensDeposit map[string][]TokensDeposit `json:"map_tokens_deposit"`
}

type AnalyticsProjectDepositExternal struct {
	Error  interface{}             `json:"error"`
	Status bool                    `json:"status"`
	Data   AnalyticsProjectDeposit `json:"data"`
}

type TokensDeposit struct {
	Name      string  `json:"name"`
	Value     string  `json:"value"`
	UsdtValue float64 `json:"usdtValue"`
}

type AnalyticsProjectDepositChan struct {
	Value *AnalyticsProjectDeposit
	Err   error
}

type WalletResponse struct {
	UserAddress  string     `json:"user_address"`
	ENS          string     `json:"ens"`
	Avatar       string     `json:"avatar"`
	Address      string     `json:"address"`
	Status       int        `json:"status"`
	Type         string     `json:"type"`
	NativeAmount []string   `json:"native_amounts"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}
