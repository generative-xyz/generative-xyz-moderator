package grpc

import (
	"context"

	"github.com/spf13/cobra"
	"go.uber.org/dig"
	"go.uber.org/zap"

	"rederinghub.io/internal/adapter"
	"rederinghub.io/internal/api/http"
	"rederinghub.io/internal/model"
	"rederinghub.io/internal/repository"
	"rederinghub.io/internal/services"
	"rederinghub.io/pkg/config"
	log "rederinghub.io/pkg/logger"
)

type Server struct {
	logger    *log.AutoLogger
	container *dig.Container
}

var ServiceCmd = &cobra.Command{
	Use:   "app",
	Short: "Rendering Hub APIs",
	Long:  "Rendering Hub APIs",
	Run: func(cmd *cobra.Command, args []string) {
		err := NewServer().Run()
		if err != nil {
			panic(any(err))
		}
	},
	Version: "1.0.0",
}

func NewServer() *Server {
	logger := log.AtLog
	if config.AppConfig().Environment == config.ProductionEnvironment {
		log.InitLoggerDefault(true)
	}
	container := dig.New()

	return &Server{
		logger:    logger,
		container: container,
	}
}

func (s *Server) Run() error {
	s.addToContainer(
		// Base
		services.Init,
		http.NewApiGateway,
		Init,
		model.NewDatabase,

		// repository
		repository.NewTemplateRepository,
		repository.NewRenderedNftRepository,

		// adapter
		adapter.NewMoralisAdapter,
		adapter.NewRenderMachineAdapter,
	)

	err := s.container.Invoke(func(server GrpcServer) {
		_ = server.Run(context.Background())
	})
	if err != nil {
		s.logger.Logger.Error("server process ends: %v", zap.Error(err))
	}
	return err
}

func (s *Server) addToContainer(in ...interface{}) {
	var err error
	for _, i := range in {
		err = s.container.Provide(i)
		if err != nil {
			s.logger.Logger.Panic(err.Error())
		}
	}
}
