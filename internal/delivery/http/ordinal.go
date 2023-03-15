package http

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"rederinghub.io/internal/delivery/http/response"
)

// UserCredits godoc
// @Summary get ordinal template
// @Description get ordinal template
// @Tags Ordinal collection
// @Accept  json
// @Produce  json
// @Success 200 {object} response.JsonResponse{data=response.CategoryResp}
// @Router /ordinal/collections/template [GET]
func (h *httpDelivery) getOrdinalTemplate(w http.ResponseWriter, r *http.Request) {
	zipfile, err := h.Usecase.GetOrdinalTemplate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	files, err := ioutil.ReadFile(zipfile.Name())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	des := fmt.Sprintf("attachment; filename=%s", zipfile.Name())
	w.Header().Set("Content-type", "application/zip")
	w.Header().Set("Content-Disposition", des)
	w.Write(files)
}

// UserCredits godoc
// @Summary Upload ordinal template
// @Description Upload ordinal template
// @Tags Ordinal collection
// @Content-Type: multipart/form-data
// @Param file formData file true "project.zip"
// @Produce  multipart/form-data
// @Success 200 {object} response.JsonResponse{data=response.CategoryResp}
// @Router /ordinal/collections [POST]
func (h *httpDelivery) uploadOrdinalTemplate(w http.ResponseWriter, r *http.Request) {
	_, err := h.Usecase.UploadOrdinalTemplate(r)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, nil, "")
}

