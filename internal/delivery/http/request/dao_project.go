package request

import (
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/constants/dao_project_voted"
)

type ListDaoProjectRequest struct {
	entity.Pagination
	Status  *int64  `query:"status"`
	Keyword *string `query:"keyword"`
	Id      *string `query:"-"`
}
type CreateDaoProjectRequest struct {
	ProjectId string `json:"project_id"`
	CreatedBy string `json:"-"`
}

type VoteDaoProjectRequest struct {
	Status        dao_project_voted.Status `json:"int64"`
	WalletAddress string                   `json:"-"`
}
