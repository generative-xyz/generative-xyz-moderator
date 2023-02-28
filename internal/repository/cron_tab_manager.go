package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/helpers"
)

func (r Repository) InsertCronJobManager(data *entity.CronJobManager) error {
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) FindCronJobManager(groupName, jobName string) (*entity.CronJobManager, error) {
	resp := &entity.CronJobManager{}
	f := bson.D{
		{Key: "group", Value: groupName},
		{Key: "job_name", Value: jobName},
	}

	usr, err := r.FilterOne(entity.CronJobManager{}.TableName(), f)
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) FindCronJobManagerByUUID(key string) (*entity.CronJobManager, error) {
	resp := &entity.CronJobManager{}
	usr, err := r.FilterOne(entity.CronJobManager{}.TableName(), bson.D{{"uuid", key}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) FindCronJobManagerByJobKey(jobKey string) ([]entity.CronJobManager, error) {
	resp := []entity.CronJobManager{}
	filter := bson.M{
		"job_key": jobKey,
	}

	cursor, err := r.DB.Collection(entity.CronJobManager{}.TableName()).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (r Repository) UpdateCronJobManager(model *entity.CronJobManager) (*mongo.UpdateResult, error) {

	filter := bson.D{{Key: "uuid", Value: model.UUID}}
	result, err := r.UpdateOne(model.TableName(), filter, model)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) UpdateCronJobManagerStatus(uuid string, status bool) (*mongo.UpdateResult, error) {
	f := bson.D{
		{Key: "uuid", Value: uuid},
	}

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}
	result, err := r.DB.Collection(entity.CronJobManager{}.TableName()).UpdateOne(context.TODO(), f, update)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (r Repository) UpdateCronJobManagerLastSatus(uuid, lastStatus string) (*mongo.UpdateResult, error) {
	f := bson.D{
		{Key: "uuid", Value: uuid},
	}

	update := bson.M{
		"$set": bson.M{
			"last_status": lastStatus,
			"updated_at":  time.Now(),
		},
	}
	result, err := r.DB.Collection(entity.CronJobManager{}.TableName()).UpdateOne(context.TODO(), f, update)
	if err != nil {
		return nil, err
	}

	return result, err
}
