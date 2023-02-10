package repository

import (
	"context"
	"encoding/json"
	"strings"

	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r Repository) FindProject( projectID string) (*entity.Projects, error) {
	resp := &entity.Projects{}
	usr, err := r.FilterOne(entity.Projects{}.TableName(), bson.D{{utils.KEY_UUID, projectID}})
	if err != nil {
		return nil, err
	}
	
	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) FindProjectByTokenID( tokenID string) (*entity.Projects, error) {
	resp := &entity.Projects{}
	usr, err := r.FilterOne(entity.Projects{}.TableName(), bson.D{{"tokenid", tokenID}})
	if err != nil {
		return nil, err
	}
	
	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) FindProjectBy( contractAddress string, tokenID string) (*entity.Projects, error) {
	resp := &entity.Projects{}
	contractAddress = strings.ToLower(contractAddress)
	go r.findProjectBy(contractAddress, tokenID)
	
	p, err := r.Cache.GetData(helpers.ProjectDetailKey(contractAddress, tokenID))
	if err != nil {
		return r.findProjectBy(contractAddress, tokenID)
	}

	bytes := []byte(*p)
	err = json.Unmarshal(bytes, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) IncreaseProjectIndex(projectID string) (error) {
	filter := bson.D{{Key: "tokenid", Value: projectID}}
	update := bson.M{"$inc": bson.M{"index": 1}}
	_, err := r.DB.Collection(utils.COLLECTION_PROJECTS).UpdateOne(context.TODO(), filter, update)
	return err
}

func (r Repository) FindProjectWithoutCache( contractAddress string, tokenID string) (*entity.Projects, error) {
	contractAddress = strings.ToLower(contractAddress)
	return  r.findProjectBy(contractAddress, tokenID)
}

func (r Repository) findProjectBy( contractAddress string, tokenID string) (*entity.Projects, error) {
	contractAddress = strings.ToLower(contractAddress)
	resp := &entity.Projects{}
	usr, err := r.FilterOne(entity.Projects{}.TableName(), bson.D{{"contractAddress", contractAddress}, {"tokenid", tokenID}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	
	r.Cache.SetData(helpers.ProjectDetailKey(contractAddress, tokenID), resp)
	return resp, nil
}

func (r Repository) CreateProject(data *entity.Projects) error {
	data.ContractAddress = strings.ToLower(data.ContractAddress)
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}

	_ = r.Cache.SetData(helpers.ProjectDetailKey(data.ContractAddress, data.TokenID), data)

	return nil
}

func (r Repository) UpdateProjectImages(ID string, images []string, processingImages []string) (*mongo.UpdateResult, error) {
	filter := bson.D{{utils.KEY_UUID, ID}}
	update := bson.M{
		"$set": bson.M{
			"images": images,
			"processingImages": processingImages,
		},
	}
	
	result, err := r.DB.Collection(utils.COLLECTION_PROJECTS).UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return nil, err	
	}
	
	return result, nil
}

func (r Repository) UpdateProject(ID string, data *entity.Projects) (*mongo.UpdateResult, error) {
	filter := bson.D{{utils.KEY_UUID, ID}}
	result, err := r.UpdateOne(entity.Projects{}.TableName(), filter, data)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r Repository) UpdateProjectMintedCount(ID string, mintedCount int32) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: utils.KEY_UUID, Value: ID}}
	update := bson.M{
		"$set": bson.M{
			"stats.minted_count": mintedCount,
		},
	}
	result, err := r.DB.Collection(utils.COLLECTION_PROJECTS).UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return nil, err	
	}
	
	return result, nil
}

func (r Repository) GetProjects(filter entity.FilterProjects) (*entity.Pagination, error)  {
	confs := []entity.Projects{}
	resp := &entity.Pagination{}
	f := r.FilterProjects(filter)
	
	
	s := r.SortProjects()
	p, err := r.Paginate(utils.COLLECTION_PROJECTS, filter.Page, filter.Limit, f, r.SelectedProjectFields(), s, &confs)
	if err != nil {
		return nil, err
	}
	
	resp.Result = confs
	resp.Page = p.Pagination.Page
	resp.Total = p.Pagination.Total
	resp.PageSize = filter.Limit
	return resp, nil
}

func (r Repository) GetAllProjects(filter entity.FilterProjects) ([]entity.Projects, error)  {
	projects := []entity.Projects{}
	f := r.FilterProjects(filter)
	cursor, err := r.DB.Collection(utils.COLLECTION_PROJECTS).Find(context.TODO(), f)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &projects); err != nil {
		return nil, err
	}

	return projects, nil
}

