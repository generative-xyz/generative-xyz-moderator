package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"rederinghub.io/internal/entity"
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

	// required token id
	if filter.TokenID == "" {
		return nil, fmt.Errorf("token id is required")
	}

	var result []entity.TokenFileFragment

	// init filter
	queryFilter := bson.M{
		"token_id": filter.TokenID,
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
	options.SetSort(bson.M{"sequence": 1})

	cursor, err := r.DB.Collection(utils.TOKEN_FILE_FRAGMENT).Find(ctx, queryFilter, options)
	if err != nil {
		return nil, err
	}

	if cursor.All(context.Background(), &result); err != nil {
		return nil, err
	}
	return result, nil
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
