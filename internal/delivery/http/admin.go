package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/utils"
)

// UserCredits godoc
// @Summary Get Redis
// @Description Get Redis
// @Tags Admin
// @Accept  json
// @Produce  json
// @Success 200 {object} response.JsonResponse{data=[]response.RedisResponse}
// @Router /admin/redis [GET]
func (h *httpDelivery) getRedisKeys(w http.ResponseWriter, r *http.Request) {
	res, err := h.Usecase.GetAllRedis()

	if err != nil {
		h.Logger.Error(err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, res, "")
}

// UserCredits godoc
// @Summary Get Redis
// @Description Get Redis
// @Tags Admin
// @Accept  json
// @Produce  json
// @Param key path string true "Redis key"
// @Success 200 {object} response.JsonResponse{data=response.RedisResponse}
// @Router /admin/redis/{key} [GET]
func (h *httpDelivery) getRedis(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	redisKey := vars["key"]
	res, err := h.Usecase.GetRedis(redisKey)

	if err != nil {
		h.Logger.Error(err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, &response.RedisResponse{Value: res}, "")
}

// UserCredits godoc
// @Summary Upsert Redis
// @Description Upsert Redis
// @Tags Admin
// @Accept  json
// @Produce  json
// @Param request body request.UpsertRedisRequest true "Upsert redis key"
// @Success 200 {object} response.JsonResponse{data=response.RedisResponse}
// @Router /admin/redis [POST]
func (h *httpDelivery) upsertRedis(w http.ResponseWriter, r *http.Request) {
	var reqBody request.UpsertRedisRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		h.Logger.Error(err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	res, err := h.Usecase.UpsertRedis(reqBody.Key, reqBody.Value)
	if err != nil {
		h.Logger.Error(err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, &response.RedisResponse{Value: res}, "")
}

// UserCredits godoc
// @Summary Delete Redis
// @Description Delete Redis
// @Tags Admin
// @Accept  json
// @Produce  json
// @Param key path string true "Redis key"
// @Success 200 {object} response.JsonResponse{data=string}
// @Router /admin/redis/{key} [DELETE]
func (h *httpDelivery) deleteRedis(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	redisKey := vars["key"]

	err = h.Usecase.DeleteRedis(redisKey)
	if err != nil {
		h.Logger.Error(err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, "", "")
}

// UserCredits godoc
// @Summary Delete Redis
// @Description Delete Redis
// @Tags Admin
// @Accept  json
// @Produce  json
// @Success 200 {object} response.JsonResponse{data=string}
// @Router /admin/redis [DELETE]
func (h *httpDelivery) deleteAllRedis(w http.ResponseWriter, r *http.Request) {
	res, err := h.Usecase.DeleteAllRedis()

	if err != nil {
		h.Logger.Error(err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, res, "")
}

//	Auto listing godoc
//
// @Summary Auto listing
// @Description  Auto listing
// @Tags Admin
// @Accept  json
// @Produce  json
// @Param request body request.ListNftIdsReq true " Auto listing"
// @Success 200 {object} response.JsonResponse{data=bool}
// @Router /admin/auto-listing [POST]
func (h *httpDelivery) autoListing(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)

	fmt.Println("iWalletAddress", iWalletAddress)

	userWalletAddr, ok := iWalletAddress.(string)
	if !ok {
		err := errors.New("Wallet address is incorect")
		h.Logger.Error("ctx.Value.Token", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	fmt.Println("userWalletAddr", userWalletAddr)

	// check admin user:
	profile, err := h.Usecase.GetUserProfileByBtcAddress(userWalletAddr)
	if err != nil {
		h.Logger.Error("h.Usecase.GetUserProfileByBtcAddress(", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	if !profile.IsAdmin {
		err := errors.New("permission denied")
		h.Logger.Error("permission", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	var reqBody request.ListNftIdsReq
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&reqBody)
	if err != nil {
		h.Logger.Error("decoder.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	res := h.Usecase.AutoListing(&reqBody)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, res, "")
}
