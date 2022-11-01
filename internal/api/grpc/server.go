package grpc

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"rederinghub.io/api"
	"rederinghub.io/internal/api/http"
	"rederinghub.io/internal/api/middleware"
	"rederinghub.io/internal/services"
	"rederinghub.io/pkg/config"
	log "rederinghub.io/pkg/logger"
)

type GrpcServer interface {
	Run(ctx context.Context) error
}

type grpcServer struct {
	ctx    context.Context
	server *grpc.Server
	apiSvc services.Service
	gw     http.ApiGateway
}

func Init(apiSvc services.Service, gw http.ApiGateway) GrpcServer {
	var g grpcServer
	g.apiSvc = apiSvc
	g.gw = gw
	return &g
}

func (g *grpcServer) Run(ctx context.Context) error {
	g.ctx = ctx
	errChan := make(chan error)
	port := fmt.Sprintf(":%d", config.ServerConfig().GRPCPort)
	grpcListener, err := net.Listen("tcp", port)
	if err != nil {
		log.AtLog.Errorf(fmt.Sprint("failed to start grpc server listener: ", err))
		return err
	}

	interceptors := middleware.NewInterceptor()
	grpc.EnableTracing = true
	baseServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.WithTimeoutInterceptor(),
			interceptors.ValidationInterceptor(),
		))

	api.RegisterApiServiceServer(baseServer, g.apiSvc)
	g.server = baseServer
	go func() {
		log.AtLog.Infof("grpc server is listening to port %v", port)
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
	log.AtLog.Print(fmt.Sprint("Service is stopped: ", err))
	g.server.GracefulStop()

	return err
}
