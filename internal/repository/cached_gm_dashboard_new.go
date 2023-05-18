package repository

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/internal/entity"
	"time"
)

func (r *Repository) AggregateGMDashboardNewCachedDataByTime(time *time.Time) (*entity.CachedGMDashBoardNew, error) {
	resp := []entity.CachedGMDashBoardNew{}
	f := bson.D{}
	cursor, err := r.DB.Collection(entity.CachedGMDashBoardNew{}.TableName()).Aggregate(context.Background(), f)
	if err != nil {
		return nil, err
	}

	if err = cursor.All((context.TODO()), &resp); err != nil {
		return nil, err
	}

	if len(resp) < 1 {
		return nil, errors.New("Cannot get backup data")
	}

	item := resp[0]
	return &item, nil
}

func (r *Repository) InsertGMDashboardNewCached(data *entity.CachedGMDashBoardNew) error {
	data.SetID()
	data.SetCreatedAt()
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetTheLatestGMDashboardNewCached() (*entity.CachedGMDashBoardNew, error) {
	resp := []entity.CachedGMDashBoardNew{}
	f := bson.A{
		bson.D{{"$sort", bson.D{{"created_at", -1}}}},
		bson.D{
			{"$project",
				bson.D{
					{"_id", 1},
					{"backup_url", 1},
					{"backup_file_path", 1},
					{"backup_file_name", 1},
				},
			},
		},
		bson.D{{"$limit", 1}},
	}

	cursor, err := r.DB.Collection(entity.CachedGMDashBoardNew{}.TableName()).Aggregate(context.Background(), f)
	if err != nil {
		return nil, err
	}

	if err = cursor.All((context.TODO()), &resp); err != nil {
		return nil, err
	}

	if len(resp) < 1 {
		return nil, errors.New("Cannot get backup data")
	}
	item := resp[0]
	return &item, nil
}
