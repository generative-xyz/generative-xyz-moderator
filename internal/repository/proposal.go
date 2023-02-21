package repository

import (
	"context"

	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"

	. "github.com/gobeam/mongo-go-pagination"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r Repository) FilterProposal(filter entity.FilterProposals) (*entity.Pagination, error) {
	pro := []entity.Proposal{}
	resp := &entity.Pagination{}

	if filter.SortBy == "" {
		filter.SortBy = "created_at"
	}
	f := r.filterProposal(filter)

	t, err := r.AggregateData(entity.Proposal{}.TableName(), filter.Page, filter.Limit, f, r.SelectedProposalFields() , r.SortProposal(filter), &pro)
	if err != nil {
		return nil, err
	}

	for _, raw := range t.Data {
		p :=&entity.QueriedProposal{}
		pResp :=&entity.Proposal{}
		marshallErr := bson.Unmarshal(raw, &p)
		if marshallErr == nil {
			err = copier.Copy(pResp, p)
			if err == nil {
				if len(p.ProposalDetail) > 0 {
					d := p.ProposalDetail[0]
					pResp.ProposalDetail = d
				}
	
				pro = append(pro, *pResp)
			}
		}

	}

	resp.Result = pro
	resp.Page = t.Pagination.Page
	resp.Total = t.Pagination.Total
	resp.PageSize = filter.Limit
	return resp, nil
}

func (r Repository)  ProcessProposals() ([]entity.Proposal, error) {
	pro := []entity.Proposal{}

	f := bson.M{}
	f[utils.KEY_DELETED_AT] = nil
	f["state"] = bson.M{"$nin": bson.A{
		entity.	PStateCanceled,
		entity.PStateDefeated,
		entity.PStateSuccesseded,
		entity.PStateExpired,
		entity.PStateExecuted,
	}}
	
	cursor, err := r.DB.Collection(utils.COLLECTION_DAO_PROPOSAL).Find(context.TODO(), f)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &pro); err != nil {
		return nil, err
	}

	return pro, nil
}

func (r Repository) AllProposals(filter entity.FilterProposals) ([]entity.Proposal, error) {
	pro := []entity.Proposal{}
	
	if filter.SortBy == "" {
		filter.SortBy = "created_at"
	}

	f := r.filterProposal(filter)
	cursor, err := r.DB.Collection(utils.COLLECTION_DAO_PROPOSAL).Find(context.TODO(), f)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &pro); err != nil {
		return nil, err
	}

	return pro, nil
}

func (r Repository) AggregateData(dbName string, page int64, limit int64, filter interface{}, selectFields interface{}, sorts []Sort, returnData interface{}) (*PaginatedData, error) {

	paginatedData := New(r.DB.Collection(dbName)).
		Context(context.TODO()).
		Limit(int64(limit)).
		Page(int64(page))

	if 	len(sorts) > 0 {
		for _, sort := range sorts {
			if sort.Sort == entity.SORT_ASC || sort.Sort == entity.SORT_DESC {
				//sortValue := bson.D{{"created_at", -1}}
				paginatedData.Sort(sort.SortBy, sort.Sort)
			}	
		}
	}else{
		paginatedData.Sort("created_at", entity.SORT_DESC)
	}

	lookUpStage := bson.M{
		"$lookup": bson.M{
			"from": "proposal_detail",
			"let": bson.M{"proposalID": "$proposalID"},
			"pipeline": bson.A{
				bson.M{
					"$match": bson.M{
						"$expr": bson.M{"$eq": bson.A{"$proposalID",  "$$proposalID"}},
					},
				},
			},
			"as": "proposalDetail",
		},
	}

	matchStage := bson.M{
		"$match" : filter,
	}

	data, err :=	paginatedData.Aggregate(lookUpStage, matchStage)

	if err != nil {
		return nil, err
	}
	
	return data, nil
}

func (r Repository) filterProposal(filter entity.FilterProposals) bson.M {
	f := bson.M{}
	f[utils.KEY_DELETED_AT] = nil

	if filter.Proposer != nil {
		if *filter.Proposer  != "" {
			f["proposer"] = *filter.Proposer
		}
	}
	
	if filter.State != nil {
		f["state"] = *filter.State
	}

	if filter.ProposalID != nil {
		if *filter.ProposalID  != "" {
			f["proposalID"] = *filter.ProposalID
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
		{"currentBlock", 1},
		{"description", 1},
		{"raw", 1},
		{"state", 1},
		{"vote", 1},
		{"proposalDetail.amount", 1},
		{"proposalDetail.receiverAddress", 1},
		{"proposalDetail.title", 1},
		{"proposalDetail.tokenType", 1},
	}
	return f
}

func (r Repository) SortProposal (filter entity.FilterProposals) []Sort {
	s := []Sort{}
	s = append(s, Sort{SortBy: filter.SortBy, Sort: filter.Sort })
	return s
}

func (r Repository) CreateProposal(obj *entity.Proposal) error {
	err := r.InsertOne(obj.TableName(), obj)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) FindProposal(proposalID string) (*entity.Proposal, error) {
	resp := &entity.Proposal{}
	usr, err := r.FilterOne(entity.Proposal{}.TableName(), bson.D{{"proposalID", proposalID}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}


func (r Repository) FindProposalByUUID(UUID string) (*entity.Proposal, error) {
	resp := &entity.Proposal{}
	usr, err := r.FilterOne(entity.Proposal{}.TableName(), bson.D{{utils.KEY_UUID, UUID}})
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