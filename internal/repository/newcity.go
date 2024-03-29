package repository

import (
	"context"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/helpers"
)

func (r Repository) InsertNewCityGm(data *entity.NewCityGm) error {

	data.Address = strings.ToLower(data.Address)
	data.UserAddress = strings.ToLower(data.UserAddress)

	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) FindNewCityGmByUserAddress(userAddress, typeReq string) (*entity.NewCityGm, error) {
	resp := &entity.NewCityGm{}

	filter := bson.D{{"user_address", strings.ToLower(userAddress)}, {"type", typeReq}}

	usr, err := r.FilterOne(entity.NewCityGm{}.TableName(), filter)
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) FindNewCityGmByType(typeReq string) ([]entity.NewCityGm, error) {
	var projects []entity.NewCityGm
	cursor, err := r.DB.Collection(entity.NewCityGm{}.TableName()).Find(context.TODO(), bson.D{{"type", typeReq}})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &projects); err != nil {
		return nil, err
	}

	return projects, nil
}

func (r Repository) ListNewCityGmByStatus(statuses []int) ([]*entity.NewCityGm, error) {
	resp := []*entity.NewCityGm{}
	filter := bson.M{
		//"status": bson.M{"$in": statuses},
	}

	// sort: created_at

	// cursor, err := r.DB.Collection(entity.NewCityGm{}.TableName()).Find(context.TODO(), filter)

	cursor, err := r.DB.Collection(entity.NewCityGm{}.TableName()).Find(context.TODO(), filter, &options.FindOptions{
		Sort: bson.D{{"created_at", 1}},
		// Limit: &limit,
		// Skip:  &offset,
	})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (r Repository) UpdateNewCityGm(newCityGm *entity.NewCityGm) (*mongo.UpdateResult, error) {
	filter := bson.D{{"uuid", newCityGm.UUID}}
	result, err := r.UpdateOne(entity.NewCityGm{}.TableName(), filter, newCityGm)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) UpdateNewCityGmENSAvatar(newCityGm *entity.NewCityGm) (*mongo.UpdateResult, error) {
	filter := bson.D{{"uuid", newCityGm.UUID}}
	result, err := r.DB.Collection(entity.NewCityGm{}.TableName()).UpdateOne(context.TODO(), filter, bson.M{
		"$set": bson.M{
			"avatar": newCityGm.Avatar,
			"ens":    newCityGm.ENS,
		},
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) FindNewCitysGmByUserAddress(address string) ([]entity.NewCityGm, error) {

	address = strings.ToLower(address)

	var projects []entity.NewCityGm
	cursor, err := r.DB.Collection(entity.NewCityGm{}.TableName()).Find(context.TODO(), bson.D{{"user_address", address}})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &projects); err != nil {
		return nil, err
	}

	return projects, nil
}

func (r Repository) SetUpdatedTimeNewCitysGm(tcAddress string) (*mongo.UpdateResult, error) {
	f := bson.D{
		{Key: "user_address", Value: tcAddress},
	}

	update := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}
	result, err := r.DB.Collection(entity.NewCityGm{}.TableName()).UpdateMany(context.TODO(), f, update)
	if err != nil {
		return nil, err
	}

	return result, err
}
