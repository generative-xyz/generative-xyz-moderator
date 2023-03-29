package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/helpers"
)

func (r Repository) CreateProjectUnzip(data *entity.ProjectZipLinks) error {
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) GetProjectUnzip(projectID string) (*entity.ProjectZipLinks, error) {
	resp := &entity.ProjectZipLinks{}
	usr, err := r.FilterOne(entity.ProjectZipLinks{}.TableName(), bson.D{{"projectID", projectID}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r Repository) GetProjectUnzips() ([]entity.ProjectZipLinks, error) {
	checkTime := time.Now().UTC().Add(-5 * time.Minute)
	zipFiles := []entity.ProjectZipLinks{}
	f := bson.M{}
	f["status"] = entity.UzipStatusFail
	f["retries"] = bson.M{"$lte": 20}
	f["updated_at"] =  bson.M{ "$lte":  primitive.NewDateTimeFromTime(checkTime) }
	//f["tokenid"] = "1001572"
	
	cursor, err := r.DB.Collection(entity.ProjectZipLinks{}.TableName()).Find(context.TODO(), f, nil)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &zipFiles); err != nil {
		return nil, err
	}

	return zipFiles, nil
}


func (r Repository) UpdateProjectUnzip(projectID string, data *entity.ProjectZipLinks) (*mongo.UpdateResult, error) {
	filter := bson.D{{"projectID", projectID}}
	result, err := r.UpdateOne(entity.ProjectZipLinks{}.TableName(), filter, data)
	if err != nil {
		return nil, err
	}
	return result, nil
}
