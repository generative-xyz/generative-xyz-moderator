package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/utils"
)

// UserCredits godoc
// @Summary Generate an api key for developer
// @Description Generate an api key for developer
// @Tags developer
// @Accept  json
// @Produce  json
// @Security Authorization
// @Param request body request.GetApiKeyReq true "request.GetApiKeyReq"
// @Success 200 {object} response.JsonResponse{}
// @Router /developer/get-api-key [POST]
func (h *httpDelivery) apiDeveloper_GenApiKey(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)
	userWalletAddr, ok := iWalletAddress.(string)
	if !ok {
		err := errors.New("Wallet address is incorect")
		h.Logger.Error("ctx.Value.Token", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	var req request.GetApiKeyReq

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		h.Logger.Error("httpDelivery.getDeveloperApiKey.NewDecoder", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	resp, err := h.Usecase.ApiDevelop_GenApiKey(userWalletAddr, &req)
	if err != nil {
		h.Logger.Error("h.Usecase.ApiDevelop_GenApiKey", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

// UserCredits godoc
// @Summary Get api key for developer
// @Description Get an api key for developer
// @Tags developer
// @Accept  json
// @Produce  json
// @Security No Authorization
// @Success 200 {object} entity.DeveloperKey
// @Router /developer/get-api-key [GET]
func (h *httpDelivery) apiDeveloper_GetApiKey(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)
	userWalletAddr, ok := iWalletAddress.(string)
	if !ok {
		err := errors.New("Wallet address is incorect")
		h.Logger.Error("ctx.Value.Token", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	resp, err := h.Usecase.ApiDevelop_GetApiKey(userWalletAddr)
	if err != nil {
		h.Logger.Error("h.Usecase.ApiDevelop_GetApiKey", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}
