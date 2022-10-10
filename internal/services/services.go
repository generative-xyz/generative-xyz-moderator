package services

import (
	"context"

	"rederinghub.io/api"
	"rederinghub.io/internal/adapter"
)

type Service interface {
	api.ApiServiceServer
}

type service struct {
	api.UnimplementedApiServiceServer
	moralisAdapter adapter.MoralisAdapter
}

func Init(moralisAdapter adapter.MoralisAdapter) Service {
	return &service{
		moralisAdapter: moralisAdapter,
	}
}

func (s *service) Live(context.Context, *api.LiveRequest) (*api.LiveResponse, error) {
	return &api.LiveResponse{
		Message: "OK",
	}, nil
}
