package usecase

import (
	"errors"

	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/mongo"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
)

func (u Usecase) CreateConfig( input structure.ConfigData) (*entity.Configs, error) {

	u.Logger.Info("input", input)
	config := &entity.Configs{
		Key: input.Key,
		Value: input.Value,
	}

	conf, err := u.Repo.FindConfig(input.Key)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err := u.Repo.InsertConfig(config)
			if err != nil {
				u.Logger.Error(err)
				return nil, err
			}
		}
	}

	conf.Value = input.Value
	updated, err := u.Repo.UpdateConfig(input.Key, conf)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}

	u.Logger.Info("updated",updated)
	return config, nil
}

func (u Usecase) UpdateConfig( input structure.ConfigData) (*entity.Configs, error) {
	return nil, nil
}

func (u Usecase) DeleteConfig( input string) error {

	u.Logger.Info("input", input)
	deleted, err := u.Repo.DeleteConfig(input)
	if err != nil {
		u.Logger.Error(err)
		return err
	}
	u.Logger.Info("deleted",deleted)

	return nil
}

func (u Usecase) GetConfig( input string) (*entity.Configs, error) {

	u.Logger.Info("input", input)

	config, err := u.Repo.FindConfig(input)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}

	return config, nil
}

func (u Usecase) GetConfigs( input structure.FilterConfigs) (*entity.Pagination, error) {
	f := &entity.FilterConfigs{}
	err := copier.Copy(f, input)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}

	confs,  err := u.Repo.ListConfigs(*f)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}

	return confs, nil

}