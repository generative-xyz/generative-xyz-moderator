package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
)

func (r Repository) InsertModular(obj *entity.ModularInscription) (*mongo.InsertOneResult, error) {
	obj.SetCreatedAt()

	inserted, err := r.DB.Collection(obj.TableName()).InsertOne(context.TODO(), obj)
	if err != nil {
		return nil, err
	}

	return inserted, nil
}

func (r Repository) SetCreatedTokenStatus(inscriptionID string, status bool) (*mongo.UpdateResult, error) {
	f := bson.D{
		{"inscription_id", inscriptionID},
	}

	updatedData := bson.D{
		{"$set", bson.D{
			{"is_created_token", status},
		}},
	}

	inserted, err := r.DB.Collection(utils.COLLECTION_MODULAR_INSCRIPTION).UpdateOne(context.TODO(), f, updatedData)
	if err != nil {
		return nil, err
	}

	return inserted, nil
}
