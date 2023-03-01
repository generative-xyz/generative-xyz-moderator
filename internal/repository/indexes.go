package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"rederinghub.io/internal/entity"
)

func (r Repository) CreateTokenURIIndexModel() ([]string, error) {
	collection := entity.TokenUri{}.TableName()

	models := []mongo.IndexModel{
		{Keys: bson.M{"gen_nft_addrress": -1}, Options: options.Index().SetName("gen_nft_addrress_desc")},
		//{ Keys: bson.M{"gen_nft_addrress":  "text",}, Options:  options.Index().SetName("gen_nft_addrress_i_text") ,} ,
		{Keys: bson.M{"project_id": -1}, Options: options.Index().SetName("project_id_desc")},
		{Keys: bson.M{"project_id_int": -1}, Options: options.Index().SetName("project_id_int_desc")},
		{Keys: bson.M{"owner_addrress": -1}, Options: options.Index().SetName("owner_addrress_desc")}, {Keys: bson.M{"creator_address": -1}, Options: options.Index().SetName("creator_address_desc")},
		{Keys: bson.M{"created_at": -1}, Options: options.Index().SetName("created_at_desc")},
		{Keys: bson.M{"minted_time": -1}, Options: options.Index().SetName("minted_time_desc")},
		{Keys: bson.M{"priority": -1}, Options: options.Index().SetName("tk_priority_desc")},
		{Keys: bson.D{{Key: "gen_nft_addrress", Value: -1}, {Key: "token_id", Value: -1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "contract_address", Value: -1}, {Key: "token_id", Value: -1}}, Options: options.Index().SetUnique(true)},
	}

	return r.CreateIndexes(collection, models)
}

func (r Repository) CreateProjectIndexModel() ([]string, error) {
	collection := entity.Projects{}.TableName()
	models := []mongo.IndexModel{
		{Keys: bson.M{"tokenid": -1}, Options: options.Index().SetUnique(true).SetName("pr_tokenID_desc")},
		{Keys: bson.M{"tokenIDInt": -1}, Options: nil},
		{Keys: bson.M{"contractAddress": -1}, Options: nil},
		{Keys: bson.M{"creatorName": -1}, Options: nil},
		{Keys: bson.M{"creatorAddress": -1}, Options: nil},
		{Keys: bson.M{"genNFTAddr": -1}, Options: nil},
		{Keys: bson.M{"created_at": -1}, Options: nil},
		{Keys: bson.M{"priority": -1}, Options: options.Index().SetName("pr_priority_desc")},
	}

	return r.CreateIndexes(collection, models)
}

func (r Repository) CreateProposalIndexModel() ([]string, error) {
	collection := entity.Proposal{}.TableName()
	models := []mongo.IndexModel{
		{Keys: bson.M{"proposalID": -1}, Options: options.Index().SetUnique(true).SetName("dao_proposalID_desc")},
	}

	return r.CreateIndexes(collection, models)
}

func (r Repository) CreateProposalVotesIndexModel() ([]string, error) {
	collection := entity.ProposalVotes{}.TableName()
	models := []mongo.IndexModel{
		{Keys: bson.M{"voter": -1}, Options: options.Index().SetName("pvotes_voter_desc")},
		{Keys: bson.M{"proposalID": -1}, Options: options.Index().SetName("pvotes_proposalID_desc")},
		{Keys: bson.M{"support": -1}, Options: options.Index().SetName("pvotes_support_desc")},
		{Keys: bson.D{{Key: "proposalID", Value: -1}, {Key: "voter", Value: -1}}, Options: options.Index().SetUnique(true)},
	}

	return r.CreateIndexes(collection, models)
}

func (r Repository) CreateMarketplaceListingsIndexModel() ([]string, error) {
	collection := entity.MarketplaceListings{}.TableName()
	models := []mongo.IndexModel{
		{Keys: bson.M{"created_at": -1}, Options: nil},
	}
	return r.CreateIndexes(collection, models)
}

func (r Repository) CreateMarketplaceOffersIndexModel() ([]string, error) {
	collection := entity.MarketplaceOffers{}.TableName()
	models := []mongo.IndexModel{
		{Keys: bson.M{"created_at": -1}, Options: nil},
	}
	return r.CreateIndexes(collection, models)
}

