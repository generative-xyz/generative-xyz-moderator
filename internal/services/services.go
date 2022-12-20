package services

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"

	"rederinghub.io/api"
	"rederinghub.io/internal/adapter"
	"rederinghub.io/internal/repository"
	"rederinghub.io/pkg/config"
	"rederinghub.io/pkg/oauth2service"
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
	userRepository repository.IUserRepository
	tokenUriRepository repository.ITokenUriRepository

	redisClient *redis.Client
	auth2Service *oauth2service.Auth2
}

func Init(moralisAdapter adapter.MoralisAdapter,
	renderMachineAdapter adapter.RenderMachineAdapter,
	templateRepository repository.TemplateRepository,
	renderedNftRepository repository.RenderedNftRepository,
	userRepository repository.IUserRepository,
	tokenUriRepository repository.ITokenUriRepository,
) Service {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.AppConfig().RedisAddr,
		Password: config.AppConfig().RedisPassword,
	})

	auth := oauth2service.NewAuth2()

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("can not connect to redis")
	}
	return &service{
		moralisAdapter:       moralisAdapter,
		renderMachineAdapter: renderMachineAdapter,

		templateRepository:    templateRepository,
		renderedNftRepository: renderedNftRepository,
		redisClient:           redisClient,
		userRepository:           userRepository,
		tokenUriRepository:           tokenUriRepository,
		auth2Service:  auth,
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
