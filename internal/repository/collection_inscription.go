package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
)

func (r Repository) InsertCollectionInscription(data *entity.CollectionInscription) error {
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) FindCollectionInscriptionByInscriptionIcon(inscriptionIcon string) ([]entity.CollectionInscription, error) {
	inscriptions := []entity.CollectionInscription{}

	f := bson.M{
		"collection_inscription_icon" : inscriptionIcon,
	}
	cursor, err := r.DB.Collection(utils.COLLECTION_COLLECTION_INSCRIPTION).Find(context.TODO(), f)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &inscriptions); err != nil {
		return nil, err
	}

	return inscriptions, nil
}
