package response

type PaginationResponse struct {
	Items       interface{} `json:"items"`
	CurrentPage int         `json:"current_page"`
	NextPage    *int        `json:"next_page"`
	PrevPage    *int        `json:"prev_page"`
	TotalItems  int64       `json:"total_items"`
	TotalPages  int64       `json:"total_pages"`
	Cursor      string      `json:"cursor"`
	// string `json:"previous"`
}