func (r Repository) GetAllProjectsWithSelectedFields() ([]entity.Projects, error)  {
	projects := []entity.Projects{}
	f := bson.M{}
	opts := options.Find().SetProjection(r.SelectedProjectFields())
	cursor, err := r.DB.Collection(utils.COLLECTION_PROJECTS).Find(context.TODO(), f, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &projects); err != nil {
		return nil, err
	}

	return projects, nil
}

func (r Repository) CountProjects(filter entity.FilterProjects) (*int64, error)  {
	//products := &entity.Products{}
	f := r.FilterProjects(filter)
	count, err := r.DB.Collection(utils.COLLECTION_PROJECTS).CountDocuments(context.TODO(), f)
	if err != nil {
		return nil, err
	}

	return &count, nil
}

func (r Repository) GetMintedOutProjects(filter entity.FilterProjects) (*entity.Pagination, error)  {
	confs := []entity.Projects{}
	resp := &entity.Pagination{}
	f := r.FilterProjects(filter)

	query := `{ "$where": "this.limitSupply == this.index + this.indexReverse " }`
	err := json. Unmarshal([]byte(query), &f)
	if err != nil {
		return nil, err
	}

	s := r.SortProjects()
	p, err := r.Paginate(utils.COLLECTION_PROJECTS, filter.Page, filter.Limit, f , r.SelectedProjectFields(), s, &confs)
	if err != nil {
		return nil, err
	}
	
	resp.Result = confs
	resp.Page = p.Pagination.Page
	resp.Total = p.Pagination.Total
	resp.PageSize = filter.Limit
	return resp, nil
}

func (r Repository) GetRecentWorksProjects(filter entity.FilterProjects) (*entity.Pagination, error)  {
	confs := []entity.Projects{}
	resp := &entity.Pagination{}
	f := r.FilterProjects(filter)

	query := `{ "$where": "this.limitSupply > this.index + this.indexReverse " }`
	err := json. Unmarshal([]byte(query), &f)
	if err != nil {
		return nil, err
	}

	s := r.SortProjects()
	p, err := r.Paginate(utils.COLLECTION_PROJECTS, filter.Page, filter.Limit, f, r.SelectedProjectFields(), s, &confs)
	if err != nil {
		return nil, err
	}
	
	resp.Result = confs
	resp.Page = p.Pagination.Page
	resp.Total = p.Pagination.Total
	resp.PageSize = filter.Limit
	return resp, nil
}

func (r Repository) FilterProjects(filter entity.FilterProjects) bson.M {
	f := bson.M{}
	f["isSynced"] = true
	f[utils.KEY_DELETED_AT] = nil
	f["isHidden"] = false

	if filter.WalletAddress != nil {
		if *filter.WalletAddress != "" {
			f["creatorAddress"] = *filter.WalletAddress
		}
	}

	if filter.Name != nil {
		if *filter.Name != "" {
			f["$text"] =  bson.M{"$search": *filter.Name}
		}
	}
	return f
}

func (r Repository) FindProjectByGenNFTAddr(genNFTAddr string) (*entity.Projects, error) {
	genNFTAddr = strings.ToLower(genNFTAddr)
	resp := &entity.Projects{}
	cached, err := r.Cache.GetData(helpers.ProjectDetailgenNftAddrrKey(genNFTAddr))
	if err == nil && cached != nil {
		err := helpers.ParseCache(cached, resp)
		if err == nil {
			return resp, nil
		}
	}
	
	prj, err := r.FilterOne(entity.Projects{}.TableName(), bson.D{{Key: "genNFTAddr", Value: genNFTAddr}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(prj, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository)SortProjects () []Sort {
	s := []Sort{}
	s = append(s, Sort{SortBy:"priority", Sort: entity.SORT_DESC })
	s = append(s, Sort{SortBy:"created_at", Sort: entity.SORT_DESC })
	return s
}

func (r Repository)SelectedProjectFields () bson.D {
	f := bson.D{
		{"id", 1},
		{"contractAddress", 1},
		{"tokenid", 1},
		{"maxSupply", 1},
		{"limitSupply", 1},
		{"mintPrice", 1},
		{"name", 1},
		{"creatorName", 1},
		{"creatorAddress", 1},
		{"thumbnail", 1},
		{"mintFee", 1},
		{"openMintUnixTimestamp", 1},
		{"closeMintUnixTimestamp", 1},
		{"genNFTAddr", 1},
		{"mintTokenAddress", 1},
		{"minted_time", 1},
		{"license", 1},
		{"description", 1},
		{"stats", 1},
		{"status", 1},
		{"tokenDescription", 1},
		{"index", 1},
		{"indexReverse", 1},
		{"creatorProfile", 1},
		{"images", 1},
		{"mintedImages", 1},
	}
	return f
}
