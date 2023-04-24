package http

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

	files, err := h.Usecase.Repo.FindPresetDatasetByName(text)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	result := []response.AISchoolPresetDataset{}
	for _, file := range files {
		result = append(result, response.AISchoolPresetDataset{
			Name:        file.Name,
			Thumbnail:   file.Thumbnail,
			UUID:        file.ID.Hex(),
			Creator:     file.Creator,
			IsPrivate:   file.IsPrivate,
			Size:        file.Size,
			NumOfAssets: file.NumOfAssets,
		})
	}
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")
}

// func (h *httpDelivery) schoolUpload(w http.ResponseWriter, r *http.Request) {

// 	address := r.URL.Query().Get("address")
// 	if address == "" {
// 		ctx := r.Context()
// 		iUserID := ctx.Value(utils.SIGNED_USER_ID)
// 		userID, ok := iUserID.(string)
// 		if !ok {
// 			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("address or accessToken cannot be empty"))
// 			return
// 		}
// 		userInfo, err := h.Usecase.UserProfile(userID)
// 		if err != nil {
// 			logger.AtLog.Logger.Error("httpDelivery.mintStatus.Usecase.UserProfile", zap.Error(err))
// 			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
// 			return
// 		}
// 		address = userInfo.WalletAddress
// 	}
// 	if address == "" {
// 		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("address or accessToken cannot be empty"))
// 		return
// 	}

// 	params := r.FormValue("params")

// 	var paramsStruct structure.AISchoolModelParams

// 	err := json.Unmarshal([]byte(params), &paramsStruct)
// 	if err != nil {
// 		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
// 		return
// 	}
// 	err = paramsStruct.SelfValidate()
// 	if err != nil {
// 		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
// 		return
// 	}
// 	usePFPDataset := false
// 	customDataset := []string{}
// 	_, datasetFile, err := r.FormFile("file")
// 	if err != nil {
// 		datasetStr := r.FormValue("datasets")
// 		if datasetStr == "" {
// 			usePFPDataset = true
// 		} else {
// 			customDataset = strings.Split(datasetStr, ",")
// 		}
// 	}
// 	var newJob entity.AISchoolJob
// 	uuid := uuid.NewString()
// 	if !usePFPDataset {
// 		dataset := datasetFile
// 		datasetSize := dataset.Size
// 		if datasetSize > 100000000 {
// 			http.Error(w, "Dataset size must be less than 100MB", http.StatusBadRequest)
// 			return
// 		}
// 		datasetName := dataset.Filename + "-" + uuid

// 		log.Println("datasetName", datasetName)
// 		log.Println("datasetSize", datasetSize)
// 		d, _ := json.MarshalIndent(paramsStruct, "", " ")
// 		log.Println("paramsStruct", string(d))
// 		file, err := h.Usecase.UploadDatasetFile(r, uuid)
// 		if err != nil {
// 			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
// 			return
// 		}
// 		log.Println("file", file)

// 		newJob = entity.AISchoolJob{
// 			JobID:       uuid,
// 			Params:      params,
// 			DatasetUUID: file.UUID,
// 			Status:      "waiting",
// 			CreatedBy:   address,
// 		}
// 	} else {
// 		newJob = entity.AISchoolJob{
// 			JobID:              uuid,
// 			Params:             params,
// 			DatasetUUID:        "",
// 			Status:             "waiting",
// 			CreatedBy:          address,
// 			UsePFPDataset:      usePFPDataset,
// 			CustomDatasetsUUID: customDataset,
// 		}
// 	}

// 	err = h.Usecase.Repo.InsertAISChoolJob(&newJob)
// 	if err != nil {
// 		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
// 		return
// 	}

// 	h.Response.RespondSuccess(w, http.StatusOK, response.Success, uuid, "")
// }

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

func (h *httpDelivery) schoolUploadDataset(w http.ResponseWriter, r *http.Request) {
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
	datasetName := r.FormValue("name")

	isPrivateStr := r.FormValue("private")
	isPrivate := false
	if isPrivateStr == "true" {
		isPrivate = true
	}

	_, datasetFile, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file cannot be empty", http.StatusBadRequest)
		return
	}
	dataset := datasetFile
	datasetSize := dataset.Size

	if datasetSize > 100000000 {
		http.Error(w, "Dataset size must be less than 100MB", http.StatusBadRequest)
		return
	}

	log.Println("datasetName", datasetName)
	log.Println("datasetSize", datasetSize)

	datasetPath := fmt.Sprintf("ai-school-dataset/%s/%s", address, datasetName)

	file, err := h.Usecase.UploadDatasetFile(r, datasetPath)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	filedata, err := dataset.Open()
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	dataseBytes := make([]byte, datasetSize)
	_, err = filedata.Read(dataseBytes)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	br := bytes.NewReader(dataseBytes)

	zr, err := zip.NewReader(br, int64(len(dataseBytes)))
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	numOfAssets := 0

	for _, fileInfo := range zr.File {
		if strings.Index(strings.ToLower(fileInfo.Name), strings.ToLower("__MACOSX")) > -1 {
			continue
		}

		if strings.Index(strings.ToLower(fileInfo.Name), strings.ToLower(".DS_Store")) > -1 {
			continue
		}
		numOfAssets++
	}

	created, err := h.Usecase.CreateDataset(file.UUID, file.FileName, datasetName, address, file.FileSize, numOfAssets, isPrivate)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, created.UUID, "")
}

func (h *httpDelivery) schoolSubmitModel(w http.ResponseWriter, r *http.Request) {
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

	datasetsStr := r.FormValue("datasets")
	datasets := []string{}
	err = json.Unmarshal([]byte(datasetsStr), &datasets)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	uuid := uuid.NewString()
	newJob := entity.AISchoolJob{
		JobID:     uuid,
		Params:    params,
		Status:    "waiting",
		CreatedBy: address,
		Datasets:  datasets,
	}

	err = h.Usecase.Repo.InsertAISChoolJob(&newJob)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, uuid, "")
}

func (h *httpDelivery) schoolDeleteDataset(w http.ResponseWriter, r *http.Request) {
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
	datasetUUID := r.URL.Query().Get("uuid")
	err := h.Usecase.DeleteDataset(datasetUUID, address)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, "ok", "")
}

func (h *httpDelivery) schoolListDataset(w http.ResponseWriter, r *http.Request) {
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
	datasets, err := h.Usecase.ListDataset(address, int64(limit), int64(offset))
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	result := []response.AISchoolPresetDataset{}
	for _, dataset := range datasets {
		result = append(result, response.AISchoolPresetDataset{
			Name:        dataset.Name,
			Thumbnail:   dataset.Thumbnail,
			UUID:        dataset.UUID,
			Creator:     dataset.Creator,
			IsPrivate:   dataset.IsPrivate,
			Size:        dataset.Size,
			NumOfAssets: dataset.NumOfAssets,
		})
	}
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")
}
