package repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
)

func (r Repository) CreateProposalVotes(obj *entity.ProposalVotes) error {
	err := r.InsertOne(obj.TableName(), obj)
	if err != nil {
		return err
	}
	
	return  nil
}

func (r Repository) FilterProposalVotes(filter entity.FilterProposalVotes) (*entity.Pagination, error) {
	pro := []entity.ProposalVotes{}
	resp := &entity.Pagination{}
	
	f := r.filterProposalVotes(filter)
	if filter.SortBy == "" {
		filter.SortBy = "created_at"
	}

	t, err := r.Paginate(entity.ProposalVotes{}.TableName(), filter.Page, filter.Limit, f, r.SelectedProposalVoteFields() , r.SortProposalVotes(filter), &pro)
	if err != nil {
		return nil, err
	}
	
	resp.Result = pro
	resp.Page = t.Pagination.Page
	resp.Total = t.Pagination.Total
	resp.PageSize = filter.Limit
	return resp, nil
}

func (r Repository) SelectedProposalVoteFields() bson.D {
	f := bson.D{
		{"uuid", 1},
		{"proposalID", 1},
		{"voter", 1},
		{"support", 1},
		{"weight", 1},
		{"created_at", 1},
		{"reason", 1},
		
	}
	return f
}

func (r Repository) SortProposalVotes (filter entity.FilterProposalVotes) []Sort {
	s := []Sort{}
	s = append(s, Sort{SortBy: filter.SortBy, Sort: filter.Sort })
	return s
}

func (r Repository) filterProposalVotes(filter entity.FilterProposalVotes) bson.M {
	f := bson.M{}
	f[utils.KEY_DELETED_AT] = nil

	if filter.Voter != nil {
		if *filter.Voter  != "" {
			f["voter"] = *filter.Voter
		}
	}
	
	if filter.ProposalID != nil {
		f["proposalID"] = *filter.ProposalID
	}

	if filter.Support != nil {
		f["support"] = *filter.Support
	}

	return f
}
