package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils/helpers"
)

type AuctionCollectionBidder struct {
	BaseEntity `bson:",inline"`

	Bidder   string `bson:"bidder" json:"bidder"`
	IsWinner bool   `bson:"isWinner" json:"isWinner"`
	Amount   string `bson:"amount" json:"amount"`
	Ens      string `bson:"ens" json:"ens"`
}

func (u AuctionCollectionBidder) TableName() string {
	return "auction_collection_bidder"
}

func (u AuctionCollectionBidder) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
