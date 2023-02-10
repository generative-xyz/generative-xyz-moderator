package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type MarketplaceBTCBuyOrder struct {
	BaseEntity    `bson:",inline"`
	OrdAddress    string    `bson:"ord_address"`
	ItemID        string    `bson:"item_id"`
	InscriptionID string    `bson:"inscriptionID"` // tokenID in btc
	SegwitAddress string    `bson:"segwit_address"`
	SegwitKey     string    `bson:"segwit_key"`
	ExpiredAt     time.Time `bson:"expired_at"`
}

func (u MarketplaceBTCBuyOrder) TableName() string {
	return utils.COLLECTION_MARKETPLACE_BTC_BUY
}

func (u MarketplaceBTCBuyOrder) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
