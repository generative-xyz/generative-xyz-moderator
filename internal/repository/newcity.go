package repository

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
		"status": bson.M{"$in": statuses},
	}

	cursor, err := r.DB.Collection(entity.NewCityGm{}.TableName()).Find(context.TODO(), filter)
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
