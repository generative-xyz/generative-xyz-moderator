package http

import (
	"encoding/json"
	"net/http"

	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/usecase/structure"
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
// @Router /files [POST]
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

// UserCredits godoc
// @Summary Upload file
// @Description Upload file
// @Tags Files
// @Content-Type: application/json
// @Security Authorization
// @Param request body structure.MinifyDataResp true "Data for minify"
// @Success 200 {object} response.JsonResponse{data=response.FileRes}
// @Router /files/minify [POST]
func (h *httpDelivery) minifyFiles(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("httpDelivery.minifyFiles", r)
	defer h.Tracer.FinishSpan(span, log )
	
	var reqBody structure.MinifyDataResp
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Error("decoder.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	
	data, err := h.Usecase.MinifyFiles(span, reqBody)
	if err != nil {
		log.Error("h.Usecase.MinifyFiles", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, data, "")
}


// UserCredits godoc
// @Summary Deflate a string
// @Description Deflate a string
// @Tags Files
// @Content-Type: application/json
// @Security Authorization
// @Param request body structure.DeflateDataResp true "Data for minify"
// @Success 200 {object} response.JsonResponse{data=structure.DeflateDataResp}
// @Router /files/deflate [POST]
func (h *httpDelivery) deflate(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("httpDelivery.deflate", r)
	defer h.Tracer.FinishSpan(span, log )
	
	reqBody := &structure.DeflateDataResp{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(reqBody)
	if err != nil {
		log.Error("decoder.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	
	err = h.Usecase.DeflateString(span, reqBody)
	if err != nil {
		log.Error(" h.Usecase.DeflateString", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, reqBody, "")
}

