package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
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

	apiKey := r.URL.Query().Get("api-key")

	developerApiKey, _ := h.Usecase.Repo.FindIDeveloperKeyByApiKey(apiKey)
	if developerApiKey == nil {
		err := errors.New("api-key not found")

		h.Logger.Error("h.developerCreateInscribe", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusUnauthorized, response.Error, err)
		return
	}

	// TODO: check request:
	now := time.Now().UTC()
	developerKeyRequests, _ := h.Usecase.Repo.FindDeveloperKeyRequests(apiKey) // TODO: load request by api endpoint
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
		now := time.Now().UTC()

		fmt.Println("now:", now)

		fmt.Println("developerKeyRequests.DayReqResetTime.Year():", developerKeyRequests.DayReqResetTime.Year())
		fmt.Println("now.Year().now.Year():", now.Year())

		fmt.Println("now.Year().now.Year():", now.Year())

		if developerKeyRequests.DayReqResetTime.Year() == now.Year() && developerKeyRequests.DayReqResetTime.YearDay() == now.YearDay() {
			if developerKeyRequests.DayReqCounter >= utils.DEVELOPER_INSCRIBE_MAX_REQUEST {
				err := errors.New("Limits reached.")
				h.Logger.Error("h.developerCreateInscribe", err.Error(), err)
				h.Response.RespondWithError(w, http.StatusTooManyRequests, response.Error, err)
				return
			} else {
				err := h.Usecase.Repo.IncreaseDeveloperReqCounter(apiKey)
				fmt.Println("IncreaseDeveloperReqCounter err: ", err)

			}
		} else {
			// reset:
			developerKeyRequests.DayReqCounter = 1
			developerKeyRequests.DayReqResetTime = &now
			developerKeyRequests.DayReqLastTime = &now
			h.Usecase.Repo.UpdateDeveloperKeyRequests(developerKeyRequests)
		}
	}

	var reqBody request.DeveloperCreateInscribeBtcReq
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusUnauthorized, response.Error, err)
		return
	}
	reqUsecase := &structure.InscribeBtcReceiveAddrRespReq{}
	err = copier.Copy(reqUsecase, reqBody)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusUnauthorized, response.Error, err)
		return
	}

	if len(reqUsecase.FileName) == 0 {
		err = errors.New("Filename is required")
		h.Logger.Error("h.developerCreateInscribe", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusUnauthorized, response.Error, err)
		return
	}

	typeFiles := strings.Split(reqUsecase.FileName, ".")
	if len(typeFiles) < 2 {
		err := errors.New("File name invalid")
		h.Logger.Error("h.developerCreateInscribe", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusUnauthorized, response.Error, err)
		return
	}

	if len(reqUsecase.WalletAddress) == 0 {
		err = errors.New("WalletAddress is required")
	}

	if ok, _ := btc.ValidateAddress("btc", reqUsecase.WalletAddress); !ok {
		err = errors.New("WalletAddress is invalid")
	}

	/*if reqUsecase.FeeRate != 15 && reqUsecase.FeeRate != 20 && reqUsecase.FeeRate != 25 {
		err = errors.New("fee rate is invalid")
	}*/

	if len(reqUsecase.File) == 0 {
		err = errors.New("file is invalid")
	}

	reqUsecase.DeveloperKeyUuid = developerApiKey.UUID

	btcWallet, err := h.Usecase.DeveloperCreateInscribe(*reqUsecase)
	if err != nil {
		logger.AtLog.Logger.Error("DeveloperCreateInscribe failed",
			zap.Any("payload", reqBody),
			zap.Error(err),
		)
		h.Logger.Error("h.developerCreateInscribe", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusUnauthorized, response.Error, err)
		return
	}
	logger.AtLog.Logger.Info("DeveloperCreateInscribe successfully", zap.Any("response", btcWallet))

	resp, err := h.developerInscribeCreatedRespResp(btcWallet)

	if err != nil {
		h.Logger.Error("h.Usecase.developerCreateInscribe.developerInscribeCreatedRespResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")

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
