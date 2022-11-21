package services

import (
	"context"
	"fmt"

	"rederinghub.io/api"
	"rederinghub.io/internal/adapter"
	"rederinghub.io/internal/repository"
	"rederinghub.io/pkg/config"

	"github.com/jellydator/ttlcache/v3"
)

type Service interface {
	api.ApiServiceServer
}

type service struct {
	api.UnimplementedApiServiceServer
	moralisAdapter       adapter.MoralisAdapter
	renderMachineAdapter adapter.RenderMachineAdapter

	templateRepository    repository.TemplateRepository
	renderedNftRepository repository.RenderedNftRepository

	cache *ttlcache.Cache[string, string]
}

func Init(moralisAdapter adapter.MoralisAdapter,
	renderMachineAdapter adapter.RenderMachineAdapter,
	templateRepository repository.TemplateRepository,
	renderedNftRepository repository.RenderedNftRepository,
) Service {
	cache := ttlcache.New[string, string]()
	return &service{
		moralisAdapter:       moralisAdapter,
		renderMachineAdapter: renderMachineAdapter,

		templateRepository:    templateRepository,
		renderedNftRepository: renderedNftRepository,
		cache: cache,
	}
}

func (s *service) Live(context.Context, *api.LiveRequest) (*api.LiveResponse, error) {
	return &api.LiveResponse{
		Message: "OK",
	}, nil
}

func (s *service) Ping(context.Context, *api.PingRequest) (*api.PingResponse, error) {
	return &api.PingResponse{Message: "OK"}, nil
}

func GetRPCURLFromChainID(chainID string) (string, bool) {
	appConfig := config.AppConfig()
	if v, ok := appConfig.ChainConfigs[chainID]; ok {
		return fmt.Sprintf("%v%v", "https://", v), true
	}

	return "", false
}
