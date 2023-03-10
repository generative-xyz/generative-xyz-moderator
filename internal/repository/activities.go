package repository

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/internal/entity"
)

func (r Repository) InsertActitvy(data *entity.Activity) error {
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return errors.Wrap(err, "collection.InsertOne")
	}
	return nil
}

// Get all BTC activity (Mint, Buy)
func (r Repository) GetRecentBTCActivity(page int64, limit int64) (*entity.Pagination, error) {
	activities := []entity.Activity{}
	resp := &entity.Pagination{}

	f := bson.M{
		"type":       bson.M{"$in": []entity.ActivityType{entity.Mint, entity.Buy}},
		"created_at": bson.M{"$gte": time.Now().Add(-14 * 24 * time.Hour).UTC()},
	}

	s := []Sort{
		{SortBy: "created_at", Sort: entity.SORT_ASC},
	}

	p, err := r.Paginate(entity.Activity{}.TableName(), page, limit, f, bson.D{}, s, &activities)
	if err != nil {
		return nil, err
	}

	resp.Result = activities
	resp.Page = p.Pagination.Page
	resp.Total = p.Pagination.Total
	resp.PageSize = limit
	return resp, nil
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

func (r Repository) CountMintActivity(projectID string) (*int64, error) {
	f := bson.M{
		"type":       entity.Mint,
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
