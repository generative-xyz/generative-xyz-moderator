package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

func (r Repository) FindAirdropByStatus(status int) ([]*entity.Airdrop, error) {
	var resp []*entity.Airdrop
	filter := bson.D{{"status", status}}
	cursor, err := r.DB.Collection(utils.COLLECTION_AIRDROP).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (r Repository) FindAirdropByTokenGatedNewUser(userUUid string) (*entity.Airdrop, error) {
	resp := &entity.Airdrop{}
	filter := bson.D{{"type", 2}, {"receiver", userUUid}}
	cursor, err := r.FilterOne(entity.Airdrop{}.TableName(), filter)
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(cursor, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

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

func (r Repository) UpdateAirdropMintInfoByUUid(uuid string, ordinalResponseAction interface{}) (*mongo.UpdateResult, error) {
	filter := bson.D{{"uuid", uuid}}
	update := bson.M{"$set": bson.M{"status": 0, "ordinalResponseAction": ordinalResponseAction}}
	result, err := r.DB.Collection(utils.COLLECTION_AIRDROP).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) UpdateAirdropStatusByTx(tx string, status int, ordinalResponseAction string) (*mongo.UpdateResult, error) {
	filter := bson.D{{"tx", tx}}
	update := bson.M{"$set": bson.M{"status": status, "ordinalResponseAction": ordinalResponseAction}}
	result, err := r.DB.Collection(utils.COLLECTION_AIRDROP).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) UpdateAirdropInscriptionByUUid(uuid string, tx string, inscriptionId string) (*mongo.UpdateResult, error) {
	filter := bson.D{{"uuid", uuid}}
	update := bson.M{"$set": bson.M{"tx": tx, "inscriptionId": inscriptionId}}
	result, err := r.DB.Collection(utils.COLLECTION_AIRDROP).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}
