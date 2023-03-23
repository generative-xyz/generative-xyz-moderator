package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
)

func (r Repository) ListSubClollectionItem(filter *structure.BaseFilters, inscriptionIds []string) ([]*entity.ItemListing, error) {
	page := filter.Page
	pageSize := filter.Limit
	result := []entity.DexVolumeInscriptionSumary{}
	pipeline := bson.A{
		bson.M{"$match": bson.M{
			"metadata.inscription_id": bson.M{"$in": inscriptionIds},
		}},
		bson.M{
			"$group": bson.M{
				"_id": "$metadata.inscription_id",
				// "total_volume": bson.M{"$sum": "$amount"},
				// "volume_1h": bson.M{
				// 	"$sum": bson.M{
				// 		"$cond": bson.A{
				// 			bson.M{"$gte": bson.A{"$timestamp", time.Now().Add(-1 * time.Hour)}},
				// 			"$amount", 0,
				// 		},
				// 	},
				// },
				// "volume_1d": bson.M{
				// 	"$sum": bson.M{
				// 		"$cond": bson.A{
				// 			bson.M{"$gte": bson.A{"$timestamp", time.Now().AddDate(0, 0, -1)}},
				// 			"$amount", 0,
				// 		},
				// 	},
				// },
				"volume_7d": bson.M{
					"$sum": bson.M{
						"$cond": bson.A{
							bson.M{"$gte": bson.A{"$timestamp", time.Now().AddDate(0, 0, -7)}},
							"$amount", 0,
						},
					},
				},
				"volume_30d": bson.M{
					"$sum": bson.M{
						"$cond": bson.A{
							bson.M{"$gte": bson.A{"$timestamp", time.Now().AddDate(0, 0, -30)}},
							"$amount", 0,
						},
					},
				},
			},
		},
		bson.M{
			"$project": bson.M{
				"inscription_id": "$_id",
				// "total_volume":   1,
				// "volume_1h":      1,
				// "volume_1d":      1,
				"volume_7d":  1,
				"volume_30d": 1,
				"_id":        0,
			},
		},
		bson.M{"$skip": (page - 1) * pageSize},
		bson.M{"$limit": pageSize},
	}

	cursor, err := r.DB.Collection(entity.DexVolumeInscription{}.TableName()).Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}

	if err = cursor.All((context.TODO()), &result); err != nil {
		return nil, err
	}

	response := []*entity.ItemListing{}
	for _, r := range result {
		data := &entity.ItemListing{
			InscriptionId: r.InscriptionId,
			Image:         fmt.Sprintf("https://generativeexplorer.com/preview/%s", r.InscriptionId),
			VolumeOneHour: &entity.VolumneObject{
				Amount: fmt.Sprintf("%d", r.Volume1h),
			},
			VolumeOneDay: &entity.VolumneObject{
				Amount: fmt.Sprintf("%d", r.Volume1d),
			},
			VolumeOneWeek: &entity.VolumneObject{
				Amount: fmt.Sprintf("%d", r.Volume7d),
			},
			VolumeOneMonth: &entity.VolumneObject{
				Amount: fmt.Sprintf("%d", r.Volume30d),
			},
		}
		response = append(response, data)
	}
	return response, nil
}

func (r Repository) ListItemListingOnSale(filter *structure.BaseFilters) ([]*entity.ItemListing, error) {
	page := filter.Page
	pageSize := filter.Limit
	result := []entity.DexBTCListing{}
	pipeline := bson.A{
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "token_uri"},
					{"localField", "inscription_id"},
					{"foreignField", "token_id"},
					{"as", "token_info"},
				},
			},
		},
		bson.D{{"$match", bson.D{{"token_info", bson.A{}}}}},
		bson.M{"$skip": (page - 1) * pageSize},
		bson.M{"$limit": pageSize},
	}

	cursor, err := r.DB.Collection(entity.DexBTCListing{}.TableName()).Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}

	if err = cursor.All((context.TODO()), &result); err != nil {
		return nil, err
	}

	response := []*entity.ItemListing{}
	addresses := []string{}
	for _, r := range result {
		data := &entity.ItemListing{
			InscriptionId: r.InscriptionID,
			Image:         fmt.Sprintf("https://generativeexplorer.com/preview/%s", r.InscriptionID),
			SellerAddress: r.SellerAddress,
		}

		addresses = append(addresses, r.SellerAddress)
		response = append(response, data)
	}

	users, err := r.FindUserByAddresses(addresses)
	if err != nil {
		return nil, err
	}

	userMap := make(map[string]entity.Users)
	for _, u := range users {
		userMap[u.WalletAddressBTCTaproot] = u
	}

	for _, r := range response {
		if u, ok := userMap[r.SellerAddress]; ok {
			r.SellerDisplayName = u.DisplayName

		}
	}
	return response, nil
}

