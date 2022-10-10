package services

import (
	"context"

	"github.com/labstack/gommon/log"
	"rederinghub.io/api"
)

func (s *service) GetTemplate(ctx context.Context, req *api.GetTemplateRequest) (*api.GetTemplateResponse, error) {
	//appConfig := config.AppConfig()
	// TODO @khoa read from config
	moralisResp, err := s.moralisAdapter.ListNFTs("0x19cbe1721a63dd4f391fc6f0a75596fe98c2301a", "goerli")
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
