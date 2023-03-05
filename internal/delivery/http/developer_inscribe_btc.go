package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/btc"
	"rederinghub.io/utils/logger"
)

// @Summary BTC Generate receive wallet address
// @Description BTC Generate receive wallet address
// @Tags Inscribe
// @Accept json
// @Produce json
// @Param request body request.DeveloperCreateInscribeBtcReq true "Create mint request for dev"
// @Success 200 {object} response.InscribeBtcResp{}
// @Router /developer/inscribe [POST]
// @Security api-key
func (h *httpDelivery) developerCreateInscribe(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, vars map[string]string) (interface{}, error) {

			apiKey := r.URL.Query().Get("api-key")

			developerApiKey, _ := h.Usecase.Repo.FindIDeveloperKeyByApiKey(apiKey)
			if developerApiKey == nil {
				err := errors.New("api-key not found")
				h.Logger.Error("h.developerCreateInscribe", err.Error(), err)
				h.Response.RespondWithError(w, http.StatusUnauthorized, response.Error, err)
				return nil, err
			}
			if developerApiKey.Status <= 0 {
				err := errors.New("The api-key is invalid. Please contact the generative for the support.")
				h.Logger.Error("h.developerCreateInscribe", err.Error(), err)
				h.Response.RespondWithError(w, http.StatusUnauthorized, response.Error, err)
				return nil, err
			}
			// TODO: check request:
			now := time.Now()
			developerKeyRequests, _ := h.Usecase.Repo.FindDeveloperKeyRequests(apiKey)
			if developerKeyRequests == nil {
				h.Usecase.Repo.InsertDeveloperKeyRequests(&entity.DeveloperKeyRequests{

					RecordId: developerApiKey.UUID,
					ApiKey:   apiKey,

					EndpointName:    "inscribe-create", // todo move const/db
					EndpointUrl:     "",
					Status:          1,
					DayReqResetTime: &now,
					DayReqLastTime:  &now,
					DayReqCounter:   1,
				})
			} else {

				// check valid by day:
				now := time.Now()
				if developerKeyRequests.DayReqResetTime.Year() == now.Year() && developerKeyRequests.DayReqResetTime.YearDay() == now.YearDay() {
					if developerKeyRequests.DayReqCounter >= utils.DEVELOPER_INSCRIBE_MAX_REQUEST {
						err := errors.New("Limits reached.")
						h.Logger.Error("h.developerCreateInscribe", err.Error(), err)
						h.Response.RespondWithError(w, http.StatusTooManyRequests, response.Error, err)
						return nil, err
					} else {
						h.Usecase.Repo.IncreaseDeveloperReqCounter(apiKey)
					}
				} else {
					// reset:
					developerKeyRequests.DayReqCounter = 0
					developerKeyRequests.DayReqResetTime = &now
					developerKeyRequests.DayReqLastTime = &now
					h.Usecase.Repo.UpdateDeveloperKeyRequests(developerKeyRequests)
				}
			}

			var reqBody request.DeveloperCreateInscribeBtcReq
			err := json.NewDecoder(r.Body).Decode(&reqBody)
			if err != nil {
				return nil, err
			}
			reqUsecase := &structure.InscribeBtcReceiveAddrRespReq{}
			err = copier.Copy(reqUsecase, reqBody)
			if err != nil {
				return nil, err
			}

			if len(reqUsecase.FileName) == 0 {
				return nil, errors.New("Filename is required")
			}

			if len(reqUsecase.WalletAddress) == 0 {
				return nil, errors.New("WalletAddress is required")
			}

			if ok, _ := btc.ValidateAddress("btc", reqUsecase.WalletAddress); !ok {
				return nil, errors.New("WalletAddress is invalid")
			}

			if reqUsecase.FeeRate != 15 && reqUsecase.FeeRate != 20 && reqUsecase.FeeRate != 25 {
				return nil, errors.New("fee rate is invalid")
			}

			if len(reqUsecase.File) == 0 {
				return nil, errors.New("file is invalid")
			}

			btcWallet, err := h.Usecase.DeveloperCreateInscribe(ctx, *reqUsecase)
			if err != nil {
				logger.AtLog.Logger.Error("DeveloperCreateInscribe failed",
					zap.Any("payload", reqBody),
					zap.Error(err),
				)
				return nil, err
			}
			logger.AtLog.Logger.Info("DeveloperCreateInscribe successfully", zap.Any("response", btcWallet))
			return h.developerInscribeCreatedRespResp(btcWallet)
		},
	).ServeHTTP(w, r)
}

