package http

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/utils"
)

// UserCredits godoc
// @Summary Create referral
// @Description Create referral
// @Tags Referral
// @Accept  json
// @Produce  json
// @Param referrerID path string true "referrerID"
// @Security Authorization
// @Success 200 {object} response.JsonResponse{data=bool}
// @Router /referral/{referrerID} [POST]
func (h *httpDelivery) createReferral(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("messages.createReferral", r)
	defer h.Tracer.FinishSpan(span, log )

	ctx := r.Context()
	iUserID := ctx.Value(utils.SIGNED_USER_ID)
	referreeID, ok := iUserID.(string)

	if !ok {
		err := errors.New( "Token is incorect")
		log.Error("ctx.Value.Token",  err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	vars := mux.Vars(r)
	referrerID := vars["referrerID"]
	span.SetTag("referrerID", referrerID)

	err := h.Usecase.CreateReferral(referrerID, referreeID)

	if err != nil {
		log.Error("h.Usecase.CreateReferral", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, true, "")
}
