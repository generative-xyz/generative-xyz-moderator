package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type FilterProposalVotes struct {
	BaseFilters
	Voter *string
	Support *int
	ProposalID *string
}

type ProposalVotes struct {
	BaseEntity `bson:",inline"`
	ProposalID string `bson:"proposalID" json:"proposalID"`
	Voter string `bson:"voter" json:"voter"`
	Support int `bson:"support" json:"support"`
	Weight uint64 `bson:"weight" json:"weight"`
	Reason string `bson:"reason" json:"reason"`
	
}

func (u ProposalVotes) TableName() string { 
	return utils.COLLECTION_DAO_PROPOSAL_VOTES
}

func (u ProposalVotes) ToBson()  (*bson.D, error) { 
	return helpers.ToDoc(u)
}