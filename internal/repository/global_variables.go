package repository

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"rederinghub.io/internal/entity"
)

func (r Repository) InsertGlobalVariableInt(key string, value int64) error {
	v := entity.GlobalVariable{
		Key: key,
		Value: value,
	}
	err := r.InsertOne(v.TableName(), &v)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (r Repository) UpdateGlobalActivity(key string, value interface{}) error {
	f := bson.D{{Key: "key", Value: key}}
	update := bson.M{
		"$set": bson.M{
			"value": value,
		},
	}
	_, err := r.DB.Collection(entity.GlobalVariable{}.TableName()).UpdateOne(context.TODO(), f, update)
	if err != nil {
		return errors.Wrap(err, "collection.InsertOne")
	}
	return nil
}

func (r Repository) GetVariableInt(key string) (*int64, error) {
	globalVariable := &entity.GlobalVariable{}
	if err := r.FindOneBy(context.TODO(), globalVariable.TableName(), bson.M{
		"key": key,
	}, globalVariable); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			var val int64
			r.InsertGlobalVariableInt(key, val)
			return &val, nil
		} else {
			return nil, err
		}
	}
	val, ok := globalVariable.Value.(int64)
	if !ok {
		return nil, errors.New("variable is not of type int64")
	}
	return &val, nil
}
