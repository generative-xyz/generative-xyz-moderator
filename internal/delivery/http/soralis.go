package http

import (
	"context"
	"net/http"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"strconv"
)

// Discord godoc
// @Summary Get all tokens
// @Description Get all tokens
// @Tags Soralis
// @Content-Type: application/json
// @Success 200 {object} response.JsonResponse
// @Security Authorization
// @Router /soralis/tokens [GET]
func (h *httpDelivery) soralisTokens(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, vars map[string]string) (interface{}, error) {
			toPageInt := 1
			perPageInt := 100
			offsetInt := 0

			toPage := r.URL.Query().Get("page")
			if toPage != "" {
				toPageI, err := strconv.Atoi(toPage)
				if err == nil {
					toPageInt = toPageI
				}
			}

			perPage := r.URL.Query().Get("limit")
			if perPage != "" {
				perPageI, err := strconv.Atoi(perPage)
				if err == nil {
					perPageInt = perPageI
				}
			}

			offset := r.URL.Query().Get("offset")
			if offset != "" {
				offsetI, err := strconv.Atoi(offset)
				if err == nil {
					offsetInt = offsetI
				}
			}

			query := request.PaginationReq{
				Page:   &toPageInt,
				Limit:  &perPageInt,
				Offset: &offsetInt,
			}
			tokens, err := h.Usecase.SoralisAllTokens(query)
			if err != nil {
				return nil, err
			}

			return tokens, nil
		},
	).ServeHTTP(w, r)
}

// Discord godoc
// @Summary Get token's min-max price
// @Description Get token's min-max price
// @Tags Soralis
// @Content-Type: application/json
// @Param tokenAddress path string true "token address"
// @Param chartType query string false "default: hour"
// @Success 200 {object} response.JsonResponse
// @Security Authorization
// @Router /soralis/tokens/{tokenAddress}/price/min-max [GET]
func (h *httpDelivery) soralisTokenMinMax(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, vars map[string]string) (interface{}, error) {
			tokenAddress := vars["tokenAddress"]
			chartType := r.URL.Query().Get("chartType")
			if chartType == "" {
				chartType = "hour"
			}

			data, err := h.Usecase.GetTokenSwapChartMinMax(tokenAddress, chartType)
			if err != nil {
				return nil, err
			}
			return data, nil
		},
	).ServeHTTP(w, r)
}

// Discord godoc
// @Summary Get token's current price
// @Description  Get token's current price
// @Tags Soralis
// @Content-Type: application/json
// @Param tokenAddress path string true "token address"
// @Success 200 {object} response.JsonResponse
// @Security Authorization
// @Router /soralis/tokens/{tokenAddress}/price/current [GET]
func (h *httpDelivery) soralisTokenCurrentPrice(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, vars map[string]string) (interface{}, error) {
			tokenAddress := vars["tokenAddress"]

			data, err := h.Usecase.GetTokenReport(tokenAddress)
			if err != nil {
				return nil, err
			}
			return data, nil
		},
	).ServeHTTP(w, r)
}

// Discord godoc
// @Summary Get token's detail
// @Description Get token's detail
// @Tags Soralis
// @Content-Type: application/json
// @Success 200 {object} response.JsonResponse
// @Param tokenAddress path string true "token address"
// @Param walletAddress path string true "wallet address"
// @Security Authorization
// @Router /soralis/tokens/{tokenAddress}/balance/{walletAddress} [GET]
func (h *httpDelivery) soralisUserTokenBalance(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, vars map[string]string) (interface{}, error) {
			tokenAddress := vars["tokenAddress"]
			walletAddress := vars["walletAddress"]

			data, err := h.Usecase.TokenHolderBalance(walletAddress, tokenAddress)
			if err != nil {
				return nil, err
			}
			return data, nil
		},
	).ServeHTTP(w, r)
}

// Discord godoc
// @Summary Get token's detail
// @Description Get token's detail
// @Tags Soralis
// @Content-Type: application/json
// @Success 200 {object} response.JsonResponse
// @Param tokenAddress path string true "token address"
// @Param walletAddress path string true "wallet address"
// @Security Authorization
// @Router /soralis/tokens/{tokenAddress}/balance/{walletAddress} [POST]
func (h *httpDelivery) soralisSnapShotUserTokenBalance(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, vars map[string]string) (interface{}, error) {
			tokenAddress := vars["tokenAddress"]
			walletAddress := vars["walletAddress"]

			data, err := h.Usecase.SoralisSnapShotUserTokenBalance(walletAddress, tokenAddress)
			if err != nil {
				return nil, err
			}
			return data, nil
		},
	).ServeHTTP(w, r)
}

// @Description Get token's detail
// @Tags Soralis
// @Content-Type: application/json
// @Success 200 {object} response.JsonResponse
// @Param tokenAddress path string true "token address"
// @Param walletAddress path string true "wallet address"
// @Security Authorization
// @Router /soralis/tokens/{tokenAddress}/balance/{walletAddress}/time-travel [GET]
func (h *httpDelivery) soralisGetSnapShotUserTokenBalance(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, vars map[string]string) (interface{}, error) {
			tokenAddress := vars["tokenAddress"]
			walletAddress := vars["walletAddress"]

			data, err := h.Usecase.SoralisGetSnapShotUserTokenBalance(walletAddress, tokenAddress)
			if err != nil {
				return nil, err
			}
			return data, nil
		},
	).ServeHTTP(w, r)
}

// Discord godoc
// @Summary Get token's time-travel
// @Description Get token's time-travel
// @Tags Soralis
// @Content-Type: application/json
// @Success 200 {object} response.JsonResponse
// @Security Authorization
// @Router /soralis/tokens/{tokenAddress}/time-travel [GET]
func (h *httpDelivery) soralisTimeTravel(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, vars map[string]string) (interface{}, error) {
			tokenAddress := vars["tokenAddress"]
			chartType := r.URL.Query().Get("chartType")
			if chartType == "" {
				chartType = "hour"
			}

			data, err := h.Usecase.GetTokenSwapChart(tokenAddress, chartType)
			if err != nil {
				return nil, err
			}
			return data, nil
		},
	).ServeHTTP(w, r)
}
