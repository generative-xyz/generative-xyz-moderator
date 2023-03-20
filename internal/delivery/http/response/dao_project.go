package response

import (
	"time"

	"rederinghub.io/utils/constants/dao_project"
	"rederinghub.io/utils/constants/dao_project_voted"
)

type DaoProject struct {
	BaseEntity      `json:",inline"`
	SeqId           uint                  `json:"seq_id"`
	CreatedBy       string                `json:"created_by"`
	User            *UserForDao           `json:"user"`
	ProjectId       string                `json:"project_id"`
	Project         *ProjectForDaoProject `json:"project"`
	ExpiredAt       time.Time             `json:"expired_at"`
	Status          dao_project.Status    `json:"status"`
	DaoProjectVoted []*DaoProjectVoted    `json:"dao_project_voted"`
	Action          *ActionDaoProject     `json:"action"`
	TotalAgainst    *int64                `json:"total_against"`
	TotalVote       *int64                `json:"total_vote"`
}

func (s DaoProject) Expired() bool {
	return s.ExpiredAt.UTC().Unix() < time.Now().UTC().Unix()
}

func (s *DaoProject) SetFields(fns ...func(*DaoProject)) {
	for _, fn := range fns {
		fn(s)
	}
}
func (s DaoProject) WithAction(action *ActionDaoProject) func(*DaoProject) {
	return func(dp *DaoProject) {
		dp.Action = action
	}
}

type ActionDaoProject struct {
	CanVote bool `json:"can_vote"`
}

type ProjectForDaoProject struct {
	BaseEntity   `json:",inline"`
	TokenID      string             `json:"token_id"`
	Name         string             `json:"name"`
	CreatorAddrr string             `json:"creator_by"`
	CreatorName  string             `json:"creator_name"`
	Thumbnail    string             `json:"thumbnail"`
	MaxSupply    int64              `json:"max_supply"`
	MintingInfo  ProjectMintingInfo `json:"minting_info"`
	IsHidden     bool               `json:"is_hidden"`
	MintPrice    string             `json:"mint_price"`
	MintPriceEth string             `json:"mint_price_eth"`
}

type ProjectMintingInfo struct {
	Index        int64 `json:"index"`
	IndexReverse int64 `json:"index_reverse"`
}

type UserForDao struct {
	BaseEntity              `json:",inline"`
	IsVerified              bool          `json:"is_verified"`
	WalletAddress           string        `json:"wallet_address"`
	WalletAddressPayment    string        `json:"wallet_address_payment"`
	WalletAddressBTC        string        `json:"wallet_address_btc"`
	WalletAddressBTCTaproot string        `json:"wallet_address_btc_taproot"`
	DisplayName             string        `json:"display_name"`
	Avatar                  string        `json:"avatar"`
	ProfileSocial           ProfileSocial `json:"profile_social"`
	Stats                   UserStats     `json:"stats"`
}

type UserStats struct {
	CollectionCreated int32 `json:"collection_created"`
	TotalMint         int64 `bson:"total_mint" json:"total_mint"`
	TotalMinted       int64 `bson:"total_minted" json:"total_minted"`
}

type DaoProjectVoted struct {
	CreatedBy    string                   `json:"created_by"`
	DisplayName  string                   `json:"display_name"`
	DaoProjectId string                   `json:"dao_project_id"`
	Status       dao_project_voted.Status `json:"status"`
}

func (s *DaoProjectVoted) SetFields(fns ...func(*DaoProjectVoted)) {
	for _, fn := range fns {
		fn(s)
	}
}
func (s DaoProjectVoted) WithDisplayName(displayName string) func(*DaoProjectVoted) {
	return func(dp *DaoProjectVoted) {
		dp.DisplayName = displayName
	}
}
