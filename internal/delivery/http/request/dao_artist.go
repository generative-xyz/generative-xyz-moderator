package request

import "rederinghub.io/internal/entity"

type ListDaoArtistRequest struct {
	entity.Pagination
	Status  *int64  `query:"status"`
	Keyword *string `query:"keyword"`
	Id      *string `query:"-"`
}
