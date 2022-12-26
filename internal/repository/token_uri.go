package repository

import (
	"strings"

	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r Repository) FindTokenBy(contractAddress string, tokenID string) (*entity.TokenUri, error) {
	resp := &entity.TokenUri{}
	contractAddress = strings.ToLower(contractAddress)
	usr, err := r.FilterOne(entity.TokenUri{}.TableName(), bson.D{{"contract_address", contractAddress}, {"token_id", tokenID}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) CreateTokenURI(data *entity.TokenUri) error {

	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}

	return nil
}


func (r Repository) UpdateTokenByID(tokenUri string, updateddUser *entity.TokenUri) (*mongo.UpdateResult, error) {
	filter := bson.D{{utils.KEY_UUID, tokenUri}}
	result, err := r.UpdateOne(updateddUser.TableName(), filter, updateddUser)
	if err != nil {
		return nil, err
	}

	return result, nil
}
