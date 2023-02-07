package response

import "time"

type ProposalResp struct {
	BaseResponse
	ProposalID string `json:"proposalID"`
	Proposer string `json:"proposer"`
	Targets []string `json:"targets"`
	Values []int64 `json:"values"`
	Signatures []string `json:"signatures"`
	Calldatas [][]byte `json:"calldatas"`
	StartBlock int64 `json:"startBlock"`
	CurrentBlock int64 `json:"currentBlock"`
	EndBlock int64 `json:"endBlock"`
	Title string `json:"title"`
	Description string `json:"description"`
	TokenType string `json:"tokenType"`
	Amount string `json:"amount"`
	Raw ProposalRaw `json:"raw"`
	State uint `json:"state"`
	ReceiverAddress string `json:"receiverAddress"`
	IsDraft bool `json:"isDraft"`
	Vote ProposalVote `json:"vote"`
	CreatedAt *time.Time `json:"createdAt"`
}

type ProposalVote struct { 
	For uint64 `json:"for"`
	Against uint64 `json:"against"`
	Abstain uint64 `json:"abstain"`

	Total uint64 `json:"total"`
	PercentFor float64 `json:"percentFor"`
	PercentAgainst float64 `json:"percentAgainst"`
	PercentAbstain float64 `json:"percentAbstain"`
}

type ProposalRaw struct {
	Address string `json:"address"`
	Topics []string `json:"topics"`
	Data []byte `json:"data"`
	BlockNumber uint64 `json:"blockNumber"`
	TransactionHash string `json:"transactionHash"`
	TransactionIndex uint `json:"transactionIndex"`
	BlockHash string `json:"blockHash"`
	LogIndex uint `json:"logIndex"`
	Removed bool `json:"removed"`
}