package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
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
	var reqBody request.CreateProjectReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		h.Logger.Error("decoder.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase := &structure.CreateProjectReq{}
	err = copier.Copy(reqUsecase, reqBody)
	if err != nil {
		h.Logger.Error("copier.Copy", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Logger.Info("reqUsecase", reqUsecase)

	message, err := h.Usecase.CreateProject(*reqUsecase)
	if err != nil {
		h.Logger.Error("h.Usecase.CreateProject", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp, err := h.projectToResp(message)
	if err != nil {
		h.Logger.Error("h.projectToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

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
	ctx := r.Context()
	iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)
	walletAddress, ok := iWalletAddress.(string)
	if !ok {
		err := errors.New("Wallet address is incorect")
		h.Logger.Error("ctx.Value.Token", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	var reqBody request.CreateBTCProjectReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		h.Logger.Error("decoder.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase := &structure.CreateBtcProjectReq{}
	err = copier.Copy(reqUsecase, reqBody)
	if err != nil {
		h.Logger.Error("copier.Copy", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	reqUsecase.CreatorAddrr = walletAddress
	h.Logger.Info("reqUsecase", reqUsecase)
	message, err := h.Usecase.CreateBTCProject(*reqUsecase)
	if err != nil {
		h.Logger.Error("h.Usecase.CreateBTCProject", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp, err := h.projectToResp(message)
	if err != nil {
		h.Logger.Error("h.projectToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

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

	vars := mux.Vars(r)
	projectID := vars["projectID"]

	ctx := r.Context()
	iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)
	walletAddress, ok := iWalletAddress.(string)
	if !ok {
		err := errors.New("Wallet address is incorect")
		h.Logger.Error("ctx.Value.Token", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	var reqBody request.UpdateBTCProjectReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		h.Logger.Error("decoder.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase := &structure.UpdateBTCProjectReq{}
	err = copier.Copy(reqUsecase, reqBody)
	if err != nil {
		h.Logger.Error("copier.Copy", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	reqUsecase.CreatetorAddress = &walletAddress
	reqUsecase.ProjectID = &projectID
	h.Logger.Info("reqUsecase", reqUsecase)
	message, err := h.Usecase.UpdateBTCProject(*reqUsecase)
	if err != nil {
		h.Logger.Error("h.Usecase.UpdateBTCProject", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp, err := h.projectToResp(message)
	if err != nil {
		h.Logger.Error("h.projectToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

// UserCredits godoc
// @Summary Delete BTC project
// @Description Delete BTC projects
// @Tags Project
// @Accept  json
// @Produce  json
// @Security Authorization
// @Param contractAddress path string true "contract adress"
// @Param projectID path string true "projectID adress"
// @Success 200 {object} response.JsonResponse{}
// @Router /project/{contractAddress}/{projectID} [DELETE]
func (h *httpDelivery) deleteBTCProject(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	projectID := vars["projectID"]

	ctx := r.Context()
	iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)
	walletAddress, ok := iWalletAddress.(string)
	if !ok {
		err := errors.New("Wallet address is incorect")
		h.Logger.Error("ctx.Value.Token", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase := &structure.UpdateBTCProjectReq{}
	reqUsecase.CreatetorAddress = &walletAddress
	reqUsecase.ProjectID = &projectID

	h.Logger.Info("reqUsecase", reqUsecase)

	message, err := h.Usecase.DeleteBTCProject(*reqUsecase)
	if err != nil {
		h.Logger.Error("h.Usecase.DeleteProject", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp, err := h.projectToResp(message)
	if err != nil {
		h.Logger.Error("h.projectToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

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

	vars := mux.Vars(r)
	contractAddress := vars["contractAddress"]

	projectID := vars["projectID"]

	project, err := h.Usecase.GetProjectDetail(structure.GetProjectDetailMessageReq{
		ContractAddress: contractAddress,
		ProjectID:       projectID,
	})

	if err != nil {
		h.Logger.Error("h.Usecase.GetToken", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp, err := h.projectToResp(project)
	if err != nil {
		h.Logger.Error(" h.projectToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	go h.Usecase.CreateViewProjectActivity(project.TokenID)

	h.Logger.Info("resp.project", resp)

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

func (h *httpDelivery) projectMarketplaceData(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	// contractAddress := vars["contractAddress"]

	projectID := vars["projectID"]

	currentListing, err := h.Usecase.Repo.ProjectGetCurrentListingNumber(projectID)
	if err != nil {
		h.Logger.Error(" h.Usecase.Repo.ProjectGetCurrentListingNumber", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	floorPrice, err := h.Usecase.Repo.RetrieveFloorPriceOfCollection(projectID)
	if err != nil {
		h.Logger.Error(" h.Usecase.Repo.RetrieveFloorPriceOfCollection", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	var result response.ProjectMarketplaceData

	result.FloorPrice = floorPrice
	result.Listed = currentListing
	result.Volume = 0

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")
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

	name := r.URL.Query().Get("name")
	categoriesRaw := r.URL.Query().Get("category")

	categoryIds := strings.Split(categoriesRaw, ",")
	if categoriesRaw == "" {
		categoryIds = []string{}
	}

	baseF, err := h.BaseFilters(r)
	if err != nil {
		h.Logger.Error("BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	f := structure.FilterProjects{}
	f.BaseFilters = *baseF
	f.Name = &name
	f.CategoryIds = categoryIds

	hidden := false
	f.IsHidden = &hidden
	uProjects, err := h.Usecase.GetProjects(f)
	if err != nil {
		h.Logger.Error("h.Usecase.GetProjects", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	pResp := []response.ProjectResp{}
	iProjects := uProjects.Result
	projects := iProjects.([]entity.Projects)
	for _, project := range projects {

		p, err := h.projectToResp(&project)
		if err != nil {
			h.Logger.Error("copier.Copy", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}

		pResp = append(pResp, *p)
	}

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
	baseF, err := h.BaseFilters(r)
	if err != nil {
		h.Logger.Error("BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	f := structure.FilterProjects{}
	f.BaseFilters = *baseF
	uProjects, err := h.Usecase.GetMintedOutProjects(f)
	if err != nil {
		h.Logger.Error("h.Usecase.GetProjects", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	pResp := []response.ProjectResp{}
	iProjects := uProjects.Result
	projects := iProjects.([]entity.Projects)
	for _, project := range projects {

		p, err := h.projectToResp(&project)
		if err != nil {
			h.Logger.Error("copier.Copy", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}

		pResp = append(pResp, *p)
	}

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
	var err error

	project, err := h.Usecase.GetRandomProject()
	if err != nil {
		h.Logger.Error(" h.GetRandomProject", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp, err := h.projectToResp(project)
	if err != nil {
		h.Logger.Error(" h.projectToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

func (h *httpDelivery) projectToResp(input *entity.Projects) (*response.ProjectResp, error) {
	resp := &response.ProjectResp{}
	social := make(map[string]string)
	response.CopyEntityToRes(resp, input)
	resp.MintPriceAddr = input.MintTokenAddress
	resp.Limit = input.LimitSupply
	resp.CreatorAddrr = input.CreatorAddrr
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
	resp.AnimationHtml = input.AnimationHtml
	resp.MaxFileSize = input.MaxFileSize

	fileExt := ""
	if len(input.Images) > 0 {
		fileExt = input.Images[0]
	} else if len(input.ProcessingImages) > 0 {
		fileExt = input.ProcessingImages[0]
	}
	spew.Dump(fileExt)
	//fileExt := strings.Split(".")

	resp.FileExtension = helpers.FileExtension(fileExt)
	if input.CatureThumbnailDelayTime == nil || *input.CatureThumbnailDelayTime == 0 {
		resp.CaptureThumbnailDelayTime = entity.DEFAULT_CAPTURE_TIME
	} else {
		resp.CaptureThumbnailDelayTime = *input.CatureThumbnailDelayTime
	}
	resp.TotalImages = len(input.Images) + len(input.ProcessingImages)
	resp.HtmlFile = input.HtmlFile
	if resp.HtmlFile == "" {
		if resp.TotalImages > 0 {
			if len(input.Images) > 0 {
				if strings.HasSuffix(input.Images[0], ".html") {
					resp.HtmlFile = input.Images[0]
				}
			} else if len(input.ProcessingImages) > 0 {
				if strings.HasSuffix(input.ProcessingImages[0], ".html") {
					resp.HtmlFile = input.ProcessingImages[0]
				}
			}
		}
	}
	resp.LimitMintPerProcess = input.LimitMintPerProcess
	if resp.LimitMintPerProcess == 0 {
		resp.LimitMintPerProcess = 100
	}

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

	profile, err := h.Usecase.UserProfileByWalletWithCache(input.CreatorAddrr)
	if err == nil && profile != nil && profile.ProfileSocial.TwitterVerified {
		profileResp, err := h.profileToResp(profile)
		if err == nil {
			resp.CreatorProfile = *profileResp
		}
	} else {
		profileResp, err := h.profileToResp(&input.CreatorProfile)
		if err == nil {
			resp.CreatorProfile = *profileResp
		}
	}

	resp.MintPriceEth = input.MintPriceEth
	resp.NetworkFeeEth = input.NetworkFeeEth
	resp.ReportUsers = []*response.ReportProject{}
	for _, r := range input.ReportUsers {
		resp.ReportUsers = append(resp.ReportUsers, &response.ReportProject{ReportUserAddress: r.ReportUserAddress, OriginalLink: r.OriginalLink})
	}

	resp.EditableIsHidden = len(input.ReportUsers) >= h.Config.MaxReportCount

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
	baseF, err := h.BaseFilters(r)
	if err != nil {
		h.Logger.Error("BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	f := structure.FilterProjects{}
	f.BaseFilters = *baseF
	uProjects, err := h.Usecase.GetRecentWorksProjects(f)
	if err != nil {
		h.Logger.Error("h.Usecase.GetProjects", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	pResp := []response.ProjectResp{}
	iProjects := uProjects.Result
	projects := iProjects.([]entity.Projects)
	for _, project := range projects {

		p, err := h.projectToResp(&project)
		if err != nil {
			h.Logger.Error("copier.Copy", err.Error(), err)
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
		h.Logger.Error("decoder.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase := &structure.UpdateProjectReq{
		ContracAddress: contractAddress,
		TokenID:        projectID,
		Priority:       reqBody.Priority,
	}

	h.Logger.Info("reqUsecase", reqUsecase)

	message, err := h.Usecase.UpdateProject(*reqUsecase)
	if err != nil {
		h.Logger.Error("h.Usecase.CreateProject", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp, err := h.projectToResp(message)
	if err != nil {
		h.Logger.Error("h.projectToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

// UserCredits godoc
// @Summary Update project
// @Description Update projects
// @Tags Project
// @Accept  json
// @Produce  json
// @Param projectID path string true "projectID adress"
// @Param request body request.ReportProjectReq true "Report Project request"
// @Success 200 {object} response.JsonResponse{}
// @Router /project/{projectID}/report [POST]
// @Security Authorization
func (h *httpDelivery) reportProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID := vars["projectID"]
	ctx := r.Context()
	iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS).(string)

	var reqBody request.ReportProjectReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		h.Logger.Error("decoder.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	message, err := h.Usecase.ReportProject(projectID, iWalletAddress, reqBody.OriginalLink)
	if err != nil {
		h.Logger.Error("h.Usecase.reportProject", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp, err := h.projectToResp(message)
	if err != nil {
		h.Logger.Error("h.projectToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

// UserCredits godoc
// @Summary Update project's categories
// @Description  Update project's categories
// @Tags Project
// @Accept  json
// @Produce  json
// @Param request body request.UpdateBTCProjectCategoriesReq true "UpdateBTCProjectCategoriesReq"
// @Param contractAddress path string true "contract adress"
// @Param projectID path string true "projectID adress"
// @Success 200 {object} response.JsonResponse{}
// @Router /project/{contractAddress}/{projectID}/categories [PUT]
func (h *httpDelivery) updateBTCProjectcategories(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	projectID := vars["projectID"]

	var reqBody request.UpdateBTCProjectCategoriesReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		h.Logger.Error("decoder.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase := &structure.UpdateBTCProjectReq{
		ProjectID:  &projectID,
		Categories: reqBody.Categories,
	}

	h.Logger.Info("reqUsecase", reqUsecase)

	message, err := h.Usecase.SetCategoriesForBTCProject(*reqUsecase)
	if err != nil {
		h.Logger.Error("h.Usecase.CreateProject", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp, err := h.projectToResp(message)
	if err != nil {
		h.Logger.Error("h.projectToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

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
	file, err := h.Usecase.UploadProjectFiles(r)
	if err != nil {
		h.Logger.ErrorAny("UploadProjectFiles", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp := &response.FileRes{}
	err = response.CopyEntityToRes(resp, file)
	if err != nil {
		h.Logger.ErrorAny("UploadProjectFiles", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

// UserCredits godoc
// @Summary get project's volumn
// @Description get project's volumn
// @Tags Project
// @Accept  json
// @Produce  json
// @Param payType query string false "payType eth|btc"
// @Param contractAddress path string true "contractAddress"
// @Param projectID path string true "token ID"
// @Success 200 {object} response.JsonResponse{}
// @Router /project/{contractAddress}/tokens/{projectID}/volumn [GET]
func (h *httpDelivery) projectVolumn(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	//contractAddress := vars["contractAddress"]
	projectID := vars["projectID"]

	paytype := r.URL.Query().Get("payType")

	v, err := h.Usecase.ProjectVolume(projectID, paytype)
	if err != nil {
		h.Logger.ErrorAny("projectVolumn", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, v, "")
}

// UserCredits godoc
// @Summary get project's random-images
// @Description get project's random-images
// @Tags Project
// @Accept  json
// @Produce  json
// @Param payType query string false "payType eth|btc"
// @Param contractAddress path string true "contractAddress"
// @Param projectID path string true "token ID"
// @Success 200 {object} response.JsonResponse{}
// @Router /project/{contractAddress}/tokens/{projectID}/random-images [GET]
func (h *httpDelivery) projectRandomImages(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	//contractAddress := vars["contractAddress"]
	projectID := vars["projectID"]

	v, err := h.Usecase.ProjectRandomImages(projectID)
	if err != nil {
		h.Logger.ErrorAny("projectVolumn", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, v, "")
}

// UserCredits godoc
// @Summary get upcomming projects
// @Description upcomming get projects
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
// @Router /project/upcomming [GET]
func (h *httpDelivery) getUpcommingProjects(w http.ResponseWriter, r *http.Request) {

	name := r.URL.Query().Get("name")
	categoriesRaw := r.URL.Query().Get("category")

	categoryIds := strings.Split(categoriesRaw, ",")
	if categoriesRaw == "" {
		categoryIds = []string{}
	}

	baseF, err := h.BaseFilters(r)
	if err != nil {
		h.Logger.Error("BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	f := structure.FilterProjects{}
	f.BaseFilters = *baseF
	f.Name = &name
	f.CategoryIds = categoryIds

	hidden := false
	f.IsHidden = &hidden
	uProjects, err := h.Usecase.GetUpcommingProjects(f)
	if err != nil {
		h.Logger.Error("h.Usecase.GetProjects", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	pResp := []response.ProjectResp{}
	iProjects := uProjects.Result
	projects := iProjects.([]entity.Projects)
	for _, project := range projects {

		p, err := h.projectToResp(&project)
		if err != nil {
			h.Logger.Error("copier.Copy", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}

		pResp = append(pResp, *p)
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, h.PaginationResp(uProjects, pResp), "")
}

// UserCredits godoc
// @Summary get project's token-traits
// @Description get project's token-traits
// @Tags Project
// @Accept  json
// @Produce  json
// @Param empty-trait query bool false "only tokens which don't have any trait are exported"
// @Param contractAddress path string true "contractAddress"
// @Param projectID path string true "token ID"
// @Success 200 {object} response.JsonResponse{}
// @Router /project/{contractAddress}/tokens/{projectID}/token-traits [GET]
func (h *httpDelivery) tokenTraits(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	//contractAddress := vars["contractAddress"]
	projectID := vars["projectID"]

	v, err := h.Usecase.ProjectTokenTraits(projectID)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, v, "")
}


// UserCredits godoc
// @Summary upload project's token-traits
// @Description upload project's token-traits
// @Tags Project
// @Content-Type: multipart/form-data
// @Param file formData file true "file"
// @Produce  multipart/form-data
// @Param contractAddress path string true "contractAddress"
// @Param projectID path string true "token ID"
// @Success 200 {object} response.JsonResponse{}
// @Router /project/{contractAddress}/tokens/{projectID}/token-traits [POST]
func (h *httpDelivery) uploadTokenTraits(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID := vars["projectID"]
	v, err := h.Usecase.UploadTokenTraits(projectID, r)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondWithoutContainer(w, http.StatusOK, v)
}
