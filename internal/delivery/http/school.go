package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
)

func (h *httpDelivery) schoolUpload(w http.ResponseWriter, r *http.Request) {

	params := r.FormValue("params")

	var paramsStruct structure.AISchoolModelParams

	err := json.Unmarshal([]byte(params), &paramsStruct)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	err = paramsStruct.SelfValidate()
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	_, datasetFile, err := r.FormFile("file")
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	dataset := datasetFile
	datasetSize := dataset.Size
	if datasetSize > 100000000 {
		http.Error(w, "Dataset size must be less than 100MB", http.StatusBadRequest)
		return
	}
	uuid := uuid.NewString()
	datasetName := dataset.Filename + "-" + uuid

	log.Println("datasetName", datasetName)
	log.Println("datasetSize", datasetSize)
	d, _ := json.MarshalIndent(paramsStruct, "", " ")
	log.Println("paramsStruct", string(d))
	file, err := h.Usecase.UploadDatasetFile(r, uuid)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	log.Println("file", file)

	newJob := entity.AISchoolJob{
		JobID:       uuid,
		Params:      params,
		DatasetUUID: file.UUID,
		Status:      "waiting",
	}

	err = h.Usecase.Repo.InsertAISChoolJob(&newJob)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, uuid, "")
}

func (h *httpDelivery) schoolCheckProgress(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("id")
	if jobID == "" {
		http.Error(w, "id cannot be empty", http.StatusBadRequest)
		return
	}
	result, err := h.Usecase.Repo.GetAISchoolJobByUUID(jobID)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, response.AISchoolJobProgress{
		result.Progress,
		result.Status,
	}, "")
}

func (h *httpDelivery) schoolDownload(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("id")
	if jobID == "" {
		http.Error(w, "id cannot be empty", http.StatusBadRequest)
		return
	}
	result, err := h.Usecase.Repo.GetAISchoolJobByUUID(jobID)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	if result.Status != "done" {
		http.Error(w, "Job is not done yet", http.StatusBadRequest)
		return
	}
	file, err := h.Usecase.Repo.GetFileByUUID(result.OutputUUID)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	http.ServeFile(w, r, file.URL)
}
