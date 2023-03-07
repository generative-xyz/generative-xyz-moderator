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

func (r Repository) FindIDeveloperKeyByApiKey(apiKey string) (*entity.DeveloperKey, error) {

	filter := bson.D{{Key: "api_key", Value: apiKey}, {Key: "status", Value: 1}}

	resp := &entity.DeveloperKey{}
	usr, err := r.FilterOne(entity.DeveloperKey{}.TableName(), filter)
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) FindIDeveloperKeyByUserID(userID string) (*entity.DeveloperKey, error) {

	resp := &entity.DeveloperKey{}
	usr, err := r.FilterOne(entity.DeveloperKey{}.TableName(), bson.D{{"user_uuid", userID}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) InsertDeveloperKey(data *entity.DeveloperKey) error {
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}

//////
func (r Repository) FindDeveloperKeyRequests(apiKey string) (*entity.DeveloperKeyRequests, error) {

	filter := bson.D{{Key: "api_key", Value: apiKey}}

	resp := &entity.DeveloperKeyRequests{}
	usr, err := r.FilterOne(entity.DeveloperKeyRequests{}.TableName(), filter)
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) InsertDeveloperKeyRequests(data *entity.DeveloperKeyRequests) error {
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) UpdateDeveloperKeyRequests(model *entity.DeveloperKeyRequests) (*mongo.UpdateResult, error) {

	filter := bson.D{{Key: "uuid", Value: model.UUID}}
	result, err := r.UpdateOne(model.TableName(), filter, model)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) IncreaseDeveloperReqCounter(apiKeyID string) error {
	filter := bson.D{{Key: "api_key", Value: apiKeyID}}

	update := bson.M{
		"$set": bson.M{
			"day_req_last_time": primitive.NewDateTimeFromTime(time.Now().UTC()),
		},
		"$inc": bson.M{"day_req_counter": 1},
	}

	_, err := r.DB.Collection(entity.DeveloperKeyRequests{}.TableName()).UpdateOne(context.TODO(), filter, update)
	return err
}

///

func (r Repository) InsertDeveloperKeyReqLogs(data *entity.DeveloperKeyReqLogs) error {
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}
