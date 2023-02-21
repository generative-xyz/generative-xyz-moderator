package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"rederinghub.io/internal/delivery/http/request"

	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/usecase/structure"
)

// UploadFile godoc
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
	defer h.Tracer.FinishSpan(span, log)
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

// CreateMultipartUpload godoc
// @Summary Create multipart upload
// @Description Create multipart upload.
// @Tags Files
// @Accept  json
// @Produce  json
// @Security Authorization
// @Param request body request.CreateMultipartUploadRequest true "Create multipart upload request"
// @Success 200 {object} response.JsonResponse{data=response.FileResponse}
// @Router /files/multipart [POST]
func (h *httpDelivery) CreateMultipartUpload(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("CreateMultipartUpload", r)
	defer h.Tracer.FinishSpan(span, log)

	ctx := r.Context()

	var reqBody request.CreateMultipartUploadRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Error("decoder.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	log.SetData("request.decoder", decoder)

	if err = reqBody.SelfValidate(); err != nil {
		log.Error("SelfValidate", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	uploadID, err := h.Usecase.CreateMultipartUpload(ctx, span, reqBody.Group, reqBody.FileName)

	if err != nil {
		log.Error("h.Usecase.CreateMultipartUpload", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	log.SetData("resp.uploadID", uploadID)
	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, response.FileResponse{UploadID: *uploadID}, "")
}

// UploadPart godoc
// @Summary Upload multipart file
// @Description Upload multipart file
// @Tags Files
// @Content-Type: multipart/form-data
// @Security Authorization
// @Produce  multipart/form-data
// @Param file formData file true "file"
// @Param uploadID path string true "upload ID"
// @Param partNumber query string  false  "part number"
// @Success 200 {object} response.JsonResponse{data=response.FileRes}
// @Router /files/multipart/{uploadID} [PUT]
func (h *httpDelivery) UploadPart(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("httpDelivery.UploadPart", r)
	defer h.Tracer.FinishSpan(span, log)

	ctx := r.Context()

	_, handler, err := r.FormFile("file")
	if err != nil {
		err = errors.Wrap(err, "error getting file from part")
		log.Error("h.Usecase.UploadFile", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	vars := mux.Vars(r)
	uploadID := vars["uploadID"]
	var partNumber int
	partNumberStr := r.URL.Query().Get("partNumber")
	if partNumberStr == "" {
		err = errors.New("missing part number")
	} else {
		partNumber, err = strconv.Atoi(partNumberStr)
	}

	if err != nil {
		err = errors.Wrap(err, "error getting partNumber from reqeust")
		log.Error("h.Usecase.UploadFile", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	data, err := handler.Open()
	if err != nil {
		err = errors.Wrap(err, "error open file from handler")
		log.Error("FileHandlerError", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	err = h.Usecase.UploadPart(ctx, span, uploadID, data, handler.Size, partNumber)
	if err != nil {
		log.Error("h.Usecase.UploadPart", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusUnprocessableEntity, response.Error, err)
		return
	}

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, map[string]interface{}{}, "")
}

// CompleteMultipartUpload godoc
// @Summary Finish multipart upload
// @Description Finish multipart upload
// @Tags Files
// @Accept  json
// @Produce  json
// @Security Authorization
// @Param uploadID path string true "upload ID"
// @Success 200 {object} response.JsonResponse{data=response.MultipartUploadResponse}
// @Router /files/multipart/{uploadID} [POST]
func (h *httpDelivery) CompleteMultipartUpload(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("httpDelivery.FinishMultipartUpload", r)
	defer h.Tracer.FinishSpan(span, log)

	ctx := r.Context()

	vars := mux.Vars(r)
	uploadID := vars["uploadID"]

	fileURL, err := h.Usecase.CompleteMultipartUpload(ctx, span, uploadID)

	if err != nil {
		log.Error("h.Usecase.CompleteMultipartUpload", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusUnprocessableEntity, response.Error, err)
		return
	}

	log.SetData("resp.fileURL", fileURL)
	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, response.MultipartUploadResponse{FileURL: *fileURL}, "")

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
	defer h.Tracer.FinishSpan(span, log)

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
	defer h.Tracer.FinishSpan(span, log)

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
