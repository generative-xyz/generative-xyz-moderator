package response

type PaginationResponse struct {
	Items       interface{} `json:"items"`
	CurrentPage int         `json:"currentPage"`
	NextPage    *int        `json:"nextPage"`
	PrevPage    *int        `json:"prevPage"`
	TotalItems  int64       `json:"totalItems"`
	TotalPages  int64       `json:"totalPages"`
	Cursor      string      `json:"cursor"`
	// string `json:"previous"`
}