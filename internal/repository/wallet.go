package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
)

func (r Repository) CreateTrackTx(trackTx *entity.WalletTrackTx) error {
	err := r.InsertOne(trackTx.TableName(), trackTx)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil
		}
		return err
	}
	return nil
}

func (r Repository) GetTrackTxs(address string, limit, offset int64) ([]entity.WalletTrackTx, error) {
	resp := []entity.WalletTrackTx{}
	filter := bson.M{
		"address": address,
	}

	cursor, err := r.DB.Collection(utils.WALLET_TRACK_TX).Find(context.TODO(), filter, &options.FindOptions{
		Sort:  bson.D{{"created_at", -1}},
		Limit: &limit,
		Skip:  &offset,
	})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