func (r Repository) FindListItemListing(filter *structure.BaseFilters) ([]*entity.ItemListing, error) {
	page := filter.Page
	pageSize := filter.Limit
	ignoreInscriptionIds := []string{"b7b65579e2dd556b83665d7a26ecb0259225dbec491a9888d4a9c1716a7f9733i0"}
	pipeline := bson.A{
		bson.M{
			"$lookup": bson.M{
				"from":         "token_uri",
				"localField":   "metadata.inscription_id",
				"foreignField": "token_id",
				"as":           "token_info",
			},
		},
		bson.M{"$match": bson.M{
			"token_info":              bson.A{},
			"metadata.inscription_id": bson.M{"$nin": ignoreInscriptionIds},
		}},
		bson.M{
			"$group": bson.M{
				"_id": "$metadata.inscription_id",
				// "total_volume": bson.M{"$sum": "$amount"},
				// "volume_1h": bson.M{
				// 	"$sum": bson.M{
				// 		"$cond": bson.A{
				// 			bson.M{"$gte": bson.A{"$timestamp", time.Now().Add(-1 * time.Hour)}},
				// 			"$amount", 0,
				// 		},
				// 	},
				// },
				// "volume_1d": bson.M{
				// 	"$sum": bson.M{
				// 		"$cond": bson.A{
				// 			bson.M{"$gte": bson.A{"$timestamp", time.Now().AddDate(0, 0, -1)}},
				// 			"$amount", 0,
				// 		},
				// 	},
				// },
				"volume_7d": bson.M{
					"$sum": bson.M{
						"$cond": bson.A{
							bson.M{"$gte": bson.A{"$timestamp", time.Now().AddDate(0, 0, -7)}},
							"$amount", 0,
						},
					},
				},
				"volume_30d": bson.M{
					"$sum": bson.M{
						"$cond": bson.A{
							bson.M{"$gte": bson.A{"$timestamp", time.Now().AddDate(0, 0, -30)}},
							"$amount", 0,
						},
					},
				},
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "dex_volume_inscription",
				"localField":   "_id",
				"foreignField": "metadata.inscription_id",
				"as":           "inscription_info",
			},
		},
		bson.M{
			"$project": bson.M{
				"inscription_id": "$_id",
				// "total_volume":   1,
				// "volume_1h":      1,
				// "volume_1d":      1,
				"volume_7d":  1,
				"volume_30d": 1,
				"dex_volume_inscription": bson.M{
					"$arrayElemAt": bson.A{
						"$inscription_info",
						0,
					},
				},
				"_id": 0,
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from": "dex_btc_listing",
				"let":  bson.M{"inscription_id": "$inscription_id"},
				"pipeline": bson.A{
					bson.M{"$match": bson.M{"$expr": bson.M{"$eq": bson.A{"$inscription_id", "$$inscription_id"}}}},
					bson.M{"$match": bson.M{"matched": true}},
				},
				"as": "dex_btc_listings",
			},
		},
		bson.M{"$sort": bson.M{"volume_30d": -1}},
		bson.M{"$skip": (page - 1) * pageSize},
		bson.M{"$limit": pageSize},
	}

	result := []entity.DexVolumeInscriptionSumary{}
	cursor, err := r.DB.Collection(entity.DexVolumeInscription{}.TableName()).Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}

	if err = cursor.All((context.TODO()), &result); err != nil {
		return nil, err
	}

	response := []*entity.ItemListing{}
	addresses := []string{}
	for _, r := range result {
		data := &entity.ItemListing{
			InscriptionId: r.InscriptionId,
			Image:         fmt.Sprintf("https://generativeexplorer.com/preview/%s", r.InscriptionId),
			VolumeOneHour: &entity.VolumneObject{
				Amount: fmt.Sprintf("%d", r.Volume1h),
			},
			VolumeOneDay: &entity.VolumneObject{
				Amount: fmt.Sprintf("%d", r.Volume1d),
			},
			VolumeOneWeek: &entity.VolumneObject{
				Amount: fmt.Sprintf("%d", r.Volume7d),
			},
			VolumeOneMonth: &entity.VolumneObject{
				Amount: fmt.Sprintf("%d", r.Volume30d),
			},
		}
		for _, d := range r.DexBTCListings {
			if d.Matched && d.SellerAddress != "" {
				data.SellerAddress = d.SellerAddress
				data.BuyerAddress = d.Buyer
				addresses = append(addresses, d.SellerAddress, d.Buyer)
			}
		}
		response = append(response, data)
	}

	users, err := r.FindUserByAddresses(addresses)
	if err != nil {
		return nil, err
	}

	userMap := make(map[string]entity.Users)
	for _, u := range users {
		userMap[u.WalletAddressBTCTaproot] = u
	}

	for _, r := range response {
		if u, ok := userMap[r.SellerAddress]; ok {
			r.SellerDisplayName = u.DisplayName
		}

		if u, ok := userMap[r.BuyerAddress]; ok {
			r.BuyerDisplayName = u.DisplayName
		}
	}
	return response, nil
}

