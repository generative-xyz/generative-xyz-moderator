package services

import (
	"rederinghub.io/api"
)

type Service interface {
	api.ApiServiceServer
}

type service struct {
	api.UnimplementedApiServiceServer
}

func Init() Service {
	return &service{}
}
