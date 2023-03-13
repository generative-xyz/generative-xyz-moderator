package http

import (
	"context"
	"net/http"

	"rederinghub.io/internal/delivery/http/response"
)

// @Summary List DAO Project
// @Description List DAO Project
// @Tags DAO Project
// @Accept json
// @Produce json
// @Param cursor query string false "Last Id"
// @Param limit query int false "Limit"
// @Success 200 {object} entity.Pagination{}
// @Router /dao-project [GET]
// @Security ApiKeyAuth
func (h *httpDelivery) listDaoProject(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, muxVars map[string]string) (interface{}, error) {
			return nil, nil
		},
	).ServeHTTP(w, r)
}

// @Summary Create DAO Project
// @Description Create DAO Project
// @Tags DAO Project
// @Accept json
// @Produce json
// @Success 200
// @Router /dao-project [POST]
// @Security ApiKeyAuth
func (h *httpDelivery) createDaoProject(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, muxVars map[string]string) (interface{}, error) {
			return nil, nil
		},
	).ServeHTTP(w, r)
}

// @Summary Get DAO Project
// @Description Get DAO Project
// @Tags DAO Project
// @Accept json
// @Produce json
// @Success 200 {object} nfts.MoralisToken{}
// @Router /dao-project/{id} [GET]
// @Security ApiKeyAuth
func (h *httpDelivery) getDaoProject(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, muxVars map[string]string) (interface{}, error) {
			return nil, nil
		},
	).ServeHTTP(w, r)
}

// @Summary Vote DAO Project
// @Description Vote DAO Project
// @Tags DAO Project
// @Accept json
// @Produce json
// @Success 200
// @Router /dao-project/{id} [PUT]
// @Security ApiKeyAuth
func (h *httpDelivery) voteDaoProject(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, muxVars map[string]string) (interface{}, error) {
			return nil, nil
		},
	).ServeHTTP(w, r)
}
