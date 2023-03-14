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
		if err := repo.CreateCompoundIndex(ctx, entity.DaoArtistVoted{}.TableName(), []string{"dao_artist_id", "created_by"}, true); err != nil {
			return err
		}
		return nil
	}, func(db *mongo.Database) error {
		return nil
	})
}
