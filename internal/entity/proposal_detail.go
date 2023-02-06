package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type ProposalDetail struct {
	BaseEntity `bson:",inline"`
	ProposalID string `bson:"proposalID" json:"proposalID"`
	ReceiverAddress string `bson:"receiverAddress" json:"receiverAddress"`
	Title string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`
	Amount string `bson:"amount" json:"amount"`
	TokenType string `bson:"tokenType" json:"tokenType"`
	IsDraft bool `bson:"isDraft" json:"isDraft"`
}

func (u ProposalDetail) TableName() string { 
	return utils.COLLECTION_DAO_PROPOSAL_DETAIL
}

func (u ProposalDetail) ToBson()  (*bson.D, error) { 
	return helpers.ToDoc(u)
}