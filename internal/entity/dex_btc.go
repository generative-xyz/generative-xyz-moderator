package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type DexBTCListing struct {
	BaseEntity    `bson:",inline"`
	RawPSBT       string `bson:"raw_psbt"`
	SplitTx       string `bson:"split_tx"`
	InscriptionID string `bson:"inscription_id"`
	Amount        uint64 `bson:"amount"`
	// InscriptionOutputValue uint64     `bson:"inscription_output_value"`
	SellerAddress string `bson:"seller_address"`
	Verified      bool   `bson:"verified"`
	// IsValid       bool       `bson:"is_valid"`
	CancelAt  *time.Time `bson:"cancel_at"`
	Cancelled bool       `bson:"cancelled"`
	CancelTx  string     `bson:"cancel_tx"`
	Inputs    []string   `bson:"inputs"`
	Matched   bool       `bson:"matched"`
	MatchedTx string     `bson:"matched_tx"`
	MatchAt   *time.Time `bson:"matched_at"`
	Buyer     string     `bson:"buyer"`
}

func (u DexBTCListing) TableName() string {
	return utils.COLLECTION_DEX_BTC_LISTING
}

func (u DexBTCListing) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
