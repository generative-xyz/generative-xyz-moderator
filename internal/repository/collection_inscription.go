package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		"collection_inscription_icon": inscriptionIcon,
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

// Only get 1000 docs each time.
func (r Repository) FindUncreatedCollectionInscription(page int) ([]entity.CollectionInscription, error) {
	inscriptions := make([]entity.CollectionInscription, 0)
	f := bson.M{
		"token_created": false,
	}
	opts := options.Find().SetLimit(1000)
	opts.SetSkip(int64(page * 1000))
	cursor, err := r.DB.Collection(entity.CollectionInscription{}.TableName()).Find(context.TODO(), f, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &inscriptions); err != nil {
		return nil, err
	}

	return inscriptions, nil
}

func (r Repository) SetTokenCreatedInscription(inscription entity.CollectionInscription) error {
	f := bson.D{
		{Key: "uuid", Value: inscription.UUID},
	}

	update := bson.M{
		"$set": bson.M{
			"token_created": true,
		},
	}

	_, err := r.DB.Collection(inscription.TableName()).UpdateOne(context.TODO(), f, update)
	if err != nil {
		return err
	}

	return err
}

func (r Repository) CountCollectionInscriptionByInscriptionIcon(inscriptionIcon string) (*int64, error) {
	f := bson.M{
		"collection_inscription_icon": inscriptionIcon,
	}
	count, err := r.DB.Collection(utils.COLLECTION_COLLECTION_INSCRIPTION).CountDocuments(context.TODO(), f)
	if err != nil {
		return nil, err
	}

	return &count, nil
}
