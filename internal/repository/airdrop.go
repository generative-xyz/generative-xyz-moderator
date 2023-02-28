package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

func (r Repository) FindAirdropByTx(tx string) (*entity.Airdrop, error) {
	resp := &entity.Airdrop{}
	usr, err := r.FilterOne(entity.Airdrop{}.TableName(), bson.D{{"tx", tx}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) InsertAirdrop(data *entity.Airdrop) error {
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) UpdateAirdropByTx(tx string, data *entity.Airdrop) (*mongo.UpdateResult, error) {
	filter := bson.D{{"tx", tx}}
	result, err := r.UpdateOne(data.TableName(), filter, data)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) UpdateAirdropStatusByTx(tx string, data *entity.Airdrop, status int) (*mongo.UpdateResult, error) {
	filter := bson.D{{"tx", tx}}
	update := bson.M{"$set": bson.M{"status": status}}
	result, err := r.DB.Collection(utils.COLLECTION_AIRDROP).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}
