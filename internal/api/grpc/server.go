package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"rederinghub.io/api"
	"rederinghub.io/internal/api/http"
	"rederinghub.io/internal/api/middleware"
	"rederinghub.io/internal/services"
	"rederinghub.io/pkg/config"
	"rederinghub.io/pkg/log"
	"syscall"
)

type GrpcServer interface {
	Run(ctx context.Context) error
}

type grpcServer struct {
	logger log.Logger
	ctx    context.Context
	server *grpc.Server
	svc    services.Service
	gw     http.ApiGateway
}

func Init(service services.Service, gw http.ApiGateway) GrpcServer {
	var g grpcServer
	g.logger = log.NewLogger("grpc_server")
	g.svc = service
	g.gw = gw
	return &g
}

func (g *grpcServer) Run(ctx context.Context) error {
	g.ctx = ctx
	errChan := make(chan error)
	port := fmt.Sprintf(":%d", config.ServerConfig().GRPCPort)
	grpcListener, err := net.Listen("tcp", port)
	if err != nil {
		g.logger.Error().Msg(fmt.Sprint("failed to start grpc server listener: ", err))
		return err
	}

	interceptors := middleware.NewInterceptor(g.logger)
	baseServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.WithTimeoutInterceptor(),
			interceptors.ValidationInterceptor(),
		))

	api.RegisterApiServiceServer(baseServer, g.svc)
	g.server = baseServer
	go func() {
		g.logger.Info().Msgf("grpc server is listening to port %v", port)
		errChan <- baseServer.Serve(grpcListener)
	}()

	//Start Gateway
	go func() {
		errChan <- g.gw.Start()
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	err = <-errChan
	g.logger.Error().Msg(fmt.Sprint("Service is stopped: ", err))
	g.server.GracefulStop()

	return err
}
