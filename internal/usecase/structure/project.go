package structure

type FilterProjects struct {
	BaseFilters
	WalletAddress *string
	Name *string
}

type FilterProposal struct {
	BaseFilters
	Proposer *string
	State *int
	ProposalID *string
}

type CreateProjectReq struct {
	ContractAddress string `json:"contractAddress"`
	TokenID string `json:"tokenID"`
	Tags []string `json:"tags"`
	Categories []string `json:"categories"`
}

type CreateProposaltReq struct {
	Title string `json:"title"`
	Description string `json:"description"`
	TokenType string `json:"tokenType"`
	Amount string `json:"amount"`
	ReceiverAddress string `json:"receiverAddress"`
}

type UpdateProjectReq struct {
	TokenID string `json:"tokenID"`
	Priority *int `json:"priority"`
	ContracAddress string `json:"contractAddress"`
}

type GetProjectReq struct {
	ContractAddr string
	TokenID string
}
