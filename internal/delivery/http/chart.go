package http

import (
	"errors"
	"net/http"
	"os"
	"rederinghub.io/external/etherscan"
	"strconv"
	"strings"

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
	useBackup := os.Getenv("API_DEPOSIT_BACKUP")
	result := &structure.AnalyticsProjectDeposit{}
	var err error

	if useBackup == "" {
		//result, err = h.Usecase.GetChartDataForGMCollection(r.URL.Query().Get("run") != "1")
		result, err = h.Usecase.GetChartDataForGMCollection(true)
		if err != nil {
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
	} else if useBackup == "true" {
		result, err = h.Usecase.GetChartDataForGMCollectionBackup()
		if err != nil {
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
	}

	result.MapItems = make(map[string]*etherscan.AddressTxItemResponse)
	for _, item := range result.Items {
		item.To = ""
	}
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")
}

func (h *httpDelivery) tryReallocate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)
	walletAddress, ok := iWalletAddress.(string)
	if !ok {
		err := errors.New("WalletAddress is incorect")
		logger.AtLog.Logger.Error("withdraw.walletAddress", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	if strings.ToLower(walletAddress) != strings.ToLower("0x07e51AEc82C7163e3237cfbf8C0E6A07413FA18E") {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("error call"))
		return
	}

	result, err := h.Usecase.ReAllocateGM()
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	result.MapItems = make(map[string]*etherscan.AddressTxItemResponse)
	for _, item := range result.Items {
		item.To = ""
	}
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")
}

func (h *httpDelivery) restoreGMDashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)
	walletAddress, ok := iWalletAddress.(string)
	if !ok {
		err := errors.New("WalletAddress is incorect")
		logger.AtLog.Logger.Error("withdraw.walletAddress", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	if strings.ToLower(walletAddress) != strings.ToLower("0x07e51AEc82C7163e3237cfbf8C0E6A07413FA18E") {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("error call"))
		return
	}

	uuid := r.URL.Query().Get("uuid")
	result, err := h.Usecase.RestoreGMDashboardCachedData(uuid)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	result.MapItems = make(map[string]*etherscan.AddressTxItemResponse)
	for _, item := range result.Items {
		item.To = ""
	}
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")
}

func (h *httpDelivery) GetPriceCoinBase(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	idInt, _ := strconv.Atoi(id)
	result, err := h.Usecase.GetPriceCoinBase(idInt)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")
}

func (h *httpDelivery) GetBitcoinBalance(w http.ResponseWriter, r *http.Request) {
	addr := r.URL.Query().Get("addr")
	result, err := h.Usecase.GetBitcoinBalance(addr)
	if err != nil {
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
// @Param address path string  false "address"
// @Success 200 {object} response.JsonResponse{}
// @Router /charts/gm-collections/extra/{address}/deposit [GET]
func (h *httpDelivery) getChartDataExtraForGMCollection(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	result := h.Usecase.GetExtraPercent(address)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")
}

func (h *httpDelivery) getListWallet(w http.ResponseWriter, r *http.Request) {
	result, err := h.Usecase.GetListWallet(r.URL.Query().Get("wallet_type"))
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")
}

// UserCredits godoc
// @Summary Chart for deposit dashboard
// @Description Chart for deposit dashboard
// @Tags Charts
// @Accept  json
// @Produce  json
// @Success 200 {object} response.JsonResponse{}
// @Router /charts/gm-collections/deposit/chart [GET]
/*func (h *httpDelivery) getChartDepositDashboard(w http.ResponseWriter, r *http.Request) {
	result, err := h.Usecase.ChartForGMDashboard()
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")
}*/

func (h *httpDelivery) GetDataOld(w http.ResponseWriter, r *http.Request) {
	result, err := h.Usecase.GetDataOld()
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")
}
