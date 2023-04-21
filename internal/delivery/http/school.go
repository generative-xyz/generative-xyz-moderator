package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/logger"
)

func (h *httpDelivery) schoolSearchDataset(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("text")
	if text == "" {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("text cannot be empty"))
		return
	}
}

func (h *httpDelivery) schoolUpload(w http.ResponseWriter, r *http.Request) {

	address := r.URL.Query().Get("address")
	if address == "" {
		ctx := r.Context()
		iUserID := ctx.Value(utils.SIGNED_USER_ID)
		userID, ok := iUserID.(string)
		if !ok {
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("address or accessToken cannot be empty"))
			return
		}
		userInfo, err := h.Usecase.UserProfile(userID)
		if err != nil {
			logger.AtLog.Logger.Error("httpDelivery.mintStatus.Usecase.UserProfile", zap.Error(err))
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
		address = userInfo.WalletAddress
	}
	if address == "" {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("address or accessToken cannot be empty"))
		return
	}

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
	usePFPDataset := false
	customDataset := []string{}
	_, datasetFile, err := r.FormFile("file")
	if err != nil {
		datasetStr := r.FormValue("datasets")
		if datasetStr == "" {
			usePFPDataset = true
		} else {
			customDataset = strings.Split(datasetStr, ",")
		}
	}
	var newJob entity.AISchoolJob
	uuid := uuid.NewString()
	if !usePFPDataset {
		dataset := datasetFile
		datasetSize := dataset.Size
		if datasetSize > 100000000 {
			http.Error(w, "Dataset size must be less than 100MB", http.StatusBadRequest)
			return
		}
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

		newJob = entity.AISchoolJob{
			JobID:       uuid,
			Params:      params,
			DatasetUUID: file.UUID,
			Status:      "waiting",
			CreatedBy:   address,
		}
	} else {
		newJob = entity.AISchoolJob{
			JobID:              uuid,
			Params:             params,
			DatasetUUID:        "",
			Status:             "waiting",
			CreatedBy:          address,
			UsePFPDataset:      usePFPDataset,
			CustomDatasetsUUID: customDataset,
		}
	}

	err = h.Usecase.Repo.InsertAISChoolJob(&newJob)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, uuid, "")
}

func (h *httpDelivery) schoolListProgress(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 20
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}

	address := r.URL.Query().Get("address")
	if address == "" {
		ctx := r.Context()
		iUserID := ctx.Value(utils.SIGNED_USER_ID)
		userID, ok := iUserID.(string)
		if !ok {
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("address or accessToken cannot be empty"))
			return
		}
		userInfo, err := h.Usecase.UserProfile(userID)
		if err != nil {
			logger.AtLog.Logger.Error("httpDelivery.mintStatus.Usecase.UserProfile", zap.Error(err))
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
		address = userInfo.WalletAddress
	}
	if address == "" {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("address or accessToken cannot be empty"))
		return
	}

	jobList, err := h.Usecase.Repo.GetAISchoolJobByCreator(address, int64(limit), int64(offset))
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	result := []response.AISchoolJobProgress{}
	for _, job := range jobList {
		params := structure.AISchoolModelParams{}
		err := json.Unmarshal([]byte(job.Params), &params)
		if err != nil {
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
		jobData := response.AISchoolJobProgress{
			JobID:       job.JobID,
			Status:      job.Status,
			Progress:    job.Progress / params.Epoch * 100,
			Output:      job.OutputLink,
			CompletedAt: job.CompletedAt,
			CreatedAt:   job.CreatedAt.Unix(),
			ModelName:   params.Name,
		}
		if job.Status == "completed" {
			jobData.Progress = 100
		}
		result = append(result, jobData)
	}
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")
}

// func (h *httpDelivery) schoolDownload(w http.ResponseWriter, r *http.Request) {
// 	jobID := r.URL.Query().Get("id")
// 	if jobID == "" {
// 		http.Error(w, "id cannot be empty", http.StatusBadRequest)
// 		return
// 	}
// 	result, err := h.Usecase.Repo.GetAISchoolJobByUUID(jobID)
// 	if err != nil {
// 		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
// 		return
// 	}
// 	if result.Status != "done" {
// 		http.Error(w, "Job is not done yet", http.StatusBadRequest)
// 		return
// 	}
// 	file, err := h.Usecase.Repo.GetFileByUUID(result.OutputUUID)
// 	if err != nil {
// 		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
// 		return
// 	}
// 	http.ServeFile(w, r, file.URL)
// }
