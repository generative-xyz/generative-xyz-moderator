package http

import (
	"net/http"

	"rederinghub.io/internal/delivery/http/response"
)

// UserCredits godoc
// @Summary Upload file
// @Description Upload file
// @Tags Files
// @Content-Type: multipart/form-data
// @Security Authorization
// @Param file formData file true "file"
// @Produce  multipart/form-data
// @Success 200 {object} response.JsonResponse{data=response.FileRes}
// @Router /v1/files [POST]
func (h *httpDelivery) UploadFile(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("httpDelivery.UploadFile", r)
	defer h.Tracer.FinishSpan(span, log )
	file, err := h.Usecase.UploadFile(span, r)
	if err != nil {
		log.Error("h.Usecase.UploadFile", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp := &response.FileRes{}
	err = response.CopyEntityToRes(resp, file)
	if err != nil {
		log.Error("response.CopyEntityToRes", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}
