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

func (r Repository) FindMintNftBtcByNftID(uuid string) (*entity.MintNftBtc, error) {

	log.Println("uuid:", uuid)

	resp := &entity.MintNftBtc{}
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
		"status":          bson.M{"$in": []entity.StatusMint{entity.StatusMint_Pending, entity.StatusMint_WaitingForConfirms}},
		"expired_at":      bson.M{"$gte": primitive.NewDateTimeFromTime(time.Now().UTC())},
		"patch_parent_id": 0, // only get parent items
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

func (r Repository) ListMintNftBtcToSendFundToMaster() ([]entity.MintNftBtc, error) {
	resp := []entity.MintNftBtc{}
	filter := bson.M{
		"status":          bson.M{"$in": []entity.StatusMint{entity.StatusMint_SentNFTToUser}},
		"patch_parent_id": 0, // only get parent items
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

// update seet filter
func (r Repository) UpdateMintNftBtcByFilter(uuid string, updateFilter bson.M) (*mongo.UpdateResult, error) {
	f := bson.D{
		{Key: "uuid", Value: uuid},
	}

	result, err := r.DB.Collection(utils.MINT_NFT_BTC).UpdateOne(context.TODO(), f, updateFilter)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (r Repository) CreateMintNftBtcLog(logs *entity.MintNftBtcLogs) error {
	err := r.InsertOne(logs.TableName(), logs)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) ListMintNftBtcByStatusAndAddress(address string, statuses []entity.StatusMint) ([]entity.MintNftBtc, error) {
	resp := []entity.MintNftBtc{}
	filter := bson.M{
		"origin_user_address": address,
		"status":              bson.M{"$in": statuses},
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

func (r Repository) UpdateTokenInscriptionIndexForMint(tokenId string, inscriptionIndex string) error {
	filter := bson.D{
		{Key: "token_id", Value: tokenId},
	}
	update := bson.M{
		"$set": bson.M{
			"inscription_index": inscriptionIndex,
		},
	}
	_, err := r.DB.Collection(utils.COLLECTION_TOKEN_URI).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return err
}

func (r Repository) UpdateCancelMintNftBtc(uuid string) error {
	filter := bson.D{
		{Key: "uuid", Value: uuid},
	}
	update := bson.M{
		"$set": bson.M{
			"status": -1,
		},
	}
	_, err := r.DB.Collection(utils.MINT_NFT_BTC).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return err
}

func (r Repository) CountBatchRecordOfItems(parentRecord string) ([]entity.MintNftBtc, error) {
	resp := []entity.MintNftBtc{}
	filter := bson.M{
		"patch_parent_id": parentRecord, // only get child items
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
