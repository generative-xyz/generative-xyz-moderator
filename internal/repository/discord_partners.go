package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/internal/entity"
)

func (r Repository) GetAllDiscordPartner() ([]entity.DiscordPartner, error) {
	partners := []entity.DiscordPartner{}
	f := bson.M{}

	cursor, err := r.DB.Collection(entity.DiscordPartner{}.TableName()).Find(context.TODO(), f)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &partners); err != nil {
		return nil, err
	}

	return partners, nil
}
