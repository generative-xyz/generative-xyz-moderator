package repository

import (
	"context"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r Repository) FindTokenUrisWithoutCache(f bson.M) ([]entity.TokenUri, error) {
	tokens := []entity.TokenUri{}

	//f[utils.KEY_DELETED_AT] = nil
	cursor, err := r.DB.Collection(utils.COLLECTION_TOKEN_URI).Find(context.TODO(), f)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &tokens); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (r Repository) FindTokenUriWithoutCache(filter bson.D) (*entity.TokenUri, error) {
	resp := &entity.TokenUri{}
	usr, err := r.FilterOne(entity.TokenUri{}.TableName(), filter)
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) FindTokenUriWithtCache(filter bson.D, cachedKey string) (*entity.TokenUri, error) {
	// always reload cache
	liveReload := func(filter bson.D, key string) (*entity.TokenUri, error) {
		token, err := r.FindTokenUriWithoutCache(filter)
		if err != nil {
			return nil, err
		}
		r.Cache.SetData(key, token)
		return token, nil
	}

	go liveReload(filter, cachedKey)

	cached, err := r.Cache.GetData(cachedKey)
	if err != nil || cached == nil {
		return liveReload(filter, cachedKey)
	}

	tok := &entity.TokenUri{}
	err = helpers.ParseCache(cached, tok)
	if err != nil {
		return nil, err
	}

	return tok, nil
}
func (r Repository) FindTokenByTokenID(tokenID string) (*entity.TokenUri, error) {
	f := bson.D{{"token_id", tokenID}}

	return r.FindTokenUriWithoutCache(f)
}

func (r Repository) FindTokenByTokenIDCustomField(tokenID string, fields []string) (*entity.TokenUri, error) {
	projectField := bson.D{
		{"_id", 1},
	}
	for _, field := range fields {
		projectField = append(projectField, bson.E{Key: field, Value: 1})
	}

	aggregates := bson.A{
		bson.D{
			{Key: "$project",
				Value: projectField,
			},
		},
		bson.D{{Key: "$match", Value: bson.D{{Key: "token_id", Value: tokenID}}}},
	}
	cursor, err := r.DB.Collection(entity.TokenUri{}.TableName()).Aggregate(context.TODO(), aggregates)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	tokenList := []entity.TokenUri{}

	if err = cursor.All((context.TODO()), &tokenList); err != nil {
		return nil, errors.WithStack(err)
	}
	if len(tokenList) > 0 {
		return &tokenList[0], nil
	}
	return nil, errors.New("token_id not found")
}

// func (r Repository) GetDexBtcsAlongWithProjectInfo(req entity.GetDexBtcListingWithProjectInfoReq) ([]entity.DexBtcListingWithProjectInfo, error) {

// }

func (r Repository) FindTokenBy(contractAddress string, tokenID string) (*entity.TokenUri, error) {
	key := helpers.TokenURIKey(contractAddress, tokenID)
	f := bson.D{{"contract_address", contractAddress}, {"token_id", tokenID}}

	return r.FindTokenUriWithtCache(f, key)
}

func (r Repository) FindTokenByWithoutCache(contractAddress string, tokenID string) (*entity.TokenUri, error) {
	f := bson.D{{"contract_address", contractAddress}, {"token_id", tokenID}}

	return r.FindTokenUriWithoutCache(f)
}

func (r Repository) FindTokenByGenNftAddr(genNftAddrr string, tokenID string) (*entity.TokenUri, error) {
	key := helpers.TokenURIByGenNftAddrKey(genNftAddrr, tokenID)
	f := bson.D{{"gen_nft_addrress", genNftAddrr}, {"token_id", tokenID}}
	return r.FindTokenUriWithtCache(f, key)
}

func (r Repository) FilterTokenUri(filter entity.FilterTokenUris) (*entity.Pagination, error) {
	tokens := []entity.TokenUri{}
	resp := &entity.Pagination{}

	f := r.filterToken(filter)
	if filter.SortBy == "" {
		filter.SortBy = "minted_time"
		filter.Sort = entity.SORT_DESC
	}

	if len(filter.Ids) != 0 {
		objectIDs, err := utils.StringsToObjects(filter.Ids)
		if err == nil {
			f["_id"] = bson.M{"$in": objectIDs}
		}
	}

	t, err := r.Paginate(entity.TokenUri{}.TableName(), filter.Page, filter.Limit, f, r.SelectedTokenFields(), r.SortToken(filter), &tokens)
	if err != nil {
		return nil, err
	}

	resp.Result = tokens
	resp.Page = t.Pagination.Page
	resp.Total = t.Pagination.Total
	resp.PageSize = filter.Limit
	//resp.PageSize = filter.Limit
	return resp, nil
}

