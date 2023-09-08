package repository

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

func (r Repository) AggregateVolumn(projectID string, payType string) ([]entity.AggregateProjectItemResp, error) {
	//resp := &entity.AggregateWalletAddres{}
	confs := []entity.AggregateProjectItemResp{}

	calculate := bson.M{"$sum": "$project_mint_price"}
	if payType == string(entity.ETH) {
		calculate = bson.M{"$sum": bson.M{
			"$multiply": bson.A{
				"$project_mint_price",
				bson.M{"$divide": bson.A{
					"$btc_rate",
					"$eth_rate",
				}},
			},
		}}
	}

	// PayType *string
	// ReferreeIDs []string
	matchStage := bson.M{"$match": bson.M{"$and": bson.A{
		bson.M{"isMinted": true},
		bson.M{"payType": payType},
		bson.M{"projectID": projectID},
	}}}

	pipeLine := bson.A{
		matchStage,
		bson.M{"$group": bson.M{"_id": bson.M{"projectID": "$projectID", "payType": "$payType"},
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
			Paytype:   res.ID.Paytype,
			BtcRate:   res.ID.BtcRate,
			EthRate:   res.ID.EthRate,
			MintPrice: res.ID.MintPrice,
			Amount:    res.Amount,
			Minted:    res.Minted,
		}
		confs = append(confs, tmp)
	}

	return confs, nil
}

func (r Repository) AggregateProjectMintPrice(projectID string, payType string) ([]entity.AggregateProjectItemResp, error) {
	//resp := &entity.AggregateWalletAddres{}
	confs := []entity.AggregateProjectItemResp{}

	calculate := bson.M{"$sum": "$project_mint_price"}
	if payType == string(entity.ETH) {
		calculate = bson.M{"$sum": bson.M{
			"$multiply": bson.A{
				"$project_mint_price",
				bson.M{"$divide": bson.A{
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
		bson.M{"projectID": projectID},
	}}}

	pipeLine := bson.A{
		matchStage,
		bson.M{"$group": bson.M{"_id": bson.M{"projectID": "$projectID", "payType": "$payType"},
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
			Paytype:   res.ID.Paytype,
			BtcRate:   res.ID.BtcRate,
			EthRate:   res.ID.EthRate,
			MintPrice: res.ID.MintPrice,
			Amount:    res.Amount,
			Minted:    res.Minted,
		}
		confs = append(confs, tmp)
	}

	return confs, nil
}

func (r Repository) AggregateAmount(filter entity.FilterVolume, groupStage bson.M) ([]entity.AggregateAmount, error) {
	//resp := &entity.AggregateWalletAddres{}
	confs := []entity.AggregateAmount{}

	f := bson.A{}

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
		f = append(f, bson.M{"$in": bson.M{"projectID": filter.ProjectIDs}})
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

func (r Repository) UpdateVolumnAmount(projectID string, payType string, amount string, earning string, gearning string) (*mongo.UpdateResult, error) {
	filter := bson.D{
		{Key: "projectID", Value: projectID},
		{Key: "payType", Value: payType},
	}

	update := bson.M{"$set": bson.M{"amount": amount, "earning": earning, "genEarning": gearning}}
	result, err := r.DB.Collection(entity.UserVolumn{}.TableName()).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r Repository) UpdateVolumnMinted(projectID string, payType string, minted int) (*mongo.UpdateResult, error) {

	filter := bson.D{
		{Key: "projectID", Value: projectID},
		{Key: "payType", Value: payType},
	}

	update := bson.M{"$set": bson.M{"minted": minted}}
	result, err := r.DB.Collection(entity.UserVolumn{}.TableName()).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r Repository) UpdateVolumMintPrice(projectID string, payType string, mintPrice int64) (*mongo.UpdateResult, error) {
	filter := bson.D{
		{Key: "projectID", Value: projectID},
		{Key: "payType", Value: payType},
	}

	update := bson.M{"$set": bson.M{"mintPrice": mintPrice}}
	result, err := r.DB.Collection(entity.UserVolumn{}.TableName()).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r Repository) AggregateBTCVolumn(projectID string) ([]entity.AggregateProjectItemResp, error) {
	//resp := &entity.AggregateWalletAddres{}
	confs := []entity.AggregateProjectItemResp{}

	calculate := bson.M{"$sum": "$project_mint_price"}
	// PayType *string
	// ReferreeIDs []string
	matchStage := bson.M{"$match": bson.M{"$and": bson.A{
		bson.M{"isMinted": true},
		bson.M{"projectID": projectID},
	}}}

	pipeLine := bson.A{
		matchStage,
		bson.M{"$group": bson.M{"_id": bson.M{"projectID": "$projectID"},
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
			Paytype:   res.ID.Paytype,
			BtcRate:   res.ID.BtcRate,
			EthRate:   res.ID.EthRate,
			MintPrice: res.ID.MintPrice,
			Amount:    res.Amount,
			Minted:    res.Minted,
		}

		confs = append(confs, tmp)
	}

	return confs, nil
}

func (r Repository) AggregateUsersVolumn(walletAddress []string) ([]entity.ReportArtist, error) {
	//resp := &entity.AggregateWalletAddres{}
	confs := []entity.ReportArtist{}

	f := bson.A{
		bson.D{
			{"$match",
				bson.D{
					{"creatorAddress",
						bson.D{
							{"$in", walletAddress},
						},
					},
				},
			},
		},
		bson.D{
			{"$group",
				bson.D{
					{"_id",
						bson.D{
							{"creatorAddress", "$creatorAddress"},
							{"payType", "$payType"},
						},
					},
					{"amount", bson.D{{"$sum", bson.D{{"$toDouble", "$amount"}}}}},
				},
			},
		},
		bson.D{
			{"$addFields",
				bson.D{
					{"payType", "$_id.payType"},
					{"walletAddress", "$_id.creatorAddress"},
				},
			},
		},
	}

	c, err := r.DB.Collection(utils.COLLECTION_USER_VOLUMN).Aggregate(context.Background(), f)
	if err != nil {
		return nil, err
	}

	err = c.All(context.Background(), &confs)
	if err != nil {
		return nil, err
	}

	return confs, nil
}

func (r Repository) AggregateMinterVolumn(walletAddress []string) ([]entity.ReportArtist, error) {
	//resp := &entity.AggregateWalletAddres{}
	confs := []entity.ReportArtist{}

	f := bson.A{
		bson.D{
			{"$match",
				bson.D{
					{"isMinted", true},
					{"origin_user_address", bson.D{{"$in", walletAddress}}},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"origin_user_address", 1},
					{"amount", 1},
					{"payType", 1},
				},
			},
		},
		bson.D{
			{"$group",
				bson.D{
					{"_id",
						bson.D{
							{"wallet_address", "$origin_user_address"},
							{"payType", "$payType"},
						},
					},
					{"amount", bson.D{{"$sum", bson.D{{"$toDouble", "$amount"}}}}},
				},
			},
		},
		bson.D{
			{"$addFields",
				bson.D{
					{"walletAddress", "$_id.wallet_address"},
					{"payType", "$_id.payType"},
				},
			},
		},
	}

	c, err := r.DB.Collection(utils.MINT_NFT_BTC).Aggregate(context.Background(), f)
	if err != nil {
		return nil, err
	}

	err = c.All(context.Background(), &confs)
	if err != nil {
		return nil, err
	}

	return confs, nil
}

func (r Repository) AggregateMinterBTCVolumnOld(walletAddress []string) ([]entity.ReportArtist, error) {
	//resp := &entity.AggregateWalletAddres{}
	confs := []entity.ReportArtist{}

	f := bson.A{
		bson.D{{"$match", bson.D{{"mintResponse.inscription", bson.D{{"$ne", ""}}}}}},
		bson.D{{"$sort", bson.D{{"created_at", 1}}}},
		bson.D{
			{"$project",
				bson.D{
					{"amount", 1},
					{"origin_user_address",
						bson.D{
							{"$ifNull",
								bson.A{
									"$origin_user_address",
									"$user_address",
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$addFields",
				bson.D{
					{"walletAddress", "$origin_user_address"},
					{"payType", "btc"},
				},
			},
		},
		bson.D{{"$match", bson.D{{"walletAddress", bson.D{{"$in", walletAddress}}}}}},
		bson.D{
			{"$group",
				bson.D{
					{"_id",
						bson.D{
							{"walletAddress", "$walletAddress"},
							{"payType", "$payType"},
						},
					},
					{"amount", bson.D{{"$sum", bson.D{{"$toDouble", "$amount"}}}}},
				},
			},
		},
		bson.D{
			{"$addFields",
				bson.D{
					{"walletAddress", "$_id.walletAddress"},
					{"payType", "$_id.payType"},
				},
			},
		},
		bson.D{{"$project", bson.D{{"_id", 0}}}},
	}

	c, err := r.DB.Collection(utils.COLLECTION_BTC_WALLET_ADDRESS).Aggregate(context.Background(), f)
	if err != nil {
		return nil, err
	}

	err = c.All(context.Background(), &confs)
	if err != nil {
		return nil, err
	}

	return confs, nil
}

func (r Repository) AggregateMinterEthVolumnOld(walletAddress []string) ([]entity.ReportArtist, error) {
	//resp := &entity.AggregateWalletAddres{}
	confs := []entity.ReportArtist{}

	f := bson.A{
		bson.D{{"$match", bson.D{{"mintResponse.inscription", bson.D{{"$ne", ""}}}}}},
		bson.D{{"$sort", bson.D{{"created_at", 1}}}},
		bson.D{
			{"$project",
				bson.D{
					{"amount", 1},
					{"origin_user_address",
						bson.D{
							{"$ifNull",
								bson.A{
									"$origin_user_address",
									"$user_address",
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$addFields",
				bson.D{
					{"walletAddress", "$origin_user_address"},
					{"payType", "btc"},
				},
			},
		},
		bson.D{
			{"$match",
				bson.D{
					{"walletAddress",
						bson.D{
							{"$in", walletAddress},
						},
					},
				},
			},
		},
		bson.D{
			{"$group",
				bson.D{
					{"_id",
						bson.D{
							{"walletAddress", "$walletAddress"},
							{"payType", "$payType"},
						},
					},
					{"amount", bson.D{{"$sum", bson.D{{"$toDouble", "$amount"}}}}},
				},
			},
		},
		bson.D{
			{"$addFields",
				bson.D{
					{"walletAddress", "$_id.walletAddress"},
					{"payType", "$_id.payType"},
				},
			},
		},
		bson.D{{"$project", bson.D{{"_id", 0}}}},
	}

	c, err := r.DB.Collection(utils.COLLECTION_ETH_WALLET_ADDRESS).Aggregate(context.Background(), f)
	if err != nil {
		return nil, err
	}

	err = c.All(context.Background(), &confs)
	if err != nil {
		return nil, err
	}

	return confs, nil
}

// userType: buyer or seller
func (r Repository) AggregateBuyer2ndSaleVolumn(walletAddress []string, userType string) ([]*entity.Report2ndSale, error) {
	//resp := &entity.AggregateWalletAddres{}
	confs := []*entity.Report2ndSale{}

	f := bson.A{
		bson.D{
			{"$match",
				bson.D{
					{userType,
						bson.D{
							{"$in", walletAddress},
						},
					},
					{"matched", true},
					{"cancelled", false},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"amount", 1},
					{userType, 1},
				},
			},
		},
		bson.D{
			{"$group",
				bson.D{
					{"_id", "$" + userType},
					{"total_amount", bson.D{{"$sum", "$amount"}}},
				},
			},
		},
		bson.D{{"$addFields", bson.D{{"walletAddressBtc", "$_id"}}}},
	}

	c, err := r.DB.Collection(utils.COLLECTION_DEX_BTC_LISTING).Aggregate(context.Background(), f)
	if err != nil {
		return nil, err
	}

	err = c.All(context.Background(), &confs)
	if err != nil {
		return nil, err
	}

	if userType != "buyer" && len(confs) > 0 {
		spew.Dump(1)
	}
	return confs, nil
}
