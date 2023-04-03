package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrNoProjectsFound = errors.New("projects: no documents in result")

func (r Repository) FindProject(projectID string) (*entity.Projects, error) {
	resp := &entity.Projects{}
	usr, err := r.FilterOne(entity.Projects{}.TableName(), bson.D{{utils.KEY_UUID, projectID}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) FindProjectsHaveMinted() ([]entity.ProjectsHaveMinted, error) {
	projects := []entity.ProjectsHaveMinted{}
	f := bson.M{}
	f["index"] = bson.M{"$gte": 1}
	//f["tokenid"] = "1001572"
	opts := options.Find().SetProjection(bson.D{
		{"tokenid", 1},
		{"name", 1},
		{"index", 1},
		{"mintpriceeth", 1},
		{"mintPrice", 1},
		{"creatorAddrr", 1},
	})
	cursor, err := r.DB.Collection(utils.COLLECTION_PROJECTS).Find(context.TODO(), f, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &projects); err != nil {
		return nil, err
	}

	return projects, nil
}

func (r Repository) FindProjectByTokenID(tokenID string) (*entity.Projects, error) {
	resp := &entity.Projects{}
	usr, err := r.FilterOne(entity.Projects{}.TableName(), bson.D{{"tokenid", tokenID}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) FindProjectByTokenIDCustomField(tokenID string, fields []string) (*entity.Projects, error) {
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
		bson.D{{Key: "$match", Value: bson.D{{Key: "tokenid", Value: tokenID}}}},
	}

	cursor, err := r.DB.Collection(entity.Projects{}.TableName()).Aggregate(context.TODO(), aggregates)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	projectList := []entity.Projects{}

	if err = cursor.All((context.TODO()), &projectList); err != nil {
		return nil, errors.WithStack(err)
	}
	if len(projectList) > 0 {
		return &projectList[0], nil
	}
	return nil, errors.New("tokenid not found")
}

func (r Repository) FindProjectByTokenIDs(tokenIds []string) ([]*entity.Projects, error) {
	resp := []*entity.Projects{}
	f := bson.M{}
	f["tokenid"] = bson.M{"$in": tokenIds}
	_, err := r.Paginate(utils.COLLECTION_PROJECTS, 1, int64(len(tokenIds)), f, r.SelectedProjectFields(), nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) FindProjectByTxHash(txHash string) (*entity.Projects, error) {
	resp := &entity.Projects{}
	usr, err := r.FilterOne(entity.Projects{}.TableName(), bson.D{{"txhash", txHash}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r Repository) FindProjectBy(contractAddress string, tokenID string) (*entity.Projects, error) {
	resp := &entity.Projects{}
	contractAddress = strings.ToLower(contractAddress)
	go r.findProjectBy(contractAddress, tokenID)

	p, err := r.Cache.GetData(helpers.ProjectDetailKey(contractAddress, tokenID))
	if err != nil {
		return r.findProjectBy(contractAddress, tokenID)
	}

	bytes := []byte(*p)
	err = json.Unmarshal(bytes, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) IncreaseProjectIndex(projectID string) error {
	filter := bson.D{{Key: "tokenid", Value: projectID}}
	update := bson.M{"$inc": bson.M{"index": 1}}
	_, err := r.DB.Collection(utils.COLLECTION_PROJECTS).UpdateOne(context.TODO(), filter, update)
	return err
}

func (r Repository) FindProjectWithoutCache(contractAddress string, tokenID string) (*entity.Projects, error) {
	contractAddress = strings.ToLower(contractAddress)
	return r.findProjectBy(contractAddress, tokenID)
}

func (r Repository) findProjectBy(contractAddress string, tokenID string) (*entity.Projects, error) {
	contractAddress = strings.ToLower(contractAddress)
	resp := &entity.Projects{}
	usr, err := r.FilterOne(entity.Projects{}.TableName(), bson.D{{"contractAddress", contractAddress}, {"tokenid", tokenID}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}

	r.Cache.SetData(helpers.ProjectDetailKey(contractAddress, tokenID), resp)
	return resp, nil
}

func (r Repository) FindProjectByProjectIdWithoutCache(tokenID string) (*entity.Projects, error) {
	resp := &entity.Projects{}
	usr, err := r.FilterOne(entity.Projects{}.TableName(), bson.D{{"tokenid", tokenID}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r Repository) FindProjectByInscriptionIcon(inscription_icon string) (*entity.Projects, error) {
	resp := &entity.Projects{}
	usr, err := r.FilterOne(entity.Projects{}.TableName(), bson.D{{"inscription_icon", inscription_icon}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r Repository) CreateProject(data *entity.Projects) error {
	data.ContractAddress = strings.ToLower(data.ContractAddress)
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}

	_ = r.Cache.SetData(helpers.ProjectDetailKey(data.ContractAddress, data.TokenID), data)

	return nil
}

func (r Repository) UpdateProjectImages(ID string, images []string, processingImages []string) (*mongo.UpdateResult, error) {
	filter := bson.D{{utils.KEY_UUID, ID}}
	update := bson.M{
		"$set": bson.M{
			"images":           images,
			"processingImages": processingImages,
		},
	}

	result, err := r.DB.Collection(utils.COLLECTION_PROJECTS).UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) UpdateProject(ID string, data *entity.Projects) (*mongo.UpdateResult, error) {
	filter := bson.D{{utils.KEY_UUID, ID}}
	result, err := r.UpdateOne(entity.Projects{}.TableName(), filter, data)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r Repository) UpdateProjectAnimationHtml(ID string, animationHtml string) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: utils.KEY_UUID, Value: ID}}
	update := bson.M{
		"$set": bson.M{
			"stats.animation_html": animationHtml,
		},
	}
	result, err := r.DB.Collection(utils.COLLECTION_PROJECTS).UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) UpdateProjectMintedCount(ID string, mintedCount int32) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: utils.KEY_UUID, Value: ID}}
	update := bson.M{
		"$set": bson.M{
			"stats.minted_count": mintedCount,
		},
	}
	result, err := r.DB.Collection(utils.COLLECTION_PROJECTS).UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) GetProjectsByWalletAddress(add string) ([]entity.Projects, error) {
	confs := []entity.Projects{}
	filter := entity.FilterProjects{}
	filter.WalletAddress = &add
	f := r.FilterProjects(filter)
	s := r.SortProjects()
	_, err := r.Paginate(utils.COLLECTION_PROJECTS, 1, 10, f, r.SelectedProjectFields(), s, &confs)
	if err != nil {
		return nil, err
	}

	return confs, nil
}

func (r Repository) GetProjects(filter entity.FilterProjects) (*entity.Pagination, error) {
	confs := []entity.Projects{}
	resp := &entity.Pagination{}
	f := r.FilterProjects(filter)

	// query := `{ "$where": "this.limitSupply > this.index + this.indexReverse " }`
	// err := json.Unmarshal([]byte(query), &f)
	// if err != nil {
	// 	return nil, err
	// }

	var s []Sort
	if filter.SortBy == "" {
		s = r.SortProjects()
	} else {
		s = []Sort{
			{SortBy: filter.SortBy, Sort: filter.Sort},
			//{SortBy: "tokenid", Sort: entity.SORT_DESC},
		}

		if filter.SortBy == "stats.trending_score" {
			s = append(s, Sort{SortBy: "priority", Sort: entity.SORT_DESC})
			s = append(s, Sort{SortBy: "stats.trending_score", Sort: entity.SORT_DESC})
		}
		s = append(s, Sort{SortBy: "tokenid", Sort: entity.SORT_ASC})
	}
	p, err := r.Paginate(utils.COLLECTION_PROJECTS, filter.Page, filter.Limit, f, r.SelectedProjectFields(), s, &confs)
	if err != nil {
		return nil, err
	}

	resp.Result = confs
	resp.Page = p.Pagination.Page
	resp.Total = p.Pagination.Total
	resp.PageSize = filter.Limit
	return resp, nil
}

func (r Repository) GetAllProjects(filter entity.FilterProjects) ([]entity.Projects, error) {
	projects := []entity.Projects{}
	f := r.FilterProjects(filter)
	cursor, err := r.DB.Collection(utils.COLLECTION_PROJECTS).Find(context.TODO(), f)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &projects); err != nil {
		return nil, err
	}

	return projects, nil
}

func (r Repository) GetAllProjectsWithSelectedFields() ([]entity.Projects, error) {
	projects := []entity.Projects{}
	f := bson.M{}
	opts := options.Find().SetProjection(r.SelectedProjectFields())
	cursor, err := r.DB.Collection(utils.COLLECTION_PROJECTS).Find(context.TODO(), f, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &projects); err != nil {
		return nil, err
	}

	return projects, nil
}

func (r Repository) CountProjects(filter entity.FilterProjects) (*int64, error) {
	//products := &entity.Products{}
	f := r.FilterProjects(filter)
	count, err := r.DB.Collection(utils.COLLECTION_PROJECTS).CountDocuments(context.TODO(), f)
	if err != nil {
		return nil, err
	}

	return &count, nil
}

func (r Repository) GetMintedOutProjects(filter entity.FilterProjects) (*entity.Pagination, error) {
	confs := []entity.Projects{}
	resp := &entity.Pagination{}
	f := r.FilterProjects(filter)

	// query := `{ "$where": "this.limitSupply == this.index + this.indexReverse " }`
	// err := json.Unmarshal([]byte(query), &f)
	// if err != nil {
	// 	return nil, err
	// }

	s := r.SortProjects()
	p, err := r.Paginate(utils.COLLECTION_PROJECTS, filter.Page, filter.Limit, f, r.SelectedProjectFields(), s, &confs)
	if err != nil {
		return nil, err
	}

	resp.Result = confs
	resp.Page = p.Pagination.Page
	resp.Total = p.Pagination.Total
	resp.PageSize = filter.Limit
	return resp, nil
}

func (r Repository) GetRecentWorksProjects(filter entity.FilterProjects) (*entity.Pagination, error) {
	confs := []entity.Projects{}
	resp := &entity.Pagination{}
	f := r.FilterProjects(filter)

	query := `{ "$where": "this.limitSupply > this.index + this.indexReverse " }`
	err := json.Unmarshal([]byte(query), &f)
	if err != nil {
		return nil, err
	}

	s := r.SortProjects()
	p, err := r.Paginate(utils.COLLECTION_PROJECTS, filter.Page, filter.Limit, f, r.SelectedProjectFields(), s, &confs)
	if err != nil {
		return nil, err
	}

	resp.Result = confs
	resp.Page = p.Pagination.Page
	resp.Total = p.Pagination.Total
	resp.PageSize = filter.Limit
	return resp, nil
}

func (r Repository) GetAllRawProjects(filter entity.FilterProjects) (*entity.Pagination, error) {
	confs := []entity.Projects{}
	resp := &entity.Pagination{}
	f := r.FilterProjectRaw(filter)
	s := r.SortProjects()
	p, err := r.Paginate(utils.COLLECTION_PROJECTS, filter.Page, filter.Limit, f, r.SelectedProjectFields(), s, &confs)
	if err != nil {
		return nil, err
	}

	resp.Result = confs
	resp.Page = p.Pagination.Page
	resp.Total = p.Pagination.Total
	resp.PageSize = filter.Limit
	return resp, nil
}

func (r Repository) FilterProjects(filter entity.FilterProjects) bson.M {
	f := r.FilterProjectRaw(filter)
	f["isSynced"] = true
	//f[utils.KEY_DELETED_AT] = nil

	//f["isHidden"] = false

	return f
}

func (r Repository) FilterProjectRaw(filter entity.FilterProjects) bson.M {
	f := bson.M{}
	
	if filter.WalletAddress != nil {
		if *filter.WalletAddress != "" {
			f["creatorAddress"] = bson.M{"$regex": primitive.Regex{
				//Pattern:  *filter.WalletAddress,
				Pattern: fmt.Sprintf(`^%s$`, *filter.WalletAddress),
				Options: "i",
			}}
		}
	}

	if len(filter.Ids) != 0 {
		objectIDs, err := utils.StringsToObjects(filter.Ids)
		if err == nil {
			f["_id"] = bson.M{"$in": objectIDs}
		}
	}

	if filter.Name != nil && len(*filter.Name) >= 3 {
		if *filter.Name != "" {
			f["$or"] = []bson.M{
				{"name": primitive.Regex{Pattern: *filter.Name, Options: "i"}},
				{"creatorProfile.display_name": primitive.Regex{Pattern: *filter.Name, Options: "i"}},
				{"creatorProfile.wallet_address": primitive.Regex{Pattern: *filter.Name, Options: "i"}},
			}
		}
	}

	if filter.CategoryIds != nil && len(filter.CategoryIds) > 0 {
		f["categories"] = bson.M{"$all": filter.CategoryIds}
	}

	if filter.IsHidden != nil {
		f["isHidden"] = *filter.IsHidden
	}

	if filter.IsSynced != nil {
		f["isSynced"] = *filter.IsSynced
	}

	if filter.Status != nil {
		f["status"] = *filter.Status
	}
	
	if filter.TxHash != nil {
		f["txhash"] = *filter.TxHash
	}
	
	if filter.TxHex != nil {
		f["txHex"] = *filter.TxHex
	}
	
	if filter.ContractAddress != nil {
		f["contractAddress"] = *filter.ContractAddress
	}
	
	if filter.CommitTxHash != nil {
		f["commitTxHash"] = *filter.CommitTxHash
	}
	
	if filter.RevealTxHash != nil {
		f["revealTxHash"] = *filter.RevealTxHash
	}

	if len(filter.CustomQueries) > 0 {
		for key, query := range filter.CustomQueries {
			f[key] = query
		}
	}

	return f

	
}

func (r Repository) FindProjectByGenNFTAddr(genNFTAddr string) (*entity.Projects, error) {
	genNFTAddr = strings.ToLower(genNFTAddr)
	resp := &entity.Projects{}
	cached, err := r.Cache.GetData(helpers.ProjectDetailgenNftAddrrKey(genNFTAddr))
	if err == nil && cached != nil {
		err := helpers.ParseCache(cached, resp)
		if err == nil {
			return resp, nil
		}
	}

	prj, err := r.FilterOne(entity.Projects{}.TableName(), bson.D{{Key: "genNFTAddr", Value: genNFTAddr}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(prj, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) UpdateTrendingScoreForProject(tokenID string, trendingScore int64) error {
	filter := bson.D{
		{Key: "tokenid", Value: tokenID},
	}
	update := bson.M{
		"$set": bson.M{
			"stats.trending_score": trendingScore,
		},
	}
	_, err := r.DB.Collection(utils.COLLECTION_PROJECTS).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return err
}

func (r Repository) GetMaxBtcProjectID() (int64, error) {
	btcID := 1000000
	btcMaxID := 1999999

	f := bson.A{
		bson.M{"$match": bson.M{"$and": bson.A{
			bson.M{"tokenIDInt": bson.M{"$gte": btcID}},
			bson.M{"tokenIDInt": bson.M{"$lte": btcMaxID}},
		}}},
		bson.M{"$group": bson.M{"_id": "$tokenIDInt"}},
		bson.M{"$sort": bson.M{"_id": -1}},
	}

	cursor, err := r.DB.Collection(utils.COLLECTION_PROJECTS).Aggregate(context.TODO(), f)
	if err != nil {
		return 0, err
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return 0, err
	}

	for _, results := range results {
		i := &entity.MaxProjectID{}
		err := helpers.Transform(results, i)
		if err != nil {
			continue
		}
		return i.ID, nil

	}

	d := int64(btcID)
	return d, nil
}

func (r Repository) SortProjects() []Sort {
	s := []Sort{}
	s = append(s, Sort{SortBy: "priority", Sort: entity.SORT_DESC})
	s = append(s, Sort{SortBy: "index", Sort: entity.SORT_DESC})
	s = append(s, Sort{SortBy: "tokenid", Sort: entity.SORT_ASC})
	return s
}

func (r Repository) SelectedProjectFields() bson.D {
	f := bson.D{
		{"id", 1},
		{"contractAddress", 1},
		{"tokenid", 1},
		{"maxSupply", 1},
		{"limitSupply", 1},
		{"mintPrice", 1},
		{"networkFee", 1},
		{"name", 1},
		{"creatorName", 1},
		{"creatorAddress", 1},
		{"categories", 1},
		{"thumbnail", 1},
		{"mintFee", 1},
		{"openMintUnixTimestamp", 1},
		{"closeMintUnixTimestamp", 1},
		{"genNFTAddr", 1},
		{"mintTokenAddress", 1},
		{"minted_time", 1},
		{"license", 1},
		{"description", 1},
		{"stats", 1},
		{"status", 1},
		{"tokenDescription", 1},
		{"index", 1},
		{"indexReverse", 1},
		{"creatorProfile", 1},
		{"images", 1},
		{"mintedImages", 1},
		{"whiteListEthContracts", 1},
		{"isFullChain", 1},
		{"reportUsers", 1},
		{"mintpriceeth", 1},
		{"fromAuthentic", 1},
		{"tokenAddress", 1},
		{"tokenId", 1},
		{"ownerOf", 1},
		{"ordinalsTx", 1},
		{"inscribedBy", 1},
		{"isSynced", 1},
		{"txhash", 1},
	}
	return f
}

func (r Repository) SetProjectInscriptionIcon(projectID string, inscriptionIcon string) error {
	f := bson.D{
		{Key: "tokenid", Value: projectID},
	}

	update := bson.M{
		"$set": bson.M{
			"inscription_icon": inscriptionIcon,
		},
	}

	_, err := r.DB.Collection(entity.Projects{}.TableName()).UpdateOne(context.TODO(), f, update)
	if err != nil {
		return err
	}

	return err
}

func (r Repository) UpdateProjectTraitStats(projectID string, traitStat []entity.TraitStat) error {
	f := bson.D{
		{Key: "tokenid", Value: projectID},
	}

	update := bson.M{
		"$set": bson.M{
			"traitsStat": traitStat,
		},
	}

	_, err := r.DB.Collection(entity.Projects{}.TableName()).UpdateOne(context.TODO(), f, update)
	if err != nil {
		return err
	}

	return err
}

func (r Repository) ProjectGetCurrentListingNumber(projectID string) (uint64, error) {
	result := []entity.TokenUriListingPage{}
	pipeline := bson.A{
		bson.D{{"$match", bson.D{{"project_id", projectID}}}},
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
					{"preserveNullAndEmptyArrays", false},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "dex_btc_buy_eth"},
					{"localField", "token_id"},
					{"foreignField", "inscription_id"},
					{"let", bson.D{{"status", "$status"}}},
					{"pipeline",
						bson.A{
							bson.D{
								{"$match",
									bson.D{
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
						},
					},
					{"as", "listing_eth"},
				},
			},
		},
		bson.D{{"$project", bson.D{{"listing_eth_size", bson.D{{"$size", "$listing_eth"}}}}}},
		bson.D{{"$match", bson.D{{"listing_eth_size", bson.D{{"$eq", 0}}}}}},
		bson.D{
			{"$facet",
				bson.D{
					{"totalCount",
						bson.A{
							bson.D{{"$count", "count"}},
						},
					},
				},
			},
		},
	}

	cursor, err := r.DB.Collection(entity.TokenUri{}.TableName()).Aggregate(context.TODO(), pipeline)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	if err = cursor.All((context.TODO()), &result); err != nil {
		return 0, errors.WithStack(err)
	}
	if len(result) > 0 {
		if len(result[0].TotalCount) > 0 {
			return uint64(result[0].TotalCount[0].Count), nil
		}
		return 0, nil
	}

	return 0, nil
}

func (r Repository) ProjectGetListingVolume(projectID string) (uint64, error) {
	result := []entity.TokenUriListingVolume{}
	pipeline := bson.A{
		bson.D{
			{"$match",
				bson.D{
					{"matched", true},
					{"cancelled", false},
					{"buyer", bson.D{{"$exists", true}}},
				},
			},
		},
		bson.D{
			{"$addFields",
				bson.D{
					{"diffbuyer",
						bson.D{
							{"$ne",
								bson.A{
									"$buyer",
									"$seller_address",
								},
							},
						},
					},
				},
			},
		},
		bson.D{{"$match", bson.D{{"diffbuyer", true}}}},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "token_uri"},
					{"localField", "inscription_id"},
					{"foreignField", "token_id"},
					{"let", bson.D{{"id", "$_id"}}},
					{"pipeline",
						bson.A{
							bson.D{{"$match", bson.D{{"project_id", projectID}}}},
						},
					},
					{"as", "collection_id"},
				},
			},
		},
		bson.D{
			{"$unwind",
				bson.D{
					{"path", "$collection_id"},
					{"preserveNullAndEmptyArrays", false},
				},
			},
		},
		bson.D{
			{"$group",
				bson.D{
					{"_id", ""},
					{"Amount", bson.D{{"$sum", "$amount"}}},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"_id", 0},
					{"totalAmount", "$Amount"},
				},
			},
		},
	}

	cursor, err := r.DB.Collection(entity.DexBTCListing{}.TableName()).Aggregate(context.TODO(), pipeline)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	if err = cursor.All((context.TODO()), &result); err != nil {
		return 0, errors.WithStack(err)
	}
	if len(result) > 0 {
		return uint64(result[0].TotalAmount), nil
	}

	return 0, nil
}

func (r Repository) ProjectGetMintVolume(projectID string) (uint64, error) {
	result := []entity.TokenUriListingVolume{}
	pipeline := bson.A{
		bson.D{
			{"$match",
				bson.D{
					{"isMinted", true},
					{"projectID", projectID},
				},
			},
		},
		bson.D{
			{"$group",
				bson.D{
					{"_id", ""},
					{"Amount", bson.D{{"$sum", "$project_mint_price"}}},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"_id", 0},
					{"totalAmount", "$Amount"},
				},
			},
		},
	}

	cursor, err := r.DB.Collection(entity.MintNftBtc{}.TableName()).Aggregate(context.TODO(), pipeline)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	if err = cursor.All((context.TODO()), &result); err != nil {
		return 0, errors.WithStack(err)
	}
	if len(result) > 0 {
		return uint64(result[0].TotalAmount), nil
	}

	return 0, nil
}

