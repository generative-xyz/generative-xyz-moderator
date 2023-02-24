package usecase

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/contracts/generative_dao"
	"rederinghub.io/utils/helpers"
)

func (u Usecase) DAOCastVote( chainLog types.Log) error {

	u.Logger.Info("chainLog", chainLog.Data)

	daoContract, err := generative_dao.NewGenerativeDao(chainLog.Address, u.Blockchain.GetClient())
	if err != nil {
		u.Logger.Error("cannot init DAO contract", err.Error(), err)
		return err
	}

	parsedCastVote, err := daoContract.ParseVoteCast(chainLog)
	if err != nil {
		u.Logger.Error("cannot parse parsedCastVote", err.Error(), err)
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

	u.Logger.Info("parsed.parsedCastVote", obj)
	err = u.Repo.CreateProposalVotes(obj)
	if err != nil {
		u.Logger.Error("cannot create CreateProposalVotes", err.Error(), err)
		return err
	}

	u.SendMessageProposalVote(*obj)
	return nil
}

func (u Usecase) DAOProposalCreated( chainLog types.Log) error {

	u.Logger.Info("chainLog", chainLog.Data)

	daoContract, err := generative_dao.NewGenerativeDao(chainLog.Address, u.Blockchain.GetClient())
	if err != nil {
		u.Logger.Error("cannot init DAO contract", err.Error(), err)
		return err
	}

	parsedProposal, err := daoContract.ParseProposalCreated(chainLog)
	if err != nil {
		u.Logger.Error("cannot parse createdProposal", err.Error(), err)
		return err
	}
	u.Logger.Info("parsed.Data", parsedProposal)
	createdProposal := u.ParseProposal(parsedProposal)

	state, err := daoContract.State(nil, parsedProposal.ProposalId)
	if err != nil {
		u.Logger.Error("daoContract.State", err.Error(), err)
	} else {
		createdProposal.State = state
	}

	err = u.Repo.CreateProposal(createdProposal)
	if err != nil {
		u.Logger.Error("cannot create CreateProposal", err.Error(), err)
		return err
	}

	u.SendMessageProposal(*createdProposal)
	u.Logger.Info("createdProposal", createdProposal)
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
		u.Logger.Error("s.Slack.SendMessageToSlack err", err.Error(), err)
	}
}

func (u Usecase) SendMessageProposalVote( createdProposalVote entity.ProposalVotes) {


	//slack
	preText := fmt.Sprintf("[Vote][Proposal %s] has been voted", createdProposalVote.ProposalID)
	content := fmt.Sprintf("Support: %d. Weight: %s", createdProposalVote.Support, createdProposalVote.Weight)
	title := fmt.Sprintf("Voter:  %s", createdProposalVote.Voter)
	//title := ""

	if _, _, err := u.Slack.SendMessageToSlack(preText, title, content); err != nil {
		u.Logger.Error("s.Slack.SendMessageToSlack err", err.Error(), err)
	}
}
