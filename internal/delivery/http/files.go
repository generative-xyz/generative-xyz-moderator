package http

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"rederinghub.io/internal/delivery/http/request"

	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/fileutil"
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
	file, err := h.Usecase.UploadFile(r)
	if err != nil {
		h.Logger.Error("h.Usecase.UploadFile", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp := &response.FileRes{}
	err = response.CopyEntityToRes(resp, file)
	if err != nil {
		h.Logger.Error("response.CopyEntityToRes", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

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

	ctx := r.Context()

	var reqBody request.CreateMultipartUploadRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		h.Logger.Error("decoder.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	h.Logger.Info("request.decoder", decoder)

	if err = reqBody.SelfValidate(); err != nil {
		h.Logger.Error("SelfValidate", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	uploadID, err := h.Usecase.CreateMultipartUpload(ctx, reqBody.Group, reqBody.FileName)

	if err != nil {
		h.Logger.Error("h.Usecase.CreateMultipartUpload", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Logger.Info("resp.uploadID", uploadID)

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

	ctx := r.Context()

	_, handler, err := r.FormFile("file")
	if err != nil {
		err = errors.Wrap(err, "error getting file from part")
		h.Logger.Error("h.Usecase.UploadFile", err.Error(), err)
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
		h.Logger.Error("h.Usecase.UploadFile", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	data, err := handler.Open()
	if err != nil {
		err = errors.Wrap(err, "error open file from handler")
		h.Logger.Error("FileHandlerError", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	err = h.Usecase.UploadPart(ctx, uploadID, data, handler.Size, partNumber)
	if err != nil {
		h.Logger.Error("h.Usecase.UploadPart", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusUnprocessableEntity, response.Error, err)
		return
	}

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

	ctx := r.Context()

	vars := mux.Vars(r)
	uploadID := vars["uploadID"]

	fileURL, err := h.Usecase.CompleteMultipartUpload(ctx, uploadID)

	if err != nil {
		h.Logger.Error("h.Usecase.CompleteMultipartUpload", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusUnprocessableEntity, response.Error, err)
		return
	}

	h.Logger.Info("resp.fileURL", fileURL)

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

	var reqBody structure.MinifyDataResp
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		h.Logger.Error("decoder.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	data, err := h.Usecase.MinifyFiles(reqBody)
	if err != nil {
		h.Logger.Error("h.Usecase.MinifyFiles", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

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
	reqBody := &structure.DeflateDataResp{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(reqBody)
	if err != nil {
		h.Logger.Error("decoder.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	err = h.Usecase.DeflateString(reqBody)
	if err != nil {
		h.Logger.Error(" h.Usecase.DeflateString", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, reqBody, "")
}

// @Summary Upload file
// @Description Upload file
// @Tags Files
// @Accept json
// @Produce json
// @Param request body request.FileResize true "Base64 File Request"
// @Success 200 {object} request.FileResize{}
// @Router /files/image/resize [POST]
func (h *httpDelivery) resizeImage(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(func(ctx context.Context, r *http.Request, vars map[string]string) (interface{}, error) {
		var reqBody request.FileResize
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			return nil, err
		}
		coI := strings.Index(reqBody.File, ",")
		dec, err := base64.StdEncoding.DecodeString(reqBody.File[coI+1:])
		if err != nil {
			return nil, err
		}
		exts := strings.Split(strings.TrimSuffix(reqBody.File[5:coI], ";base64"), "/")
		if len(exts) < 2 {
			return nil, errors.New("image not support")
		}
		imgByte, err := fileutil.ResizeImage(dec, exts[1], fileutil.MaxImageByteSize)
		if err != nil {
			return nil, err
		}
		return &request.FileResize{
			File: reqBody.File[:coI+1] + base64.StdEncoding.EncodeToString(imgByte),
		}, nil
	}).ServeHTTP(w, r)
}
