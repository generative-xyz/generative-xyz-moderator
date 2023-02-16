package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

func (r Repository) GetInscribeInfo(id string) (*entity.InscribeInfo, error) {
	resp := &entity.InscribeInfo{}

	f := bson.D{
		{Key: "id", Value: id},
	}

	inscribeInfo, err := r.FilterOne(utils.INSCRIBE_INFO, f, &options.FindOneOptions{
		Sort: bson.D{{Key: "created_at", Value: -1}},
	})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(inscribeInfo, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) CreateInscribeInfo(inscribeInfo *entity.InscribeInfo) error {
	err := r.InsertOne(inscribeInfo.TableName(), inscribeInfo)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) UpsertTokenUri(id string, inputData *entity.InscribeInfo) (*mongo.UpdateResult, error) {
	inputData.SetUpdatedAt()
	inputData.SetCreatedAt()
	bData, _ := inputData.ToBson()
	filter := bson.D{{"id", id}}
	update := bson.D{{"$set", bData}}
	updateOpts := options.Update().SetUpsert(true)
	//indexOpts := options.CreateIndexes().SetMaxTime(10 * time.Second)

	//id := fmt.Sprintf("%s%s", contractAddress, tokenID)
	result, err := r.DB.Collection(inputData.TableName()).UpdateOne(context.TODO(), filter, update, updateOpts)
	if err != nil {
		return nil, err
	}

	return result, nil
}
