package http

import (
	"context"
	"net/http"

	"rederinghub.io/internal/delivery/http/response"
)

// @Summary List DAO Artist
// @Description List DAO Artist
// @Tags DAO Artist
// @Accept json
// @Produce json
// @Param keyword query string false "Keyword"
// @Param status query int false "Status"
// @Param cursor query string false "Last Id"
// @Param limit query int false "Limit"
// @Success 200 {object} entity.Pagination{}
// @Router /dao-artist [GET]
// @Security ApiKeyAuth
func (h *httpDelivery) listDaoArtist(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, muxVars map[string]string) (interface{}, error) {
			return nil, nil
		},
	).ServeHTTP(w, r)
}

// @Summary Create DAO Artist
// @Description Create DAO Artist
// @Tags DAO Artist
// @Accept json
// @Produce json
// @Success 200
// @Router /dao-artist [POST]
// @Security ApiKeyAuth
func (h *httpDelivery) createDaoArtist(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, muxVars map[string]string) (interface{}, error) {
			return nil, nil
		},
	).ServeHTTP(w, r)
}

// @Summary Get DAO Artist
// @Description Get DAO Artist
// @Tags DAO Artist
// @Accept json
// @Produce json
// @Param id path string true "DAO Artist Id"
// @Success 200
// @Router /dao-artist/{id} [GET]
// @Security ApiKeyAuth
func (h *httpDelivery) getDaoArtist(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, muxVars map[string]string) (interface{}, error) {
			return nil, nil
		},
	).ServeHTTP(w, r)
}

// @Summary Vote DAO Artist
// @Description Vote DAO Artist
// @Tags DAO Artist
// @Accept json
// @Produce json
// @Param id path string true "DAO Artist Id"
// @Success 200
// @Router /dao-artist/{id} [PUT]
// @Security ApiKeyAuth
func (h *httpDelivery) voteDaoArtist(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, muxVars map[string]string) (interface{}, error) {
			return nil, nil
		},
	).ServeHTTP(w, r)
}
