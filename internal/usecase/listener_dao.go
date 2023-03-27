package usecase

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/contracts/generative_dao"
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/logger"
)

func (u Usecase) DAOCastVote( chainLog types.Log) error {

	logger.AtLog.Logger.Info("chainLog", zap.Any("chainLog.Data", chainLog.Data))

	daoContract, err := generative_dao.NewGenerativeDao(chainLog.Address, u.Blockchain.GetClient())
	if err != nil {
		logger.AtLog.Logger.Error("cannot init DAO contract", zap.Error(err))
		return err
	}

	parsedCastVote, err := daoContract.ParseVoteCast(chainLog)
	if err != nil {
		logger.AtLog.Logger.Error("cannot parse parsedCastVote", zap.Error(err))
		return err
	}

	obj := &entity.ProposalVotes{
		ProposalID: parsedCastVote.ProposalId.String(),
		Voter:      strings.ToLower(parsedCastVote.Voter.String()),
		Support:    int(parsedCastVote.Support),
		WeightNum:  helpers.ParseBigToFloat(parsedCastVote.Weight),
		Weight:     parsedCastVote.Weight.String(),
		Reason:     parsedCastVote.Reason,
	}

	logger.AtLog.Logger.Info("parsed.parsedCastVote", zap.Any("obj", obj))
	err = u.Repo.CreateProposalVotes(obj)
	if err != nil {
		logger.AtLog.Logger.Error("cannot create CreateProposalVotes", zap.Error(err))
		return err
	}

	u.SendMessageProposalVote(*obj)
	return nil
}

func (u Usecase) DAOProposalCreated( chainLog types.Log) error {

	logger.AtLog.Logger.Info("chainLog", zap.Any("chainLog.Data", chainLog.Data))

	daoContract, err := generative_dao.NewGenerativeDao(chainLog.Address, u.Blockchain.GetClient())
	if err != nil {
		logger.AtLog.Logger.Error("cannot init DAO contract", zap.Error(err))
		return err
	}

	parsedProposal, err := daoContract.ParseProposalCreated(chainLog)
	if err != nil {
		logger.AtLog.Logger.Error("cannot parse createdProposal", zap.Error(err))
		return err
	}
	logger.AtLog.Logger.Info("parsed.Data", zap.Any("parsedProposal", parsedProposal))
	createdProposal := u.ParseProposal(parsedProposal)

	state, err := daoContract.State(nil, parsedProposal.ProposalId)
	if err != nil {
		logger.AtLog.Logger.Error("daoContract.State", zap.Error(err))
	} else {
		createdProposal.State = state
	}

	err = u.Repo.CreateProposal(createdProposal)
	if err != nil {
		logger.AtLog.Logger.Error("cannot create CreateProposal", zap.Error(err))
		return err
	}

	u.SendMessageProposal(*createdProposal)
	logger.AtLog.Logger.Info("createdProposal", zap.Any("createdProposal", createdProposal))
	return nil
}

func (u Usecase) ParseProposal(input *generative_dao.GenerativeDaoProposalCreated) *entity.Proposal {

	targets := []string{}
	for _, target := range input.Targets {
		targets = append(targets, strings.ToLower(target.String()))
	}

	values := []int64{}
	for _, value := range input.Values {
		values = append(values, value.Int64())
	}

	createdProposal := &entity.Proposal{
		ProposalID:      input.ProposalId.String(),
		Proposer:        strings.ToLower(input.Proposer.String()),
		StartBlock:      input.StartBlock.Int64(),
		EndBlock:        input.EndBlock.Int64(),
		Title:           input.Description,
		Targets:         targets,
		Values:          values,
		Signatures:      input.Signatures,
		Calldatas:       input.Calldatas,
		Raw:             u.ParseRaw(input.Raw),
		Amount:          "0",
		TokenType:       "NATIVE",
		ReceiverAddress: strings.ToLower(input.Proposer.String()),
	}
	return createdProposal
}

func (u Usecase) ParseRaw(input types.Log) entity.ProposalRaw {
	r := entity.ProposalRaw{}
	r.Address = input.Address.String()
	r.Data = input.Data
	r.BlockNumber = input.BlockNumber
	r.TransactionHash = strings.ToLower(input.TxHash.String())
	r.TransactionIndex = input.TxIndex
	r.BlockHash = strings.ToLower(input.BlockHash.String())
	r.LogIndex = input.Index
	r.Removed = input.Removed
	return r
}

func (u Usecase) SendMessageProposal( createdProposal entity.Proposal) {


	//slack
	preText := fmt.Sprintf("[Proposal %s] has been created by %s", createdProposal.ProposalID, createdProposal.Proposer)
	//content := fmt.Sprintf("Title: %s. Token: %s", helpers.CreateProfileLink(owner,  profile.DisplayName),  helpers.CreateTokenLink( token.ProjectID, token.TokenID,  token.Name))
	content := ""
	title := ""
	//title := fmt.Sprintf("Proposal:  %s is %s", createdProposal.ProposalID, event)

	if _, _, err := u.Slack.SendMessageToSlack(preText, title, content); err != nil {
		logger.AtLog.Logger.Error("s.Slack.SendMessageToSlack err", zap.Error(err))
	}
}

func (u Usecase) SendMessageProposalVote( createdProposalVote entity.ProposalVotes) {


	//slack
	preText := fmt.Sprintf("[Vote][Proposal %s] has been voted", createdProposalVote.ProposalID)
	content := fmt.Sprintf("Support: %d. Weight: %s", createdProposalVote.Support, createdProposalVote.Weight)
	title := fmt.Sprintf("Voter:  %s", createdProposalVote.Voter)
	//title := ""

	if _, _, err := u.Slack.SendMessageToSlack(preText, title, content); err != nil {
		logger.AtLog.Logger.Error("s.Slack.SendMessageToSlack err", zap.Error(err))
	}
}
