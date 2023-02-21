package repository

import (
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
