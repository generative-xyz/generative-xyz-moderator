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

func (r Repository) UnCreatedModularInscriptions(ctx context.Context, offset, limit int) ([]entity.ModularInscription, error) {
	f := bson.A{
		bson.D{{"$match", bson.D{{"is_created_token", false}}}},
		bson.D{{"$sort", bson.D{{"_id", 1}}}}, //first int first out
		bson.D{{"$skip", offset}},
		bson.D{{"$limit", limit}},
	}

	cursor, err := r.DB.Collection(utils.COLLECTION_MODULAR_INSCRIPTION).Aggregate(ctx, f)
	if err != nil {
		return nil, err
	}

	aggregation := []entity.ModularInscription{}
	if err = cursor.All(ctx, &aggregation); err != nil {
		return nil, err
	}

	return aggregation, nil
}
