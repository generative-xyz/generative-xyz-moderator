package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/usecase/structure"
)

// UserCredits godoc
// @Summary get token uri data
// @Description get token uri data
// @Tags token_uri
// @Accept  json
// @Produce  json
// @Param contractAddress path string true "contract address"
// @Param tokenID path string true "token ID"
// @Success 200 {object} response.JsonResponse{data=response.GeneratedMessage}
// @Router /token/{contractAddress}/{tokenID} [GET]
func (h *httpDelivery) tokenURI(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("tokenURI", r)
	defer h.Tracer.FinishSpan(span, log )

	vars := mux.Vars(r)
	contractAddress := vars["contractAddress"]
	span.SetTag("contractAddress", contractAddress)
	
	tokenID := vars["tokenID"]
	span.SetTag("tokenID", tokenID)

	message, err := h.Usecase.GetToken(span, structure.GetTokenMessageReq{
		ContractAddress: contractAddress,
		TokenID: tokenID,
	})

	if err != nil {
		log.Error("h.Usecase.GetToken", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	log.SetData("resp.message", message)
	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, message, "")
}
