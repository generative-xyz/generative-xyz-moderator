package repository

import (
	"context"

	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r Repository) FilterProposal(filter entity.FilterProposals) (*entity.Pagination, error) {
	pro := []entity.Proposal{}
	resp := &entity.Pagination{}
	
	f := r.filterProposal(filter)
	if filter.SortBy == "" {
		filter.SortBy = "created_at"
	}

	t, err := r.Paginate(entity.Proposal{}.TableName(), filter.Page, filter.Limit, f, r.SelectedProposalFields() , r.SortProposal(filter), &pro)
	if err != nil {
		return nil, err
	}
	
	resp.Result = pro
	resp.Page = t.Pagination.Page
	resp.Total = t.Pagination.Total
	resp.PageSize = filter.Limit
	return resp, nil
}

func (r Repository) filterProposal(filter entity.FilterProposals) bson.M {
	f := bson.M{}
	f[utils.KEY_DELETED_AT] = nil

	if filter.Proposer != nil {
		if *filter.Proposer  != "" {
			f["proposer"] = *filter.Proposer
		}
	}

	return f
}

func (r Repository) SelectedProposalFields() bson.D {
	f := bson.D{
		{"uuid", 1},
		{"proposalID", 1},
		{"proposer", 1},
		{"targets", 1},
		{"values", 1},
		{"signatures", 1},
		{"calldatas", 1},
		{"startBlock", 1},
		{"endBlock", 1},
		{"description", 1},
		{"raw", 1},
		{"state", 1},
	}
	return f
}

func (r Repository) SortProposal (filter entity.FilterProposals) []Sort {
	s := []Sort{}
	s = append(s, Sort{SortBy: filter.SortBy, Sort: filter.Sort })
	return s
}

func (r Repository) CreateProposal(obj *entity.Proposal) error {

	obj.SetID()
	obj.SetCreatedAt()
	obj.ProposalID = obj.UUID


	bData, err := obj.ToBson()
	if err != nil {
		return  err
	}

	_ , err = r.DB.Collection(obj.TableName()).InsertOne(context.TODO(), &bData)
	if err != nil {
		return err
	}

	err = obj.Decode(bData)
	if err != nil {
		return  err
	}

	return  nil
}

func (r Repository) FindProposal(pID string) (*entity.Proposal, error) {
	resp := &entity.Proposal{}
	usr, err := r.FilterOne(entity.Proposal{}.TableName(), bson.D{{"proposalID", pID}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}


func (r Repository) FindProposalByID(pID string) (*entity.Proposal, error) {
	resp := &entity.Proposal{}
	usr, err := r.FilterOne(entity.Proposal{}.TableName(), bson.D{{utils.KEY_UUID, pID}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) UpdateProposal(ID string, data *entity.Proposal) (*mongo.UpdateResult, error) {
	filter := bson.D{{utils.KEY_UUID, ID}}
	result, err := r.UpdateOne(entity.Proposal{}.TableName(), filter, data)
	if err != nil {
		return nil, err
	}
	return result, nil
}