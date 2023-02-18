package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
)

func (r Repository) CreateTokenUriHistory(data *entity.TokenUriHistories) error {
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}


func (r Repository) ListTokenURIHistory(projectID string) ([]entity.TokenUriHistories, error) {
	// confs := []entity.TokenUriHistoriesArr{}

	// lookUpStage := bson.M{
	// 	"$lookup": bson.M{"from": utils.COLLECTION_TOKEN_URI, "localField": "token_id", "foreignField": "token_id", "as": "token"},
	// }
	
	// projectStage := bson.M{
	// 	"$project": bson.M{"processID": 1, "action": 1, "projectID": 1, "token": 1, "tokenID": 1},
	// }

	// f := bson.A{
	// 	bson.M{"$match": bson.M{"$and": bson.A{
	// 		bson.M{"projectID": projectID},
	// 		bson.M{"action": entity.SENT},
	// 	}}},
	// 	lookUpStage,
	// 	projectStage,
	// 	//bson.M{"$sort": bson.M{"_id": -1}},
	// }

	// cursor, err := r.DB.Collection(utils.COLLECTION_TOKEN_URI_HISTORIES).Aggregate(context.TODO(), f)
	// if err != nil {
	// 	return nil, err
	// }
	// var results []bson.M
	// if err = cursor.All(context.TODO(), &results); err != nil {
	// 	return nil, err
	// }

	// for _, results := range results {
	// 	i := &entity.TokenUriHistoriesArr{}
	// 	err := helpers.Transform(results, i)
	// 	if err != nil {
	// 		continue
	// 	}
	// 	confs =  append(confs, *i)

	// }

	confs := []entity.TokenUriHistories{}

	f := bson.M{}
	f["projectID"] = projectID
	f["action"] = entity.SENT
	
	opts := options.Find()
	cursor, err := r.DB.Collection(utils.COLLECTION_TOKEN_URI_HISTORIES).Find(context.TODO(), f, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &confs); err != nil {
		return nil, err
	}

	
	return confs, nil

}

func (r Repository) ListTokensIn(tokenIDs []string) ([]entity.TokenUri, error) {

	confs := []entity.TokenUri{}

	f := bson.M{}
	f["token_id"] = bson.M{"$in": tokenIDs}
	opts := options.Find()
	cursor, err := r.DB.Collection(entity.TokenUri{}.TableName()).Find(context.TODO(), f, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &confs); err != nil {
		return nil, err
	}

	
	return confs, nil

}