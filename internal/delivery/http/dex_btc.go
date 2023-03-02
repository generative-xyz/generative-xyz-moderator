package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/utils"
)

func (h *httpDelivery) dexBTCListing(w http.ResponseWriter, r *http.Request) {
	var reqBody request.CreateDexBTCListing
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		h.Logger.Error("httpDelivery.dexBTCListing.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	var ok bool
	ctx := r.Context()
	iUserID := ctx.Value(utils.SIGNED_USER_ID)
	userID, ok := iUserID.(string)
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
	address := userInfo.WalletAddressBTCTaproot
	err = h.Usecase.DexBTCListing(address, reqBody.RawPSBT, reqBody.InscriptionID)
	if err != nil {
		h.Logger.Error("httpDelivery.dexBTCListing.Usecase.DexBTCListing", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, "ok", "")
}

func (h *httpDelivery) cancelBTCListing(w http.ResponseWriter, r *http.Request) {
	var reqBody request.CancelDexBTCListing
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		h.Logger.Error("httpDelivery.dexBTCListing.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	var ok bool
	ctx := r.Context()
	iUserID := ctx.Value(utils.SIGNED_USER_ID)
	userID, ok := iUserID.(string)
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
	address := userInfo.WalletAddressBTCTaproot
	//find order by inscription_id and user_address
	err = h.Usecase.CancelDexBTCListing(reqBody.Txhash, address, reqBody.InscriptionID, reqBody.OrderID)
	if err != nil {
		h.Logger.Error("httpDelivery.dexBTCListing.Usecase.DexBTCListing", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, "ok", "")
}

func (h *httpDelivery) retrieveBTCListingOrderInfo(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("order_id")
	if orderID == "" {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("orderID cannot be empty"))
		return
	}
	orderInfo, err := h.Usecase.Repo.GetDexBTCListingOrderByID(orderID)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("get order info failed"))
		return
	}
	result := response.DexBTCListingOrderInfo{
		RawPSBT: orderInfo.RawPSBT,
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")
}