func (h *httpDelivery) developerInscribeBtcCreatedRespResp(input *entity.DeveloperInscribe) (*response.InscribeBtcResp, error) {
	resp := &response.InscribeBtcResp{}
	resp.UserAddress = input.UserAddress
	resp.Amount = input.Amount
	resp.MintFee = input.MintFee
	resp.SentTokenFee = input.SentTokenFee
	resp.OrdAddress = input.OrdAddress
	resp.FileURI = input.FileURI
	resp.IsConfirm = input.IsConfirm
	resp.InscriptionID = input.InscriptionID
	resp.Balance = input.Balance
	resp.TimeoutAt = fmt.Sprintf("%d", time.Now().Add(time.Hour*1).Unix()) // return FE in 1h. //TODO: need update
	resp.SegwitAddress = input.SegwitAddress
	return resp, nil
}

// @Summary BTC List Inscribe
// @Description BTC List Inscribe
// @Tags Inscribe
// @Accept json
// @Produce json
// @Success 200 {object} entity.Pagination{}
// @Router /developer/inscribe [GET]
// @Security ApiKeyAuth
func (h *httpDelivery) developerInscribeList(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, muxVars map[string]string) (interface{}, error) {
			userUuid := ctx.Value(utils.SIGNED_USER_ID).(string)
			page := entity.GetPagination(r)
			req := &entity.FilterDeveloperInscribeBT{
				BaseFilters: entity.BaseFilters{
					Limit: page.PageSize,
					Page:  page.Page,
				},
				UserUuid: &userUuid,
			}
			return h.Usecase.ListDeveloperInscribeBTC(req)
		},
	).ServeHTTP(w, r)
}

// @Summary BTC NFT Detail Inscribe
// @Description BTC NFT Detail Inscribe
// @Tags Inscribe
// @Accept json
// @Produce json
// @Param ID path string true "inscribe ID"
// @Success 200 {object} entity.InscribeBTCResp{}
// @Router /developer/inscribe/{ID} [GET]
// @Security ApiKeyAuth
func (h *httpDelivery) developerDetailInscribe(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uuid := vars["ID"]

	result, err := h.Usecase.DetailDeveloperInscribeBTC(uuid)
	if err != nil {
		h.Logger.Error("h.Usecase.DetailInscribeBTC", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")

}

// @Summary BTC Retry Inscribe
// @Description BTC Retry Inscribe
// @Tags Inscribe
// @Accept json
// @Produce json
// @Param ID path string true "inscribe ID"
// @Success 200
// @Router developer/inscribe/retry/{ID} [POST]
// @Security ApiKeyAuth
func (h *httpDelivery) developerRetryInscribeBTC(w http.ResponseWriter, r *http.Request) {

	// vars := mux.Vars(r)
	// id := vars["ID"]

	// err := h.Usecase.DeveloperRetryInscribeBTC(id)
	// if err != nil {
	// 	h.Logger.Error("h.Usecase.RetryInscribeBTC", err.Error(), err)
	// 	h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
	// 	return
	// }

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, true, "")

}

func (h *httpDelivery) developerInscribeCreatedRespResp(input *entity.DeveloperInscribe) (*response.InscribeBtcResp, error) {
	resp := &response.InscribeBtcResp{}

	resp.ID = input.UUID

	resp.UserAddress = input.UserAddress
	resp.Amount = input.Amount
	resp.MintFee = input.MintFee
	resp.SentTokenFee = input.SentTokenFee
	resp.OrdAddress = input.OrdAddress
	resp.FileURI = input.FileURI
	resp.IsConfirm = input.IsConfirm
	resp.InscriptionID = input.InscriptionID
	resp.Balance = input.Balance
	resp.TimeoutAt = fmt.Sprintf("%d", time.Now().Add(time.Hour*1).Unix()) // return FE in 1h. //TODO: need update
	resp.SegwitAddress = input.SegwitAddress
	return resp, nil
}
