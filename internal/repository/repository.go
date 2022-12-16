package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"rederinghub.io/internal/model"
	"rederinghub.io/pkg/drivers/mongodb"
)

type IRepository interface {
	mongodb.Repository
	CreateOne(ctx context.Context, model model.IEntity, opts ...*options.InsertOneOptions) (primitive.ObjectID, error)
	UpdateOneByID(ctx context.Context, model model.IEntity, id  primitive.ObjectID , opts ...*options.FindOneAndReplaceOptions) error
}

type repository struct {
	mongodb.BaseRepository
}

func (b *repository) CreateOne(ctx context.Context, model model.IEntity, opts ...*options.InsertOneOptions) (primitive.ObjectID, error) {
	model.SetID()
	model.SetCreatedAt()
	return b.Create(ctx, model, opts...)
}

func (b *repository) UpdateOneByID(ctx context.Context, model model.IEntity, id  primitive.ObjectID , opts ...*options.FindOneAndReplaceOptions) error {
	model.SetUpdatedAt()
	now := time.Now().UTC()
	return b.Update(ctx, model,id, now, opts...)
}