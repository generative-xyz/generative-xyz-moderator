package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils/helpers"
)

type Faucet struct {
	BaseEntity  `bson:",inline"`
	Tx          string `bson:"tx"`
	BtcTx       string `bson:"btc_tx"`
	Nonce       int64  `bson:"nonce"`
	Address     string `bson:"address"`
	TwitterName string `bson:"twitter_name"`
	Status      int    `bson:"status"` // 0 pending, 1 have tx tc, 2 have tx btc, 3 success, 4 false.
	Amount      string `bson:"amount"`
	TwShareID   string `bson:"twitter_share_id"`
	ErrLogs     string `bson:"err_logs"`
}

func (u Faucet) TableName() string {
	return "faucets"
}

func (u Faucet) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