func (r Repository) FilterTokenUriNew(filter entity.FilterTokenUris) (*entity.Pagination, error) {
	tokens := []entity.TokenUriListingPage{}
	resp := &entity.Pagination{}

	f := r.filterToken(filter)
	if filter.SortBy == "" {
		filter.SortBy = "priceBTC"
	}

	if len(filter.Ids) != 0 {
		objectIDs, err := utils.StringsToObjects(filter.Ids)
		if err == nil {
			f["_id"] = bson.M{"$in": objectIDs}
		}
	}

	listingAmountDefault := -1
	if filter.SortBy == "priceBTC" && filter.Sort == 1 {
		listingAmountDefault = 99999999999999
	}

	dexBtcMatch := bson.D{
		{"matched", false},
		{"cancelled", false},
	}

	addNoneBuyItems := true
	filterPendingBuyETH := bson.D{{"$gte", 0}}
	filterPendingBuyBTC := bson.D{{"$gte", 0}}
	if filter.IsBuynow != nil {
		if *filter.IsBuynow == true {
			addNoneBuyItems = false
		}
	}

	priceFilter := bson.A{}

	isFilterPrice := false
	if filter.FromPrice != nil {
		isFilterPrice = true
		priceFilter = append(priceFilter, bson.D{{"amount", bson.D{{"$gte", *filter.FromPrice}}}})
	}

	if filter.ToPrice != nil {
		isFilterPrice = true
		priceFilter = append(priceFilter, bson.D{{"amount", bson.D{{"$lte", *filter.ToPrice}}}})
	}

	if isFilterPrice {
		dexBtcMatchAnd := bson.E{"$and", priceFilter}
		dexBtcMatch = append(dexBtcMatch, dexBtcMatchAnd)
		addNoneBuyItems = false
	}
	//buyable only
	if !addNoneBuyItems {
		filterPendingBuyETH = bson.D{{"$eq", 0}}
		filterPendingBuyBTC = bson.D{{"$eq", 0}}
	}

	f2 := bson.A{
		bson.D{{"$match", f}},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "users"},
					{"localField", "owner_addrress"},
					{"foreignField", "wallet_address_btc_taproot"},
					{"as", "owner_object"},
					{"let",
						bson.D{
							{"wallet_address_btc_taproot", "$wallet_address_btc_taproot"},
						},
					},
					{"pipeline",
						bson.A{
							bson.D{
								{"$match",
									bson.D{
										{"wallet_address_btc_taproot", bson.D{{"$ne", ""}, {"$exists", true}}},
									},
								},
							},
						},
					},
				},
			},
		},

		bson.D{
			{"$unwind",
				bson.D{
					{"path", "$owner_object"},
					{"preserveNullAndEmptyArrays", true},
				},
			},
		},

		bson.D{
			{"$lookup",
				bson.D{
					{"from", "dex_btc_listing"},
					{"localField", "token_id"},
					{"foreignField", "inscription_id"},
					{"let",
						bson.D{
							{"cancelled", "$cancelled"},
							{"matched", "$matched"},
						},
					},
					{"pipeline",
						bson.A{
							bson.D{
								{"$match",
									dexBtcMatch,
								},
							},
						},
					},
					{"as", "listing"},
				},
			},
		},
		bson.D{
			{"$unwind",
				bson.D{
					{"path", "$listing"},
					{"preserveNullAndEmptyArrays", addNoneBuyItems},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "dex_btc_buy_eth"},
					{"let", bson.D{{"token_id", "$token_id"}}},
					{"pipeline",
						bson.A{
							bson.D{
								{"$addFields",
									bson.D{
										{"matched",
											bson.D{
												{"$eq",
													bson.A{
														"$inscription_id",
														"$$token_id",
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
										{"matched_multi",
											bson.D{
												{"$reduce",
													bson.D{
														{"input", "$inscription_list"},
														{"initialValue", false},
														{"in",
															bson.D{
																{"$cond",
																	bson.A{
																		bson.D{
																			{"$eq",
																				bson.A{
																					"$$this",
																					"$$token_id",
																				},
																			},
																		},
																		true,
																		"$$value",
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
							bson.D{
								{"$match",
									bson.D{
										{"$or",
											bson.A{
												bson.D{{"matched", true}},
												bson.D{{"matched_multi", true}},
											},
										},
										{"status",
											bson.D{
												{"$in",
													bson.A{
														1,
														2,
													},
												},
											},
										},
									},
								},
							},
							bson.D{{"$limit", 1}},
						},
					},
					{"as", "listing_eth"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "btc_tx_submit"},
					{"let", bson.D{{"token_id", "$token_id"}}},
					{"pipeline",
						bson.A{
							bson.D{
								{"$addFields",
									bson.D{
										{"matched_multi",
											bson.D{
												{"$reduce",
													bson.D{
														{"input", "$related_inscriptions"},
														{"initialValue", false},
														{"in",
															bson.D{
																{"$cond",
																	bson.A{
																		bson.D{
																			{"$eq",
																				bson.A{
																					"$$this",
																					"$$token_id",
																				},
																			},
																		},
																		true,
																		"$$value",
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
							bson.D{
								{"$match",
									bson.D{
										{"$or",
											bson.A{
												bson.D{{"matched_multi", true}},
											},
										},
										{"status",
											bson.D{
												{"$in",
													bson.A{
														0,
														1,
													},
												},
											},
										},
									},
								},
							},
							bson.D{{"$limit", 1}},
						},
					},
					{"as", "buying_btc"},
				},
			},
		},
		bson.D{
			{"$addFields",
				bson.D{
					{"buyable",
						bson.D{
							{"$cond",
								bson.D{
									{"if",
										bson.D{
											{"$or",
												bson.A{
													bson.D{
														{"$eq",
															bson.A{
																bson.D{
																	{"$ifNull",
																		bson.A{
																			"$listing",
																			0,
																		},
																	},
																},
																0,
															},
														},
													},
													bson.D{
														{"$gt",
															bson.A{
																bson.D{{"$size", "$listing_eth"}},
																0,
															},
														},
													},
													bson.D{
														{"$gt",
															bson.A{
																bson.D{{"$size", "$buying_btc"}},
																0,
															},
														},
													},
												},
											},
										},
									},
									{"then", false},
									{"else", true},
								},
							},
						},
					},
					{"priceBTC",
						bson.D{
							{"$cond",
								bson.D{
									{"if",
										bson.D{
											{"$or",
												bson.A{
													bson.D{
														{"$eq",
															bson.A{
																bson.D{
																	{"$ifNull",
																		bson.A{
																			"$listing",
																			0,
																		},
																	},
																},
																0,
															},
														},
													},
													bson.D{
														{"$gt",
															bson.A{
																bson.D{{"$size", "$listing_eth"}},
																0,
															},
														},
													},
													bson.D{
														{"$gt",
															bson.A{
																bson.D{{"$size", "$buying_btc"}},
																0,
															},
														},
													},
												},
											},
										},
									},
									{"then", listingAmountDefault},
									{"else", "$listing.amount"},
								},
							},
						},
					},
					{"orderID", "$listing._id"},
					{"sell_verified", "$listing.verified"},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"_id", 1},
					{"token_id", 1},
					{"name", 1},
					{"gen_nft_addrress", 1},
					{"contract_address", 1},
					{"project_id", 1},
					{"image", 1},
					{"priority", 1},
					{"inscription_index", 1},
					{"order_inscription_index", 1},
					{"token_id_int", 1},
					{"sell_verified", 1},
					{"thumbnail", 1},
					{"buyable", 1},
					{"priceBTC", 1},
					{"orderID", 1},
					{"nftTokenId", 1},
					{"project.tokenid", 1},
					{"project.royalty", 1},
					{"owner_addrress", 1},
					{"owner", 1},
					{"owner_object.wallet_address", 1},
					{"owner_object.wallet_address_btc_taproot", 1},
					{"owner_object.avatar", 1},
					{"owner_object.display_name", 1},
					{"listing_eth_size", bson.D{{"$size", "$listing_eth"}}},
					{"buying_btc_size", bson.D{{"$size", "$buying_btc"}}},
				},
			},
		},
		bson.D{
			{"$match",
				bson.D{
					{"$and",
						bson.A{
							bson.D{{"listing_eth_size", filterPendingBuyETH}},
							bson.D{{"buying_btc_size", filterPendingBuyBTC}},
						},
					},
				},
			},
		},
		bson.D{{"$sort", bson.D{{filter.SortBy, filter.Sort}, {"order_inscription_index", 1}, {"token_id_int", 1}}}},
		bson.D{
			{"$facet",
				bson.D{
					{"totalData",
						bson.A{
							bson.D{{"$match", bson.D{}}},
							bson.D{{"$skip", (filter.Page - 1) * filter.Limit}},
							bson.D{{"$limit", filter.Limit}},
						},
					},
					{"totalCount",
						bson.A{
							bson.D{{"$count", "count"}},
						},
					},
				},
			},
		},
	}

	// t, err := r.Aggregate(entity.TokenUri{}.TableName(), filter.Page, filter.Limit, f2, r.SelectedTokenFieldsNew(), r.SortToken(filter), &tokens)
	// if err != nil {
	// 	return nil, err
	// }
	cursor, err := r.DB.Collection(entity.TokenUri{}.TableName()).Aggregate(context.TODO(), f2)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err = cursor.All((context.TODO()), &tokens); err != nil {
		return nil, errors.WithStack(err)
	}
	if len(tokens) > 0 {
		//log.Println("len(tokens)", len(tokens[0].TotalCount))

		resp.Result = tokens[0].TotalData
		resp.Page = filter.Page
		if len(tokens[0].TotalCount) > 0 {
			resp.Total = tokens[0].TotalCount[0].Count
		}
		resp.TotalPage = resp.Total / filter.Limit
		resp.PageSize = filter.Limit
	}

	//resp.PageSize = filter.Limit
	return resp, nil
}

func (r Repository) FilterTokenUriTCNew(filter entity.FilterTokenUris) (*entity.Pagination, error) {
	tokens := []entity.TokenUriListingPage{}
	resp := &entity.Pagination{}

	f := r.filterToken(filter)
	if filter.SortBy == "" {
		filter.SortBy = "priceBTC"
	}

	if len(filter.Ids) != 0 {
		objectIDs, err := utils.StringsToObjects(filter.Ids)
		if err == nil {
			f["_id"] = bson.M{"$in": objectIDs}
		}
	}

	listingAmountDefault := -1
	if filter.SortBy == "priceBTC" && filter.Sort == 1 {
		listingAmountDefault = 99999999999999
	}

	dexBtcMatch := bson.D{
		{"matched", false},
		{"cancelled", false},
	}

	addNoneBuyItems := true
	filterPendingBuyETH := bson.D{{"$gte", 0}}
	filterPendingBuyBTC := bson.D{{"$gte", 0}}
	if filter.IsBuynow != nil {
		if *filter.IsBuynow == true {
			addNoneBuyItems = false
		}
	}

	priceFilter := bson.A{}

	isFilterPrice := false
	if filter.FromPrice != nil {
		isFilterPrice = true
		priceFilter = append(priceFilter, bson.D{{"amount", bson.D{{"$gte", *filter.FromPrice}}}})
	}

	if filter.ToPrice != nil {
		isFilterPrice = true
		priceFilter = append(priceFilter, bson.D{{"amount", bson.D{{"$lte", *filter.ToPrice}}}})
	}

	if isFilterPrice {
		dexBtcMatchAnd := bson.E{"$and", priceFilter}
		dexBtcMatch = append(dexBtcMatch, dexBtcMatchAnd)
		addNoneBuyItems = false
	}
	//buyable only
	if !addNoneBuyItems {
		filterPendingBuyETH = bson.D{{"$eq", 0}}
		filterPendingBuyBTC = bson.D{{"$eq", 0}}
	}

	f2 := bson.A{
		bson.D{{"$match", f}},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "users"},
					{"localField", "owner_addrress"},
					{"foreignField", "wallet_address_btc_taproot"},
					{"as", "owner_object"},
					{"let",
						bson.D{
							{"wallet_address_btc_taproot", "$wallet_address_btc_taproot"},
						},
					},
					{"pipeline",
						bson.A{
							bson.D{
								{"$match",
									bson.D{
										{"wallet_address_btc_taproot", bson.D{{"$ne", ""}, {"$exists", true}}},
									},
								},
							},
						},
					},
				},
			},
		},

		bson.D{
			{"$unwind",
				bson.D{
					{"path", "$owner_object"},
					{"preserveNullAndEmptyArrays", true},
				},
			},
		},

		bson.D{
			{"$lookup",
				bson.D{
					{"from", "dex_btc_listing"},
					{"localField", "token_id"},
					{"foreignField", "inscription_id"},
					{"let",
						bson.D{
							{"cancelled", "$cancelled"},
							{"matched", "$matched"},
						},
					},
					{"pipeline",
						bson.A{
							bson.D{
								{"$match",
									dexBtcMatch,
								},
							},
						},
					},
					{"as", "listing"},
				},
			},
		},
		bson.D{
			{"$unwind",
				bson.D{
					{"path", "$listing"},
					{"preserveNullAndEmptyArrays", addNoneBuyItems},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "dex_btc_buy_eth"},
					{"let", bson.D{{"token_id", "$token_id"}}},
					{"pipeline",
						bson.A{
							bson.D{
								{"$addFields",
									bson.D{
										{"matched",
											bson.D{
												{"$eq",
													bson.A{
														"$inscription_id",
														"$$token_id",
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
										{"matched_multi",
											bson.D{
												{"$reduce",
													bson.D{
														{"input", "$inscription_list"},
														{"initialValue", false},
														{"in",
															bson.D{
																{"$cond",
																	bson.A{
																		bson.D{
																			{"$eq",
																				bson.A{
																					"$$this",
																					"$$token_id",
																				},
																			},
																		},
																		true,
																		"$$value",
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
							bson.D{
								{"$match",
									bson.D{
										{"$or",
											bson.A{
												bson.D{{"matched", true}},
												bson.D{{"matched_multi", true}},
											},
										},
										{"status",
											bson.D{
												{"$in",
													bson.A{
														1,
														2,
													},
												},
											},
										},
									},
								},
							},
							bson.D{{"$limit", 1}},
						},
					},
					{"as", "listing_eth"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "btc_tx_submit"},
					{"let", bson.D{{"token_id", "$token_id"}}},
					{"pipeline",
						bson.A{
							bson.D{
								{"$addFields",
									bson.D{
										{"matched_multi",
											bson.D{
												{"$reduce",
													bson.D{
														{"input", "$related_inscriptions"},
														{"initialValue", false},
														{"in",
															bson.D{
																{"$cond",
																	bson.A{
																		bson.D{
																			{"$eq",
																				bson.A{
																					"$$this",
																					"$$token_id",
																				},
																			},
																		},
																		true,
																		"$$value",
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
							bson.D{
								{"$match",
									bson.D{
										{"$or",
											bson.A{
												bson.D{{"matched_multi", true}},
											},
										},
										{"status",
											bson.D{
												{"$in",
													bson.A{
														0,
														1,
													},
												},
											},
										},
									},
								},
							},
							bson.D{{"$limit", 1}},
						},
					},
					{"as", "buying_btc"},
				},
			},
		},
		bson.D{
			{"$addFields",
				bson.D{
					{"buyable",
						bson.D{
							{"$cond",
								bson.D{
									{"if",
										bson.D{
											{"$or",
												bson.A{
													bson.D{
														{"$eq",
															bson.A{
																bson.D{
																	{"$ifNull",
																		bson.A{
																			"$listing",
																			0,
																		},
																	},
																},
																0,
															},
														},
													},
													bson.D{
														{"$gt",
															bson.A{
																bson.D{{"$size", "$listing_eth"}},
																0,
															},
														},
													},
													bson.D{
														{"$gt",
															bson.A{
																bson.D{{"$size", "$buying_btc"}},
																0,
															},
														},
													},
												},
											},
										},
									},
									{"then", false},
									{"else", true},
								},
							},
						},
					},
					{"priceBTC",
						bson.D{
							{"$cond",
								bson.D{
									{"if",
										bson.D{
											{"$or",
												bson.A{
													bson.D{
														{"$eq",
															bson.A{
																bson.D{
																	{"$ifNull",
																		bson.A{
																			"$listing",
																			0,
																		},
																	},
																},
																0,
															},
														},
													},
													bson.D{
														{"$gt",
															bson.A{
																bson.D{{"$size", "$listing_eth"}},
																0,
															},
														},
													},
													bson.D{
														{"$gt",
															bson.A{
																bson.D{{"$size", "$buying_btc"}},
																0,
															},
														},
													},
												},
											},
										},
									},
									{"then", listingAmountDefault},
									{"else", "$listing.amount"},
								},
							},
						},
					},
					{"orderID", "$listing._id"},
					{"sell_verified", "$listing.verified"},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"_id", 1},
					{"token_id", 1},
					{"name", 1},
					{"gen_nft_addrress", 1},
					{"contract_address", 1},
					{"project_id", 1},
					{"image", 1},
					{"priority", 1},
					{"inscription_index", 1},
					{"order_inscription_index", 1},
					{"token_id_int", 1},
					{"sell_verified", 1},
					{"thumbnail", 1},
					{"buyable", 1},
					{"priceBTC", 1},
					{"orderID", 1},
					{"nftTokenId", 1},
					{"project.tokenid", 1},
					{"project.royalty", 1},
					{"owner_addrress", 1},
					{"owner", 1},
					{"owner_object.wallet_address", 1},
					{"owner_object.wallet_address_btc_taproot", 1},
					{"owner_object.avatar", 1},
					{"owner_object.display_name", 1},
					{"listing_eth_size", bson.D{{"$size", "$listing_eth"}}},
					{"buying_btc_size", bson.D{{"$size", "$buying_btc"}}},
				},
			},
		},
		bson.D{
			{"$match",
				bson.D{
					{"$and",
						bson.A{
							bson.D{{"listing_eth_size", filterPendingBuyETH}},
							bson.D{{"buying_btc_size", filterPendingBuyBTC}},
						},
					},
				},
			},
		},
		bson.D{{"$sort", bson.D{{filter.SortBy, filter.Sort}, {"order_inscription_index", 1}, {"token_id_int", 1}}}},
		bson.D{
			{"$facet",
				bson.D{
					{"totalData",
						bson.A{
							bson.D{{"$match", bson.D{}}},
							bson.D{{"$skip", (filter.Page - 1) * filter.Limit}},
							bson.D{{"$limit", filter.Limit}},
						},
					},
					{"totalCount",
						bson.A{
							bson.D{{"$count", "count"}},
						},
					},
				},
			},
		},
	}

	// t, err := r.Aggregate(entity.TokenUri{}.TableName(), filter.Page, filter.Limit, f2, r.SelectedTokenFieldsNew(), r.SortToken(filter), &tokens)
	// if err != nil {
	// 	return nil, err
	// }
	cursor, err := r.DB.Collection(entity.TokenUri{}.TableName()).Aggregate(context.TODO(), f2)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err = cursor.All((context.TODO()), &tokens); err != nil {
		return nil, errors.WithStack(err)
	}
	if len(tokens) > 0 {
		//log.Println("len(tokens)", len(tokens[0].TotalCount))

		resp.Result = tokens[0].TotalData
		resp.Page = filter.Page
		if len(tokens[0].TotalCount) > 0 {
			resp.Total = tokens[0].TotalCount[0].Count
		}
		resp.TotalPage = resp.Total / filter.Limit
		resp.PageSize = filter.Limit
	}

	//resp.PageSize = filter.Limit
	return resp, nil
}

func (r Repository) filterToken(filter entity.FilterTokenUris) bson.M {
	f := bson.M{}
	//f[utils.KEY_DELETED_AT] = nil

	if filter.CreatorAddr != nil {
		if *filter.CreatorAddr != "" {
			f["creator_address"] = *filter.CreatorAddr
		}
	}

	if filter.Keyword != nil {
		if *filter.Keyword != "" {
			kwInt, err := strconv.Atoi(*filter.Keyword)
			if err == nil {
				f["token_id_mini"] = kwInt
			} else {
				f["token_id_mini"] = *filter.Keyword
			}

		}
	}

	if filter.Search != nil {
		searchInt, err := strconv.Atoi(*filter.Search)
		if err == nil {
			f["$or"] = []bson.M{
				{"inscription_index": primitive.Regex{Pattern: *filter.Search, Options: "i"}},
				{"order_inscription_index": searchInt},
			}
		} else {
			f["$or"] = []bson.M{
				{"token_id": primitive.Regex{Pattern: *filter.Search, Options: "i"}},
				{"inscription_index": primitive.Regex{Pattern: *filter.Search, Options: "i"}},
			}

		}
	}

	if filter.OwnerAddr != nil {
		if *filter.OwnerAddr != "" {
			f["owner_addrress"] = *filter.OwnerAddr
		}
	}

	if filter.GenNFTAddr != nil {
		if *filter.GenNFTAddr != "" {
			f["gen_nft_addrress"] = *filter.GenNFTAddr
		}
	}

	if filter.ContractAddress != nil {
		if *filter.ContractAddress != "" {
			f["contract_address"] = *filter.ContractAddress
		}
	}

	if len(filter.CollectionIDs) > 0 {
		f["gen_nft_addrress"] = bson.D{{Key: "$in", Value: filter.CollectionIDs}}
	}

	if len(filter.TokenIDs) > 0 {
		f["token_id"] = bson.D{{Key: "$in", Value: filter.TokenIDs}}
	}

	// if filter.HasPrice != nil || filter.FromPrice != nil || filter.ToPrice != nil {
	// 	priceFilter := bson.M{}
	// 	if filter.HasPrice != nil {
	// 		priceFilter["$exists"] = *filter.HasPrice
	// 		priceFilter["$ne"] = nil
	// 	}
	// 	if filter.FromPrice != nil {
	// 		priceFilter["$gte"] = *filter.FromPrice
	// 	}
	// 	if filter.ToPrice != nil {
	// 		priceFilter["$lte"] = *filter.ToPrice
	// 	}
	// 	f["stats.price_int"] = priceFilter
	// }

	andFilters := []bson.M{
		f,
	}

	if filter.Attributes != nil && len(filter.Attributes) > 0 {
		for _, attribute := range filter.Attributes {
			andFilters = append(andFilters, bson.M{
				"parsed_attributes_str": bson.M{
					"$elemMatch": bson.M{
						"trait_type": attribute.TraitType,
						"value": bson.M{
							"$in": attribute.Values,
						},
					},
				},
			})
		}
	}

	if filter.RarityAttributes != nil && len(filter.RarityAttributes) > 0 {
		traits := []string{}
		values := []string{}
		for _, attribute := range filter.RarityAttributes {
			traits = append(traits, attribute.TraitType)
			for _, value := range attribute.Values {
				values = append(values, value)
			}
		}

		andFilters = append(andFilters, bson.M{
			"parsed_attributes_str": bson.M{
				"$elemMatch": bson.M{
					"trait_type": bson.M{
						"$in": traits,
					},
					"value": bson.M{
						"$in": values,
					},
				},
			},
		})

	}

	return bson.M{
		"$and": andFilters,
	}
}

func (r Repository) GetAllTokens() ([]entity.TokenUri, error) {
	tokens := []entity.TokenUri{}

	f := bson.M{}
	//f[utils.KEY_DELETED_AT] = nil
	cursor, err := r.DB.Collection(utils.COLLECTION_TOKEN_URI).Find(context.TODO(), f)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &tokens); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (r Repository) GetOwnerTokens(ownerAddress string) ([]entity.TokenUri, error) {
	tokens := []entity.TokenUri{}

	f := bson.D{{Key: "owner_addrress", Value: strings.ToLower(ownerAddress)}}

	cursor, err := r.DB.Collection(utils.COLLECTION_TOKEN_URI).Find(context.TODO(), f)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &tokens); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (r Repository) GetAllTokensSeletedFields() ([]entity.TokenUri, error) {
	tokens := []entity.TokenUri{}

	f := bson.M{}
	//f[utils.KEY_DELETED_AT] = nil
	opts := options.Find().SetProjection(r.SelectedTokenFields())
	cursor, err := r.DB.Collection(utils.COLLECTION_TOKEN_URI).Find(context.TODO(), f, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &tokens); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (r Repository) GetAllTokensHasTraitSeletedFields() ([]entity.TokenUri, error) {
	tokens := []entity.TokenUri{}

	f := bson.M{
		"parsed_attributes.0": bson.M{"$exists": true},
	}
	//f[utils.KEY_DELETED_AT] = nil
	opts := options.Find().SetProjection(r.SelectedTokenFields())
	cursor, err := r.DB.Collection(utils.COLLECTION_TOKEN_URI).Find(context.TODO(), f, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &tokens); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (r Repository) UpdateOrInsertTokenUri(contractAddress string, tokenID string, inputData *entity.TokenUri) (*mongo.UpdateResult, error) {
	inputData.SetUpdatedAt()
	inputData.SetCreatedAt()
	bData, _ := inputData.ToBson()
	filter := bson.D{{"contract_address", contractAddress}, {"token_id", tokenID}}
	update := bson.D{{"$set", bData}}
	updateOpts := options.Update().SetUpsert(true)
	//indexOpts := options.CreateIndexes().SetMaxTime(10 * time.Second)

	//id := fmt.Sprintf("%s%s", contractAddress, tokenID)
	result, err := r.DB.Collection(inputData.TableName()).UpdateOne(context.TODO(), filter, update, updateOpts)
	if err != nil {
		return nil, err
	}

	key1 := helpers.TokenURIKey(contractAddress, tokenID)
	key2 := helpers.TokenURIByGenNftAddrKey(inputData.GenNFTAddr, tokenID)

	r.Cache.SetData(key1, inputData)
	r.Cache.SetData(key2, inputData)
	return result, nil
}

func (r Repository) GetAllTokensByProjectID(projectID string) ([]entity.TokenUri, error) {
	tokens := []entity.TokenUri{}
	f := bson.D{{
		Key:   utils.KEY_PROJECT_ID,
		Value: projectID,
	}}

	cursor, err := r.DB.Collection(utils.COLLECTION_TOKEN_URI).Find(context.TODO(), f)
	if err != nil {
		return nil, err
	}

	if err = cursor.All((context.TODO()), &tokens); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (r Repository) GetAllTokenTraitsByProjectID(projectID string) ([]entity.AggregateTokenUriTraits, error) {
	tokens := []entity.AggregateTokenUriTraits{}
	matchStage := bson.D{{
		Key:   utils.KEY_PROJECT_ID,
		Value: projectID,
	}}

	pipeLine := bson.A{
		bson.D{
			{"$unwind", bson.D{
				{"path", "$parsed_attributes_str"},
				{"preserveNullAndEmptyArrays", true},
			}},
		},
		bson.D{
			{"$match", matchStage},
		},
		bson.D{
			{"$group",
				bson.D{
					{"_id",
						bson.D{
							{"project_id", "$projectID"},
							{"token_id", "$token_id"},
						},
					},
					{"parsed_attributes_str", bson.D{{"$push", "$parsed_attributes_str"}}},
					{"size", bson.D{{"$sum", 1}}},
				},
			},
		},
		bson.D{
			{"$sort", bson.M{"size": -1}},
		},
	}

	cursor, err := r.DB.Collection(utils.COLLECTION_TOKEN_URI).Aggregate(context.TODO(), pipeLine)
	if err != nil {
		return nil, err
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	//spew.Dump(results)

	for _, results := range results {
		i := &entity.AggregateTokenUriTraits{}
		err := helpers.Transform(results, i)
		if err != nil {
			continue
		}
		tokens = append(tokens, *i)

	}
	return tokens, nil
}

func (r Repository) SelectedTokenFields() bson.D {
	f := bson.D{
		{"token_id", 1},
		{"gen_nft_addrress", 1},
		{"contract_address", 1},
		{"thumbnail", 1},
		{"description", 1},
		{"name", 1},
		{"price", 1},
		{"owner_addrress", 1},
		{"creator_address", 1},
		{"project_id", 1},
		{"minted_time", 1},
		{"priority", 1},
		{"image", 1},
		{"project.tokenid", 1},
		{"project.tokenIDInt", 1},
		{"project.contractAddress", 1},
		{"project.name", 1},
		//{"project", 1},
		{"owner.wallet_address", 1},
		{"owner.display_name", 1},
		{"owner.avatar", 1},
		{"creator.wallet_address", 1},
		{"creator.display_name", 1},
		{"creator.avatar", 1},
		{"stats.price_int", 1},
		{"minter_address", 1},
		{"inscription_index", 1},
		{"order_inscription_index", 1},
	}
	return f
}

func (r Repository) SelectedTokenFieldsNew() bson.D {
	f := bson.D{
		{"token_id", 1},
		{"gen_nft_addrress", 1},
		{"contract_address", 1},
		{"thumbnail", 1},
		{"description", 1},
		{"name", 1},
		{"price", 1},
		{"owner_addrress", 1},
		{"creator_address", 1},
		{"project_id", 1},
		{"minted_time", 1},
		{"priority", 1},
		{"image", 1},
		{"project.tokenid", 1},
		{"project.tokenIDInt", 1},
		{"project.contractAddress", 1},
		{"project.name", 1},
		//{"project", 1},
		{"owner.wallet_address", 1},
		{"owner.display_name", 1},
		{"owner.avatar", 1},
		{"creator.wallet_address", 1},
		{"creator.display_name", 1},
		{"creator.avatar", 1},
		{"stats.price_int", 1},
		{"minter_address", 1},
		{"inscription_index", 1},
		{"order_inscription_index", 1},
		{"created_at", 1},
	}
	return f
}

func (r Repository) SortToken(filter entity.FilterTokenUris) []Sort {
	s := []Sort{}
	s = append(s, Sort{SortBy: filter.SortBy, Sort: filter.Sort})
	return s
}

func (r Repository) CreateToken(data *entity.TokenUri) error {
	data.SetCreatedAt()
	bData, err := data.ToBson()
	if err != nil {
		return err
	}

	inserted, err := r.DB.Collection(data.TableName()).InsertOne(context.TODO(), &bData)
	if err != nil {
		return err
	}

	_ = inserted
	return nil
}

func (r Repository) UpdateTokenPriceByTokenId(tokenId string, price int64) error {
	filter := bson.D{
		{Key: "token_id", Value: tokenId},
	}
	update := bson.M{
		"$set": bson.M{
			"stats.price_int": price,
		},
	}
	_, err := r.DB.Collection(utils.COLLECTION_TOKEN_URI).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return err
}

func (r Repository) UnsetTokenPriceByTokenId(tokenId string) error {
	filter := bson.D{
		{Key: "token_id", Value: tokenId},
	}
	update := bson.M{
		"$unset": bson.M{
			"stats.price_int": true,
		},
	}
	_, err := r.DB.Collection(utils.COLLECTION_TOKEN_URI).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return err
}

func (r Repository) UpdateTokenOnchainStatusByTokenId(tokenId string) error {
	filter := bson.D{
		{Key: "token_id", Value: tokenId},
	}
	update := bson.M{
		"$set": bson.M{
			"isOnchain": true,
		},
	}
	_, err := r.DB.Collection(utils.COLLECTION_TOKEN_URI).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return err
}

func (r Repository) GetNotSyncInscriptionIndexToken() ([]entity.TokenUri, error) {
	tokens := []entity.TokenUri{}

	aggregates := bson.A{
		bson.D{
			{"$match",
				bson.D{
					{"project_id_int", bson.D{{"$gt", 1000000}}},
					{"synced_inscription_info", bson.D{{"$ne", true}}},
				},
			},
		},
		bson.D{{"$sample", bson.D{{"size", 100}}}},
		bson.D{{"$project", r.SelectedTokenFieldsNew()}},
	}

	cursor, err := r.DB.Collection(entity.TokenUri{}.TableName()).Aggregate(context.TODO(), aggregates)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err = cursor.All((context.TODO()), &tokens); err != nil {
		return nil, errors.WithStack(err)
	}

	return tokens, nil
}

func (r Repository) UpdateTokenInscriptionIndex(tokenId string, inscriptionIndex string) error {
	filter := bson.D{
		{Key: "token_id", Value: tokenId},
	}
	update := bson.M{
		"$set": bson.M{
			"inscription_index":       inscriptionIndex,
			"synced_inscription_info": true,
		},
	}
	_, err := r.DB.Collection(utils.COLLECTION_TOKEN_URI).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return err
}

func (r Repository) UpdateTokenOwner(tokenId string, owner *entity.Users) error {
	filter := bson.D{
		{Key: "token_id", Value: tokenId},
	}
	update := bson.M{
		"$set": bson.M{
			"owner_addrress": owner.WalletAddressBTC,
			"owner":          owner,
		},
	}
	_, err := r.DB.Collection(utils.COLLECTION_TOKEN_URI).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return err
}

func (r Repository) UpdateTokenOwnerAddr(tokenId string, addr string) error {
	filter := bson.D{
		{Key: "token_id", Value: tokenId},
	}
	update := bson.M{
		"$set": bson.M{
			"owner_addrress": addr,
		},
	}
	_, err := r.DB.Collection(utils.COLLECTION_TOKEN_URI).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return err
}

func (r Repository) FindOneTokenByListOfTokenIds(tokenIds []string) (*entity.TokenUri, error) {
	resp := &entity.TokenUri{}
	filter := bson.D{
		{Key: "token_id", Value: bson.M{"$in": tokenIds}},
	}
	usr, err := r.FilterOne(entity.TokenUri{}.TableName(), filter)
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) CountTokenUriByProjectId(projectId string) (*int64, error) {
	f := bson.M{
		"project_id": projectId,
	}
	count, err := r.DB.Collection(utils.COLLECTION_TOKEN_URI).CountDocuments(context.TODO(), f)
	if err != nil {
		return nil, err
	}

	return &count, nil
}

func (r Repository) GetNotCreatedActivitiesToken(page int64, limit int64) (*entity.Pagination, error) {
	confs := []entity.TokenUri{}
	resp := &entity.Pagination{}
	f := bson.M{"created_mint_activity": bson.M{"$ne": true}}
	s := []Sort{{SortBy: "created_at", Sort: entity.SORT_ASC}}
	p, err := r.Paginate(entity.TokenUri{}.TableName(), page, limit, f, r.SelectedTokenFieldsNew(), s, &confs)
	if err != nil {
		return nil, err
	}

	resp.Result = confs
	resp.Page = p.Pagination.Page
	resp.Total = p.Pagination.Total
	resp.PageSize = limit
	return resp, nil
}

func (r Repository) GetNotCreatedTxToken(page int64, limit int64) (*entity.Pagination, error) {
	confs := []entity.TokenUri{}
	resp := &entity.Pagination{}
	f := bson.M{"created_token_tx": bson.M{"$ne": true}}
	// hardcode for product environment
	// if os.Getenv("ENVIRONMENT") == "production" {
	// 	f["project_id"] = bson.M{
	// 		"$in": []string{"1000466", "1002270"},
	// 	}
	// }
	// f["token_id"] = {}
	s := []Sort{{SortBy: "created_at", Sort: entity.SORT_ASC}}
	p, err := r.Paginate(entity.TokenUri{}.TableName(), page, limit, f, r.SelectedTokenFieldsNew(), s, &confs)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp.Result = confs
	resp.Page = p.Pagination.Page
	resp.Total = p.Pagination.Total
	resp.PageSize = limit
	return resp, nil
}

func (r Repository) UpdateTokenCreatedMintActivity(tokenID string) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: "token_id", Value: tokenID}}
	update := bson.M{
		"$set": bson.M{"created_mint_activity": true},
	}

	result, err := r.DB.Collection(entity.TokenUri{}.TableName()).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (r Repository) UpdateTokenCreatedTokenTx(tokenID string) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: "token_id", Value: tokenID}}
	update := bson.M{
		"$set": bson.M{"created_token_tx": true},
	}

	result, err := r.DB.Collection(entity.TokenUri{}.TableName()).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (r Repository) FindTokenByTokenIds(tokenIDs []string) ([]entity.TokenUri, error) {
	tokens := []entity.TokenUri{}
	f := bson.M{
		"token_id": bson.M{
			"$in": tokenIDs,
		},
	}
	opts := options.Find().SetProjection(r.SelectedTokenFields())
	cursor, err := r.DB.Collection(entity.TokenUri{}.TableName()).Find(context.TODO(), f, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &tokens); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (r Repository) UpdateTokenUriCreatorByUuid(uuid string, user *entity.Users) error {
	filter := bson.D{
		{Key: "uuid", Value: uuid},
	}
	update := bson.M{
		"$set": bson.M{
			"creator": user,
		},
	}
	_, err := r.DB.Collection(utils.COLLECTION_TOKEN_URI).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return err
}

func (r Repository) AnalyticsTokenUriOwner(f entity.FilterTokenUris) ([]*entity.TokenUriOnwer, error) {
	tokens := []*entity.TokenUriOnwer{}
	//offset := (f.Page - 1) * f.Limit

	filter := bson.A{
		bson.D{
			{"$match",
				bson.D{
					{"project_id", strings.ToLower(*f.GenNFTAddr)},
					//{"token_id", "2d37dbe24f059cbe5004f76df10c8c1bebe3d88adc7229db4d462c05e42fd406i0"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "users"},
					{"localField", "owner_addrress"},
					{"foreignField", "wallet_address_btc"},
					{"as", "owner"},
				},
			},
		},
		bson.D{
			{"$unwind",
				bson.D{
					{"path", "$owner"},
					{"includeArrayIndex", "string"},
					{"preserveNullAndEmptyArrays", true},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"token_id", 1},
					{"owner_addrress", 1},
					{"owner", 1},
				},
			},
		},
	}

	cursor, err := r.DB.Collection(utils.COLLECTION_TOKEN_URI).Aggregate(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All((context.TODO()), &tokens)
	if err != nil {
		return nil, err
	}

	return tokens, err
}
