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

	"github.com/davecgh/go-spew/spew"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"rederinghub.io/api"
	"rederinghub.io/internal/api/middleware"
	"rederinghub.io/pkg/config"
	log "rederinghub.io/pkg/logger"
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
	rootPath, _ := os.Getwd()
	spew.Dump(rootPath)
	spew.Dump(http.Dir(rootPath + "/swaggerUI"))
	
	fs := http.FileServer(http.Dir(rootPath + "/swaggerUI"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	httpServer := &http.Server{
		Addr:        fmt.Sprintf("%s:%s", server, httpPort),
		Handler:     mux,
		BaseContext: func(_ net.Listener) context.Context { return ctx },
	}
	// Run server
	go func() {
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			// it is fine to use Fatal here because it is not main goroutine
			log.AtLog.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	log.AtLog.Infof("http server is listening to port %v", httpPort)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan,
		syscall.SIGHUP,  // kill -SIGHUP XXXX
		syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
	)

	<-signalChan
	log.AtLog.Infof("os.Interrupt - shutting down...\n")

	go func() {
		<-signalChan
		log.AtLog.Infof("os.Kill - terminating...\n")
	}()

	graceFullCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := httpServer.Shutdown(graceFullCtx); err != nil {
		log.AtLog.Infof("shutdown error: %v\n", err)
		defer os.Exit(1)
		return err
	}
	log.AtLog.Infof("gracefully stopped\n")

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
