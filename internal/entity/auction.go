package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils/helpers"
)

type AuctionCollectionBidder struct {
	BaseEntity `bson:",inline" json:"-"`

	Bidder   string `bson:"bidder" json:"bidder"`
	IsWinner bool   `bson:"isWinner" json:"isWinner"`
	Amount   string `bson:"amount" json:"amount"`

	UnitPrice string `bson:"unitPrice" json:"unitPrice"`
	Quantity  int    `bson:"quantity" json:"quantity"`

	Contract string `bson:"contract" json:"contract"`

	Ens string `bson:"ens" json:"ens"`
}

func (u AuctionCollectionBidder) TableName() string {
	return "auction_collection_bidder"
}

func (u AuctionCollectionBidder) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}

type AuctionCollectionBidderShort struct {
	Bidder   string `bson:"bidder" json:"bidder"`
	IsWinner bool   `bson:"isWinner" json:"isWinner"`
	Amount   string `bson:"amount" json:"amount"`

	UnitPrice string `bson:"unitPrice" json:"unitPrice"`
	Quantity  int    `bson:"quantity" json:"quantity"`

	Contract string `bson:"contract" json:"contract"`

	Ens string `bson:"ens" json:"ens"`
}

type AuctionWinnerList struct {
	Address   string `bson:"address" json:"address"`
	Quantity  int    `bson:"quantity" json:"quantity"`
	MintPrice int    `bson:"mintPrice" json:"mintPrice"`
}
