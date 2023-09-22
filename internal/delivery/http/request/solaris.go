package request

import (
	"fmt"
	"net/url"
)

type PaginationReq struct {
	Limit  *int
	Page   *int
	Offset *int
	SortBy *string
	Sort   *int
}

func (pq PaginationReq) ToNFTServiceUrlQuery() url.Values {
	q := url.Values{}

	if pq.Limit != nil && *pq.Limit != 0 {
		q.Set("limit", fmt.Sprintf("%d", *pq.Limit))
	}

	if pq.Offset != nil {
		q.Set("offset", fmt.Sprintf("%d", *pq.Offset))
	}

	return q
}
