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
		"status":     bson.M{"$in": []entity.StatusMint{entity.StatusMint_Pending, entity.StatusMint_WaitingForConfirms}},
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
	_, err := r.DB.Collection(utils.COLLECTION_TOKEN_URI).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return err
}
type  AggregateWalletAddressItem struct {
	ID AggregateItemID `bson:"_id" json:"id"`
	Amount float64 `bson:"amount" json:"amount"`
}

type  AggregateItemID struct {
	ProjectID string `bson:"projectID" json:"projectID"`
	Paytype string `bson:"payType" json:"payType"`
}

//analytis
type  AggregateWalletAddres struct {
	Items []AggregateWalletAddressItem  `json:"items"`
	TotalBTC float64 `json:"totalAmountBTC"`
	TotalETH float64 `json:"totalAmountETH"`
}

func (r Repository) VolumeByProjectIDs(projectIDs []string, amountType string) (*AggregateWalletAddres, error) {
	resp := &AggregateWalletAddres{}
	totalBTC := 0.0
	totalETH := 0.0
	confs := []AggregateWalletAddressItem{}
	pipeLine := bson.A{
		bson.M{"$match": bson.M{"$and": bson.A{
			bson.M{"projectID": bson.M{"$in": projectIDs}},
			bson.M{"status": entity.StatusMint_SentFundToMaster},
		}}},
		bson.M{"$group": bson.M{"_id": 
			bson.M{ "projectID": "$projectID", "payType": "$payType" }, 
			"amount": bson.M{"$sum": bson.M{"$toDouble": "$amount"}},
		}},
		bson.M{"$sort": bson.M{"_id": -1}},
	}
	
	cursor, err := r.DB.Collection(entity.MintNftBtc{}.TableName()).Aggregate(context.TODO(), pipeLine, nil)
	if err != nil {
		return nil, err
	}

	// display the results
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	for _, item := range results {
		res := &AggregateWalletAddressItem{}
		err = helpers.Transform(item, res)
		if err != nil {
			return nil, err
		}
		confs = append(confs, *res)
		if res.ID.Paytype  == string(entity.BIT) {
			totalBTC += res.Amount
		}else{
			totalETH += res.Amount
		}
		
	}
	resp.TotalBTC = totalBTC
	resp.TotalETH = totalETH
	resp.Items = confs
	return resp, nil
}

