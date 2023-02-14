package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type InscribeInfo struct {
	BaseEntity `bson:",inline"`
	ID                 string `bson:"id" json:"id"`
	Index              string `bson:"index" json:"index"`
	Address            string `bson:"address" json:"address"`
	OutputValue        string `bson:"outputValue" json:"outputValue"`
	Sat                string `bson:"sat" json:"sat"`
	Preview            string `bson:"preview" json:"preview"`
	Content            string `bson:"content" json:"content"`
	ContentLength      string `bson:"contentLength" json:"contentLength"`
	ContentType        string `bson:"contentType" json:"contentType"`
	Timestamp          string `bson:"timestamp" json:"timestamp"`
	GenesisHeight      string `bson:"genesisHeight" json:"genesisHeight"`
	GenesisFee         string `bson:"genesisFee" json:"genesisFee"`
	GenesisTransaction string `bson:"genesisTransaction" json:"genesisTransaction"`
	Location           string `bson:"location" json:"location"`
	Output             string `bson:"output" json:"output"`
	Offset             string `bson:"offset" json:"offset"`
}

func (u InscribeInfo) TableName() string {
	return utils.INSCRIBE_INFO
}

func (u InscribeInfo) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
