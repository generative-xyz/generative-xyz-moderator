package http

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/utils/logger"
)

func (h *httpDelivery) getNftsByAddress(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	address := vars["address"]

	if address == "" {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("address cannot be empty"))
		return
	}
	result, err := h.Usecase.GetNftsByAddress(address)
	if err != nil {
		logger.AtLog.Logger.Error("httpDelivery.walletInfo.Usecase.GetNftsByAddress", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")
}
