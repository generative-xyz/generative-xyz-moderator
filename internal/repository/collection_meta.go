package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/helpers"
)

func (r Repository) InsertCollectionMeta(data *entity.CollectionMeta) error {
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) FindCollectionMetaByInscriptionIcon(inscriptionIcon string) (*entity.CollectionMeta, error) {
	resp := &entity.CollectionMeta{}
	usr, err := r.FilterOne(entity.CollectionMeta{}.TableName(), bson.D{{Key: "inscription_icon", Value: inscriptionIcon}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r Repository) FindUncreatedCollectionMeta() ([]entity.CollectionMeta, error) {	
	metas := []entity.CollectionMeta{}
	f := bson.M{
		"project_created" : bson.M{"$ne": true},
	}
	cursor, err := r.DB.Collection(entity.CollectionMeta{}.TableName()).Find(context.TODO(), f)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &metas); err != nil {
		return nil, err
	}

	return metas, nil
}

func (r Repository) SetProjectCreatedMeta(meta entity.CollectionMeta) error {
	f := bson.D{
		{Key: "uuid", Value: meta.UUID,},
	}

	update := bson.M{
		"$set": bson.M{
			"project_created": true,
		},
	}

	_, err := r.DB.Collection(meta.TableName()).UpdateOne(context.TODO(), f, update)
	if err != nil {
		return err
	}

	return err
} 

