package services

import (
	"context"

	"rederinghub.io/api"
	"rederinghub.io/internal/adapter"
	"rederinghub.io/internal/repository"
)

type Service interface {
	api.ApiServiceServer
}

type service struct {
	api.UnimplementedApiServiceServer
	moralisAdapter adapter.MoralisAdapter

	templateRepository repository.TemplateRepository
}

func Init(moralisAdapter adapter.MoralisAdapter, templateRepository repository.TemplateRepository) Service {
	return &service{
		moralisAdapter:     moralisAdapter,
		templateRepository: templateRepository,
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
