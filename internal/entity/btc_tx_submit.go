package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type BTCTransactionSubmit struct {
	BaseEntity          `bson:",inline"`
	Txhash              string                     `bson:"txhash" json:"txhash"`
	Raw                 string                     `bson:"raw" json:"raw"`
	RelatedInscriptions []string                   `bson:"related_inscriptions" json:"related_inscriptions"`
	Status              BTCTransactionSubmitStatus `bson:"status" json:"status"`
}

func (u BTCTransactionSubmit) TableName() string {
	return utils.COLLECTION_BTC_TX_SUBMIT
}

func (u BTCTransactionSubmit) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}

type BTCTransactionSubmitStatus int

const (
	StatusBTCTransactionSubmit_Waiting BTCTransactionSubmitStatus = iota // 0: waiting
	StatusBTCTransactionSubmit_Pending                                   // 1: pending
	StatusBTCTransactionSubmit_Success                                   // 2: successful
	StatusBTCTransactionSubmit_Failed                                    // 3: failed
)
