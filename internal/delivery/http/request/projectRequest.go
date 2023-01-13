package request

type CreateProjectReq struct {
	ContractAddress string `json:"contractAddress"`
	TokenID string `json:"tokenID"`
	Tags []string `json:"tags"`
	Categories []string `json:"categories"`
}

type UpdateProjectReq struct {
	Priority *int `json:"priority"`
}
