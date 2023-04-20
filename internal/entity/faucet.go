package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils/helpers"
)

type Faucet struct {
	BaseEntity  `bson:",inline"`
	Tx          string `bson:"tx" json:"tcTx"`
	BtcTx       string `bson:"btc_tx" json:"btcTx"`
	Nonce       int64  `bson:"nonce" json:"-"`
	Address     string `bson:"address" json:"address"`
	TwitterName string `bson:"twitter_name" json:"twitterName"`
	Status      int    `bson:"status" json:"status"` // 0 pending, 1 have tx tc, 2 have tx btc, 3 success, 4 false.
	Amount      string `bson:"amount" json:"amount"`
	TwShareID   string `bson:"twitter_share_id" json:"twitterShareId"`
	FaucetType  string `bson:"faucet_type" json:"faucetType"`
	UserTx      string `bson:"user_tx" json:"userTx"`
	SharedLink  string `bson:"shared_link" json:"shared_link"`

	ErrLogs string `bson:"err_logs" json:"-"`

	StatusStr string `bson:"-" json:"status_str"`
}

func (u Faucet) TableName() string {
	return "faucets"
}

func (u Faucet) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
