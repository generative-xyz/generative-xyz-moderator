package repository

import (
	"context"
	"fmt"

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

func (r Repository) FindCollectionMetaByProjectId(projectID string) (*entity.CollectionMeta, error) {
	resp := &entity.CollectionMeta{}
	usr, err := r.FilterOne(entity.CollectionMeta{}.TableName(), bson.D{{Key: "project_id", Value: projectID}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r Repository) UpdateCollectionMeta(uuid string, updatedData map[string]interface{}) error {
	collection := entity.CollectionMeta{}.TableName()
	filter := bson.D{{"uuid", uuid}}
	updatingFields := bson.D{}
	for k, v := range updatedData {
		updatingFields = append(updatingFields, bson.E{Key: k, Value: v})
	}
	result, err := r.DB.Collection(collection).UpdateOne(context.TODO(), filter, bson.D{{"$set", updatingFields}})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("no document matched the filter")
	}
	return nil
}

func (r Repository) FindUncreatedCollectionMeta() ([]entity.CollectionMeta, error) {
	metas := []entity.CollectionMeta{}
	f := bson.M{
		"project_created": bson.M{"$ne": true},
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
		{Key: "uuid", Value: meta.UUID},
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

func (r Repository) SetMetaMappedProjectID(meta entity.CollectionMeta, projectID string) error {
	f := bson.D{
		{Key: "uuid", Value: meta.UUID},
	}
	update := bson.M{
		"$set": bson.M{
			"project_id": projectID,
		},
	}

	_, err := r.DB.Collection(meta.TableName()).UpdateOne(context.TODO(), f, update)
	if err != nil {
		return err
	}

	return err
}

func (r Repository) SetMetaProjectExisted(meta entity.CollectionMeta, existed bool) error {
	f := bson.D{
		{Key: "uuid", Value: meta.UUID},
	}
	update := bson.M{
		"$set": bson.M{
			"project_existed": existed,
		},
	}

	_, err := r.DB.Collection(meta.TableName()).UpdateOne(context.TODO(), f, update)
	if err != nil {
		return err
	}

	return err
}
