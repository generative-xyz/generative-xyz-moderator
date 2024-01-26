package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"time"
)

type TokenFileFragmentFileter struct {
	TokenID  string
	Sequence int
	Status   entity.FileFragmentStatus
	Page     int
	PageSize int
}

type AggregateTokenMintingInfo struct {
	TokenID string `bson:"token_id" json:"token_id"`
	All     int    `bson:"all" json:"all"`
	Pending int    `bson:"pending" json:"pending"`
	Done    int    `bson:"done" json:"done"`
}

type TokenFragmentJobFilter struct {
	Status   entity.TokenFragmentJobStatus
	Page     int
	PageSize int
}

func (r Repository) FindTokenFileFragment(ctx context.Context, tokenID string, sequence int) (*entity.TokenFileFragment, error) {
	var file entity.TokenFileFragment
	err := r.DB.Collection(file.TableName()).FindOne(ctx, bson.M{"token_id": tokenID, "sequence": sequence}).Decode(&file)

	if err != nil {
		return nil, err
	}

	return &file, nil
}

func (r Repository) InsertFileFragment(ctx context.Context, file *entity.TokenFileFragment) error {
	id := primitive.NewObjectID()
	file.CreatedAt = time.Now()
	file.BaseEntity = entity.BaseEntity{
		ID:   id,
		UUID: id.Hex(),
	}
	_, err := r.DB.Collection(file.TableName()).InsertOne(ctx, file)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) FindTokenFileFragments(ctx context.Context, filter TokenFileFragmentFileter) ([]entity.TokenFileFragment, error) {

	// check limit , override if 0
	limit := filter.PageSize
	if limit == 0 {
		limit = 10
	}
	if filter.Page == 0 {
		filter.Page++
	}

	var result []entity.TokenFileFragment

	queryFilter := bson.M{}
	// init filter
	if filter.TokenID != "" {
		queryFilter["token_id"] = filter.TokenID
	}

	if filter.Sequence > 0 {
		queryFilter["sequence"] = filter.Sequence
	}
	if filter.Status > 0 {
		queryFilter["status"] = filter.Status
	}

	// init options
	options := options.Find()
	options.SetSkip(int64((filter.Page - 1) * limit))
	options.SetLimit(int64(limit))
	options.SetSort(bson.M{"created_at": 1})

	cursor, err := r.DB.Collection(utils.TOKEN_FILE_FRAGMENT).Find(ctx, queryFilter, options)
	if err != nil {
		return nil, err
	}

	if cursor.All(context.Background(), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (r Repository) UpdateFileFragmentStatus(ctx context.Context, id string, updateFields map[string]interface{}) error {
	filter := bson.M{"uuid": id}
	updatedQuery := bson.M{}
	for k, v := range updateFields {
		updatedQuery[k] = v
	}

	_, err := r.DB.Collection(utils.TOKEN_FILE_FRAGMENT).UpdateOne(ctx, filter, bson.M{"$set": updatedQuery})
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) CreateFragmentJob(ctx context.Context, job *entity.TokenFragmentJob) error {

	if job.TokenId == "" {
		return fmt.Errorf("token id is required")
	}

	if job.FilePath == "" {
		return fmt.Errorf("file path is required")
	}

	id := primitive.NewObjectID()
	job.CreatedAt = time.Now()

	job.BaseEntity = entity.BaseEntity{
		ID:   id,
		UUID: id.Hex(),
	}

	job.Status = entity.FragmentJobStatusPending
	_, err := r.DB.Collection(job.TableName()).InsertOne(ctx, job)

	if err != nil {
		return err
	}

	return nil
}

func (r Repository) UpdateFragmentJobStatus(ctx context.Context, uuid string, status entity.TokenFragmentJobStatus, note string) error {
	_, err := r.DB.Collection(utils.TOKEN_FILE_FRAGMENT_JOB).UpdateOne(ctx, bson.M{"uuid": uuid}, bson.M{"$set": bson.M{"status": status, "note": note}})
	return err
}

func (r Repository) FindFragmentJobs(ctx context.Context, filter TokenFragmentJobFilter) ([]entity.TokenFragmentJob, error) {
	var jobs []entity.TokenFragmentJob

	limit := filter.PageSize
	if limit == 0 {
		limit = 5
	}
	if filter.Page == 0 {
		filter.Page++
	}

	queryFilter := bson.M{}

	if filter.Status > 0 {
		queryFilter["status"] = filter.Status
	}

	options := options.Find()
	options.SetSkip(int64((filter.Page - 1) * limit))
	options.SetLimit(int64(limit))
	options.SetSort(bson.M{"created_at": 1})

	cursor, err := r.DB.Collection(utils.TOKEN_FILE_FRAGMENT_JOB).Find(ctx, queryFilter, options)
	if err != nil {
		return nil, err
	}

	if cursor.All(context.Background(), &jobs); err != nil {
		return nil, err
	}

	return jobs, nil

}

func (r Repository) GetStoreWallet() (*entity.StoreFileWallet, error) {
	var wallet *entity.StoreFileWallet
	err := r.DB.Collection(entity.StoreFileWallet{}.TableName()).FindOne(context.Background(), bson.M{}).Decode(&wallet)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (r Repository) AggregateMintingInfo(ctx context.Context, tokenID string) ([]AggregateTokenMintingInfo, error) {
	f := bson.A{
		bson.D{{"$match", bson.D{{"token_id", tokenID}}}},
		bson.D{
			{"$project",
				bson.D{
					{"token_id", 1},
					{"status", 1},
					{"pending",
						bson.D{
							{"$cond",
								bson.A{
									bson.D{
										{"$eq",
											bson.A{
												"$status",
												1,
											},
										},
									},
									1,
									0,
								},
							},
						},
					},
					{"done",
						bson.D{
							{"$cond",
								bson.A{
									bson.D{
										{"$eq",
											bson.A{
												"$status",
												2,
											},
										},
									},
									1,
									0,
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$group",
				bson.D{
					{"_id", bson.D{{"token_id", "$token_id"}}},
					{"all", bson.D{{"$sum", 1}}},
					{"pending", bson.D{{"$sum", "$pending"}}},
					{"done", bson.D{{"$sum", "$done"}}},
				},
			},
		},
		bson.D{{"$addFields", bson.D{{"token_id", "$_id.token_id"}}}},
	}

	cursor, err := r.DB.Collection(entity.TokenFileFragment{}.TableName()).Aggregate(ctx, f)
	if err != nil {
		return nil, err
	}

	aggregation := []AggregateTokenMintingInfo{}
	if err = cursor.All(ctx, &aggregation); err != nil {
		return nil, err
	}

	return aggregation, nil
}

func (r Repository) AggregateModularInscriptions(ctx context.Context, projectID string, offset, limit int) ([]entity.TokenUri, error) {
	f := bson.A{
		bson.D{{"$match", bson.D{{"project_id", projectID}}}},
		bson.D{{"$project", bson.D{
			{"_id", 1},
			{"token_id", 1},
			{"owner_addrress", 1},
		}}},
		bson.D{{"$skip", offset}},
		bson.D{{"$limit", limit}},
	}

	cursor, err := r.DB.Collection(entity.TokenUri{}.TableName()).Aggregate(ctx, f)
	if err != nil {
		return nil, err
	}

	aggregation := []entity.TokenUri{}
	if err = cursor.All(ctx, &aggregation); err != nil {
		return nil, err
	}

	return aggregation, nil
}

func (r Repository) AggregateListModularInscriptions(ctx context.Context, filter structure.FilterTokens) (*entity.Pagination, error) {
	_match := bson.D{}

	if filter.OwnerAddr != nil && *filter.OwnerAddr != "" {
		_match = append(_match, bson.E{"owner_addrress", filter.OwnerAddr})
	}

	if filter.GenNFTAddr != nil && *filter.GenNFTAddr != "" {
		_match = append(_match, bson.E{"project_id", filter.GenNFTAddr})
	}

	limit := int64(20)
	page := int64(1)

	if filter.Page > 0 {
		page = filter.Page
	}

	if filter.Limit > 0 {
		limit = filter.Limit
	}

	offset := (page - 1) * limit
	f := bson.A{
		bson.D{{"$match", _match}},
		bson.D{{"$sort", bson.D{{"_id", -1}}}},
		//bson.D{{"$project", bson.D{
		//	{"_id", 1},
		//	{"token_id", 1},
		//	{"owner_addrress", 1},
		//}}},
		bson.D{{"$skip", offset}},
		bson.D{{"$limit", limit}},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "projects"},
					{"localField", "project_id"},
					{"foreignField", "tokenid"},
					{"as", "project"},
				},
			},
		},
		bson.D{
			{"$unwind",
				bson.D{
					{"path", "$project"},
					{"preserveNullAndEmptyArrays", true},
				},
			},
		},

		bson.D{
			{"$lookup",
				bson.D{
					{"from", "users"},
					{"localField", "owner_addrress"},
					{"foreignField", "wallet_address_btc_taproot"},
					{"as", "owner"},
				},
			},
		},
		bson.D{
			{"$unwind",
				bson.D{
					{"path", "$owner"},
					{"preserveNullAndEmptyArrays", true},
				},
			},
		},
	}

	cursor, err := r.DB.Collection(entity.TokenUri{}.TableName()).Aggregate(ctx, f)
	if err != nil {
		return nil, err
	}

	aggregation := []entity.ModularTokenUri{}
	if err = cursor.All(ctx, &aggregation); err != nil {
		return nil, err
	}

	countF := bson.A{
		bson.D{{"$match", _match}},
		bson.D{
			{"$group",
				bson.D{
					{"_id", "_id"},
					{"total", bson.D{{"$sum", 1}}},
				},
			},
		},
	}
	cursor1, err := r.DB.Collection(entity.TokenUri{}.TableName()).Aggregate(ctx, countF)
	if err != nil {
		return nil, err
	}

	aggregation1 := []entity.Total{}
	if err = cursor1.All(ctx, &aggregation1); err != nil {
		return nil, err
	}

	ap := entity.Total{}
	if len(aggregation1) > 0 {
		ap = aggregation1[0]
	}
	p := entity.Pagination{
		Page:      page,
		Total:     ap.Total,
		TotalPage: int64(math.Ceil(float64(ap.Total) / float64(limit))),
		PageSize:  limit,
		Result:    aggregation,
	}
	return &p, nil
}

func (r Repository) GroupModularInscByAttr(ctx context.Context, filter structure.FilterTokens) ([]*entity.ModularTokenAttr, error) {
	_match := bson.D{}

	if filter.OwnerAddr != nil && *filter.OwnerAddr != "" {
		_match = append(_match, bson.E{"owner_addrress", filter.OwnerAddr})
	}

	if filter.GenNFTAddr != nil && *filter.GenNFTAddr != "" {
		_match = append(_match, bson.E{"project_id", filter.GenNFTAddr})
	}

	skip := int64(0)
	limit := int64(20)

	if filter.Limit > 0 {
		limit = filter.Limit
	}

	if filter.Page > 0 {
		skip = (filter.Page - 1) * limit
	}

	f := bson.A{
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "token_uri"},
					{"localField", "inscription_id"},
					{"foreignField", "token_id"},
					{"pipeline",
						bson.A{
							bson.D{
								{"$match", _match},
							},
						},
					},
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
			{"$group",
				bson.D{
					{"_id", "$attribute"},
					{"total", bson.D{{"$sum", 1}}},
				},
			},
		},
		bson.D{{"$sort", bson.D{{"total", -1}}}},
		bson.D{{"$skip", skip}},
		bson.D{{"$limit", limit}},
		bson.D{
			{"$graphLookup",
				bson.D{
					{"from", "modular_inscription_attribute"},
					{"startWith", "$_id"},
					{"connectFromField", "attribute"},
					{"connectToField", "attribute"},
					{"as", "string"},
					{"maxDepth", 1},
					{"depthField", "string"},
					{"as", "attr"},
				},
			}},
	}

	cursor, err := r.DB.Collection(utils.COLLECTION_MODULAR_INSCRIPTION_ATTRIBUTE).Aggregate(ctx, f)
	if err != nil {
		return nil, err
	}

	aggregation := []*entity.ModularTokenAttr{}
	if err = cursor.All(ctx, &aggregation); err != nil {
		return nil, err
	}

	return aggregation, nil
}

func (r Repository) AggregateListModularInscriptionsByTokenIDs(ctx context.Context, tokenIDs []string) ([]entity.ModularTokenUri, error) {
	_match := bson.D{}
	_match = append(_match, bson.E{"token_id", bson.D{{"$in", tokenIDs}}})

	f := bson.A{
		bson.D{{"$match", _match}},
		bson.D{{"$sort", bson.D{{"_id", -1}}}},
		//bson.D{{"$project", bson.D{
		//	{"_id", 1},
		//	{"token_id", 1},
		//	{"owner_addrress", 1},
		//}}},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "projects"},
					{"localField", "project_id"},
					{"foreignField", "tokenid"},
					{"as", "project"},
				},
			},
		},
		bson.D{
			{"$unwind",
				bson.D{
					{"path", "$project"},
					{"preserveNullAndEmptyArrays", true},
				},
			},
		},

		bson.D{
			{"$lookup",
				bson.D{
					{"from", "users"},
					{"localField", "owner_addrress"},
					{"foreignField", "wallet_address_btc_taproot"},
					{"as", "owner"},
				},
			},
		},
		bson.D{
			{"$unwind",
				bson.D{
					{"path", "$owner"},
					{"preserveNullAndEmptyArrays", true},
				},
			},
		},
	}

	cursor, err := r.DB.Collection(entity.TokenUri{}.TableName()).Aggregate(ctx, f)
	if err != nil {
		return nil, err
	}

	aggregation := []entity.ModularTokenUri{}
	if err = cursor.All(ctx, &aggregation); err != nil {
		return nil, err
	}

	return aggregation, nil
}
