package http

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"math/big"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/logger"
)

var specialProjectVolume = map[string]uint64{
	"1002573": 263000000,
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
		logger.AtLog.Logger.Error("ctx.Value.Token", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	var reqBody request.CreateBTCProjectReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		logger.AtLog.Logger.Error("decoder.Decode", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase := &structure.CreateBtcProjectReq{}
	err = copier.Copy(reqUsecase, reqBody)
	if err != nil {
		logger.AtLog.Logger.Error("copier.Copy", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	reqUsecase.CreatorAddrr = walletAddress
	logger.AtLog.Logger.Info("reqUsecase", zap.Any("reqUsecase", reqUsecase))
	message, err := h.Usecase.CreateBTCProject(*reqUsecase)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.CreateBTCProject", zap.Error(err))
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
		err := errors.New("wallet address is incorect")
		logger.AtLog.Logger.Error("updateBTCProject", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	var reqBody request.UpdateBTCProjectReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		logger.AtLog.Logger.Error("updateBTCProject", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase := &structure.UpdateBTCProjectReq{}
	err = copier.Copy(reqUsecase, reqBody)
	if err != nil {
		logger.AtLog.Logger.Error("updateBTCProject", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase.CreatetorAddress = &walletAddress
	reqUsecase.ProjectID = &projectID
	logger.AtLog.Logger.Info("reqUsecase", zap.Any("reqUsecase", reqUsecase))
	message, err := h.Usecase.UpdateBTCProject(*reqUsecase)
	if err != nil {
		logger.AtLog.Logger.Error("updateBTCProject", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp, err := h.projectToResp(message)
	if err != nil {
		logger.AtLog.Logger.Error("updateBTCProject", zap.Error(err))
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
		logger.AtLog.Logger.Error("ctx.Value.Token", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase := &structure.UpdateBTCProjectReq{}
	reqUsecase.CreatetorAddress = &walletAddress
	reqUsecase.ProjectID = &projectID

	logger.AtLog.Logger.Info("reqUsecase", zap.Any("reqUsecase", reqUsecase))

	message, err := h.Usecase.DeleteBTCProject(*reqUsecase)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.DeleteProject", zap.Error(err))
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
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, vars map[string]string) (interface{}, error) {
			contractAddress := vars["contractAddress"]
			projectID := vars["projectID"]
			userAddress := r.URL.Query().Get("userAddress")
			project, err := h.Usecase.GetProjectDetailWithFeeInfo(structure.GetProjectDetailMessageReq{
				ContractAddress:            contractAddress,
				ProjectID:                  projectID,
				UserAddressToCheckDiscount: userAddress,
			})
			if err != nil {
				logger.AtLog.Logger.Error("GetProjectDetailWithFeeInfo failed", zap.Error(err))
				return nil, err
			}

			if r.URL.Query().Get("isEdit") == "" {
				project.CurrentLoginUserID = vars[utils.SIGNED_USER_ID]
			}

			resp, err := h.projectToResp(project)
			if err != nil {
				return nil, err
			}
			resp.IsReviewing = h.Usecase.IsProjectReviewing(ctx, project.ID.Hex())
			if project.CreatorAddrr != "" && strings.EqualFold(vars[utils.SIGNED_WALLET_ADDRESS], project.CreatorAddrr) {
				if resp.IsHidden {
					daoProject, exists := h.Usecase.CheckDAOProjectAvailableByUser(ctx, project.CreatorAddrr, project.ID)
					resp.CanCreateProposal = !exists
					if daoProject != nil {
						resp.ProposalSeqId = &daoProject.SeqId
					}
				}
			} else {
				if resp.IsHidden {
					if daoProject, err := h.Usecase.GetLastDAOProjectByProjectId(ctx, project.ID); err == nil {
						resp.ProposalSeqId = &daoProject.SeqId
					}
				}
			}
			go h.Usecase.CreateViewProjectActivity(project.TokenID)
			return resp, nil
		},
	).ServeHTTP(w, r)
}

func (h *httpDelivery) projectMarketplaceData(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	// contractAddress := vars["contractAddress"]

	projectID := vars["projectID"]
	ignoreCache := r.URL.Query().Get("refresh") == "true"
	cached, err := h.Cache.GetData(helpers.GenerateMKPDataKey(projectID))
	if ignoreCache || err != nil || cached == nil {
		currentListing, err := h.Usecase.Repo.ProjectGetCurrentListingNumber(projectID)
		if err != nil {
			logger.AtLog.Logger.Error(" h.Usecase.Repo.ProjectGetCurrentListingNumber", zap.Error(err))
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}

		floorPrice, err := h.Usecase.Repo.RetrieveFloorPriceOfCollection(projectID)
		if err != nil {
			logger.AtLog.Logger.Error(" h.Usecase.Repo.RetrieveFloorPriceOfCollection", zap.Error(err))
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}

		volume, err := h.Usecase.Repo.ProjectGetListingVolume(projectID)
		if err != nil {
			logger.AtLog.Logger.Error(" h.Usecase.Repo.ProjectGetListingVolume", zap.Error(err))
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}

		volumeCEX, err := h.Usecase.Repo.ProjectGetCEXVolume(projectID)
		if err != nil {
			logger.AtLog.Logger.Error(" h.Usecase.Repo.ProjectGetListingVolume", zap.Error(err))
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
		// projectInfo, err := h.Usecase.Repo.FindProjectByTokenID(projectID)
		// if err != nil {
		// 	logger.AtLog.Logger.Error(" h.Usecase.Repo.ProjectGetListingVolume", zap.Error(err))
		// 	h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		// 	return
		// }
		// minted := projectInfo.MintingInfo.Index
		// mintPrice, _ := strconv.Atoi(projectInfo.MintPrice)
		// // mintFee, _ := strconv.Atoi(projectInfo.NetworkFee)
		// mintVolume := uint64(minted) * uint64(mintPrice)

		mintVolume, err := h.Usecase.Repo.ProjectGetMintVolume(projectID)
		if err != nil {
			logger.AtLog.Logger.Error(" h.Usecase.Repo.ProjectGetListingVolume", zap.Error(err))
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
		var result response.ProjectMarketplaceData

		result.FloorPrice = floorPrice
		result.Listed = currentListing
		result.TotalVolume = volume + volumeCEX
		result.MintVolume = mintVolume
		result.CEXVolume = volumeCEX

		if additionalAmount, has := specialProjectVolume[projectID]; has {
			result.TotalVolume += additionalAmount
		}

		h.Cache.SetDataWithExpireTime(helpers.GenerateMKPDataKey(projectID), result, 600) // 10 min
		h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")
		return
	}
	var result response.ProjectMarketplaceData
	bytes := []byte(*cached)

	err = json.Unmarshal(bytes, &result)
	if err != nil {
		logger.AtLog.Logger.Error(" h.Usecase.Repo.ProjectGetListingVolume", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

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
	creatorAddress := r.URL.Query().Get("creatorAddress")

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
	f.CategoryIds = categoryIds
	f.WalletAddress = &creatorAddress

	hidden := false
	f.IsHidden = &hidden
	uProjects, err := h.Usecase.GetProjects(f)
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

	projectFloorPriceMap, err := h.Usecase.GetProjectsFloorPrice(projectToGetFloorPrice)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.GetProjectsFloorPrice", zap.Error(err))
	}
	for _, project := range projects {
		p, err := h.projectToResp(&project)
		if err != nil {
			logger.AtLog.Logger.Error("copier.Copy", zap.Error(err))
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
		floorPrice, ok := projectFloorPriceMap[p.TokenID]
		if ok {
			p.BtcFloorPrice = floorPrice
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
		logger.AtLog.Logger.Error("BaseFilters", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	f := structure.FilterProjects{}
	f.BaseFilters = *baseF
	uProjects, err := h.Usecase.GetMintedOutProjects(f)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.GetProjects", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	pResp := []response.ProjectResp{}
	iProjects := uProjects.Result
	projects := iProjects.([]entity.Projects)
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
		logger.AtLog.Logger.Error(" h.GetRandomProject", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp, err := h.projectToResp(project)
	if err != nil {
		logger.AtLog.Logger.Error(" h.projectToResp", zap.Error(err))
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
	resp.IsSynced = input.IsSynced
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
	resp.TxHash = input.TxHash
	resp.TxHex = input.TxHex

	fileExt := ""
	if len(input.Images) > 0 {
		fileExt = input.Images[0]
	} else if len(input.ProcessingImages) > 0 {
		fileExt = input.ProcessingImages[0]
	}
	//spew.Dump(fileExt)
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
				r := vv.RarityF
				if r == 0 {
					r = float64(vv.Rarity)
				}

				traitValueStats = append(traitValueStats, response.TraitValueStat{
					Value:  vv.Value,
					Rarity: math.Floor(r*100) / 100,
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
	if err == nil && profile != nil {
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

	// check is generative code
	resp.IsGenerative = true
	if input.Source == "" {
		// -> from generative.xyz
		if resp.TotalImages != 0 {
			if len(input.Images) > 0 {
				if !strings.HasSuffix(input.Images[0], ".html") {
					resp.IsGenerative = false
				}
			}
			if resp.IsGenerative && len(input.ProcessingImages) > 0 {
				if !strings.HasSuffix(input.ProcessingImages[0], ".html") {
					resp.IsGenerative = false
				}
			}
		}
	} else {
		// crawler
		resp.IsGenerative = false
	}

	if input.IsSupportGMHolder && input.CurrentLoginUserID != "" {
		minimumGMSupport, ok := new(big.Int).SetString(input.MinimumGMSupport, 10)
		if ok && minimumGMSupport.Cmp(new(big.Int).SetUint64(0)) > 0 {
			currentUser, err := h.Usecase.UserProfile(input.CurrentLoginUserID)
			if err == nil {
				gmBalance, err := h.Usecase.GetGMBalance(currentUser.WalletAddress)
				if err == nil {
					resp.CurrentUserBalanceGM = gmBalance.String()
					if gmBalance.Cmp(minimumGMSupport) >= 0 {
						resp.MintPrice = "0"
						resp.MintPriceEth = "0"
					}
				}
			}
		}
	}

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
		logger.AtLog.Logger.Error("BaseFilters", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	f := structure.FilterProjects{}
	f.BaseFilters = *baseF
	uProjects, err := h.Usecase.GetRecentWorksProjects(f)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.GetProjects", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	pResp := []response.ProjectResp{}
	iProjects := uProjects.Result
	projects := iProjects.([]entity.Projects)
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
// @Summary Update project's hash
// @Description Update project's hash via txHasg
// @Tags Project
// @Accept  json
// @Produce  json
// @Param request body structure.UpdateProjectHash true "Request body"
// @Param contractAddress path string true "contract adress"
// @Param txHash path string true "txHash adress"
// @Success 200 {object} response.JsonResponse{}
// @Router /project/{contractAddress}/tx-hash/{txHash} [PUT]
func (h *httpDelivery) updateProjectHash(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	txHash := vars["txHash"]
	contractAddress := vars["contractAddress"]

	var reqBody structure.UpdateProjectHash
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		logger.AtLog.Logger.Error("decoder.Decode", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	reqBody.TxHash = &txHash
	reqBody.ContractAddress = &contractAddress
	message, err := h.Usecase.UpdateProjectHash(reqBody)
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
		logger.AtLog.Logger.Error("decoder.Decode", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	message, err := h.Usecase.ReportProject(projectID, iWalletAddress, reqBody.OriginalLink)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.reportProject", zap.Error(err))
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
		logger.AtLog.Logger.Error("decoder.Decode", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase := &structure.UpdateBTCProjectReq{
		ProjectID:  &projectID,
		Categories: reqBody.Categories,
	}

	logger.AtLog.Logger.Info("reqUsecase", zap.Any("reqUsecase", reqUsecase))

	message, err := h.Usecase.SetCategoriesForBTCProject(*reqUsecase)
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
		logger.AtLog.Logger.Error("UploadProjectFiles", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp := &response.FileRes{}
	err = response.CopyEntityToRes(resp, file)
	if err != nil {
		logger.AtLog.Logger.Error("UploadProjectFiles", zap.Error(err))
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
		logger.AtLog.Logger.Error("projectVolumn", zap.Error(err))
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
		logger.AtLog.Logger.Error("projectVolumn", zap.Error(err))
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
		logger.AtLog.Logger.Error("BaseFilters", zap.Error(err))
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
		logger.AtLog.Logger.Error("h.Usecase.GetProjects", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	pResp := []response.ProjectResp{}
	iProjects := uProjects.Result
	projects := iProjects.([]entity.Projects)
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

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, v, "")
}

func (h *httpDelivery) triggerPubsubTokenThumbnail(w http.ResponseWriter, r *http.Request) {
	tokenId := r.URL.Query().Get("tokenId")
	ctx := r.Context()
	iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)
	userWalletAddr, ok := iWalletAddress.(string)
	if !ok || strings.ToLower(userWalletAddr) != strings.ToLower("0x668ea0470396138acd0B9cCf6FBdb8a845B717B0") {
		err := errors.New("wallet address is incorect")
		logger.AtLog.Logger.Error("ctx.Value.Token", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	if len(tokenId) == 0 {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("token empty"))
		return
	}
	token, err := h.Usecase.TriggerPubsubTokenThumbnail(tokenId)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("token empty"))
		return
	}
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, token.TokenID, "")
}

// UserCredits godoc
// @Summary Create project's allow list
// @Description Create project's allow list
// @Tags Project
// @Accept  json
// @Produce  json
// @Security Authorization
// @Param contractAddress path string true "contractAddress request"
// @Param projectID path string true "projectID request"
// @Success 200 {object} response.JsonResponse{}
// @Router /project/{contractAddress}/{projectID}/allow-list [POST]
func (h *httpDelivery) createProjectAllowList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID := vars["projectID"]
	var err error
	var walletAddress *string
	var resp *entity.ProjectAllowList

	defer func() {
		if err != nil {
			logger.AtLog.Logger.Error("createProjectAllowList", zap.String("projectID", projectID), zap.Any("walletAddress", walletAddress), zap.Error(err))
		}
		logger.AtLog.Logger.Info("createProjectAllowList", zap.String("projectID", projectID), zap.Any("walletAddress", walletAddress), zap.Any("resp", zap.Any("resp)", resp)))
	}()

	ctx := r.Context()
	iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)
	wa, ok := iWalletAddress.(string)
	if !ok {
		err := errors.New("walletAddress is incorect")
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	walletAddress = &wa

	reqUsecase := &structure.CreateProjectAllowListReq{
		ProjectID:         &projectID,
		UserWalletAddress: walletAddress,
	}

	resp, err = h.Usecase.CreateProjectAllowList(*reqUsecase)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

// UserCredits godoc
// @Summary Create project's allow list
// @Description Create project's allow list
// @Tags Project
// @Accept  json
// @Produce  json
// @Security Authorization
// @Param contractAddress path string true "contractAddress request"
// @Param projectID path string true "projectID request"
// @Param walletAddress path string true "walletAddress request"
// @Success 200 {object} response.JsonResponse{}
// @Router /project/{contractAddress}/{projectID}/{walletAddress}/allow-list-gm [POST]
func (h *httpDelivery) createProjectAllowListGM(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID := vars["projectID"]
	if projectID != "999998" {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("invalid project"))
		return
	}
	walletAddress := vars["walletAddress"]
	walletAddress = strings.ToLower(walletAddress)
	var err error
	var resp *entity.ProjectAllowList

	defer func() {
		if err != nil {
			logger.AtLog.Logger.Error("createProjectAllowList", zap.String("projectID", projectID), zap.Any("walletAddress", walletAddress), zap.Error(err))
		}
		logger.AtLog.Logger.Info("createProjectAllowList", zap.String("projectID", projectID), zap.Any("walletAddress", walletAddress), zap.Any("resp", zap.Any("resp)", resp)))
	}()

	reqUsecase := &structure.CreateProjectAllowListReq{
		ProjectID:         &projectID,
		UserWalletAddress: &walletAddress,
	}

	resp, err = h.Usecase.CreateProjectAllowList(*reqUsecase)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

// UserCredits godoc
// @Summary Check project's allow list
// @Description Check project's allow list
// @Tags Project
// @Accept  json
// @Produce  json
// @Security Authorization
// @Param contractAddress path string true "contractAddress request"
// @Param projectID path string true "projectID request"
// @Success 200 {object} response.JsonResponse{}
// @Router /project/{contractAddress}/{projectID}/allow-list [GET]
func (h *httpDelivery) getProjectAllowList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID := vars["projectID"]
	var err error
	var walletAddress *string
	var resp *entity.ProjectAllowList

	defer func() {
		if err != nil {
			logger.AtLog.Logger.Error("createProjectAllowList", zap.String("projectID", projectID), zap.Any("walletAddress", walletAddress), zap.Error(err))
		}
		logger.AtLog.Logger.Info("createProjectAllowList", zap.String("projectID", projectID), zap.Any("walletAddress", walletAddress), zap.Any("resp", zap.Any("resp)", resp)))
	}()

	ctx := r.Context()
	iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)
	wa, ok := iWalletAddress.(string)
	if !ok {
		err := errors.New("walletAddress is incorect")
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	walletAddress = &wa

	reqUsecase := &structure.CreateProjectAllowListReq{
		ProjectID:         &projectID,
		UserWalletAddress: walletAddress,
	}

	existed, allowedBy := h.Usecase.CheckExistedProjectAllowList(*reqUsecase)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, response.ExistedInAllowList{
		Existed:   existed,
		AllowedBy: allowedBy,
	}, "")
}

func (h *httpDelivery) getCountingAllowList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID := vars["projectID"]
	public, al, err := h.Usecase.CountingProjectAllowList(projectID)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
	} else {
		h.Response.RespondSuccess(w, http.StatusOK, response.Success, response.CountingAllowList{
			Public:    public,
			AllowList: al,
		}, "")
	}
}

// AnalyticsTokenUriOwner godoc
// @Summary Calculate Token's Onwers by project
// @Description  Calculate Token's Onwers by project
// @Tags Project
// @Accept  json
// @Produce  json
// @Param page query string false "page"
// @Param limit query int false "limit"
// @Param search query string false "search"
// @Param projectID path string true "projectID request"
// @Param contractAddress path string true "contractAddress request"
// @Success 200 {object} response.JsonResponse{}
// @Router /project/{contractAddress}/{projectID}/token-onwers/analytics [GET]
func (h *httpDelivery) AnalyticsTokenUriOwner(w http.ResponseWriter, r *http.Request) {
	f := structure.FilterTokens{}
	err := f.CreateFilter(r)
	if err != nil {
		logger.AtLog.Logger.Error("AnalyticsTokenUriOwner", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	vars := mux.Vars(r)
	projectID := vars["projectID"]
	f.GenNFTAddr = &projectID
	bf, err := h.BaseFilters(r)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.getProfileNfts.BaseFilters", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	f.BaseFilters = *bf
	resp, err := h.Usecase.AnalyticsTokenUriOwner(f)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.getProfileNfts.getTokens", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	//
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}
