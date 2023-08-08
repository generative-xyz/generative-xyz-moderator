package repository

import (
	"context"
	"errors"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/internal/entity"
)

func (r Repository) GetCutomTypesByProjectID(projectID int) (*entity.DiscordCustomTypes, error) {
	resp := &entity.DiscordCustomTypes{}
	f := bson.D{
		{"project_id", projectID},
	}

	cursor := r.DB.Collection(entity.DiscordCustomTypes{}.TableName()).FindOne(context.TODO(), f)
	if err := cursor.Decode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (r Repository) GetCutomTypeByProjectID(projectID string, key string) (*string, error) {
	projectIDInt, err := strconv.Atoi(projectID)
	if err != nil {
		return nil, err
	}

	w, err := r.GetCutomTypesByProjectID(projectIDInt)
	if err != nil {
		return nil, err
	}

	k, ok := w.CustomTypes[key]
	if !ok {
		return nil, errors.New("Cannot get key for webhook")
	}

	return &k, nil
}

func (r Repository) InsertTypeByProjectID(obj *entity.DiscordCustomTypes) error {
	err := r.InsertOne(entity.DiscordCustomTypes{}.TableName(), obj)
	if err != nil {
		return err
	}
	return nil
}
