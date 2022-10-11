package services

import (
	"context"
	"fmt"

	"github.com/labstack/gommon/log"
	"rederinghub.io/api"
	"rederinghub.io/pkg/config"
)

func (s *service) GetTemplate(ctx context.Context, req *api.GetTemplateRequest) (*api.GetTemplateResponse, error) {
	appConfig := config.AppConfig()
	moralisResp, err := s.moralisAdapter.ListNFTs(appConfig.GenerativeBoilerplateContract, appConfig.ChainID)
	if err != nil {
		log.Errorf("moralis get nft error", err)
		return nil, err
	}
	resp := api.GetTemplateResponse{
		Template: make([]*api.Template, 0, len(moralisResp.Result)),
	}
	for _, nft := range moralisResp.Result {
		resp.Template = append(resp.Template, &api.Template{
			Name:    nft.Name,
			TokenId: nft.TokenID,
			Symbol:  nft.Symbol,
		})

	}
	return &resp, nil
}

func (s *service) GetTemplateDetail(ctx context.Context, req *api.GetTemplateDetailRequest) (*api.GetTemplateDetailResponse, error) {
	fmt.Println(req.Name)
	return &api.GetTemplateDetailResponse{
		Code:         "abc",
		TemplateType: 1,
	}, nil
}
