package http

import (
	"net/http"
	"os"
	"strconv"

	"rederinghub.io/docs"
	_ "rederinghub.io/docs"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/tracer"

	"github.com/opentracing/opentracing-go"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (h *httpDelivery) registerRoutes() {
	h.RegisterDocumentRoutes()
	h.RegisterV1Routes()
}

func (h *httpDelivery) RegisterV1Routes() {
	h.Handler.Use(h.MiddleWare.Tracer)
	h.Handler.Use(h.MiddleWare.LoggingMiddleware)
	h.Handler.HandleFunc("/", h.healthCheck).Methods("GET")

	//api
	api := h.Handler.PathPrefix("/api").Subrouter()
	api.HandleFunc("/token/{contractAddress}/{tokenID}", h.tokenURI).Methods("GET")
	api.HandleFunc("/trait/{contractAddress}/{tokenID}", h.tokenTrait).Methods("GET")
	
	
	//v1 := api.PathPrefix("/v1").Subrouter()
	api.HandleFunc("/", h.healthCheck).Methods("GET")

	//auth
	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/nonce", h.generateMessage).Methods("POST")
	auth.HandleFunc("/nonce/verify", h.verifyMessage).Methods("POST")

	files := api.PathPrefix("/files").Subrouter()
	// files.Use(h.MiddleWare.AccessToken)
	files.HandleFunc("", h.UploadFile).Methods("POST")
	

	//profile
	singedIn := api.PathPrefix("/profile").Subrouter()
	singedIn.Use(h.MiddleWare.AccessToken)
	singedIn.HandleFunc("", h.profile).Methods("GET")
	singedIn.HandleFunc("", h.updateProfile).Methods("PUT")
	singedIn.HandleFunc("/logout", h.logout).Methods("PUT")

	//project
	project := api.PathPrefix("/project").Subrouter()
	project.HandleFunc("", h.getProjects).Methods("GET")
	project.HandleFunc("", h.createProjects).Methods("POST")
	project.HandleFunc("/{contractAddress}/tokens/{projectID}", h.projectDetail).Methods("GET")
	project.HandleFunc("/{contractAddress}/tokens", h.projectTokens).Methods("GET")
	
	
	//project
	config := api.PathPrefix("/configs").Subrouter()
	config.HandleFunc("", h.getConfigs).Methods("GET")
	config.HandleFunc("", h.createConfig).Methods("POST")
	config.HandleFunc("/{key}", h.getConfig).Methods("GET")
	config.HandleFunc("/{key}", h.deleteConfig).Methods("DELETE")
}

func (h *httpDelivery) RegisterDocumentRoutes() {
	documentUrl := `/swagger/`
	domain := os.Getenv("swagger_domain")
	docs.SwaggerInfo.Host = domain
	docs.SwaggerInfo.BasePath = "/api"
	swaggerURL := documentUrl + "swagger/doc.json"
	h.Handler.PathPrefix(documentUrl).Handler(httpSwagger.Handler(
		httpSwagger.URL(swaggerURL), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		//httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("#swagger-ui"),
	))
}

func (h *httpDelivery) StartSpan(name string, r *http.Request) (opentracing.Span, *tracer.TraceLog) {
	span := h.Tracer.StartSpanFromHeaderInjection(r.Header, name)
	log := tracer.NewTraceLog()
	return span, log
}

func (h *httpDelivery) healthCheck(w http.ResponseWriter, r *http.Request) {
	span := h.Tracer.StartSpan("healthCheck")
	h.Response.SetTrace(h.Tracer)
	h.Response.SetSpan(span)
	defer span.Finish()
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, "It work!", "")
}

func (h *httpDelivery) PaginationResp(data *entity.Pagination, items interface{}) response.PaginationResponse {
	resp := response.PaginationResponse{}
	resp.Result = items
	resp.Currsor = data.Currsor
	resp.Total = data.Total
	resp.Page = data.Page
	resp.PageSize = data.PageSize
	return resp
}

func (h *httpDelivery) BaseFilters(r *http.Request) (*structure.BaseFilters, error) {
	f := &structure.BaseFilters{}

	limitInt := 10
	pageInt := 1
	var err error

	limit := r.URL.Query().Get("limit")
	if limit != "" {
		limitInt, err = strconv.Atoi(limit)
		if err != nil {
			return nil, err
		}
	}

	page := r.URL.Query().Get("page")
	if page != "" {
		pageInt, err = strconv.Atoi(page)
		if err != nil {
			return nil, err
		}
	}

	sort := r.URL.Query().Get("sort")
	if sort != "" {
		sortInt, err := strconv.Atoi(sort)
		if err != nil {
			return nil, err
		}
		f.Sort = sortInt
	}

	f.SortBy = r.URL.Query().Get("sort_by")

	f.Page = int64(pageInt)
	f.Limit = int64(limitInt)

	return f, nil
}
