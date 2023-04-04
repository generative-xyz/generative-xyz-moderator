package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils/helpers"
)

type AuctionCollectionBidder struct {
	BaseEntity `bson:",inline"`

	Bidder   string `bson:"bidder"`
	IsWinner bool   `bson:"isWinner"`
	Amount   string `bson:"amount"`
	Ens      string `bson:"ens"`
}

func (u AuctionCollectionBidder) TableName() string {
	return "auction_collection_bidder"
}

func (u AuctionCollectionBidder) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
