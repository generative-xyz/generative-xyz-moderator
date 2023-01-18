package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"rederinghub.io/internal/delivery/http/middleware"
	_httpResponse "rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils/config"
	"rederinghub.io/utils/global"
	_logger "rederinghub.io/utils/logger"
	_redis "rederinghub.io/utils/redis"
	_tracer "rederinghub.io/utils/tracer"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type IServer interface {
	StartServer()
}

type deliveryConfig struct {
	Handler  *mux.Router
	Config   *config.Config
	Response _httpResponse.IHttpResponse
	Logger   _logger.Ilogger
	Tracer   _tracer.ITracer
	Cache    _redis.IRedisCache
}

func (dc *deliveryConfig) LoadConfig(g *global.Global) {
	dc.Handler = g.MuxRouter
	dc.Config = g.Conf
	dc.Response = _httpResponse.NewHttpResponse()
	dc.Logger = g.Logger
	dc.Tracer = g.Tracer
	dc.Cache = g.Cache
}

type httpDelivery struct {
	deliveryConfig
	Usecase    usecase.Usecase
	MiddleWare middleware.IMiddleware
}

func NewHandler(global *global.Global, uc usecase.Usecase) (*httpDelivery, error) {
	h := new(httpDelivery)
	h.LoadConfig(global)
	m := middleware.NewMiddleware(uc, global)
	h.Usecase = uc
	h.MiddleWare = m
	return h, nil
}

func (h *httpDelivery) StartServer() {
	var wait time.Duration
	h.Logger.Info("httpDelivery.StartServer - Starting http-server")

	h.registerRoutes()
	h.Handler.NotFoundHandler = h.Handler.NewRoute().HandlerFunc(http.NotFound).GetHandler()
	credentials := handlers.AllowCredentials()
	methods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	headers := handlers.AllowedHeaders([]string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "X-Requested-With", "param"})
	hCORS := handlers.CORS(credentials, methods, origins, headers)(h.Handler)

	timeOut := h.Config.Context.TimeOut * 10
	
	srv := &http.Server{
		Handler: hCORS,
		Addr:    fmt.Sprintf(":%s",h.Config.ServicePort),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: time.Duration(timeOut) * time.Second,
		ReadTimeout: time.Duration(timeOut) * time.Second,
	}

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		h.Logger.Info(fmt.Sprintf("Server is listening at port %s ...",h.Config.ServicePort))
		if err := srv.ListenAndServe(); err != nil {
			h.Logger.Error("httpDelivery.StartServer - Can not start http server", err)
		}
	}()

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	err := srv.Shutdown(ctx)
	if err != nil {
		h.Logger.Error("httpDelivery.StartServer - Server can not shutdown", err)
		return
	}
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	h.Logger.Warning("httpDelivery.StartServer - server is shutting down")
	os.Exit(0)
}
