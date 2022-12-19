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

	tokenID := new(big.Int)
	tokenID, ok = tokenID.SetString(req.GetTokenId(), 10)
	if !ok {
		err := errors.New("Cannot convert tokenID")
		return nil, err
	}
	// call to contract to get emotion
	client, err := ethclient.Dial(chainURL)
	if err != nil {
		return nil, err
	}
	addr := common.HexToAddress(req.ContractAddr)

	gProject, err := generative_project_contract.NewGenerativeProjectContract(addr, client)
	if err != nil {
		return nil, err
	}

	projectID := new(big.Int).Div(tokenID, big.NewInt(1000000))
	proDetail, err := gProject.ProjectDetails(nil,  projectID)
	if err != nil {
		return nil, err
	}
	spew.Dump(proDetail.GenNFTAddr)

	gNft, err := generative_nft_contract.NewGenerativeNftContract(addr, client)
	if err != nil {
		return nil, err
	}

	value, err := gNft.TokenGenerativeURI(nil, tokenID)
	if err !=nil {
		return nil, err
	}
	spew.Dump(value)

	resp := &api.GetTokenMessageResp{}
	return resp, nil
}
