package usecase

import (
	"github.com/jinzhu/copier"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
)

func (u Usecase) CreateCategory( input structure.CategoryData) (*entity.Categories, error) {

	u.Logger.Info("input", input)
	category := &entity.Categories{
		Name: input.Name,
	}

	err := u.Repo.InsertCategory(category)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}

	u.Logger.Info("inserted",category)
	return category, nil
}

func (u Usecase) UpdateCategory( input structure.UpdateCategoryData) (*entity.Categories, error) {
u.Logger.Info("input", input)
	cat, err := u.Repo.FindCategory(*input.ID)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}

	if input.Name != nil {
		cat.Name = *input.Name
	}
updated, err := u.Repo.UpdateCategory(*input.ID, cat)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}

	u.Logger.Info("updated", updated)
	return cat, nil
}

func (u Usecase) DeleteCategory( input string) error {
u.Logger.Info("input", input)
	deleted, err := u.Repo.DeleteCategory(input)
	if err != nil {
		u.Logger.Error(err)
		return err
	}
	u.Logger.Info("deleted",deleted)

	return nil
}

func (u Usecase) GetCategory( input string) (*entity.Categories, error) {
u.Logger.Info("input", input)
	category, err := u.Repo.FindCategory(input)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}

	return category, nil
}

func (u Usecase) GetCategories( input structure.FilterCategories) (*entity.Pagination, error) {
f := &entity.FilterCategories{}
	err := copier.Copy(f, input)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}

	confs,  err := u.Repo.ListCategories(*f)
	if err != nil {
		u.Logger.Error(err)
		return nil, err
	}

	return confs, nil

}