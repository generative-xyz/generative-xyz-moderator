package request

import (
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/constants/dao_artist_voted"
)

type ListDaoArtistRequest struct {
	*entity.Pagination
	Status  *int64  `query:"status"`
	Keyword *string `query:"keyword"`
	Id      *string `query:"-"`
}

type VoteDaoArtistRequest struct {
	Status dao_artist_voted.Status `json:"status"`
}
