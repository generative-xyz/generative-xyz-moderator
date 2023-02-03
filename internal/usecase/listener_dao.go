package usecase

import (
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/opentracing/opentracing-go"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/contracts/generative_dao"
)


func (u Usecase) DAOProposalCreated(rootSpan opentracing.Span, chainLog types.Log) error {
	span, log := u.StartSpan("DAOProposalCreated", rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	log.SetData("chainLog", chainLog.Data)

	daoContract, err := generative_dao.NewGenerativeDao(chainLog.Address, u.Blockchain.GetClient())
	if  err != nil {
		log.Error("cannot init DAO contract", err.Error(), err)
		return err
	}

	parsedProposal, err := daoContract.ParseProposalCreated(chainLog)
	if err != nil {
		log.Error("cannot parse createdProposal", err.Error(), err)
		return err
	}
	log.SetData("parsed.Data", parsedProposal)

	

	createdProposal := u.ParseProposal(parsedProposal)	

	state, err := daoContract.State(nil, parsedProposal.ProposalId)
	if err != nil {
		log.Error("daoContract.State", err.Error(), err)
	}else{
		createdProposal.State = state
	}

	err = u.Repo.CreateProposal(createdProposal)
	if  err != nil {
		log.Error("cannot create CreateProposal", err.Error(), err)
		return err
	}

	log.SetData("createdProposal", createdProposal)
	return nil
}

func (u Usecase) ParseProposal(input  *generative_dao.GenerativeDaoProposalCreated) *entity.Proposal {
	
	targets := []string{}
	for _, target := range  input.Targets {
		targets = append(targets, strings.ToLower(target.String()))
	}
	
	values := []int64{}
	for _, value := range  input.Values {
		values = append(values, value.Int64())
	}

	createdProposal := &entity.Proposal{
		ProposalID: input.ProposalId.String(),
		Proposer: strings.ToLower(input.Proposer.String()),
		StartBlock: input.StartBlock.Int64(),
		EndBlock: input.EndBlock.Int64(),
		Description: input.Description,
		Targets: targets,
		Values: values,
		Signatures: input.Signatures,
		Calldatas: input.Calldatas,
		Raw:  u.ParseRaw(input.Raw),
	}
	return createdProposal
}

func (u Usecase) ParseRaw(input  types.Log) entity.ProposalRaw {
		r :=  entity.ProposalRaw {}
		r.Address = input.Address.String()
		r.Data = input.Data
		r.BlockNumber = input.BlockNumber
		r.TransactionHash =  strings.ToLower(input.TxHash.String())
		r.TransactionIndex = input.TxIndex
		r.BlockHash = strings.ToLower(input.BlockHash.String())
		r.LogIndex = input.Index
		r.Removed = input.Removed
		return r
}