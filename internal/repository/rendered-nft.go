package repository

import (
	"rederinghub.io/internal/model"
	"rederinghub.io/pkg/drivers/mongodb"
)

type RenderedNftRepository interface {
	mongodb.Repository
}

type renderedNftRepository struct {
	mongodb.BaseRepository
}

func NewRenderedNftRepository(db model.Database) RenderedNftRepository {
	return &renderedNftRepository{
		BaseRepository: mongodb.BaseRepository{
			CollectionName: model.RenderedNft{}.CollectionName(),
			DB:             db.DB(),
		},
	}
}
