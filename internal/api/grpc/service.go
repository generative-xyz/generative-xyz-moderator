package grpc

import (
	"context"
	"rederinghub.io/internal/api/http"
	"rederinghub.io/internal/services"
	"rederinghub.io/pkg/log"

	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

type Server struct {
	logger    log.Logger
	container *dig.Container
}

var ServiceCmd = &cobra.Command{
	Use:   "app",
	Short: "Rendering Hub APIs",
	Long:  `Rendering Hub APIs`,
	Run: func(cmd *cobra.Command, args []string) {
		err := NewServer().Run()
		if err != nil {
			panic(any(err))
		}
	},
	Version: "1.0.0",
}

func NewServer() *Server {
	logger := log.NewLogger("grpc_server")
	container := dig.New()

	return &Server{
		logger:    logger,
		container: container,
	}
}

func (s *Server) Run() error {
	s.addToContainer(
		services.Init,
		http.NewApiGateway,
		Init,
	)

	err := s.container.Invoke(func(server GrpcServer) {
		_ = server.Run(context.Background())
	})
	if err != nil {
		s.logger.Error().Msgf("server process ends: %v", err)
	}
	return err
}

func (s *Server) addToContainer(in ...interface{}) {
	var err error
	for _, i := range in {
		err = s.container.Provide(i)
		if err != nil {
			s.logger.Panic().Msg(err.Error())
		}
	}
}
