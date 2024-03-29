package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
	"gopkg.in/ezzarghili/recaptcha-go.v4"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	faucetconst "rederinghub.io/utils/constants/faucet"
	"rederinghub.io/utils/logger"
)

func (h *httpDelivery) requestFaucet(w http.ResponseWriter, r *http.Request) {

	var reqBody request.FaucetReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
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
	if reqBody.RecaptchaResponse == "" {

		err = errors.New("the recaptcha is required.")
		logger.AtLog.Logger.Error("h.requestFaucet", zap.String("err", err.Error()))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	if len(h.Usecase.Config.CaptcharSecret) > 0 {
		captcha, _ := recaptcha.NewReCAPTCHA(h.Usecase.Config.CaptcharSecret, recaptcha.V3, 10*time.Second) // for v2 API get your secret from https://www.google.com/recaptcha/admin

		err = captcha.Verify(reqBody.RecaptchaResponse)
		if err != nil {
			//logger.AtLog.Logger.Error("h.requestFaucet.recaptcha.Verify", zap.String("err", err.Error()))
			//h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			//return
		}
	}
	result, err := h.Usecase.ApiCreateFaucet(reqBody.Address, reqBody.Url, reqBody.Txhash, reqBody.Type, reqBody.Source)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.GetFaucetPaymentInfo", zap.String("err", err.Error()))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")
}

func (h *httpDelivery) listFaucet(w http.ResponseWriter, r *http.Request) {

	address := r.URL.Query().Get("address")

	result, err := h.Usecase.ApiListCheckFaucet(address)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.ApiListCheckFaucet", zap.String("err", err.Error()))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")
}

func (h *httpDelivery) getFaucetConfig(w http.ResponseWriter, r *http.Request) {
	result := response.FaucetConfigRes{
		FaucetAmounts: make(map[string]string),
	}

	result.FaucetAmounts["normal"] = fmt.Sprintf("%.1f", float64(faucetconst.SpecialFaucetAmount)/float64(1e18))
	result.FaucetAmounts["bns"] = fmt.Sprintf("%.1f", float64(faucetconst.BNSFaucetAmount)/float64(1e18))
	result.FaucetAmounts["artifact"] = fmt.Sprintf("%.1f", float64(faucetconst.ArtifactFaucetAmount)/float64(1e18))
	result.FaucetAmounts["special"] = fmt.Sprintf("%.1f", float64(faucetconst.SpecialFaucetAmount)/float64(1e18))

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")
}

func (h *httpDelivery) getCurrentFaucetStep(w http.ResponseWriter, r *http.Request) {

	address := r.URL.Query().Get("address")

	faucetItems, err := h.Usecase.Repo.FindFaucetByTwitterNameOrAddress(address, address)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.getCurrentFaucetStep", zap.String("err", err.Error()))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	faucetStatus := make(map[string]string)

	for _, item := range faucetItems {
		// resItem := response.FaucetStatusRes{
		// 	CreatedAt: item.CreatedAt.Unix(),
		// 	Txhash:    item.Tx,
		// 	Status:    item.StatusStr,
		// }
		item.StatusStr = "Pending"
		if item.Status == 2 {
			item.StatusStr = "Processing"
		} else if item.Status == 3 {
			item.StatusStr = "Success"
		}
		if item.FaucetType == "" {
			faucetStatus["normal"] = item.StatusStr
		} else {
			if item.FaucetType == "dapps" {
				item.FaucetType = "bns"
			}
			faucetStatus[item.FaucetType] = item.StatusStr
		}
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, faucetStatus, "")
}

func (h *httpDelivery) getNonces(w http.ResponseWriter, r *http.Request) {

	address := r.URL.Query().Get("address")

	result, err := h.Usecase.ApiFaucetGetNonce(address)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.ApiFaucetGetNonce", zap.String("err", err.Error()))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")
}
