package migrate

import (
	"context"

	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/mongo"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/repository"
)

func init() {
	_ = migrate.Register(func(db *mongo.Database) error {
		repo := &repository.Repository{
			DB: db,
		}
		ctx := context.Background()
		_, err := repo.DB.Collection(entity.DaoProject{}.TableName()).Indexes().DropAll(ctx)
		if err != nil {
			return err
		}
		if err := repo.CreateIndices(ctx, entity.DaoProject{}.TableName(), []string{"project_id", "created_by"}, false); err != nil {
			return err
		}
		return nil
	}, func(db *mongo.Database) error {
		return nil
	})
}
