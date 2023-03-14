package response

import (
	"time"

	"rederinghub.io/utils/constants/dao_project"
)

type DaoProject struct {
	BaseEntity `json:",inline"`
	CreatedBy  string                `json:"created_by"`
	User       *UserForDaoProject    `json:"user"`
	ProjectId  string                `json:"project_id"`
	Project    *ProjectForDaoProject `json:"project"`
	ExpiredAt  time.Time             `json:"expired_at"`
	Status     dao_project.Status    `json:"status"`
	Action     *ActionDaoProject     `json:"action"`
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
	BaseEntity  `json:",inline"`
	Name        string             `json:"name"`
	CreatorName string             `json:"creator_name"`
	Thumbnail   string             `json:"thumbnail"`
	MintingInfo ProjectMintingInfo `json:",inline"`
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
