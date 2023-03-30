package usecase

import (
	"errors"

	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/logger"
)

func (u Usecase) CreateConfig( input structure.ConfigData) (*entity.Configs, error) {

	logger.AtLog.Logger.Info("input", zap.Any("input", input))
	config := &entity.Configs{
		Key: input.Key,
		Value: input.Value,
	}

	conf, err := u.Repo.FindConfig(input.Key)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err := u.Repo.InsertConfig(config)
			if err != nil {
				logger.AtLog.Logger.Error("err", zap.Error(err))
				return nil, err
			}
		}
	}

	conf.Value = input.Value
	updated, err := u.Repo.UpdateConfig(input.Key, conf)
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		return nil, err
	}

	logger.AtLog.Logger.Info("updated",zap.Any("updated", updated))
	return config, nil
}

func (u Usecase) UpdateConfig( input structure.ConfigData) (*entity.Configs, error) {
	return nil, nil
}

func (u Usecase) DeleteConfig( input string) error {

	logger.AtLog.Logger.Info("input", zap.Any("input", input))
	deleted, err := u.Repo.DeleteConfig(input)
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		return err
	}
	logger.AtLog.Logger.Info("deleted",zap.Any("deleted", deleted))

	return nil
}

func (u Usecase) GetConfig( input string) (*entity.Configs, error) {

	logger.AtLog.Logger.Info("input", zap.Any("input", input))

	config, err := u.Repo.FindConfig(input)
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		return nil, err
	}

	return config, nil
}

func (u Usecase) GetConfigs( input structure.FilterConfigs) (*entity.Pagination, error) {
	f := &entity.FilterConfigs{}
	err := copier.Copy(f, input)
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		return nil, err
	}

	confs,  err := u.Repo.ListConfigs(*f)
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		return nil, err
	}

	return confs, nil

}