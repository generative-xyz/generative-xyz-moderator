package repository

import (
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
	usr, err := r.FilterOne(entity.ProjectAllowList{}.TableName(), bson.D{{"projectID", walletAddress}, {"userWalletAddress", walletAddress}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}