package repository

import (
	"context"
	"time"

	"github.com/davecgh/go-spew/spew"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"rederinghub.io/internal/entity"
)
func (r Repository) CreateTokenURIIndexModel() ([]string, error) {
	collection := entity.TokenUri{}.TableName()

	

 	models :=  []mongo.IndexModel{
		{ Keys: bson.M{"gen_nft_addrress":  -1,}, Options:  options.Index().SetName("gen_nft_addrress_desc") ,} ,
		{ Keys: bson.M{"gen_nft_addrress":  "text",}, Options:  options.Index().SetName("gen_nft_addrress_i_text") ,} ,
		
		{ Keys: bson.M{"project_id":  -1,}, Options:  options.Index().SetName("project_id_desc") ,} ,
		{ Keys: bson.M{"project_id_int":  -1,}, Options:  options.Index().SetName("project_id_int_desc") ,} ,
		{ Keys: bson.M{"owner_addrress":  -1,}, Options:  options.Index().SetName("owner_addrress_desc") ,} ,		
		{ Keys: bson.M{"creator_address":  -1,}, Options:  options.Index().SetName("creator_address_desc") ,} ,
		
		{ Keys: bson.D{ {Key: "gen_nft_addrress", Value: -1}, {Key: "token_id", Value: -1} }, Options: options.Index().SetUnique(true),} ,
		{ Keys: bson.D{ {Key: "contract_address", Value: -1}, {Key: "token_id", Value: -1} }, Options: options.Index().SetUnique(true),} ,
	}

	return r.CreateIndexes(collection, models)
}

func (r Repository) CreateProjectIndexModel() ([]string, error) {
	collection := entity.Projects{}.TableName()
 	models :=  []mongo.IndexModel{
		{ Keys: bson.M{"tokenid": -1,}, Options: nil,} ,
		{ Keys: bson.M{"tokenIDInt": -1,}, Options: nil,} ,
		{ Keys: bson.M{"contractAddress":  -1,}, Options: nil,} ,
		{ Keys: bson.M{"creatorName": -1,}, Options: nil,} ,
		{ Keys: bson.M{"creatorAddress": -1,}, Options: nil,} ,
		{ Keys: bson.M{"genNFTAddr": -1,}, Options: nil,} ,
	}

	return r.CreateIndexes(collection, models)
}

func (r Repository) CreateIndexes(collectionName string, models []mongo.IndexModel) ([]string, error) {
 	col := r.DB.Collection(collectionName)
	
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	ind, err := col.Indexes().CreateMany(ctx, models, opts)
	if err != nil {
		return nil, err
	}

	spew.Dump(ind)
	return ind, nil
}