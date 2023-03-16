package http

import (
	"context"
	"encoding/json"
	"net/http"

	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
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
			req := &request.ListDaoArtistRequest{}
			if err := utils.QueryParser(r, req); err != nil {
				return nil, err
			}
			req.Pagination = entity.GetPagination(r)
			userWallet := muxVars[utils.SIGNED_WALLET_ADDRESS]
			return h.Usecase.ListDAOArtist(ctx, userWallet, req)
		},
	).ServeHTTP(w, r)
}

// @Summary Create DAO Artist
// @Description Create DAO Artist
// @Tags DAO Artist
// @Accept json
// @Produce json
// @Param request body request.CreateDaoArtistRequest true "Create Dao Artist Request"
// @Success 200
// @Router /dao-artist [POST]
// @Security ApiKeyAuth
func (h *httpDelivery) createDaoArtist(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, muxVars map[string]string) (interface{}, error) {
			var reqBody request.CreateDaoArtistRequest
			err := json.NewDecoder(r.Body).Decode(&reqBody)
			if err != nil {
				return nil, err
			}
			return h.Usecase.CreateDAOArtist(ctx, muxVars[utils.SIGNED_WALLET_ADDRESS], &reqBody)
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
			return h.Usecase.GetDAOArtist(ctx, muxVars["id"], muxVars[utils.SIGNED_WALLET_ADDRESS])
		},
	).ServeHTTP(w, r)
}

// @Summary Vote DAO Artist
// @Description Vote DAO Artist
// @Tags DAO Artist
// @Accept json
// @Produce json
// @Param id path string true "DAO Artist Id"
// @Param request body request.VoteDaoArtistRequest true "Vote Dao Artist Request"
// @Success 200
// @Router /dao-artist/{id} [PUT]
// @Security ApiKeyAuth
func (h *httpDelivery) voteDaoArtist(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, muxVars map[string]string) (interface{}, error) {
			var reqBody request.VoteDaoArtistRequest
			err := json.NewDecoder(r.Body).Decode(&reqBody)
			if err != nil {
				return nil, err
			}
			return nil, h.Usecase.VoteDAOArtist(ctx, muxVars["id"], muxVars[utils.SIGNED_WALLET_ADDRESS], &reqBody)
		},
	).ServeHTTP(w, r)
}
