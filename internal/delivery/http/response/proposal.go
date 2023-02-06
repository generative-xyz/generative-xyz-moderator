package response

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
}

type ProposalVote struct { 
	For int `json:"for"`
	Against int `json:"against"`
	Total int `json:"total"`
	PercentFor float32 `json:"percentFor"`
	PercentAgainst float32 `json:"percentAgainst"`
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