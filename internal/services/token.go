package services

import (
	"context"
	"errors"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"rederinghub.io/api"
	"rederinghub.io/pkg/contracts/generative_nft"
	"rederinghub.io/pkg/logger"
)

func (s *service) GetToken(ctx context.Context, req *api.GetTokenMessageReq) (*api.GetTokenMessageResp, error) {
	logger.AtLog.Infof("Handle [GetToken] %s %s", req.ContractAddr, req.TokenId)
		chainURL, ok := GetRPCURLFromChainID("80001")
		if !ok {
			return nil, errors.New("missing config chain_config from server")
		}

		// call to contract to get emotion
		client, err := ethclient.Dial(chainURL)
		if err != nil {
			return nil, err
		}
		addr := common.HexToAddress(req.ContractAddr)

		instance, err := generative_nft.NewGenerativeNft(addr, client)
		if err != nil {
			return nil, err
		}
		
		spew.Dump(instance)

	resp := &api.GetTokenMessageResp{}
	return resp, nil
}
