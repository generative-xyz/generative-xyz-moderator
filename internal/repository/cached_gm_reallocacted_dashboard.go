package repository

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/internal/entity"
)

func (r *Repository) InsertReAllocated(data *entity.CachedGMReAllocatedDashBoard) error {
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetTheLatestReAllocated() (*entity.CachedGMReAllocatedDashBoard, error) {
	resp := []entity.CachedGMReAllocatedDashBoard{}
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

	cursor, err := r.DB.Collection(entity.CachedGMReAllocatedDashBoard{}.TableName()).Aggregate(context.Background(), f)
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
