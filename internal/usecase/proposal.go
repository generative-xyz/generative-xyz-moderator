package usecase

import (
	"github.com/jinzhu/copier"
	"github.com/opentracing/opentracing-go"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
)

func (u Usecase) CreateDraftProposal(rootSpan opentracing.Span, req structure.CreateProposaltReq) (*entity.ProposalDetail, error) {
	span, log := u.StartSpan("CreateDraftProposal", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	pe := &entity.ProposalDetail{}
	err := copier.Copy(pe, req)
	if err != nil {
		log.Error("copier.Copy", err.Error(), err)
		return nil, err
	}

	pe.IsDraft = true
	err = u.Repo.CreateProposalDetail(pe)
	if err != nil {
		log.Error("u.Repo.CreateProject", err.Error(), err)
		return nil, err
	}

	log.SetData("pe", pe)
	return pe, nil
}

func (u Usecase) MapOffToOnChainProposal(rootSpan opentracing.Span, ID string, proposalID string) (*entity.ProposalDetail, error) {
	span, log := u.StartSpan("MapOffToOnChainProposal", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	
	pD, err := u.Repo.FindProposalDetailByUUID(ID)
	if err != nil {
		log.Error("MapOffToOnChainProposal.FindProposalByID", err.Error(), err)
		return nil, err
	}
	
	pD.ProposalID = proposalID
	pD.IsDraft = false
	updated, err := u.Repo.UpdateProposalDetail(ID, pD)
	if err != nil {
		log.Error("MapOffToOnChainProposal.UpdateProposalDetail", err.Error(), err)
		return nil, err
	}

	log.SetData("updated", updated)	
	return pD, nil
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

	log.SetData("proposals", proposals.Total)
	return proposals, nil
}

func (u Usecase) GetProposal(rootSpan opentracing.Span, proposalID string) (*entity.Proposal, error) {
	span, log := u.StartSpan("GetProposal", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	
	
	proposal, err := u.Repo.FindProposal(proposalID)
	if err != nil {
		log.Error("u.Repo.FilterProposal", err.Error(), err)
		return nil, err
	}

	pDetail, err := u.Repo.FindProposalDetail(proposalID)
	if err == nil {
		proposal.ProposalDetail = *pDetail
	}else{
		log.Error("u.Repo.FilterProposal", err.Error(), err)
	}

	log.SetData("proposal", proposal)
	return proposal, nil
}

func (u Usecase) GetProposalVotes(rootSpan opentracing.Span, filter structure.FilterProposalVote) (*entity.Pagination, error) {
	span, log := u.StartSpan("GetProposalVotes", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	
	f := &entity.FilterProposalVotes{}
	err := copier.Copy(f, filter)
	if err != nil {
		log.Error("filterProposalVotes.copier.Copy", err.Error(), err)
		return nil, err
	}
	
	proposalVotes, err := u.Repo.FilterProposalVotes(*f)
	if err != nil {
		log.Error("u.Repo.FilterProposalVotes", err.Error(), err)
		return nil, err
	}


	log.SetData("proposalVotes", proposalVotes)
	return proposalVotes, nil
}