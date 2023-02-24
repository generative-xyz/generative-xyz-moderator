package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
)

// func (h *httpDelivery) inscriptionByOutput(w http.ResponseWriter, r *http.Request) {

// 	var reqBody request.InscriptionByOutput
// 	decoder := json.NewDecoder(r.Body)
// 	err := decoder.Decode(&reqBody)
// 	if err != nil {
// 		h.Logger.Error("httpDelivery.inscriptionByOutput.Decode", err.Error(), err)
// 		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
// 		return
// 	}

// 	result, _, _, err := h.Usecase.InscriptionsByOutputs(reqBody.Outputs)
// 	if err != nil {
// 		h.Logger.Error("httpDelivery.inscriptionByOutput.Usecase.InscriptionsByOutputs", err.Error(), err)
// 		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
// 		return
// 	}

// 	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")
// }

func (h *httpDelivery) walletInfo(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")

	if address == "" {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("address cannot be empty"))
		return
	}
	result, err := h.Usecase.GetBTCWalletInfo(address)
	if err != nil {
		h.Logger.Error("httpDelivery.walletInfo.Usecase.GetBTCWalletInfo", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")
}

func (h *httpDelivery) mintStatus(w http.ResponseWriter, r *http.Request) {

	address := r.URL.Query().Get("address")
	userID := ""
	if address == "" {
		var ok bool
		ctx := r.Context()
		iUserID := ctx.Value(utils.SIGNED_USER_ID)
		userID, ok = iUserID.(string)
		if !ok {
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("address or accessToken cannot be empty"))
			return
		}
		userInfo, err := h.Usecase.UserProfile(userID)
		if err != nil {
			h.Logger.Error("httpDelivery.mintStatus.Usecase.UserProfile", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
		address = userInfo.WalletAddressBTCTaproot
	}

	result, err := h.Usecase.GetCurrentMintingByWalletAddress(address)
	if err != nil {
		h.Logger.Error("httpDelivery.mintStatus.Usecase.GetBTCWalletInfo", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")
}

func (h *httpDelivery) trackTx(w http.ResponseWriter, r *http.Request) {
	var reqBody request.TrackTx
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		h.Logger.Error("httpDelivery.trackTx.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	if reqBody.Address == "" || reqBody.Txhash == "" {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("address nor txhash cannot be empty"))
		return
	}

	err = h.Usecase.TrackWalletTx(reqBody.Address, structure.WalletTrackTx{Txhash: reqBody.Txhash, Type: reqBody.Type, Amount: reqBody.Amount, InscriptionID: reqBody.InscriptionID})
	if err != nil {
		h.Logger.Error("httpDelivery.trackTx.TrackWalletTx", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, "ok", "")
}

func (h *httpDelivery) walletTrackedTx(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 20
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}

	address := r.URL.Query().Get("address")
	userID := ""
	if address == "" {
		var ok bool
		ctx := r.Context()
		iUserID := ctx.Value(utils.SIGNED_USER_ID)
		userID, ok = iUserID.(string)
		if !ok {
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("address or accessToken cannot be empty"))
			return
		}
		userInfo, err := h.Usecase.UserProfile(userID)
		if err != nil {
			h.Logger.Error("httpDelivery.walletTrackedTx.Usecase.UserProfile", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
		address = userInfo.WalletAddressBTCTaproot
	}

	txList, err := h.Usecase.GetWalletTrackTxs(address, int64(limit), int64(offset))
	if err != nil {
		h.Logger.Error("httpDelivery.walletTrackedTx.GetWalletTrackTxs", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, txList, "")
}
