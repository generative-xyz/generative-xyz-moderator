package repository

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)


func (r Repository) AggregateVolumn(payType string) ([]entity.AggregateProjectItemResp, error) {
	//resp := &entity.AggregateWalletAddres{}
	confs := []entity.AggregateProjectItemResp{}

	calculate := bson.M{"$sum": "$project_mint_price"}
	if payType == string(entity.ETH) {
		calculate = bson.M{"$sum": bson.M{
			"$multiply": bson.A{
				"$project_mint_price",
				 bson.M{ "$divide": bson.A{
					"$btc_rate",
					"$eth_rate",
				 }},
			},
		}}
	}

	// PayType *string
	// ReferreeIDs []string
	matchStage := bson.M{"$match": bson.M{"$and": bson.A{
		bson.M{"status": entity.StatusMint_SentFundToMaster},
		bson.M{"payType": payType},
		bson.M{"$and": bson.A{
				bson.M{"eth_rate": bson.M{"$gt": 0}},
				bson.M{"btc_rate": bson.M{"$gt": 0}},
			},
		},
	}}}

	pipeLine := bson.A{
		matchStage,
		bson.M{"$group": bson.M{"_id": 
			bson.M{ "projectID": "$projectID", "mintPrice": "$project_mint_price", "payType": "$payType" }, 
			"amount": calculate,
			"minted": bson.M{"$sum": 1},
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
		res := &entity.AggregateProjectItem{}
		err = helpers.Transform(item, res)
		if err != nil {
			return nil, err
		}
		tmp := entity.AggregateProjectItemResp{
			ProjectID: res.ID.ProjectID,
			Paytype: res.ID.Paytype,
			BtcRate: res.ID.BtcRate,
			EthRate: res.ID.EthRate,
			MintPrice: res.ID.MintPrice,
			Amount: res.Amount,
			Minted: res.Minted,
		}
		confs = append(confs, tmp)
	}
	
	return confs, nil
}

func (r Repository) AggregateAmount(filter entity.FilterVolume, groupStage bson.M) ([]entity.AggregateAmount, error) {
	//resp := &entity.AggregateWalletAddres{}
	confs := []entity.AggregateAmount{}

	f :=  bson.A{}
	
	if filter.AmountType != nil && *filter.AmountType != "" {
		f = append(f, bson.M{"payType": *filter.AmountType})
	}
	
	if filter.CreatorAddress != nil && *filter.CreatorAddress != "" {
		f = append(f, bson.M{"creatorAddress": *filter.CreatorAddress})
	}
	
	if filter.ProjectID != nil && *filter.ProjectID != "" {
		f = append(f, bson.M{"projectID": *filter.ProjectID})
	}
	
	if len(filter.ProjectIDs) > 0 {
		f = append(f, bson.M{"$in": bson.M{"projectID": filter.ProjectIDs } })
	}

	// PayType *string
	// ReferreeIDs []string
	matchStage := bson.M{"$match": bson.M{"$and": f}}

	pipeLine := bson.A{
		matchStage,
		groupStage,
		bson.M{"$sort": bson.M{"_id": -1}},
	}
	
	cursor, err := r.DB.Collection(entity.UserVolumn{}.TableName()).Aggregate(context.TODO(), pipeLine, nil)
	if err != nil {
		return nil, err
	}

	// display the results
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	for _, item := range results {
		
		tmp := &entity.AggregateAmount{}
		err = helpers.Transform(item, tmp)
		confs = append(confs, *tmp)
	}
	
	return confs, nil
}


func (r Repository) FindVolumn(projectID string, amountType string) (*entity.UserVolumn, error) {
	projectID = strings.ToLower(projectID)
	resp := &entity.UserVolumn{}
	usr, err := r.FilterOne(entity.UserVolumn{}.TableName(), bson.D{{"projectID", projectID}, {"payType", amountType}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r Repository) FindVolumnByWalletAddress(creatorAddress string, amountType string) (*entity.UserVolumn, error) {
	creatorAddress = strings.ToLower(creatorAddress)
	resp := &entity.UserVolumn{}
	usr, err := r.FilterOne(entity.UserVolumn{}.TableName(), bson.D{{"creatorAddress", creatorAddress}, {"payType", amountType}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r Repository) CreateVolumn(input *entity.UserVolumn) error {
	err := r.InsertOne(input.TableName(), input)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) UpdateVolumn(ID string, data *entity.UserVolumn) (*mongo.UpdateResult, error) {
	filter := bson.D{{utils.KEY_UUID, ID}}
	result, err := r.UpdateOne(entity.UserVolumn{}.TableName(), filter, data)
	if err != nil {
		return nil, err
	}
	return result, nil
}


func (r Repository) UpdateVolumnAmount(ID string, amount string,  earning string, gearning string) (*mongo.UpdateResult, error) {
	filter := bson.D{{utils.KEY_UUID, ID}}
	update := bson.M{"$set": bson.M{"amount": amount, "earning": earning, "genEarning":  gearning}}
	result, err := r.DB.Collection(entity.UserVolumn{}.TableName()).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r Repository) UpdateVolumnMinted(ID string, minted int) (*mongo.UpdateResult, error) {
	filter := bson.D{{utils.KEY_UUID, ID}}
	update := bson.M{"$set": bson.M{"minted": minted}}
	result, err := r.DB.Collection(entity.UserVolumn{}.TableName()).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}


func (r Repository) UpdateVolumMintPrice(ID string, mintPrice int64) (*mongo.UpdateResult, error) {
	filter := bson.D{{utils.KEY_UUID, ID}}
	update := bson.M{"$set": bson.M{"mintPrice": mintPrice}}
	result, err := r.DB.Collection(entity.UserVolumn{}.TableName()).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}
