package usecase

import (
	"github.com/jinzhu/copier"
	"github.com/opentracing/opentracing-go"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
)

func (u Usecase) CreateCategory(rootSpan opentracing.Span, input structure.CategoryData) (*entity.Categories, error) {
	span, log := u.StartSpan("CreateCategory", rootSpan)
	defer u.Tracer.FinishSpan(span, log )

	log.SetData("input", input)
	category := &entity.Categories{
		Name: input.Name,
	}

	err := u.Repo.InsertCategory(category)
	if err != nil {
		log.Error(" u.Repo.InsertCategory", err.Error(), err)
		return nil, err
	}

	log.SetData("inserted",category)
	return category, nil
}

func (u Usecase) UpdateCategory(rootSpan opentracing.Span, input structure.UpdateCategoryData) (*entity.Categories, error) {
	span, log := u.StartSpan("UpdateCategory", rootSpan)
	defer u.Tracer.FinishSpan(span, log )

	log.SetData("input", input)
	cat, err := u.Repo.FindCategory(*input.ID)
	if err != nil {
		log.Error(" u.Repo.FindCategory", err.Error(), err)
		return nil, err
	}

	if input.Name != nil {
		cat.Name = *input.Name
	}
	
	updated, err := u.Repo.UpdateCategory(*input.ID, cat)
	if err != nil {
		log.Error(" u.Repo.UpdateCategory", err.Error(), err)
		return nil, err
	}

	log.SetData("updated", updated)
	return cat, nil
}

func (u Usecase) DeleteCategory(rootSpan opentracing.Span, input string) error {
	span, log := u.StartSpan("GetCategory", rootSpan)
	defer u.Tracer.FinishSpan(span, log )

	log.SetData("input", input)
	deleted, err := u.Repo.DeleteCategory(input)
	if err != nil {
		log.Error(" u.Repo.DeleteCategory", err.Error(), err)
		return err
	}
	log.SetData("deleted",deleted)

	return nil
}

func (u Usecase) GetCategory(rootSpan opentracing.Span, input string) (*entity.Categories, error) {
	span, log := u.StartSpan("GetCategory", rootSpan)
	defer u.Tracer.FinishSpan(span, log )

	log.SetData("input", input)

	category, err := u.Repo.FindCategory(input)
	if err != nil {
		log.Error(" u.Repo.FindCategory", err.Error(), err)
		return nil, err
	}

	return category, nil
}

func (u Usecase) GetCategories(rootSpan opentracing.Span, input structure.FilterCategories) (*entity.Pagination, error) {
	span, log := u.StartSpan("GetCategories", rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	f := &entity.FilterCategories{}
	err := copier.Copy(f, input)
	if err != nil {
		log.Error("copier.Copy", err.Error(), err)
		return nil, err
	}

	confs,  err := u.Repo.ListCategories(*f)
	if err != nil {
		log.Error(" u.Repo.FindCategory", err.Error(), err)
		return nil, err
	}

	return confs, nil

}