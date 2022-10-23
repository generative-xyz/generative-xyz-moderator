package repository

import (
	"rederinghub.io/internal/model"
	"rederinghub.io/pkg/drivers/mongodb"
)

type TemplateRepository interface {
	mongodb.Repository
}

type templateRepository struct {
	mongodb.BaseRepository
}

func NewTemplateRepository(db model.Database) TemplateRepository {
	return &templateRepository{
		BaseRepository: mongodb.BaseRepository{
			CollectionName: model.Template{}.CollectionName(),
			DB:             db.DB(),
		},
	}
}
