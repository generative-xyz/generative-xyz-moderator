package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type DexBTCListing struct {
	BaseEntity             `bson:",inline"`
	RawPSBT                string     `bson:"raw_psbt"`
	InscriptionID          string     `bson:"inscription_id"`
	Amount                 uint64     `bson:"amount"`
	InscriptionOutputValue uint64     `bson:"inscription_output_value"`
	SellerAddress          string     `bson:"seller_address"`
	Verified               bool       `bson:"verified"`
	CancelAt               *time.Time `bson:"cancel_at"`
	Cancelled              bool       `bson:"cancelled"`
	CancelTx               string     `bson:"cancel_tx"`
	Inputs                 []string   `bson:"inputs"`
	Matched                bool       `bson:"matched"`
	MatchedTx              string     `bson:"matched_tx"`
	MatchAt                *time.Time `bson:"matched_at"`
}

func (u DexBTCListing) TableName() string {
	return utils.COLLECTION_DEX_BTC_LISTING
}

func (u DexBTCListing) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
