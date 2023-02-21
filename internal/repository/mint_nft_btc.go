package repository

import (
	"context"
	"log"
	"time"

	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r Repository) FindMintNftBtc(key string) (*entity.MintNftBtc, error) {
	resp := &entity.MintNftBtc{}
	usr, err := r.FilterOne(entity.MintNftBtc{}.TableName(), bson.D{{utils.KEY_UUID, key}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) FindMintNftBtcByNftID(uuid string) (*entity.MintNftBtcResp, error) {

	log.Println("uuid:", uuid)

	resp := &entity.MintNftBtcResp{}
	usr, err := r.FilterOne(entity.MintNftBtc{}.TableName(), bson.D{{"uuid", uuid}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) InsertMintNftBtc(data *entity.MintNftBtc) error {
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}

// new:
func (r Repository) ListMintNftBtcPending() ([]entity.MintNftBtc, error) {
	resp := []entity.MintNftBtc{}
	filter := bson.M{
		"status":     entity.StatusMint_Pending,
		"expired_at": bson.M{"$gte": primitive.NewDateTimeFromTime(time.Now().UTC())},
	}

	cursor, err := r.DB.Collection(utils.MINT_NFT_BTC).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
func (r Repository) ListMintNftBtcByStatus(statuses []entity.StatusMint) ([]entity.MintNftBtc, error) {
	resp := []entity.MintNftBtc{}
	filter := bson.M{
		"status": bson.M{"$in": statuses},
	}

	cursor, err := r.DB.Collection(utils.MINT_NFT_BTC).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (r Repository) UpdateMintNftBtc(model *entity.MintNftBtc) (*mongo.UpdateResult, error) {

	filter := bson.D{{Key: "uuid", Value: model.UUID}}
	result, err := r.UpdateOne(model.TableName(), filter, model)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) CreateMintNftBtcLog(logs *entity.MintNftBtcLogs) error {
	err := r.InsertOne(logs.TableName(), logs)
	if err != nil {
		return err
	}
	return nil
}
