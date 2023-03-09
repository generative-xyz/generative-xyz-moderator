package repository

import (
	"context"
	"log"
	"strconv"

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
	tokens := entity.TokenUriListingPage{}
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

	f2 := bson.A{
		bson.D{{"$match", f}},
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
									bson.D{
										{"matched", false},
										{"cancelled", false},
									},
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
					{"preserveNullAndEmptyArrays", true},
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
									},
									{"then", listingAmountDefault},
									{"else", "$listing.amount"},
								},
							},
						},
					},
					{"orderID", "$listing._id"},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"_id", 1},
					{"token_id", 1},
					{"gen_nft_addrress", 1},
					{"contract_address", 1},
					{"project_id", 1},
					{"image", 1},
					{"priority", 1},
					{"inscription_index", 1},
					{"order_inscription_index", 1},
					{"thumbnail", 1},
					{"buyable", 1},
					{"priceBTC", 1},
					{"orderID", 1},
					{"project.tokenid", 1},
				},
			},
		},
		bson.D{{"$sort", bson.D{{filter.SortBy, filter.Sort}}}},
		bson.D{{"$skip", (filter.Page - 1) * filter.Limit}},
		bson.D{{"$limit", filter.Limit}},
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

	log.Println("len(tokens)", len(tokens.TotalCount))

	resp.Result = tokens.TotalData
	resp.Page = filter.Page
	if len(tokens.TotalCount) > 0 {
		resp.Total = tokens.TotalCount[0].Count
	}
	resp.PageSize = filter.Limit
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

	if filter.Search != nil && len(*filter.Search) >= 3 {
		f["$or"] = []bson.M{
			{"token_id": primitive.Regex{Pattern: *filter.Search, Options: "i"}},
			{"inscription_index": primitive.Regex{Pattern: *filter.Search, Options: "i"}},
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

	if filter.HasPrice != nil || filter.FromPrice != nil || filter.ToPrice != nil {
		priceFilter := bson.M{}
		if filter.HasPrice != nil {
			priceFilter["$exists"] = *filter.HasPrice
			priceFilter["$ne"] = nil
		}
		if filter.FromPrice != nil {
			priceFilter["$gte"] = *filter.FromPrice
		}
		if filter.ToPrice != nil {
			priceFilter["$lte"] = *filter.ToPrice
		}
		f["stats.price_int"] = priceFilter
	}

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

func (r Repository) GetAllNotSyncInscriptionIndexToken() ([]entity.TokenUri, error) {
	tokens := []entity.TokenUri{}

	f := bson.M{
		"project_id_int":          bson.M{"$gt": 1000000},
		"synced_inscription_info": bson.M{"$ne": true},
	}
	//f[utils.KEY_DELETED_AT] = nil
	opts := options.Find().SetProjection(r.SelectedTokenFields()).SetLimit(100)
	cursor, err := r.DB.Collection(utils.COLLECTION_TOKEN_URI).Find(context.TODO(), f, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &tokens); err != nil {
		return nil, err
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
