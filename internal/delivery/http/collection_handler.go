package http

import (
	"net/http"

	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/algolia"
)

// UserCredits godoc
// @Summary CollectionListing
// @Description get list CollectionListing
// @Tags CollectionListing
// @Accept  json
// @Produce  json
// @Param page query string false "page"
// @Param limit query int false "limit"
// @Success 200 {object} response.JsonResponse{}
// @Router /collections [GET]
func (h *httpDelivery) getCollectionListing(w http.ResponseWriter, r *http.Request) {
	bf, err := h.BaseAlgoliaFilters(r)
	if err != nil {
		h.Logger.Error("h.Usecase.getCollectionListing.BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	filter := &algolia.AlgoliaFilter{
		Page: int(bf.Page), Limit: int(bf.Limit),
	}

	dataResp, t, tp, err := h.Usecase.AlgoliaSearchProjectListing(filter)
	result := &entity.Pagination{}
	result.Result = dataResp
	result.Page = int64(filter.Page)
	result.PageSize = int64(filter.Limit)
	result.TotalPage = int64(tp)
	result.Total = int64(t)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, h.PaginationResp(result, dataResp), "")

	// baseF, err := h.BaseFilters(r)
	// if err != nil {
	// 	h.Logger.Error("BaseFilters", err.Error(), err)
	// 	h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
	// 	return
	// }

	// bytes, _ := json.Marshal(*baseF)
	// key := helpers.GenerateMd5String(string(bytes))
	// cached, err := h.Usecase.Cache.GetData(key)
	// if err == nil && cached != nil {
	// 	resp := response.PaginationResponse{}
	// 	err = json.Unmarshal([]byte(*cached), &resp)
	// 	if err == nil {
	// 		h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
	// 		return
	// 	}
	// }

	// f := structure.FilterProjects{}
	// f.BaseFilters = *baseF
	// hidden := false
	// f.IsHidden = &hidden
	// f.Sort = -1
	// f.SortBy = "stats.trending_score"

	// uProjects, err := h.Usecase.GetProjects(f)
	// if err != nil {
	// 	h.Logger.Error("h.Usecase.GetProjects", err.Error(), err)
	// 	h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
	// 	return
	// }

	// iProjects := uProjects.Result
	// projects := iProjects.([]entity.Projects)
	// listings := make(map[string]*response.ProjectListing)
	// mainW := &sync.WaitGroup{}

	// address := []string{}
	// mapProject := []string{}
	// for _, project := range projects {
	// 	mapProject = append(mapProject, project.ID.Hex())
	// 	if project.CreatorAddrr == "" {
	// 		continue
	// 	}
	// 	address = append(address, project.CreatorAddrr)
	// }

	// users, _ := h.Usecase.Repo.FindUserByAddresses(address)
	// usersMap := make(map[string]entity.Users)
	// for _, u := range users {
	// 	usersMap[u.WalletAddress] = u
	// }

	// for _, project := range projects {
	// 	mainW.Add(1)
	// 	go func(w *sync.WaitGroup, p entity.Projects) {
	// 		defer w.Done()

	// 		projectID := p.TokenID
	// 		floorPrice, err := h.Usecase.Repo.RetrieveFloorPriceOfCollection(projectID)
	// 		if err != nil {
	// 			h.Logger.Error(" h.Usecase.Repo.RetrieveFloorPriceOfCollection", err.Error(), err)
	// 			return
	// 		}

	// 		if floorPrice <= 0 && p.MintingInfo.Index < p.MaxSupply {
	// 			return
	// 		}

	// 		currentListing, err := h.Usecase.Repo.ProjectGetCurrentListingNumber(projectID)
	// 		if err != nil {
	// 			h.Logger.Error(" h.Usecase.Repo.ProjectGetCurrentListingNumber", err.Error(), err)
	// 			return
	// 		}

	// 		volume, err := h.Usecase.Repo.ProjectGetListingVolume(projectID)
	// 		if err != nil {
	// 			h.Logger.Error(" h.Usecase.Repo.ProjectGetListingVolume", err.Error(), err)
	// 			return
	// 		}

	// 		volumeCEX, err := h.Usecase.Repo.ProjectGetCEXVolume(projectID)
	// 		if err != nil {
	// 			h.Logger.Error(" h.Usecase.Repo.ProjectGetListingVolume", err.Error(), err)
	// 			return
	// 		}

	// 		mintVolume, err := h.Usecase.Repo.ProjectGetMintVolume(projectID)
	// 		if err != nil {
	// 			h.Logger.Error(" h.Usecase.Repo.ProjectGetMintVolume", err.Error(), err)
	// 			return
	// 		}

	// 		// tokens, err := h.Usecase.Repo.GetAllTokensByProjectID(projectID)
	// 		// if err != nil {
	// 		// 	h.Logger.Error(" h.Usecase.Repo.GetAllTokensByProjectID", err.Error(), err)
	// 		// 	return
	// 		// }

	// 		// checkers := []string{}
	// 		// for _, t := range tokens {
	// 		// 	checkers = append(checkers, t.OwnerAddr)
	// 		// }

	// 		var result response.ProjectMarketplaceData
	// 		result.FloorPrice = floorPrice
	// 		result.Listed = currentListing
	// 		result.TotalVolume = volume + mintVolume + volumeCEX
	// 		result.MintVolume = mintVolume

	// 		data := &response.ProjectListing{
	// 			// NumberOwners: int64(len(utils.StringUnique(checkers))),
	// 			Project: &response.ProjectInfo{
	// 				Name:            p.Name,
	// 				TokenId:         projectID,
	// 				Thumbnail:       p.Thumbnail,
	// 				ContractAddress: p.ContractAddress,
	// 				CreatorAddress:  p.CreatorAddrr,
	// 				MaxSupply:       p.MaxSupply,
	// 				MintingInfo: response.ProjectMintingInfo{
	// 					Index:        p.MintingInfo.Index,
	// 					IndexReverse: p.MintingInfo.IndexReverse,
	// 				},
	// 				IsMintedOut: p.MintingInfo.Index == p.MaxSupply,
	// 			},
	// 			ProjectMarketplaceData: &result,
	// 		}

	// 		if user, ok := usersMap[p.CreatorAddrr]; ok {
	// 			data.Owner = &response.OwnerInfo{
	// 				DisplayName:             user.DisplayName,
	// 				WalletAddress:           user.WalletAddress,
	// 				WalletAddressPayment:    user.WalletAddressPayment,
	// 				WalletAddressBTC:        user.WalletAddressBTC,
	// 				WalletAddressBTCTaproot: user.WalletAddressBTCTaproot,
	// 				Avatar:                  user.Avatar,
	// 			}
	// 		}
	// 		listings[p.ID.Hex()] = data
	// 	}(mainW, project)
	// }
	// mainW.Wait()
	// data := []*response.ProjectListing{}
	// for _, k := range mapProject {
	// 	if d, ok := listings[k]; ok {
	// 		data = append(data, d)
	// 	}
	// }

	// h.Usecase.Cache.SetDataWithExpireTime(key, h.PaginationResp(uProjects, data), 15*60)
	// h.Response.RespondSuccess(w, http.StatusOK, response.Success, h.PaginationResp(uProjects, data), "")
}
