package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type FilterProposals struct {
	BaseFilters
	Proposer *string
}

const (
	StatePending = iota
	StateActivate
	Canceled
	Defeated
	Successeded
	Queued
	Expired
	Executed
)

type Proposal struct {
	BaseEntity `bson:",inline"`
	ProposalID string `bson:"proposalID" json:"proposalID"`
	Proposer string `bson:"proposer" json:"proposer"`
	ReceiverAddress string `bson:"receiverAddress" json:"receiverAddress"`
	Targets []string `bson:"targets" json:"targets"`
	Values []int64 `bson:"values" json:"values"`
	Signatures []string `bson:"signatures" json:"signatures"`
	Calldatas [][]byte `bson:"calldatas" json:"calldatas"`
	StartBlock int64 `bson:"startBlock" json:"startBlock"`
	EndBlock int64 `bson:"endBlock" json:"endBlock"`
	Title string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`
	Amount string `bson:"amount" json:"amount"`
	TokenType string `bson:"tokenType" json:"tokenType"`
	Raw ProposalRaw `bson:"raw" json:"raw"`
	State uint8 `bson:"state" json:"state"`	
}

type ProposalRaw struct {
	Address string `bson:"address" json:"address"`
	Topics []string `bson:"topics" json:"topics"`
	Data []byte `bson:"data" json:"data"`
	BlockNumber uint64 `bson:"blockNumber" json:"blockNumber"`
	TransactionHash string `bson:"transactionHash" json:"transactionHash"`
	TransactionIndex uint `bson:"transactionIndex" json:"transactionIndex"`
	BlockHash string `bson:"blockHash" json:"blockHash"`
	LogIndex uint `bson:"logIndex" json:"logIndex"`
	Removed bool `bson:"removed" json:"removed"`
}


func (u Proposal) TableName() string { 
	return utils.COLLECTION_DAO_PROPOSAL
}

func (u Proposal) ToBson()  (*bson.D, error) { 
	return helpers.ToDoc(u)
}