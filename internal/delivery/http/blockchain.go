package http

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/logger"
)

// UserCredits godoc
// @Summary get nft transactions
// @Description get nft transactions
// @Tags Blockchain
// @Accept  json
// @Produce  json
// @Param contractAddress path string true "contract address"
// @Param tokenID path string true "token ID"
// @Success 200 {object} response.JsonResponse{}
// @Router /nfts/{contractAddress}/transactions/{tokenID} [GET]
func (h *httpDelivery) getNftTransactions(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	contractAddress := vars["contractAddress"]
	tokenID := vars["tokenID"]

	covalentResp, err := h.Usecase.GetNftTransactions(structure.GetNftTransactionsReq{
		ContractAddress: contractAddress,
		TokenID: tokenID,
	})

	if err != nil {
		logger.AtLog.Logger.Error("getNftTransactions", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp := covalentResp.Data

	logger.AtLog.Logger.Info("resp", zap.Any("resp", resp));
	
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")	
}

// UserCredits godoc
// @Summary get token holder
// @Description get token holder
// @Tags Blockchain
// @Accept  json
// @Produce  json
// @Param contractAddress path string true "contract address"
// @Param page query string false "page"
// @Param limit query string false "limit"
// @Success 200 {object} response.JsonResponse{}
// @Router /nfts/{contractAddress}/nft_holders [GET]
func (h *httpDelivery) getTokenHolder(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	contractAddress := vars["contractAddress"]

	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		logger.AtLog.Logger.Error("parse page param to int", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		logger.AtLog.Logger.Error("parse limit param to int", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp, err := h.Usecase.GetTokenHolders(structure.GetTokenHolderRequest{
		ContractAddress: contractAddress,
		Page: int32(page),
		Limit: int32(limit),
	})

	if err != nil {
		logger.AtLog.Logger.Error("parse limit param to int", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	pagResp := h.PaginationResp(resp, resp.Result)

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, pagResp, "")
}
