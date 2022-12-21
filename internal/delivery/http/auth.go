package http

import (
	"encoding/json"
	"net/http"

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
// @Router /v1/auth/nonce [POST]
func (h *httpDelivery) generateMessage(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("generateMessage", r)
	defer h.Tracer.FinishSpan(span, log )

	var reqBody request.GenerateMessageRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Error("decoder.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	log.SetData("request.decoder", decoder)
	message, err := h.Usecase.GenerateMessage(span, structure.GenerateMessage{
		Address: reqBody.Address,
	})

	if err != nil {
		log.Error("h.Usecase.GenerateMessage(", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	log.SetData("resp.message", message)
	h.Response.SetLog(h.Tracer, span)
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
// @Router /v1/auth/nonce/verify [POST]
func (h *httpDelivery) verifyMessage(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("messages.Verify", r)
	defer h.Tracer.FinishSpan(span, log )

	var reqBody request.VerifyMessageRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Error("decoder.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	
	log.SetData("request.decoder", decoder)
	verified, err := h.Usecase.VerifyMessage(span, structure.VerifyMessage{
		Signature: reqBody.Sinature,
		Address: reqBody.Address,
	})

	log.SetData("verified", verified)
	if err != nil {
		log.Error("h.Usecase.GenerateMessage(", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	resp := response.VerifyResponse{
		IsVerified: verified.IsVerified,
		Token: verified.Token,
		RefreshToken: verified.RefreshToken,
	}

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}
