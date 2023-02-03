package usecase

import (
	"github.com/jinzhu/copier"
	"github.com/opentracing/opentracing-go"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
)

func (u Usecase) CreateDraftProposal(rootSpan opentracing.Span, req structure.CreateProposaltReq) (*entity.Proposal, error) {
	span, log := u.StartSpan("CreateDraftProposal", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	pe := &entity.Proposal{}
	err := copier.Copy(pe, req)
	if err != nil {
		log.Error("copier.Copy", err.Error(), err)
		return nil, err
	}

	pe.IsDraft = true
	err = u.Repo.CreateProposal(pe)
	if err != nil {
		log.Error("u.Repo.CreateProject", err.Error(), err)
		return nil, err
	}

	log.SetData("pe", pe)
	return pe, nil
}

func (u Usecase) GetProposals(rootSpan opentracing.Span, req structure.FilterProposal) (*entity.Pagination, error) {
	span, log := u.StartSpan("GetProposals", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	
	pe := &entity.FilterProposals{}
	err := copier.Copy(pe, req)
	if err != nil {
		log.Error("copier.Copy", err.Error(), err)
		return nil, err
	}
	
	proposals, err := u.Repo.FilterProposal(*pe)
	if err != nil {
		log.Error("u.Repo.FilterProposal", err.Error(), err)
		return nil, err
	}

	log.SetData("projects", proposals)
	return proposals, nil
}


func (u Usecase) GetProposal(rootSpan opentracing.Span, proposalID string) (*entity.Proposal, error) {
	span, log := u.StartSpan("GetProposal", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	
	
	proposals, err := u.Repo.FindProposal(proposalID)
	if err != nil {
		log.Error("u.Repo.FilterProposal", err.Error(), err)
		return nil, err
	}

	log.SetData("projects", proposals)
	return proposals, nil
}