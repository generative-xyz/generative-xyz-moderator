package repository

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
)

func (r Repository) InsertActitvy(data *entity.Activity) error {
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return errors.Wrap(err, "collection.InsertOne")
	}
	return nil
}

// Get all BTC activity (Mint, Buy)
func (r Repository) GetRecentBTCActivity() ([]entity.Activity, error) {
	activities := []entity.Activity{}

	f := bson.M{
		"type":       bson.M{"$in": []entity.ActivityType{entity.Mint, entity.Buy}},
		"created_at": bson.M{"$gte": time.Now().Add(-24 * time.Hour).UTC()},
	}

	cursor, err := r.DB.Collection(utils.COLLECTION_ACTIVITIES).Find(context.TODO(), f)
	if err != nil {
		return nil, errors.Wrap(err, "collection.Find")
	}

	if err = cursor.All(context.TODO(), &activities); err != nil {
		return nil, errors.Wrap(err, "collection.Cursor")
	}

	return activities, nil
}

func (r Repository) CountViewActivity(projectID string) (*int64, error) {
	f := bson.M{
		"type":       entity.View,
		"created_at": bson.M{"$gte": time.Now().Add(-24 * time.Hour).UTC()},
		"project_id": projectID,
	}

	count, err := r.DB.Collection(entity.Activity{}.TableName()).CountDocuments(context.TODO(), f)
	if err != nil {
		return nil, errors.Wrap(err, "collection.CountDocuments")
	}

	return &count, nil
}

func (r Repository) JobDeleteOldActivities() error {
	f := bson.M{
		"created_at": bson.M{"$lte": time.Now().Add(-7 * 24 * time.Hour)},
	}

	_, err := r.DB.Collection(entity.Activity{}.TableName()).DeleteMany(context.TODO(), f)
	if err != nil {
		return errors.Wrap(err, "collection.DeleteMany")
	}
	return nil
}
