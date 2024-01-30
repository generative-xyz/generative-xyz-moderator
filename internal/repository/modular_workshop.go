package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"rederinghub.io/internal/entity"
)

func (r Repository) GetListModularWorkShopByAddress(ctx context.Context, ownerAddr string, offset, limit int64) ([]entity.ModularWorkshopEntity, error) {
	filter := bson.M{
		"delete_at": bson.M{"$exists": false},
	}
	if len(ownerAddr) > 0 {
		filter["owner_addr"] = ownerAddr
	}
	projection := bson.M{"name": 1, "created_at": 1}
	options := options.Find().SetSort(bson.M{"created_at": -1})
	options.SetProjection(projection)
	options.SetSkip(offset)
	options.SetLimit(limit)
	cursor, err := r.DB.Collection(entity.ModularWorkshopEntity{}.TableName()).Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	data := []entity.ModularWorkshopEntity{}
	if err = cursor.All(ctx, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (r Repository) GetTotalModularWorkShopByAddress(ctx context.Context, ownerAddr string) (int64, error) {
	filter := bson.M{"owner_addr": ownerAddr,
		"delete_at": bson.M{"$exists": false},
	}
	total, err := r.DB.Collection(entity.ModularWorkshopEntity{}.TableName()).CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r Repository) SaveModularWorkshop(ctx context.Context, data *entity.ModularWorkshopEntity) (primitive.ObjectID, error) {
	result, err := r.DB.Collection(entity.ModularWorkshopEntity{}.TableName()).InsertOne(ctx, data)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func (r Repository) UpdateModularWorkshop(ctx context.Context, data *entity.ModularWorkshopEntity) error {
	filter := bson.M{"_id": data.ID,
		"owner_addr": data.OwnerAddr,
	}
	update := bson.M{
		"$set": bson.M{
			"update_at": data.UpdatedAt,
			"meta_data": data.MetaData,
			"name":      data.Name,
		},
	}
	_, err := r.DB.Collection(entity.ModularWorkshopEntity{}.TableName()).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
func (r Repository) RemoveModularWorkshop(ctx context.Context, id, ownerAddr string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{
		"_id":        objectID,
		"owner_addr": ownerAddr,
		"delete_at":  bson.M{"$exists": false},
	}
	_, err = r.DB.Collection(entity.ModularWorkshopEntity{}.TableName()).DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) GetModularWorkshopById(ctx context.Context, id string) (*entity.ModularWorkshopEntity, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objectID}
	var data entity.ModularWorkshopEntity
	err = r.DB.Collection(entity.ModularWorkshopEntity{}.TableName()).FindOne(ctx, filter).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