func (r Repository) CreateBTCWalletIndexModel() ([]string, error) {
	collection := entity.BTCWalletAddress{}.TableName()

	models := []mongo.IndexModel{
		{Keys: bson.M{"user_address": -1}, Options: options.Index().SetName("btc_user_address_desc")},
		{Keys: bson.M{"ordAddress": -1}, Options: options.Index().SetName("btc_ordAddress_desc")},
	}

	return r.CreateIndexes(collection, models)
}

func (r Repository) CreateMintBTCCIndexModel() ([]string, error) {
	collection := entity.MintNftBtc{}.TableName()
	models := []mongo.IndexModel{
		{Keys: bson.M{"projectID": -1}, Options: options.Index().SetName("mbtc_projectID_desc")},
		{Keys: bson.M{"payType": -1}, Options: options.Index().SetName("mbtc_paytype_desc")},
		{Keys: bson.M{"inscriptionID": -1}, Options: options.Index().SetName("mbtc_ins_desc")},
		{Keys: bson.M{"user_address": -1}, Options: options.Index().SetName("mbtc_uaddr_desc")},
	}

	return r.CreateIndexes(collection, models)
}

func (r Repository) CreateWalletTrackTxIndexModel() ([]string, error) {
	collection := entity.WalletTrackTx{}.TableName()
	models := []mongo.IndexModel{
		{Keys: bson.M{"txhash": 1}, Options: options.Index().SetName("txhash_asc").SetUnique(true)},
		{Keys: bson.M{"address": 1}, Options: options.Index().SetName("address_asc")},
	}

	return r.CreateIndexes(collection, models)
}

func (r Repository) CreateReferalIndexModel() ([]string, error) {
	collection := entity.Referral{}.TableName()
	models := []mongo.IndexModel{
		{Keys: bson.M{"referrer_id": -1}, Options: options.Index().SetName("ref_referrer_id_desc")},
		{Keys: bson.M{"referree_id": -1}, Options: options.Index().SetName("ref_referree_id_desc")},
		{Keys: bson.D{{Key: "referrer_id", Value: -1}, {Key: "referree_id", Value: -1}}, Options: options.Index().SetUnique(true)},
	}

	return r.CreateIndexes(collection, models)
}

func (r Repository) CreateVolumnIndexModel() ([]string, error) {
	collection := entity.UserVolumn{}.TableName()
	models := []mongo.IndexModel{
		{Keys: bson.M{"amountType": -1}, Options: options.Index().SetName("vln_amountType_desc")},
		{Keys: bson.M{"creatorAddress": -1}, Options: options.Index().SetName("vln_userID_desc")},
		{Keys: bson.M{"projectID": -1}, Options: options.Index().SetName("vln_projectID_desc")},
		{Keys: bson.D{{Key: "amountType", Value: -1}, {Key: "creatorAddress", Value: -1}, {Key: "projectID", Value: -1}}, Options: options.Index().SetUnique(true)},
	}

	return r.CreateIndexes(collection, models)
}

// func (r Repository) CreateCategoryIndexModel() ([]string, error) {
// 	collection := entity.Categories{}.TableName()
// 	models := []mongo.IndexModel{
// 		{Keys: bson.M{"priority": -1}, Options: options.Index().SetName("cat_priority_desc")},

// 	}

// 	return r.CreateIndexes(collection, models)
// }

func (r Repository) CreateIndexes(collectionName string, models []mongo.IndexModel) ([]string, error) {
	col := r.DB.Collection(collectionName)
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	ind, err := col.Indexes().CreateMany(ctx, models, opts)
	if err != nil {
		return nil, err
	}
	return ind, nil
}

func (r Repository) CreateIndices(ctx context.Context, collectionName string, indices []string, unique bool, opts ...*options.CreateIndexesOptions) error {
	var indexModels []mongo.IndexModel
	for _, index := range indices {
		indexModel := mongo.IndexModel{
			Keys: bsonx.Doc{{Key: index,
				Value: bsonx.Int32(1)}},
			Options: options.Index().SetUnique(unique),
		}
		indexModels = append(indexModels, indexModel)
	}
	_, err := r.DB.Collection(collectionName).Indexes().CreateMany(ctx, indexModels, opts...)
	return err
}
