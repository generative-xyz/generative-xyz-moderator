package usecase

import (
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/logger"
)

func (u Usecase) CreateDraftProposal( req structure.CreateProposaltReq) (*entity.ProposalDetail, error) {

	pe := &entity.ProposalDetail{}
	err := copier.Copy(pe, req)
	if err != nil {
		logger.AtLog.Logger.Error("copier.Copy", zap.Error(err))
		return nil, err
	}

	pe.IsDraft = true
	err = u.Repo.CreateProposalDetail(pe)
	if err != nil {
		logger.AtLog.Logger.Error("u.Repo.CreateProject", zap.Error(err))
		return nil, err
	}

	logger.AtLog.Logger.Info("pe", zap.Any("pe", pe))
	return pe, nil
}

func (u Usecase) MapOffToOnChainProposal( ID string, proposalID string) (*entity.ProposalDetail, error) {

pD, err := u.Repo.FindProposalDetailByUUID(ID)
	if err != nil {
		logger.AtLog.Logger.Error("MapOffToOnChainProposal.FindProposalByID", zap.Error(err))
		return nil, err
	}
pD.ProposalID = proposalID
	pD.IsDraft = false
	updated, err := u.Repo.UpdateProposalDetail(ID, pD)
	if err != nil {
		logger.AtLog.Logger.Error("MapOffToOnChainProposal.UpdateProposalDetail", zap.Error(err))
		return nil, err
	}

	logger.AtLog.Logger.Info("updated", zap.Any("updated", updated))
	return pD, nil
}

func (u Usecase) GetProposals( req structure.FilterProposal) (*entity.Pagination, error) {

pe := &entity.FilterProposals{}
	err := copier.Copy(pe, req)
	if err != nil {
		logger.AtLog.Logger.Error("copier.Copy", zap.Error(err))
		return nil, err
	}
proposals, err := u.Repo.FilterProposal(*pe)
	if err != nil {
		logger.AtLog.Logger.Error("u.Repo.FilterProposal", zap.Error(err))
		return nil, err
	}

	logger.AtLog.Logger.Info("proposals", zap.Any("proposals.Total", proposals.Total))
	return proposals, nil
}

func (u Usecase) GetProposal( proposalID string) (*entity.Proposal, error) {

	proposal, err := u.Repo.FindProposal(proposalID)
	if err != nil {
		logger.AtLog.Logger.Error("u.Repo.FilterProposal", zap.Error(err))
		return nil, err
	}

	pDetail, err := u.Repo.FindProposalDetail(proposalID)
	if err == nil {
		proposal.ProposalDetail = *pDetail
	}else{
		logger.AtLog.Logger.Error("u.Repo.FilterProposal", zap.Error(err))
	}

	logger.AtLog.Logger.Info("proposal", zap.Any("proposal", proposal))
	return proposal, nil
}

func (u Usecase) GetProposalVotes( filter structure.FilterProposalVote) (*entity.Pagination, error) {

f := &entity.FilterProposalVotes{}
	err := copier.Copy(f, filter)
	if err != nil {
		logger.AtLog.Logger.Error("filterProposalVotes.copier.Copy", zap.Error(err))
		return nil, err
	}
proposalVotes, err := u.Repo.FilterProposalVotes(*f)
	if err != nil {
		logger.AtLog.Logger.Error("u.Repo.FilterProposalVotes", zap.Error(err))
		return nil, err
	}


	logger.AtLog.Logger.Info("proposalVotes", zap.Any("proposalVotes", proposalVotes))
	return proposalVotes, nil
}