package http

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/utils/logger"
)

func (h *httpDelivery) requestGM(w http.ResponseWriter, r *http.Request) {

	var reqBody request.FaucetReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		logger.AtLog.Logger.Error("httpDelivery.requestGM.Decode", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	result, err := h.Usecase.ApiCreateNewGM(reqBody.Address, reqBody.Type)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.requestGM", zap.String("err", err.Error()))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")
}
