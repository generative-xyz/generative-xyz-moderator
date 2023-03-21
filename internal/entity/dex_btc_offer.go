package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type DexBTCOffer struct {
	BaseEntity     `bson:",inline"`
	RawPSBT        string     `bson:"raw_psbt"`
	SplitTx        string     `bson:"split_tx"`
	InscriptionID  string     `bson:"inscription_id"`
	Amount         uint64     `bson:"amount"`
	OffererAddress string     `bson:"offerer_address"`
	Verified       bool       `bson:"verified"`
	CancelAt       *time.Time `bson:"cancel_at"`
	Cancelled      bool       `bson:"cancelled"`
	CancelTx       string     `bson:"cancel_tx"`
	Inputs         []string   `bson:"inputs"`
	Matched        bool       `bson:"matched"`
	MatchedTx      string     `bson:"matched_tx"`
	MatchAt        *time.Time `bson:"matched_at"`
	Seller         string     `bson:"seller"`
	InvalidMatch   bool       `bson:"invalid_match"`
	InvalidMatchTx string     `bson:"invalid_match_tx"`
	ExpiredAt      *time.Time `bson:"expired_at"`

	CreatedVerifiedActivity  bool `bson:"created_verified_activity"`
	CreatedCancelledActivity bool `bson:"created_cancelled_activity"`
	CreatedMatchedActivity   bool `bson:"created_matched_activity"`

	IsTimeSeriesData bool `json:"is_time_series_data"`

	FromOtherMkp bool `bson:"from_other_mkp"`
}

func (u DexBTCOffer) TableName() string {
	return utils.COLLECTION_DEX_BTC_OFFER
}

func (u DexBTCOffer) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