func (r Repository) ProjectGetCEXVolume(projectID string) (uint64, error) {
	result := []entity.TokenUriListingVolume{}
	pipeline := bson.A{
		bson.D{{"$match", bson.D{{"isSold", true}}}},
		bson.D{{"$addFields", bson.D{{"price", bson.D{{"$toDouble", "$amount"}}}}}},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "token_uri"},
					{"localField", "inscriptionID"},
					{"foreignField", "token_id"},
					{"let", bson.D{{"id", "$_id"}}},
					{"pipeline",
						bson.A{
							bson.D{{"$match", bson.D{{"project_id", projectID}}}},
						},
					},
					{"as", "collection_id"},
				},
			},
		},
		bson.D{
			{"$unwind",
				bson.D{
					{"path", "$collection_id"},
					{"preserveNullAndEmptyArrays", false},
				},
			},
		},
		bson.D{
			{"$group",
				bson.D{
					{"_id", ""},
					{"Amount", bson.D{{"$sum", "$price"}}},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"_id", 0},
					{"totalAmount", "$Amount"},
				},
			},
		},
	}

	cursor, err := r.DB.Collection(entity.MarketplaceBTCListing{}.TableName()).Aggregate(context.TODO(), pipeline)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	if err = cursor.All((context.TODO()), &result); err != nil {
		return 0, errors.WithStack(err)
	}
	if len(result) > 0 {
		return uint64(result[0].TotalAmount), nil
	}

	return 0, nil
}

