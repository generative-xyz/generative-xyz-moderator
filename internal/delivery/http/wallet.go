package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
)

func (h *httpDelivery) inscriptionByOutput(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("httpDelivery.inscriptionByOutput", r)
	defer h.Tracer.FinishSpan(span, log)
	h.Response.SetLog(h.Tracer, span)

	var reqBody request.InscriptionByOutput
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Error("httpDelivery.inscriptionByOutput.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	result, err := h.Usecase.InscriptionsByOutputs(reqBody.Outputs)
	if err != nil {
		log.Error("httpDelivery.inscriptionByOutput.Usecase.InscriptionsByOutputs", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")
}

func (h *httpDelivery) walletInfo(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("httpDelivery.walletInfo", r)
	defer h.Tracer.FinishSpan(span, log)
	h.Response.SetLog(h.Tracer, span)

	address := r.URL.Query().Get("address")

	if address == "" {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("address cannot be empty"))
		return
	}
	result, err := h.Usecase.GetBTCWalletInfo(address)
	if err != nil {
		log.Error("httpDelivery.walletInfo.Usecase.GetBTCWalletInfo", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")
}
