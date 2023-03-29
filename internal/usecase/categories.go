package usecase

import (
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/logger"
)

func (u Usecase) CreateCategory( input structure.CategoryData) (*entity.Categories, error) {

	logger.AtLog.Logger.Info("input", zap.Any("input", input))
	category := &entity.Categories{
		Name: input.Name,
	}

	err := u.Repo.InsertCategory(category)
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		return nil, err
	}

	logger.AtLog.Logger.Info("inserted", zap.Any("category", category))
	return category, nil
}

func (u Usecase) UpdateCategory( input structure.UpdateCategoryData) (*entity.Categories, error) {
logger.AtLog.Logger.Info("input", zap.Any("input", input))
	cat, err := u.Repo.FindCategory(*input.ID)
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		return nil, err
	}

	if input.Name != nil {
		cat.Name = *input.Name
	}
updated, err := u.Repo.UpdateCategory(*input.ID, cat)
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		return nil, err
	}

	logger.AtLog.Logger.Info("updated", zap.Any("updated", updated))
	return cat, nil
}

func (u Usecase) DeleteCategory( input string) error {
logger.AtLog.Logger.Info("input", zap.Any("input", input))
	deleted, err := u.Repo.DeleteCategory(input)
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		return err
	}
	logger.AtLog.Logger.Info("inserted", zap.Any("deleted", deleted))
	return nil
}

func (u Usecase) GetCategory( input string) (*entity.Categories, error) {
logger.AtLog.Logger.Info("input", zap.Any("input", input))
	category, err := u.Repo.FindCategory(input)
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		return nil, err
	}

	return category, nil
}

func (u Usecase) GetCategories( input structure.FilterCategories) (*entity.Pagination, error) {
f := &entity.FilterCategories{}
	err := copier.Copy(f, input)
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		return nil, err
	}

	confs,  err := u.Repo.ListCategories(*f)
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		return nil, err
	}

	return confs, nil

}