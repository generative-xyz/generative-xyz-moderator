package repository

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/helpers"
)

func (r Repository) FindAuctionCollectionBidderByAddress(bidder string) (*entity.AuctionCollectionBidder, error) {
	bidder = strings.ToLower(bidder)
	resp := &entity.AuctionCollectionBidder{}
	usr, err := r.FilterOne(entity.AuctionCollectionBidder{}.TableName(), bson.D{{"bidder", bidder}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) InsertAuctionCollectionBidder(data *entity.AuctionCollectionBidder) error {
	data.Bidder = strings.ToLower(data.Bidder)
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) UpdateAuctionCollectionBidder(model *entity.AuctionCollectionBidder) (*mongo.UpdateResult, error) {

	filter := bson.D{{"uuid", model.UUID}}
	result, err := r.UpdateOne(model.TableName(), filter, model)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) ListAuctionCollectionBidder() ([]entity.AuctionCollectionBidder, error) {
	confs := []entity.AuctionCollectionBidder{}
	f := bson.M{}

	opts := options.Find()
	cursor, err := r.DB.Collection(entity.AuctionCollectionBidder{}.TableName()).Find(context.TODO(), f, opts)
	if err != nil {
		return confs, err
	}

	if err = cursor.All(context.TODO(), &confs); err != nil {
		return confs, err
	}

	return confs, nil
}
