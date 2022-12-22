package response

type PaginationResponse struct {
	Result interface{} `json:"result"`
	Page int64 `json:"page"`
	PageSize int64 `json:"pageSize"`
	//TotalPage int64 `json:"pageSize"`
	// Next int64
	// Prev int64
	//Limit int64
	Total int64 `json:"total"`
	Currsor string `json:"cursor"`
}