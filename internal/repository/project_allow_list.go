package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/helpers"
)

func (r Repository) CreateProjectAllowList(data *entity.ProjectAllowList) error {
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) GetProjectAllowList(projectID string, walletAddress string) (*entity.ProjectAllowList, error) {
	resp := &entity.ProjectAllowList{}
	usr, err := r.FilterOne(entity.ProjectAllowList{}.TableName(), bson.D{{"projectID", projectID}, {"userWalletAddress", walletAddress}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r Repository) GetProjectAllowListTotal(projectID string) (int64, error) {
	coll := r.DB.Collection(entity.ProjectAllowList{}.TableName())
	estCount, err := coll.EstimatedDocumentCount(context.TODO())
	if err != nil {
		return 0, err
	}
	_ = estCount
	count, err := coll.CountDocuments(context.TODO(), bson.D{{"projectID", projectID}})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r Repository) GetProjectAllowListTotalByTyppe(projectID string, allowType string) (int64, error) {
	coll := r.DB.Collection(entity.ProjectAllowList{}.TableName())
	estCount, err := coll.EstimatedDocumentCount(context.TODO())
	if err != nil {
		return 0, err
	}
	_ = estCount
	count, err := coll.CountDocuments(context.TODO(), bson.D{{"projectID", projectID}, {"allowedBy", allowType}})
	if err != nil {
		return 0, err
	}
	return count, nil
}
