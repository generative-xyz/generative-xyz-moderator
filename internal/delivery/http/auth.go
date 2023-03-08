package http

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/usecase/structure"

	"rederinghub.io/internal/delivery/http/request"
)

// UserCredits godoc
// @Summary Generate a message
// @Description Generate a message for user's wallet
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param request body request.GenerateMessageRequest true "Generate message request"
// @Success 200 {object} response.JsonResponse{data=response.GeneratedMessage}
// @Router /auth/nonce [POST]
func (h *httpDelivery) generateMessage(w http.ResponseWriter, r *http.Request) {

	var reqBody request.GenerateMessageRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		h.Logger.Error(err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	err = reqBody.SelfValidate()
	if err != nil {
		h.Logger.Error(err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Logger.LogAny("generateMessage", zap.Any("reqBody", reqBody))
	message, err := h.Usecase.GenerateMessage(structure.GenerateMessage{
		Address: *reqBody.Address,
	})

	if err != nil {
		h.Logger.Error(err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Logger.Info("resp.message", message)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, response.GeneratedMessage{Message: *message}, "")
}

// UserCredits godoc
// @Summary Verified the generated message
// @Description Verified the generated message
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param request body request.VerifyMessageRequest true "Verify message request"
// @Success 200 {object} response.JsonResponse{data=response.VerifyResponse}
// @Router /auth/nonce/verify [POST]
func (h *httpDelivery) verifyMessage(w http.ResponseWriter, r *http.Request) {

	var reqBody request.VerifyMessageRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		h.Logger.Error(err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	err = reqBody.SelfValidate()
	if err != nil {
		h.Logger.Error(err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Logger.Info("request.decoder", decoder)
	verifyMessage := structure.VerifyMessage{
		ETHSignature:     *reqBody.ETHSinature,
		Signature:        *reqBody.Sinature,
		Address:          *reqBody.Address,         //eth
		AddressBTC:       reqBody.AddressBTC,       //btc taproot addree -> use for transfer nft
		AddressBTCSegwit: reqBody.AddressBTCSegwit, //btc segwit address -> use for verify signature
		MessagePrefix:    reqBody.MessagePrefix,    //btc prefix message
	}
	verified, err := h.Usecase.VerifyMessage(verifyMessage)

	h.Logger.Info("verified", verified)
	if err != nil {
		h.Logger.Error(err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp := response.VerifyResponse{
		IsVerified:   verified.IsVerified,
		Token:        verified.Token,
		RefreshToken: verified.RefreshToken,
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}