func (r Repository) UpdateProjectIndexAndMaxSupply(projectID string, maxSupply int64, index int64) error {
	f := bson.D{
		{Key: "tokenid", Value: projectID},
	}

	update := bson.M{
		"$set": bson.M{
			"maxSupply": maxSupply,
			"index":     index,
		},
	}

	_, err := r.DB.Collection(entity.Projects{}.TableName()).UpdateOne(context.TODO(), f, update)
	if err != nil {
		return err
	}

	return err
}

func (r Repository) GetProjectTrendingScore(projectID string) (int64, error) {
	var trendingScore int64
	resp := &entity.Projects{}
	opts := options.FindOne().SetProjection(bson.D{
		{Key: "stats.trending_score", Value: 1},
	})
	project, err := r.FilterOne(entity.Projects{}.TableName(), bson.D{{Key: "tokenid", Value: projectID}}, opts)
	if err != nil {
		return trendingScore, errors.WithStack(err)
	}

	err = helpers.Transform(project, resp)
	if err != nil {
		return trendingScore, errors.WithStack(err)
	}
	trendingScore = resp.Stats.TrendingScore
	return trendingScore, nil
}

func (r Repository) AggregateProjectsFloorPrice(projectIDs []string) ([]structure.ProjectFloorPrice, error) {

	result := []structure.ProjectFloorPrice{}

	projectsBsonA := bson.A{}
	for _, v := range projectIDs {
		projectsBsonA = append(projectsBsonA, v)
	}

	pipeLine := bson.A{
		bson.D{
			{"$match",
				bson.D{
					{"matched", false},
					{"cancelled", false},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "token_uri"},
					{"localField", "inscription_id"},
					{"foreignField", "token_id"},
					{"let", bson.D{}},
					{"pipeline",
						bson.A{
							bson.D{
								{"$match",
									bson.D{
										{"$expr",
											bson.D{
												{"$in",
													bson.A{
														"$project_id",
														projectsBsonA,
													}},
											},
										},
									},
								},
							},
						},
					},
					{"as", "collection_id"},
				},
			},
		},
		bson.D{
			{"$unwind",
				bson.D{
					{"path", "$collection_id"},
					{"preserveNullAndEmptyArrays", false},
				},
			},
		},
		bson.D{
			{"$group",
				bson.D{
					{"_id", "$collection_id.project_id"},
					{"floor", bson.D{{"$min", "$amount"}}},
				},
			},
		},
	}

	cursor, err := r.DB.Collection(entity.DexBTCListing{}.TableName()).Aggregate(context.TODO(), pipeLine, nil)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) SetProjectIndex(projectID string, index int) error {
	f := bson.D{
		{Key: "tokenid", Value: projectID},
	}

	update := bson.M{
		"$set": bson.M{
			"index": index,
		},
	}

	_, err := r.DB.Collection(entity.Projects{}.TableName()).UpdateOne(context.TODO(), f, update)
	if err != nil {
		return err
	}

	return err
}
