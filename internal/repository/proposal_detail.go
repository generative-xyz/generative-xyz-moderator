package repository

import (
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r Repository) FilterProposalDetail(filter entity.FilterProposals) (*entity.Pagination, error) {
	pro := []entity.ProposalDetail{}
	resp := &entity.Pagination{}
	
	f := r.filterProposal(filter)
	if filter.SortBy == "" {
		filter.SortBy = "created_at"
	}

	t, err := r.Paginate(entity.ProposalDetail{}.TableName(), filter.Page, filter.Limit, f, r.SelectedProposalFields() , r.SortProposal(filter), &pro)
	if err != nil {
		return nil, err
	}
	
	resp.Result = pro
	resp.Page = t.Pagination.Page
	resp.Total = t.Pagination.Total
	resp.PageSize = filter.Limit
	return resp, nil
}

func (r Repository) filterProposalDetail(filter entity.FilterProposals) bson.M {
	f := bson.M{}
	f[utils.KEY_DELETED_AT] = nil

	if filter.Proposer != nil {
		if *filter.Proposer  != "" {
			f["proposer"] = *filter.Proposer
		}
	}

	return f
}

func (r Repository) SelectedProposalDetailFields() bson.D {
	f := bson.D{
		{"uuid", 1},
		{"raw", 1},
		{"state", 1},
	}
	return f
}

func (r Repository) SortProposalDetail (filter entity.FilterProposals) []Sort {
	s := []Sort{}
	s = append(s, Sort{SortBy: filter.SortBy, Sort: filter.Sort })
	return s
}

func (r Repository) CreateProposalDetail(obj *entity.ProposalDetail) error {
	err := r.InsertOne(obj.TableName(), obj)
	if err != nil {
		return err
	}
	
	return  nil
}

func (r Repository) FindProposalDetail(proposalID string) (*entity.ProposalDetail, error) {
	resp := &entity.ProposalDetail{}
	usr, err := r.FilterOne(entity.ProposalDetail{}.TableName(), bson.D{{"proposalID", proposalID}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}


func (r Repository) FindProposalDetailByUUID(UUID string) (*entity.ProposalDetail, error) {
	resp := &entity.ProposalDetail{}
	usr, err := r.FilterOne(entity.ProposalDetail{}.TableName(), bson.D{{utils.KEY_UUID, UUID}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) UpdateProposalDetail(ID string, data *entity.ProposalDetail) (*mongo.UpdateResult, error) {
	filter := bson.D{{utils.KEY_UUID, ID}}
	result, err := r.UpdateOne(entity.ProposalDetail{}.TableName(), filter, data)
	if err != nil {
		return nil, err
	}
	return result, nil
}