package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
	"time"
)

type SoralisSnapShotBalance struct {
	BaseEntity    `bson:",inline" json:"-"`
	WalletAddress string `bson:"walletAddress" json:"walletAddress"`
	TokenAddress  string `bson:"tokenAddress" json:"tokenAddress"`
	Balance       string `bson:"balance" json:"balance"`
	Decimal       int    `bson:"decimal" json:"decimal"`
}

type FilteredSoralisSnapShotBalance struct {
	WalletAddress string     `bson:"walletAddress" json:"walletAddress"`
	TokenAddress  string     `bson:"tokenAddress" json:"tokenAddress"`
	Balance       string     `bson:"balance" json:"balance"`
	Decimal       int        `bson:"decimal" json:"decimal"`
	CreatedAt     *time.Time `bson:"created_at" json:"createdAt"`
}

func (job SoralisSnapShotBalance) TableName() string {
	return utils.COLLECTION_SORALIS_SNAPSHOT_BALANCE
}

func (job SoralisSnapShotBalance) ToBson() (*bson.D, error) {
	return helpers.ToDoc(job)
}
