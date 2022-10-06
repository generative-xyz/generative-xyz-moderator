package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc/credentials/insecure"
	"rederinghub.io/api"
	"rederinghub.io/internal/api/middleware"
	"rederinghub.io/pkg/config"
	"rederinghub.io/pkg/log"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type ApiGateway interface {
	Start() error
}

type apiGateway struct {
}

func NewApiGateway() ApiGateway {
	return &apiGateway{}
}

func (a *apiGateway) run() error {
	logger := log.NewLogger("http_server")
	server := config.ServerConfig().SERVERUrl
	grpcPort := fmt.Sprintf("%d", config.ServerConfig().GRPCPort)
	httpPort := fmt.Sprintf("%d", config.ServerConfig().HTTPPort)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := http.NewServeMux()
	gwMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := api.RegisterApiServiceHandlerFromEndpoint(ctx, gwMux, fmt.Sprintf("%s:%s", server, grpcPort), opts)
	if err != nil {
		return err
	}

	mux.Handle("/", middleware.AllowCORS(gwMux))

	httpServer := &http.Server{
		Addr:        fmt.Sprintf("%s:%s", server, httpPort),
		Handler:     mux,
		BaseContext: func(_ net.Listener) context.Context { return ctx },
	}
	// Run server
	go func() {
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			// it is fine to use Fatal here because it is not main goroutine
			logger.Fatal().Msgf("HTTP server ListenAndServe: %v", err)
		}
	}()

	logger.Info().Msgf("http server is listening to port %v", httpPort)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan,
		syscall.SIGHUP,  // kill -SIGHUP XXXX
		syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
	)

	<-signalChan
	logger.Info().Msgf("os.Interrupt - shutting down...\n")

	go func() {
		<-signalChan
		logger.Fatal().Msgf("os.Kill - terminating...\n")
	}()

	graceFullCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := httpServer.Shutdown(graceFullCtx); err != nil {
		logger.Info().Msgf("shutdown error: %v\n", err)
		defer os.Exit(1)
		return err
	} else {
		logger.Info().Msgf("gracefully stopped\n")
	}

	// manually cancel context if not using httpServer.RegisterOnShutdown(cancel)
	cancel()

	defer os.Exit(0)
	return nil
}

func (a *apiGateway) Start() error {
	if err := a.run(); err != nil {
		return err
	}
	return nil
}
