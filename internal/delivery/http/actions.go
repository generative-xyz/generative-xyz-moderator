package http

import (
	"context"
	"net/http"
	"rederinghub.io/internal/delivery/http/response"
)

// UserCredits godoc
// @Summary like project
// @Description like project
// @Tags Project
// @Accept  json
// @Produce  json
// @Param contractAddress path string true "contract address"
// @Param projectID path string true "token ID"
// @Success 200 {object} response.JsonResponse{}
// @Router /project/{contractAddress}/tokens/{projectID}/like [GET]
func (h *httpDelivery) LikeProject(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, vars map[string]string) (interface{}, error) {
			contractAddress := vars["contractAddress"]
			projectID := vars["projectID"]
			return h.Usecase.LikeProject(contractAddress, projectID)
		},
	).ServeHTTP(w, r)
}

// UserCredits godoc
// @Summary dislike project
// @Description dislike project
// @Tags Like & dislike
// @Accept  json
// @Produce  json
// @Param contractAddress path string true "contract address"
// @Param projectID path string true "token ID"
// @Success 200 {object} response.JsonResponse{}
// @Router /project/{contractAddress}/tokens/{projectID}/dislike [GET]
func (h *httpDelivery) DisLikeProject(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, vars map[string]string) (interface{}, error) {
			contractAddress := vars["contractAddress"]
			projectID := vars["projectID"]
			return h.Usecase.DisLikeProject(contractAddress, projectID)
		},
	).ServeHTTP(w, r)
}

// UserCredits godoc
// @Summary like token
// @Description like token
// @Tags Project
// @Accept  json
// @Produce  json
// @Param contractAddress path string true "contract address"
// @Param projectID path string true "token ID"
// @Success 200 {object} response.JsonResponse{}
// @Router /tokens/{contractAddress}/{tokenID}/like [POST]
func (h *httpDelivery) LikeTokenURI(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, vars map[string]string) (interface{}, error) {
			contractAddress := vars["contractAddress"]
			projectID := vars["projectID"]
			return h.Usecase.LikeToken(contractAddress, projectID)
		},
	).ServeHTTP(w, r)
}

// UserCredits godoc
// @Summary dislike token
// @Description dislike token
// @Tags Like & dislike
// @Accept  json
// @Produce  json
// @Param contractAddress path string true "contract address"
// @Param projectID path string true "token ID"
// @Success 200 {object} response.JsonResponse{}
// @Router /tokens/{contractAddress}/{tokenID}/dislike [POST]
func (h *httpDelivery) DisLikeTokenURI(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, vars map[string]string) (interface{}, error) {
			contractAddress := vars["contractAddress"]
			projectID := vars["projectID"]
			return h.Usecase.DisLikeToken(contractAddress, projectID)
		},
	).ServeHTTP(w, r)
}
