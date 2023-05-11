package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/logger"
)

// UserCredits godoc
// @Summary TokenURI's chart
// @Description TokenURI's chart
// @Tags Charts
// @Accept  json
// @Produce  json
// @Param tokenID path string  true "tokenID"
// @Param dateRange query string false "dateRange"
// @Success 200 {object} response.JsonResponse{}
// @Router /charts/tokens/{tokenID}/charts [GET]
func (h *httpDelivery) getChartDataFoTokenURI(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// contractAddress := vars["contractAddress"]
	tokenID := vars["tokenID"]
	dateRange := r.URL.Query().Get("dateRange")
	f := utils.ParseAggregation(dateRange)
	filter := &structure.AggerateChartForToken{
		TokenID:  &tokenID,
		FromDate: &f.FromDate,
		ToDate:   &f.ToDate,
	}

	logger.AtLog.Logger.Info("getChartDataFoTokenURI.Filter", zap.Any("filter", zap.Any("filter)", filter)))
	result, err := h.Usecase.GetChartDataOFTokens(*filter)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.getCollectionListing", zap.String("err", err.Error()))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")
}

// UserCredits godoc
// @Summary GM deposit analytics
// @Description GM deposit analytics
// @Tags CollectionListing
// @Accept  json
// @Produce  json
// @Success 200 {object} response.JsonResponse{}
// @Router /charts/collections/{projectID}/charts [GET]
func (h *httpDelivery) getChartDataForCollection(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// contractAddress := vars["contractAddress"]
	projectID := vars["projectID"]
	dateRange := r.URL.Query().Get("dateRange")
	f := utils.ParseAggregation(dateRange)
	filter := &structure.AggerateChartForProject{
		ProjectID: &projectID,
		FromDate:  &f.FromDate,
		ToDate:    &f.ToDate,
	}

	logger.AtLog.Logger.Info("getChartDataForCollection.Filter", zap.Any("filter", zap.Any("filter)", filter)))
	result, err := h.Usecase.GetChartDataOFProject(*filter)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.getCollectionListing", zap.String("err", err.Error()))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")
}

// UserCredits godoc
// @Summary CollectionListing
// @Description get list CollectionListing
// @Tags Charts
// @Accept  json
// @Produce  json
// @Param projectID path string  false "projectID"
// @Param dateRange query string false "dateRange"
// @Success 200 {object} response.JsonResponse{}
// @Router /charts/gm-collections/deposit [GET]
func (h *httpDelivery) getChartDataForGMCollection(w http.ResponseWriter, r *http.Request) {
	result, err := h.Usecase.GetChartDataForGMCollection()
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")
}
