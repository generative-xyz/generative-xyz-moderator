package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
)

func (r Repository) InsertActitvy(data *entity.Activity) error {
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}

// Get all BTC activity (Mint, Buy)
func (r Repository) GetRecentBTCActivity() ([]entity.Activity, error) {
	activities := []entity.Activity{}

	f := bson.M{
		"type": bson.M{"$in": []entity.ActivityType{entity.Mint, entity.Buy}},
		"created_at": bson.M{"$gte":  time.Now().Add(-24*time.Hour).UTC()},
	}

	cursor, err := r.DB.Collection(utils.COLLECTION_ACTIVITIES).Find(context.TODO(), f)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &activities); err != nil {
		return nil, err
	}

	return activities, nil
}

func (r Repository) CountViewActivity(projectID string) (*int64, error) {
	f := bson.M{
		"type": entity.View,
		"created_at": bson.M{"$gte":  time.Now().Add(-24*time.Hour).UTC()},
		"project_id": projectID,
	}

	count, err := r.DB.Collection(utils.COLLECTION_ACTIVITIES).CountDocuments(context.TODO(), f)
	if err != nil {
		return nil, err
	}

	return &count, nil
}
