package request

import (
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/constants/dao_project_voted"
)

type ListDaoProjectRequest struct {
	*entity.Pagination
	Status  *int64  `query:"status"`
	Keyword *string `query:"keyword"`
	SeqId   *string `query:"seq_id"`
	Id      *string `query:"id"`
}
type CreateDaoProjectRequest struct {
	ProjectIds []string `json:"project_ids" validate:"required"`
	CreatedBy  string   `json:"-"`
}

type VoteDaoProjectRequest struct {
	Status dao_project_voted.Status `json:"status"`
}

type ListProjectHiddenRequest struct {
	*entity.Pagination
	Keyword *string `query:"keyword"`
}
