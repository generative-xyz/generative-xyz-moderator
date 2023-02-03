package request

type CreateProposalReq struct {
	Title string `json:"title"`
	Description string `json:"description"`
	TokenType string `json:"tokenType"`
	Amount string `json:"amount"`
	ReceiverAddress string `json:"receiverAddress"`
}

type UpdateProposalReq struct {
	ID string `json:"ID"`
	ProposalID string `json:"proposalID"`
}
