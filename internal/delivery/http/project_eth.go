package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/logger"
)

// UserCredits godoc
// @Summary Create project
// @Description Create projects
// @Tags ETH Project
// @Accept  json
// @Produce  json
// @Security Authorization
// @Param request body request.CreateETHProjectReq true "Create eth-project request"
// @Success 200 {object} response.JsonResponse{}
// @Router /project [POST]
func (h *httpDelivery) createEthProjects(w http.ResponseWriter, r *http.Request) {
	var reqBody request.CreateETHProjectReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		logger.AtLog.Logger.Error("decoder.Decode", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	ctx := r.Context()
	iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)
	walletAddress, ok := iWalletAddress.(string)
	if !ok {
		err := errors.New("Wallet address is incorect")
		logger.AtLog.Logger.Error("ctx.Value.Token", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase := &structure.CreateProjectReq{}
	err = copier.Copy(reqUsecase, reqBody)
	if err != nil {
		logger.AtLog.Logger.Error("copier.Copy", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase.CreatorAddrr = walletAddress
	message, err := h.Usecase.CreateProject(*reqUsecase)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.CreateProject", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp, err := h.projectToResp(message)
	if err != nil {
		logger.AtLog.Logger.Error("h.projectToResp", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

// UserCredits godoc
// @Summary get all projects
// @Description get all projects without project's status
// @Tags Project
// @Accept  json
// @Produce  json
// @Param contractAddress query string false "Filter project via contract address"
// @Param name query string false "filter project via name"
// @Param category query string false "filter project via category ids"
// @Param txHash query string false "txHash"
// @Param commitTxHash query string false "commitTxHash"
// @Param txHex query string false "txHex"
// @Param revealTxHash query string false "revealTxHash"
// @Param walletAddress query string false "walletAddress"
// @Param status query bool false "status"
// @Param isSynced query bool false "isSynced"
// @Param isHidden query bool false "isHidden"
// @Param limit query int false "limit"
// @Param page query int false "page"
// @Param sort query string false "newest, oldest, priority-asc, priority-desc, trending-score"
// @Param cursor query string false "The cursor returned in the previous response (used for getting the next page)."
// @Success 200 {object} response.JsonResponse{}
// @Router /project/all [GET]
func (h *httpDelivery) getAllProjects(w http.ResponseWriter, r *http.Request) {

	name := r.URL.Query().Get("name")
	categoriesRaw := r.URL.Query().Get("category")
	txHash := r.URL.Query().Get("txHash")
	commitTxHash := r.URL.Query().Get("commitTxHash")
	revealTxHash := r.URL.Query().Get("revealTxHash")
	txHex := r.URL.Query().Get("txHex")
	contractAddress := r.URL.Query().Get("contractAddress")
	walletAddress := r.URL.Query().Get("walletAddress")
	if walletAddress == "" {
		walletAddress = r.URL.Query().Get("creatorAddress")
	}
	isSyncedStr := r.URL.Query().Get("isSynced")
	isHiddenStr := r.URL.Query().Get("isHidden")
	statusStr := r.URL.Query().Get("status")

	categoryIds := strings.Split(categoriesRaw, ",")
	if categoriesRaw == "" {
		categoryIds = []string{}
	}

	baseF, err := h.BaseFilters(r)
	if err != nil {
		logger.AtLog.Logger.Error("BaseFilters", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	f := structure.FilterProjects{}
	f.BaseFilters = *baseF
	f.Name = &name
	f.TxHash = &txHash
	f.CommitTxHash = &commitTxHash
	f.RevealTxHash = &revealTxHash
	f.TxHex = &txHex
	f.CategoryIds = categoryIds
	f.ContractAddress = &contractAddress
	f.WalletAddress = &walletAddress

	if isHiddenStr != "" {
		hidden, err := strconv.ParseBool(isHiddenStr)
		if err == nil {
			f.IsHidden = &hidden
		}
	}

	if isSyncedStr != "" {
		isSynced, err := strconv.ParseBool(isSyncedStr)
		if err == nil {
			f.IsSynced = &isSynced
		}
	}

	if statusStr != "" {
		status, err := strconv.ParseBool(statusStr)
		if err == nil {
			f.Status = &status
		}
	}

	uProjects, err := h.Usecase.GetAllProjects(f)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.GetProjects", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	pResp := []response.ProjectResp{}
	iProjects := uProjects.Result
	projects := iProjects.([]entity.Projects)
	projectToGetFloorPrice := []string{}
	for _, project := range projects {
		projectToGetFloorPrice = append(projectToGetFloorPrice, project.TokenID)
	}

	for _, project := range projects {
		p, err := h.projectToResp(&project)
		if err != nil {
			logger.AtLog.Logger.Error("copier.Copy", zap.Error(err))
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
	
		pResp = append(pResp, *p)
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, h.PaginationResp(uProjects, pResp), "")
}

// UserCredits godoc
// @Summary Update project
// @Description Update projects
// @Tags Project
// @Accept  json
// @Produce  json
// @Param request body request.UpdateProjectReq true "Create profile request"
// @Param contractAddress path string true "contract adress"
// @Param projectID path string true "projectID adress"
// @Success 200 {object} response.JsonResponse{}
// @Router /project/{contractAddress}/{projectID} [PUT]
func (h *httpDelivery) updateProject(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	projectID := vars["projectID"]
	contractAddress := vars["contractAddress"]

	var reqBody request.UpdateProjectReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		logger.AtLog.Logger.Error("decoder.Decode", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase := &structure.UpdateProjectReq{
		ContracAddress: contractAddress,
		TokenID:        projectID,
		Priority:       reqBody.Priority,
	}

	logger.AtLog.Logger.Info("reqUsecase", zap.Any("reqUsecase", reqUsecase))

	message, err := h.Usecase.UpdateProject(*reqUsecase)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.CreateProject", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp, err := h.projectToResp(message)
	if err != nil {
		logger.AtLog.Logger.Error("h.projectToResp", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}