package http

import (
	"fmt"
	"net/http"
	"time"

	"rederinghub.io/internal/delivery/http/middleware"
	_httpResponse "rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils/config"
	"rederinghub.io/utils/global"
	_logger "rederinghub.io/utils/logger"
	_redis "rederinghub.io/utils/redis"

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
	Cache    _redis.IRedisCache
}

func (dc *deliveryConfig) LoadConfig(g *global.Global) {
	dc.Handler = g.MuxRouter
	dc.Config = g.Conf
	dc.Response = _httpResponse.NewHttpResponse()
	dc.Logger = g.Logger
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
		Handler: handlers.CompressHandler(hCORS),
		Addr:    fmt.Sprintf(":%s", h.Config.ServicePort),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: time.Duration(timeOut) * time.Second,
		ReadTimeout:  time.Duration(timeOut) * time.Second,
	}

	h.Logger.Info(fmt.Sprintf("Server is listening at port %s ...", h.Config.ServicePort))
	if err := srv.ListenAndServe(); err != nil {
		h.Logger.Error("httpDelivery.StartServer - Can not start http server", err)
	}

}
