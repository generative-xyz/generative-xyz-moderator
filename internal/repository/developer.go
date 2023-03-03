package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/helpers"
)

func (r Repository) FindIDeveloperKeyByApiKey(apiKey string) (*entity.DeveloperKey, error) {

	resp := &entity.DeveloperKey{}
	usr, err := r.FilterOne(entity.DeveloperKey{}.TableName(), bson.D{{"api_key", apiKey}})
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

	resp := &entity.DeveloperKeyRequests{}
	usr, err := r.FilterOne(entity.DeveloperKeyRequests{}.TableName(), bson.D{{"api_key", apiKey}})
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

func (r Repository) IncreaseDeveloperReqCounter(apiKeyID string) error {
	filter := bson.D{{Key: "api_key", Value: apiKeyID}}
	update := bson.M{"$inc": bson.M{"day_req_counter": 1}, "day_req_last_time": time.Now().UTC()}
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
