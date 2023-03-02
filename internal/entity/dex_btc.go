package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type DexBTCListing struct {
	BaseEntity             `bson:",inline"`
	RawPSBT                string   `bson:"raw_psbt"`
	InscriptionID          string   `bson:"inscription_id"`
	Amount                 uint64   `bson:"amount"`
	InscriptionOutputValue uint64   `bson:"inscription_output_value"`
	SellerAddress          string   `bson:"seller_address"`
	Verified               bool     `bson:"verified"`
	Cancelled              bool     `bson:"cancelled"`
	CancelTx               string   `bson:"cancel_tx"`
	Inputs                 []string `bson:"inputs"`
	Matched                bool     `bson:"matched"`
}

func (u DexBTCListing) TableName() string {
	return utils.COLLECTION_DEX_BTC_LISTING
}

func (u DexBTCListing) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
