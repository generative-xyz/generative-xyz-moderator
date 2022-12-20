package services

import (
	"context"
	"errors"
	"math/big"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"rederinghub.io/api"
	"rederinghub.io/pkg/contracts/generative_nft_contract"
	"rederinghub.io/pkg/contracts/generative_project_contract"
	"rederinghub.io/pkg/logger"
)

func (s *service) GetToken(ctx context.Context, req *api.GetTokenMessageReq) (*api.GetTokenMessageResp, error) {
	logger.AtLog.Infof("Handle [GetToken] %s %s", req.ContractAddr, req.TokenId)
	chainURL, ok := GetRPCURLFromChainID("80001")
	if !ok {
		return nil, errors.New("missing config chain_config from server")
	}

	addr := common.HexToAddress(req.GetContractAddr())

	// call to contract to get emotion
	client, err := ethclient.Dial(chainURL)
	if err != nil {
		return nil, err
	}

	nftProjectDetail, err := s.getNftContractDetail(client, addr, req.TokenId)
	if err != nil {
		return nil, err
	}

	parentAddr := nftProjectDetail.GenNFTAddr
	tokenUriData, err := s.getNftProjectTokenUri(client, parentAddr,  req.TokenId)
	spew.Dump(tokenUriData)
	
	resp := &api.GetTokenMessageResp{}
	return resp, nil
}

func (s *service) getNftContractDetail(client *ethclient.Client, contractAddr common.Address, tokenIDStr string) (*generative_project_contract.NFTProjectProject, error) {
	tokenID := new(big.Int)
	tokenID, ok := tokenID.SetString(tokenIDStr, 10)
	if !ok {
		err := errors.New("Cannot convert tokenID")
		return nil, err
	}

	gProject, err := generative_project_contract.NewGenerativeProjectContract(contractAddr, client)
	if err != nil {
		return nil, err
	}

	projectID := new(big.Int).Div(tokenID, big.NewInt(1000000))
	proDetail, err := gProject.ProjectDetails(nil,  projectID)
	if err != nil {
		return nil, err
	}
	return &proDetail, nil
}

func (s *service) getNftProjectTokenUri(client *ethclient.Client, contractAddr common.Address, tokenIDStr string) (*string, error) {
	tokenID := new(big.Int)
	tokenID, ok := tokenID.SetString(tokenIDStr, 10)
	if !ok {
		err := errors.New("Cannot convert tokenID")
		return nil, err
	}

	gNft, err := generative_nft_contract.NewGenerativeNftContract(contractAddr, client)
	if err != nil {
		return nil, err
	}

	value, err := gNft.TokenGenerativeURI(nil, tokenID)
	if err !=nil {
		return nil, err
	}

	return &value, nil
}