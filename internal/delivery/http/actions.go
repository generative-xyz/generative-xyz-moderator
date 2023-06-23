package http

import (
	"context"
	"errors"
	"net/http"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/utils"
)

// UserCredits godoc
// @Summary like project
// @Description like project
// @Tags Like & dislike
// @Accept  json
// @Produce  json
// @Security Authorization
// @Param projectID path string true "projectID"
// @Success 200 {object} response.JsonResponse{}
// @Router /action/project/{projectID}/like [POST]
func (h *httpDelivery) LikeProject(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, vars map[string]string) (interface{}, error) {
			iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)
			walletAddress, ok := iWalletAddress.(string)
			if !ok {
				return nil, errors.New("Cannot get wallet address")
			}

			projectID := vars["projectID"]
			return h.Usecase.LikeProject(projectID, walletAddress)
		},
	).ServeHTTP(w, r)
}

// UserCredits godoc
// @Summary dislike project
// @Description dislike project
// @Tags Like & dislike
// @Accept  json
// @Produce  json
// @Security Authorization
// @Param projectID path string true "projectID"
// @Success 200 {object} response.JsonResponse{}
// @Router /action/project/{projectID}/dislike [POST]
func (h *httpDelivery) DisLikeProject(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, vars map[string]string) (interface{}, error) {

			iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)
			walletAddress, ok := iWalletAddress.(string)
			if !ok {
				return nil, errors.New("Cannot get wallet address")
			}

			projectID := vars["projectID"]
			return h.Usecase.DisLikeProject(projectID, walletAddress)
		},
	).ServeHTTP(w, r)
}

// UserCredits godoc
// @Summary like token
// @Description like token
// @Tags Like & dislike
// @Accept  json
// @Produce  json
// @Security Authorization
// @Param tokenID path string true "token ID"
// @Success 200 {object} response.JsonResponse{}
// @Router /action/tokens/{tokenID}/like [POST]
func (h *httpDelivery) LikeTokenURI(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, vars map[string]string) (interface{}, error) {

			iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)
			walletAddress, ok := iWalletAddress.(string)
			if !ok {
				return nil, errors.New("Cannot get wallet address")
			}

			tokenID := vars["tokenID"]
			return h.Usecase.LikeToken(tokenID, walletAddress)
		},
	).ServeHTTP(w, r)
}

// UserCredits godoc
// @Summary dislike token
// @Description dislike token
// @Tags Like & dislike
// @Accept  json
// @Produce  json
// @Security Authorization
// @Param tokenID path string true "token ID"
// @Success 200 {object} response.JsonResponse{}
// @Router  /action/tokens/{tokenID}/dislike [POST]
func (h *httpDelivery) DisLikeTokenURI(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, vars map[string]string) (interface{}, error) {
			iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)
			walletAddress, ok := iWalletAddress.(string)
			if !ok {
				return nil, errors.New("Cannot get wallet address")
			}

			tokenID := vars["tokenID"]
			return h.Usecase.DisLikeToken(tokenID, walletAddress)
		},
	).ServeHTTP(w, r)
}
