package usecase

import (
	"github.com/jinzhu/copier"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
)

func (u Usecase) CreateDraftProposal( req structure.CreateProposaltReq) (*entity.ProposalDetail, error) {

	pe := &entity.ProposalDetail{}
	err := copier.Copy(pe, req)
	if err != nil {
		u.Logger.Error("copier.Copy", err.Error(), err)
		return nil, err
	}

	pe.IsDraft = true
	err = u.Repo.CreateProposalDetail(pe)
	if err != nil {
		u.Logger.Error("u.Repo.CreateProject", err.Error(), err)
		return nil, err
	}

	u.Logger.Info("pe", pe)
	return pe, nil
}

func (u Usecase) MapOffToOnChainProposal( ID string, proposalID string) (*entity.ProposalDetail, error) {

pD, err := u.Repo.FindProposalDetailByUUID(ID)
	if err != nil {
		u.Logger.Error("MapOffToOnChainProposal.FindProposalByID", err.Error(), err)
		return nil, err
	}
pD.ProposalID = proposalID
	pD.IsDraft = false
	updated, err := u.Repo.UpdateProposalDetail(ID, pD)
	if err != nil {
		u.Logger.Error("MapOffToOnChainProposal.UpdateProposalDetail", err.Error(), err)
		return nil, err
	}

	u.Logger.Info("updated", updated)
	return pD, nil
}

func (u Usecase) GetProposals( req structure.FilterProposal) (*entity.Pagination, error) {

pe := &entity.FilterProposals{}
	err := copier.Copy(pe, req)
	if err != nil {
		u.Logger.Error("copier.Copy", err.Error(), err)
		return nil, err
	}
proposals, err := u.Repo.FilterProposal(*pe)
	if err != nil {
		u.Logger.Error("u.Repo.FilterProposal", err.Error(), err)
		return nil, err
	}

	u.Logger.Info("proposals", proposals.Total)
	return proposals, nil
}

func (u Usecase) GetProposal( proposalID string) (*entity.Proposal, error) {

	proposal, err := u.Repo.FindProposal(proposalID)
	if err != nil {
		u.Logger.Error("u.Repo.FilterProposal", err.Error(), err)
		return nil, err
	}

	pDetail, err := u.Repo.FindProposalDetail(proposalID)
	if err == nil {
		proposal.ProposalDetail = *pDetail
	}else{
		u.Logger.Error("u.Repo.FilterProposal", err.Error(), err)
	}

	u.Logger.Info("proposal", proposal)
	return proposal, nil
}

func (u Usecase) GetProposalVotes( filter structure.FilterProposalVote) (*entity.Pagination, error) {

f := &entity.FilterProposalVotes{}
	err := copier.Copy(f, filter)
	if err != nil {
		u.Logger.Error("filterProposalVotes.copier.Copy", err.Error(), err)
		return nil, err
	}
proposalVotes, err := u.Repo.FilterProposalVotes(*f)
	if err != nil {
		u.Logger.Error("u.Repo.FilterProposalVotes", err.Error(), err)
		return nil, err
	}


	u.Logger.Info("proposalVotes", proposalVotes)
	return proposalVotes, nil
}