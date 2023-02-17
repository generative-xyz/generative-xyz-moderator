package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
)

// UserCredits godoc
// @Summary Create project
// @Description Create projects
// @Tags Project
// @Accept  json
// @Produce  json
// @Param request body request.CreateProjectReq true "Create profile request"
// @Success 200 {object} response.JsonResponse{}
// @Router /project [POST]
func (h *httpDelivery) createProjects(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("createProjects", r)
	defer h.Tracer.FinishSpan(span, log)

	var reqBody request.CreateProjectReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Error("decoder.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase := &structure.CreateProjectReq{}
	err = copier.Copy(reqUsecase, reqBody)
	if err != nil {
		log.Error("copier.Copy", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	log.SetData("reqUsecase", reqUsecase)
	log.SetTag("contractAddress", reqUsecase.ContractAddress)
	log.SetTag("tokenID", reqUsecase.TokenID)

	message, err := h.Usecase.CreateProject(span, *reqUsecase)
	if err != nil {
		log.Error("h.Usecase.CreateProject", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp, err := h.projectToResp(message)
	if err != nil {
		log.Error("h.projectToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

// UserCredits godoc
// @Summary Create btc project
// @Description Create btc project
// @Tags Project
// @Accept  json
// @Produce  json
// @Security Authorization
// @Param request body request.CreateBTCProjectReq true "Create profile request"
// @Success 200 {object} response.JsonResponse{}
// @Router /project/btc [POST]
func (h *httpDelivery) createBTCProject(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("createBTCProject", r)
	defer h.Tracer.FinishSpan(span, log)

	ctx := r.Context()
	iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)
	walletAddress, ok := iWalletAddress.(string)
	if !ok {
		err := errors.New("Wallet address is incorect")
		log.Error("ctx.Value.Token", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	var reqBody request.CreateBTCProjectReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Error("decoder.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase := &structure.CreateBtcProjectReq{}
	err = copier.Copy(reqUsecase, reqBody)
	if err != nil {
		log.Error("copier.Copy", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	reqUsecase.CreatorAddrr = walletAddress
	log.SetData("reqUsecase", reqUsecase)
	message, err := h.Usecase.CreateBTCProject(span, *reqUsecase)
	if err != nil {
		log.Error("h.Usecase.CreateBTCProject", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp, err := h.projectToResp(message)
	if err != nil {
		log.Error("h.projectToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

// UserCredits godoc
// @Summary update a btc project
// @Description update btc project
// @Tags Project
// @Accept  json
// @Produce  json
// @Security Authorization
// @Param request body request.UpdateBTCProjectReq true "Update project request"
// @Success 200 {object} response.JsonResponse{}
// @Router /project/{contractAddress}/tokens/{projectID} [PUT]
func (h *httpDelivery) updateBTCProject(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("updateBTCProject", r)
	defer h.Tracer.FinishSpan(span, log)

	vars := mux.Vars(r)
	projectID := vars["projectID"]
	span.SetTag("projectID", projectID)

	ctx := r.Context()
	iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)
	walletAddress, ok := iWalletAddress.(string)
	if !ok {
		err := errors.New("Wallet address is incorect")
		log.Error("ctx.Value.Token", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	var reqBody request.UpdateBTCProjectReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Error("decoder.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase := &structure.UpdateBTCProjectReq{}
	err = copier.Copy(reqUsecase, reqBody)
	if err != nil {
		log.Error("copier.Copy", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	reqUsecase.CreatetorAddress = &walletAddress
	reqUsecase.ProjectID = &projectID
	log.SetData("reqUsecase", reqUsecase)
	message, err := h.Usecase.UpdateBTCProject(span, *reqUsecase)
	if err != nil {
		log.Error("h.Usecase.UpdateBTCProject", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp, err := h.projectToResp(message)
	if err != nil {
		log.Error("h.projectToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

// UserCredits godoc
// @Summary get project's detail
// @Description get project's detail
// @Tags Project
// @Accept  json
// @Produce  json
// @Param contractAddress path string true "contract address"
// @Param projectID path string true "token ID"
// @Success 200 {object} response.JsonResponse{}
// @Router /project/{contractAddress}/tokens/{projectID} [GET]
func (h *httpDelivery) projectDetail(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("projectDetail", r)
	defer h.Tracer.FinishSpan(span, log)

	vars := mux.Vars(r)
	contractAddress := vars["contractAddress"]
	span.SetTag("contractAddress", contractAddress)

	projectID := vars["projectID"]
	span.SetTag("projectID", projectID)

	project, err := h.Usecase.GetProjectDetail(span, structure.GetProjectDetailMessageReq{
		ContractAddress: contractAddress,
		ProjectID:       projectID,
	})

	if err != nil {
		log.Error("h.Usecase.GetToken", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp, err := h.projectToResp(project)
	if err != nil {
		log.Error(" h.projectToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	go h.Usecase.CreateViewProjectActivity(project.TokenID)

	log.SetData("resp.project", resp)
	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

// UserCredits godoc
// @Summary get projects
// @Description get projects
// @Tags Project
// @Accept  json
// @Produce  json
// @Param contractAddress query string false "Filter project via contract address"
// @Param name query string false "filter project via name"
// @Param category query string false "filter project via category ids"
// @Param limit query int false "limit"
// @Param page query int false "limit"
// @Param sort query string false "newest, oldest, priority-asc, priority-desc, trending-score"
// @Param cursor query string false "The cursor returned in the previous response (used for getting the next page)."
// @Success 200 {object} response.JsonResponse{}
// @Router /project [GET]
func (h *httpDelivery) getProjects(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("projects", r)
	defer h.Tracer.FinishSpan(span, log)
	var err error
	vars := mux.Vars(r)
	contractAddress := vars["contractAddress"]
	span.SetTag("contractAddress", contractAddress)

	name := r.URL.Query().Get("name")

	categoriesRaw := r.URL.Query().Get("category")

	categoryIds := strings.Split(categoriesRaw, ",")
	if categoriesRaw == "" {
		categoryIds = []string{}
	}

	baseF, err := h.BaseFilters(r)
	if err != nil {
		log.Error("BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	f := structure.FilterProjects{}
	f.BaseFilters = *baseF
	f.Name = &name
	f.CategoryIds = categoryIds
	uProjects, err := h.Usecase.GetProjects(span, f)
	if err != nil {
		log.Error("h.Usecase.GetProjects", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	pResp := []response.ProjectResp{}
	iProjects := uProjects.Result
	projects := iProjects.([]entity.Projects)
	for _, project := range projects {

		p, err := h.projectToResp(&project)
		if err != nil {
			log.Error("copier.Copy", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}

		pResp = append(pResp, *p)
	}

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, h.PaginationResp(uProjects, pResp), "")
}

// UserCredits godoc
// @Summary get minted out projects
// @Description  get minted out projects
// @Tags Project
// @Accept  json
// @Produce  json
// @Param contractAddress query string false "Filter project via contract address"
// @Param limit query int false "limit"
// @Param cursor query string false "The cursor returned in the previous response (used for getting the next page)."
// @Success 200 {object} response.JsonResponse{}
// @Router /project/minted-out [GET]
func (h *httpDelivery) getMintedOutProjects(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("projects", r)
	defer h.Tracer.FinishSpan(span, log)
	var err error
	vars := mux.Vars(r)
	contractAddress := vars["contractAddress"]
	span.SetTag("contractAddress", contractAddress)

	baseF, err := h.BaseFilters(r)
	if err != nil {
		log.Error("BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	f := structure.FilterProjects{}
	f.BaseFilters = *baseF
	uProjects, err := h.Usecase.GetMintedOutProjects(span, f)
	if err != nil {
		log.Error("h.Usecase.GetProjects", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	pResp := []response.ProjectResp{}
	iProjects := uProjects.Result
	projects := iProjects.([]entity.Projects)
	for _, project := range projects {

		p, err := h.projectToResp(&project)
		if err != nil {
			log.Error("copier.Copy", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}

		pResp = append(pResp, *p)
	}

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, h.PaginationResp(uProjects, pResp), "")
}

// UserCredits godoc
// @Summary get the random projects
// @Description get the random projects
// @Tags Project
// @Accept  json
// @Produce  json
// @Success 200 {object} response.JsonResponse{}
// @Router /project/random [GET]
func (h *httpDelivery) getRandomProject(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("getRandomProject", r)
	defer h.Tracer.FinishSpan(span, log)
	var err error

	project, err := h.Usecase.GetRandomProject(span)
	if err != nil {
		log.Error(" h.GetRandomProject", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp, err := h.projectToResp(project)
	if err != nil {
		log.Error(" h.projectToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

func (h *httpDelivery) projectToResp(input *entity.Projects) (*response.ProjectResp, error) {
	resp := &response.ProjectResp{}
	social := make(map[string]string)
	response.CopyEntityToRes(resp, input)
	resp.MintPriceAddr = input.MintTokenAddress
	resp.Limit = input.LimitSupply
	resp.CreatorAddr = input.CreatorAddrr
	resp.Desc = input.Description
	resp.ItemDesc = input.Description
	resp.License = input.License
	resp.Image = input.Thumbnail
	resp.ScriptType = input.ThirdPartyScripts
	resp.NftTokenURI = input.NftTokenUri
	resp.NetworkFee = input.NetworkFee
	resp.Categories = input.Categories
	social["web"] = input.SocialWeb
	social["twitter"] = input.SocialTwitter
	social["discord"] = input.SocialDiscord
	social["medium"] = input.SocialMedium
	social["instagram"] = input.SocialInstagram
	resp.MintingInfo = response.NftMintingDetail{
		Index:        input.MintingInfo.Index,
		IndexReserve: input.MintingInfo.IndexReverse,
	}
	resp.Royalty = input.Royalty
	resp.Reservers = input.Reservers
	resp.CompleteTime = input.CompleteTime
	resp.BlockNumberMinted = input.BlockNumberMinted
	resp.MintedTime = input.MintedTime
	resp.IsFullChain = input.IsFullChain
	resp.IsHidden = input.IsHidden
	resp.CreatorAddrrBTC = input.CreatorAddrrBTC
	resp.TotalImages = len(input.Images)
	resp.Stats = response.ProjectStatResp{
		UniqueOwnerCount:   input.Stats.UniqueOwnerCount,
		TotalTradingVolumn: input.Stats.TotalTradingVolumn,
		FloorPrice:         input.Stats.FloorPrice,
		BestMakeOfferPrice: input.Stats.BestMakeOfferPrice,
		ListedPercent:      input.Stats.ListedPercent,
	}
	resp.Categories = input.Categories
	if input.TraitsStat != nil {
		traitStat := make([]response.TraitStat, 0)
		for _, v := range input.TraitsStat {
			traitValueStats := make([]response.TraitValueStat, 0)
			for _, vv := range v.TraitValuesStat {
				traitValueStats = append(traitValueStats, response.TraitValueStat{
					Value:  vv.Value,
					Rarity: vv.Rarity,
				})
			}
			traitStat = append(traitStat, response.TraitStat{
				TraitName:       v.TraitName,
				TraitValuesStat: traitValueStats,
			})
		}
		resp.TraitStat = traitStat
	}

	profileResp, err := h.profileToResp(&input.CreatorProfile)
	if err == nil {
		resp.CreatorProfile = *profileResp
	}

	resp.MintPriceEth = input.MintPriceEth
	resp.NetworkFeeEth = input.NetworkFeeEth

	return resp, nil
}

// UserCredits godoc
// @Summary get the recent work projects
// @Description  get the recent work projects (posible of minted out)
// @Tags Project
// @Accept  json
// @Produce  json
// @Param contractAddress query string false "Filter project via contract address"
// @Param limit query int false "limit"
// @Param cursor query string false "The cursor returned in the previous response (used for getting the next page)."
// @Success 200 {object} response.JsonResponse{}
// @Router /project/recent-works [GET]
func (h *httpDelivery) getRecentWorksProjects(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("projects", r)
	defer h.Tracer.FinishSpan(span, log)
	var err error
	vars := mux.Vars(r)
	contractAddress := vars["contractAddress"]
	span.SetTag("contractAddress", contractAddress)

	baseF, err := h.BaseFilters(r)
	if err != nil {
		log.Error("BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	f := structure.FilterProjects{}
	f.BaseFilters = *baseF
	uProjects, err := h.Usecase.GetRecentWorksProjects(span, f)
	if err != nil {
		log.Error("h.Usecase.GetProjects", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	pResp := []response.ProjectResp{}
	iProjects := uProjects.Result
	projects := iProjects.([]entity.Projects)
	for _, project := range projects {

		p, err := h.projectToResp(&project)
		if err != nil {
			log.Error("copier.Copy", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}

		pResp = append(pResp, *p)
	}

	h.Response.SetLog(h.Tracer, span)
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
	span, log := h.StartSpan("updateProject", r)
	defer h.Tracer.FinishSpan(span, log)

	vars := mux.Vars(r)
	projectID := vars["projectID"]
	contractAddress := vars["contractAddress"]

	var reqBody request.UpdateProjectReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Error("decoder.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase := &structure.UpdateProjectReq{
		ContracAddress: contractAddress,
		TokenID:        projectID,
		Priority:       reqBody.Priority,
	}

	log.SetData("reqUsecase", reqUsecase)
	log.SetTag("projectID", projectID)
	log.SetTag("contractAddress", contractAddress)

	message, err := h.Usecase.UpdateProject(span, *reqUsecase)
	if err != nil {
		log.Error("h.Usecase.CreateProject", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp, err := h.projectToResp(message)
	if err != nil {
		log.Error("h.projectToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

// UserCredits godoc
// @Summary Upload images for a project
// @Description Upload images for a project
// @Tags Project
// @Content-Type: multipart/form-data
// @Param projectName formData string true "Project's name"
// @Param file formData file true "file"
// @Produce  multipart/form-data
// @Success 200 {object} response.JsonResponse{data=response.FileRes}
// @Router /project/btc/files [POST]
func (h *httpDelivery) UploadProjectFiles(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("httpDelivery.UploadProjectFiles", r)
	defer h.Tracer.FinishSpan(span, log)
	file, err := h.Usecase.UploadProjectFiles(span, r)
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
