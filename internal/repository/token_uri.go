package repository

import (
	"context"
	"errors"
	"strings"

	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (r Repository) FilterTokenUri(filter entity.FilterTokenUris) (*entity.Pagination, error) {
	tokens := []entity.TokenUri{}
	resp := &entity.Pagination{}
	
	f := r.filterToken(filter)
	
	t, err := r.Paginate(entity.TokenUri{}.TableName(), filter.Page, filter.Limit, f, filter.SortBy, filter.Sort, &tokens)
	if err != nil {
		return nil, err
	}
	
	resp.Result = tokens
	resp.Page = t.Pagination.Page
	resp.Total = t.Pagination.Total
	return resp, nil
}

func (r Repository) filterToken(filter entity.FilterTokenUris) bson.M {
	f := bson.M{}
	f[utils.KEY_DELETED_AT] = nil


	if filter.CreatorAddr != nil {
		if *filter.CreatorAddr != "" {
			f["creator_address"] = *filter.CreatorAddr
		}
 	}
	
	if filter.OwnerAddr != nil {
		if *filter.OwnerAddr != "" {
			f["owner_addrress"] = *filter.OwnerAddr
		}
 	}
	
	if filter.GenNFTAddr != nil {
		if *filter.GenNFTAddr != "" {
			f["gen_nft_addrress"] = *filter.GenNFTAddr
		}
 	}
	
	if filter.ContractAddress != nil {
		if *filter.ContractAddress != "" {
			f["contract_address"] = *filter.ContractAddress
		}
 	}

	return f
}

func (r Repository) CreateTokenURI(data *entity.TokenUri) error {
	t, err := r.FindTokenBy(data.ContractAddress, data.TokenID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments)  {
			
			err = r.InsertOne(data.TableName(), data)
			if err != nil {
				return err
			}

		}
		
		return  err
	}

	data =  t
	return nil
}

func (r Repository) UpdateTokenByID(tokenID string, inputData *entity.TokenUri) (*mongo.UpdateResult, error) {
	filter := bson.D{{utils.KEY_UUID, tokenID}}
	result, err := r.UpdateOne(inputData.TableName(), filter, inputData)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) GetAllTokens() ([]entity.TokenUri, error)  {
	tokens := []entity.TokenUri{}
	
	f := bson.M{}
	f[utils.KEY_DELETED_AT] = nil
	cursor, err := r.DB.Collection(utils.COLLECTION_TOKEN_URI).Find(context.TODO(), f)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &tokens); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (r Repository) UpdateOrInsertTokenUri(contractAddress string, tokenID string, inputData *entity.TokenUri) (*mongo.UpdateResult, error) {
	if inputData.UUID == ""  {
		inputData.SetID()
		inputData.SetCreatedAt()
	}

	inputData.SetUpdatedAt()
	bData, _ := inputData.ToBson()
	filter := bson.D{{"contract_address", contractAddress}, {"token_id", tokenID}}
	update := bson.D{{"$set",  bData}}
	opts := options.Update().SetUpsert(true)

	//id := fmt.Sprintf("%s%s", contractAddress, tokenID)
	result, err := r.DB.Collection(inputData.TableName()).UpdateOne( context.TODO(), filter, update, opts)
	if err != nil {
		return nil, err
	}
	return result, nil
}


func (r Repository) GetAllTokensByProjectID(projectID string) ([]entity.TokenUri, error) {
	tokens := []entity.TokenUri{}
	f := bson.D{{
		Key: utils.KEY_PROJECT_ID, 
		Value: projectID,
	}}

	cursor, err := r.DB.Collection(utils.COLLECTION_TOKEN_URI).Find(context.TODO(), f)
	if err != nil {
		return nil, err
	}

	if err = cursor.All((context.TODO()), &tokens); err != nil {
		return nil, err
	}

	return tokens, nil
}
