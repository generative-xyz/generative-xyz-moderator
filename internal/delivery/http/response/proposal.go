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
	EndBlock int64 `json:"endBlock"`
	Description string `json:"description"`
	Raw ProposalRaw `json:"raw"`
	State uint `json:"state"`
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