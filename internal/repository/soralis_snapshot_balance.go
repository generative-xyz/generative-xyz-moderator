package repository

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/internal/entity"
	"strings"
)

func (r Repository) InsertSnapshot(data *entity.SoralisSnapShotBalance) error {
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return errors.Wrap(err, "collection.InsertOne")
	}
	return nil
}

func (r Repository) GetSnapshotByWalletAddress(walletAddress string, tokenAddress string) ([]entity.FilteredSoralisSnapShotBalance, error) {
	snapshot := []entity.FilteredSoralisSnapShotBalance{}

	f := bson.A{
		bson.D{
			{"$match",
				bson.D{
					{"walletAddress", strings.ToLower(walletAddress)},
					{"tokenAddress", strings.ToLower(tokenAddress)},
				},
			},
		},
		bson.D{{"$sort", bson.D{{"created_at", 1}}}},
	}

	cursor, err := r.DB.Collection(entity.SoralisSnapShotBalance{}.TableName()).Aggregate(context.TODO(), f)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err = cursor.All((context.TODO()), &snapshot); err != nil {
		return nil, errors.WithStack(err)
	}

	return snapshot, nil
}
