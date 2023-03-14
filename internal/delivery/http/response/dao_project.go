package response

import (
	"time"

	"rederinghub.io/utils/constants/dao_project"
	"rederinghub.io/utils/constants/dao_project_voted"
)

type DaoProject struct {
	BaseEntity      `json:",inline"`
	CreatedBy       string                `json:"created_by"`
	User            *UserForDaoProject    `json:"user"`
	ProjectId       string                `json:"project_id"`
	Project         *ProjectForDaoProject `json:"project"`
	ExpiredAt       time.Time             `json:"expired_at"`
	Status          dao_project.Status    `json:"status"`
	DaoProjectVoted []*DaoProjectVoted    `json:"dao_project_voted"`
	Action          *ActionDaoProject     `json:"action"`
	TotalAgainst    int64                 `json:"total_against"`
	TotalVote       int64                 `json:"total_vote"`
}

func (s DaoProject) Expired() bool {
	return s.ExpiredAt.UTC().Unix() > time.Now().UTC().Unix()
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
	BaseEntity  `json:",inline,omitempty"`
	TokenID     string             `json:"token_id,omitempty"`
	Name        string             `json:"name,omitempty"`
	CreatorName string             `json:"creator_name,omitempty"`
	Thumbnail   string             `json:"thumbnail,omitempty"`
	MintingInfo ProjectMintingInfo `json:"minting_info,omitempty"`
}

type ProjectMintingInfo struct {
	Index        int64 `json:"index"`
	IndexReverse int64 `json:"index_reverse"`
}

type UserForDaoProject struct {
	BaseEntity    `json:",inline,omitempty"`
	IsVerified    bool   `json:"is_verified,omitempty"`
	WalletAddress string `json:"wallet_address,omitempty"`
	DisplayName   string `json:"display_name,omitempty"`
	Avatar        string `json:"avatar,omitempty"`
}

type DaoProjectVoted struct {
	CreatedBy    string                   `json:"created_by,omitempty" bson:"created_by"`
	DisplayName  string                   `json:"display_name,omitempty"`
	DaoProjectId string                   `json:"dao_project_id,omitempty" bson:"dao_project_id"`
	Status       dao_project_voted.Status `json:"status,omitempty" bson:"status"`
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
