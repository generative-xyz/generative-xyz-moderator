package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/utils"
	"rederinghub.io/utils/logger"
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
		logger.AtLog.Logger.Error("err", zap.Error(err))
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
// @Success 200 {object} response.JsonResponse{data=[]response.RedisResponse}
// @Router /admin-test [GET]
func (h *httpDelivery) adminTest(w http.ResponseWriter, r *http.Request) {
	h.Usecase.JobCrawlTokenTxNotFromTokenUri()
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, "", "")
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
		logger.AtLog.Logger.Error("err", zap.Error(err))
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
		logger.AtLog.Logger.Error("err", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	res, err := h.Usecase.UpsertRedis(reqBody.Key, reqBody.Value)
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
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
		logger.AtLog.Logger.Error("err", zap.Error(err))
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
		logger.AtLog.Logger.Error("err", zap.Error(err))
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
		logger.AtLog.Logger.Error("ctx.Value.Token", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	fmt.Println("userWalletAddr", userWalletAddr)

	// check admin user:
	profile, err := h.Usecase.GetUserProfileByWalletAddress(userWalletAddr)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.GetUserProfileByWalletAddress(", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	if !profile.IsAdmin {
		err := errors.New("permission denied")
		logger.AtLog.Logger.Error("permission", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	var reqBody request.ListNftIdsReq
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&reqBody)
	if err != nil {
		logger.AtLog.Logger.Error("decoder.Decode", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	res := h.Usecase.AutoListing(&reqBody)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, res, "")
}
func (h *httpDelivery) checkRefundMintBtc(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)

	fmt.Println("iWalletAddress", iWalletAddress)

	userWalletAddr, ok := iWalletAddress.(string)
	if !ok {
		err := errors.New("Wallet address is incorect")
		logger.AtLog.Logger.Error("ctx.Value.Token", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	fmt.Println("userWalletAddr", userWalletAddr)

	// check admin user:
	profile, err := h.Usecase.GetUserProfileByWalletAddress(userWalletAddr)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.GetUserProfileByWalletAddress(", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	if !profile.IsAdmin {
		err := errors.New("permission denied")
		logger.AtLog.Logger.Error("permission", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	// res := h.Usecase.CheckRefundNftBtc()
	// h.Response.RespondSuccess(w, http.StatusOK, response.Success, res, "")
	// res := h.Usecase.CheckRefundNftBtc()
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, true, "")
}

func (h *httpDelivery) getMintFreeTemAddress(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)

	fmt.Println("iWalletAddress", iWalletAddress)

	userWalletAddr, ok := iWalletAddress.(string)
	if !ok {
		err := errors.New("Wallet address is incorect")
		logger.AtLog.Logger.Error("ctx.Value.Token", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	fmt.Println("userWalletAddr", userWalletAddr)

	// check admin user:
	profile, err := h.Usecase.GetUserProfileByWalletAddress(userWalletAddr)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.GetUserProfileByWalletAddress(", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	if !profile.IsAdmin {
		err := errors.New("permission denied")
		logger.AtLog.Logger.Error("permission", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	tx, err := h.Usecase.GenMintFreeTemAddress()

	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, tx, "")

}

func (h *httpDelivery) updateDeclaredNow(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)

	fmt.Println("iWalletAddress", iWalletAddress)

	userWalletAddr, ok := iWalletAddress.(string)
	if !ok {
		err := errors.New("Wallet address is incorect")
		logger.AtLog.Logger.Error("ctx.Value.Token", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	fmt.Println("userWalletAddr", userWalletAddr)

	// check admin user:
	profile, err := h.Usecase.GetUserProfileByWalletAddress(userWalletAddr)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.GetUserProfileByWalletAddress(", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	if !profile.IsAdmin {
		err := errors.New("permission denied")
		logger.AtLog.Logger.Error("permission", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	err = h.Usecase.APIAuctionDeclaredNow()

	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, true, "")

}

func (h *httpDelivery) updateWinnerFromContract(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)

	fmt.Println("iWalletAddress", iWalletAddress)

	userWalletAddr, ok := iWalletAddress.(string)
	if !ok {
		err := errors.New("Wallet address is incorect")
		logger.AtLog.Logger.Error("ctx.Value.Token", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	fmt.Println("userWalletAddr", userWalletAddr)

	// check admin user:
	profile, err := h.Usecase.GetUserProfileByWalletAddress(userWalletAddr)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.GetUserProfileByWalletAddress(", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	if !profile.IsAdmin {
		err := errors.New("permission denied")
		logger.AtLog.Logger.Error("permission", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	err = h.Usecase.APIAuctionCrawlWinnerNow()

	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, true, "")

}

// crontab on/off:
func (h *httpDelivery) updateEnabledJob(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)

	fmt.Println("iWalletAddress", iWalletAddress)

	userWalletAddr, ok := iWalletAddress.(string)
	if !ok {
		err := errors.New("Wallet address is incorect")
		logger.AtLog.Logger.Error("ctx.Value.Token", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	fmt.Println("userWalletAddr", userWalletAddr)

	// check admin user:
	profile, err := h.Usecase.GetUserProfileByWalletAddress(userWalletAddr)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.GetUserProfileByWalletAddress(", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	if !profile.IsAdmin {
		err := errors.New("permission denied")
		logger.AtLog.Logger.Error("permission", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	jobKey := r.URL.Query().Get("jobKey")

	status := true

	statusStr := r.URL.Query().Get("status")

	if statusStr == "0" {
		status = false
	}

	if len(jobKey) > 0 {

		fmt.Println("jobKey", jobKey)
		fmt.Println("status", status)

		_, err = h.Usecase.Repo.UpdateCronJobManagerStatusByJobKey(jobKey, status)

		fmt.Println("err UpdateCronJobManagerStatusByJobKey", err)
	} else {
		err = errors.New("jobKey empty")
	}

	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, true, "")

}

func (h *httpDelivery) requestFaucetAdmin(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)

	fmt.Println("iWalletAddress", iWalletAddress)

	userWalletAddr, ok := iWalletAddress.(string)
	if !ok {
		err := errors.New("Wallet address is incorect")
		logger.AtLog.Logger.Error("ctx.Value.Token", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	fmt.Println("userWalletAddr", userWalletAddr)

	// check admin user:
	profile, err := h.Usecase.GetUserProfileByWalletAddress(userWalletAddr)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.GetUserProfileByWalletAddress(", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	if !profile.IsAdmin {
		err := errors.New("permission denied")
		logger.AtLog.Logger.Error("permission", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	var reqBody request.FaucetAdminReq
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&reqBody)
	if err != nil {
		logger.AtLog.Logger.Error("httpDelivery.requestFaucet.Decode", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	if len(reqBody.Url) == 0 {
		err = errors.New("url invalid")
		logger.AtLog.Logger.Error("h.requestFaucet", zap.String("err", err.Error()))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	if len(reqBody.ListAddress) > 0 {
		result, err := h.Usecase.ApiAdminCreateBatchFaucet(reqBody.ListAddress, reqBody.Url, reqBody.Type, reqBody.Amount)
		if err != nil {
			logger.AtLog.Logger.Error("h.Usecase.GetFaucetPaymentInfo", zap.String("err", err.Error()))
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
		h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")

	} else if len(reqBody.MapAddress) > 0 {
		result, err := h.Usecase.ApiAdminCreateMapFaucet(reqBody.MapAddress, reqBody.Url, reqBody.Type)
		if err != nil {
			logger.AtLog.Logger.Error("h.Usecase.GetFaucetPaymentInfo", zap.String("err", err.Error()))
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
		h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")

	} else {
		result, err := h.Usecase.ApiAdminCreateFaucet(reqBody.Address, reqBody.Url, reqBody.Txhash, reqBody.Type, reqBody.Source)
		if err != nil {
			logger.AtLog.Logger.Error("h.Usecase.GetFaucetPaymentInfo", zap.String("err", err.Error()))
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}

		h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")
	}

}

func (h *httpDelivery) withdrawNewCityFunds(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)

	fmt.Println("iWalletAddress", iWalletAddress)

	userWalletAddr, ok := iWalletAddress.(string)
	if !ok {
		err := errors.New("Wallet address is incorect")
		logger.AtLog.Logger.Error("ctx.Value.Token", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	fmt.Println("userWalletAddr", userWalletAddr)

	// check admin user:
	profile, err := h.Usecase.GetUserProfileByWalletAddress(userWalletAddr)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.GetUserProfileByWalletAddress(", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	if !profile.IsAdmin {
		err := errors.New("permission denied")
		logger.AtLog.Logger.Error("permission", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	data, err := h.Usecase.ApiAdminCrawlFunds()

	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, data, "")

}
