package repository

import (
	"context"
	"rederinghub.io/utils"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
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
	cursor, err := r.DB.Collection(utils.COLLECTION_PROJECTS).Find(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &projects); err != nil {
		return nil, err
	}

	return projects, nil
}
