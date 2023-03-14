package http

import (
	"context"
	"encoding/json"
	"net/http"

	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/utils"
)

// @Summary List DAO Project
// @Description List DAO Project
// @Tags DAO Project
// @Accept json
// @Produce json
// @Param keyword query string false "Keyword"
// @Param status query int false "Status"
// @Param cursor query string false "Last Id"
// @Param limit query int false "Limit"
// @Success 200 {object} entity.Pagination{}
// @Router /dao-project [GET]
// @Security ApiKeyAuth
func (h *httpDelivery) listDaoProject(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, muxVars map[string]string) (interface{}, error) {
			req := &request.ListDaoProjectRequest{}
			if err := utils.QueryParser(r, req); err != nil {
				return nil, err
			}
			userWallet := muxVars[utils.SIGNED_WALLET_ADDRESS]
			return h.Usecase.ListDAOProject(ctx, userWallet, req)
		},
	).ServeHTTP(w, r)
}

// @Summary Create DAO Project
// @Description Create DAO Project
// @Tags DAO Project
// @Accept json
// @Produce json
// @Param request body request.CreateDaoProjectRequest true "Create Dao Project Request"
// @Success 200
// @Router /dao-project [POST]
// @Security ApiKeyAuth
func (h *httpDelivery) createDaoProject(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, muxVars map[string]string) (interface{}, error) {
			var reqBody request.CreateDaoProjectRequest
			err := json.NewDecoder(r.Body).Decode(&reqBody)
			if err != nil {
				return nil, err
			}
			reqBody.CreatedBy = muxVars[utils.SIGNED_WALLET_ADDRESS]
			return h.Usecase.CreateDAOProject(ctx, &reqBody)
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