func (r Repository) FindDexVolumeInscription(filter *structure.DexVolumeInscritionFilter) ([]entity.DexVolumeInscription, error) {
	return nil, nil
}

func (r Repository) InsertDexVolumeInscription(data *entity.DexVolumeInscription) error {
	if data == nil {
		return errors.New("insertDexVolumeInscription Invalid data")
	}
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) AggregateVolumnCollection(filter *entity.AggerateChartForProject) ([]entity.AggragetedProject, error) {
	f := bson.A{
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "token_uri"},
					{"localField", "metadata.inscription_id"},
					{"foreignField", "token_id"},
					{"as", "token_uri"},
				},
			},
		},
		bson.D{
			{"$unwind",
				bson.D{
					{"path", "$token_uri"},
					{"preserveNullAndEmptyArrays", false},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "projects"},
					{"localField", "token_uri.project_id"},
					{"foreignField", "tokenid"},
					{"as", "projects"},
				},
			},
		},
		bson.D{
			{"$unwind",
				bson.D{
					{"path", "$projects"},
					{"preserveNullAndEmptyArrays", false},
				},
			},
		},
		bson.D{
			{"$match",
				bson.D{
					{"token_uri.project_id", bson.D{{"$eq", filter.ProjectID}}},
					{"timestamp", bson.D{{"$gte", filter.FromDate}}},
					{"timestamp", bson.D{{"$lte", filter.ToDate}}},
				},
			},
		},
		bson.D{
			{"$group",
				bson.D{
					{"_id",
						bson.D{
							{"projectID", "$projects.tokenid"},
							{"projectName", "$projects.name"},
							{"timestamp",
								bson.D{
									{"$dateToString",
										bson.D{
											{"format", "%Y-%m-%d"},
											{"date", "$timestamp"},
										},
									},
								},
							},
						},
					},
					{"amount", bson.D{{"$sum", "$amount"}}},
				},
			},
		},
		bson.D{
			{"$sort",
				bson.D{
					{"_id.timestamp", 1},
					{"amount", -1},
				},
			},
		},
	}

	cursor, err := r.DB.Collection(entity.DexVolumeInscription{}.TableName()).Aggregate(context.TODO(), f)
	if err != nil {
		return nil, err
	}

	result := []entity.AggragetedProject{}
	if err = cursor.All((context.TODO()), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) AggregateVolumnToken(filter *entity.AggerateChartForToken) ([]entity.AggragetedToken, error) {
	f := bson.A{
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "token_uri"},
					{"localField", "metadata.inscription_id"},
					{"foreignField", "token_id"},
					{"as", "token_uri"},
				},
			},
		},
		bson.D{
			{"$unwind",
				bson.D{
					{"path", "$token_uri"},
					{"preserveNullAndEmptyArrays", false},
				},
			},
		},
		bson.D{
			{"$match",
				bson.D{
					{"token_uri.token_id", filter.TokenID},
					{"timestamp", bson.D{{"$gte", filter.FromDate}}},
					{"timestamp", bson.D{{"$lte", filter.ToDate}}},
				},
			},
		},
		bson.D{
			{"$group",
				bson.D{
					{"_id",
						bson.D{
							{"tokenID", "$token_uri.token_id"},
							{"timestamp",
								bson.D{
									{"$dateToString",
										bson.D{
											{"format", "%Y-%m-%d"},
											{"date", "$timestamp"},
										},
									},
								},
							},
						},
					},
					{"amount", bson.D{{"$sum", "$amount"}}},
				},
			},
		},
		bson.D{
			{"$sort",
				bson.D{
					{"_id.timestamp", 1},
					{"amount", -1},
				},
			},
		},
	}

	cursor, err := r.DB.Collection(entity.DexVolumeInscription{}.TableName()).Aggregate(context.TODO(), f)
	if err != nil {
		return nil, err
	}

	result := []entity.AggragetedToken{}
	if err = cursor.All((context.TODO()), &result); err != nil {
		return nil, err
	}

	return result, nil
}
