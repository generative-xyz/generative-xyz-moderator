package repository

import (
	"context"
	"strconv"

	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r Repository) FindTokenUriWithoutCache(filter bson.D) (*entity.TokenUri, error) {
	resp := &entity.TokenUri{}
	usr, err := r.FilterOne(entity.TokenUri{}.TableName(), filter)
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) FindTokenUriWithtCache(filter bson.D, cachedKey string) (*entity.TokenUri, error) {
	// always reload cache
	liveReload := func (filter bson.D, key string) (*entity.TokenUri, error) {
		token, err := r.FindTokenUriWithoutCache(filter)
		if err != nil {
			return nil, err
		}
		//spew.Dump("FindTokenUriWithtCache", filter, token)
		r.Cache.SetData(key, token)
		return token, nil
	}

	go liveReload(filter, cachedKey)
	
	cached, err := r.Cache.GetData(cachedKey)
	if err != nil  || cached == nil{
		return liveReload(filter, cachedKey)
	}

	tok := &entity.TokenUri{}
	err = helpers.ParseCache(cached, tok)
	if err != nil {
		return nil, err
	}

	return  tok, nil
}

func (r Repository) FindTokenBy(contractAddress string, tokenID string) (*entity.TokenUri, error) {
	key := helpers.TokenURIKey(contractAddress, tokenID)
	f := bson.D{{"contract_address", contractAddress}, {"token_id", tokenID}}

	return r.FindTokenUriWithtCache(f, key)
}

func (r Repository) FindTokenByWithoutCache(contractAddress string, tokenID string) (*entity.TokenUri, error) {
	f := bson.D{{"contract_address", contractAddress}, {"token_id", tokenID}}

	return r.FindTokenUriWithoutCache(f)
}

func (r Repository) FindTokenByGenNftAddr(genNftAddrr string, tokenID string) (*entity.TokenUri, error) {
	key := helpers.TokenURIByGenNftAddrKey(genNftAddrr, tokenID)
	f := bson.D{{"gen_nft_addrress", genNftAddrr}, {"token_id", tokenID}}
	return r.FindTokenUriWithtCache(f, key)
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
	
	if filter.Keyword != nil {
		if *filter.Keyword != "" {
			kwInt, err := strconv.Atoi(*filter.Keyword)
			if err == nil {
				f["token_id_mini"] = kwInt
			}else{
				f["token_id_mini"] = *filter.Keyword
			}
			
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
	
	if len(filter.CollectionIDs) > 0 {
		f["gen_nft_addrress"] = bson.D{ {"$in", filter.CollectionIDs} }
 	}
	
	if len(filter.TokenIDs) > 0 {
		f["token_id"] =  bson.D{ {"$in", filter.TokenIDs} }
 	}

	return f
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
	inputData.SetUpdatedAt()
	inputData.SetCreatedAt()
	bData, _ := inputData.ToBson()
	filter := bson.D{{"contract_address", contractAddress}, {"token_id", tokenID}}
	update := bson.D{{"$set",  bData}}
	updateOpts := options.Update().SetUpsert(true)
	//indexOpts := options.CreateIndexes().SetMaxTime(10 * time.Second)


	//id := fmt.Sprintf("%s%s", contractAddress, tokenID)
	result, err := r.DB.Collection(inputData.TableName()).UpdateOne( context.TODO(), filter, update, updateOpts)
	if err != nil {
		return nil, err
	}

	key1 := helpers.TokenURIKey(contractAddress, tokenID)
	key2 := helpers.TokenURIByGenNftAddrKey(inputData.GenNFTAddr, tokenID)
	
	r.Cache.SetData(key1, inputData)
	r.Cache.SetData(key2, inputData)
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