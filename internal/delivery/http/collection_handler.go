package http

import (
	"net/http"
	"sync"

	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
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
	// bf, err := h.BaseAlgoliaFilters(r)
	// if err != nil {
	// 	h.Logger.Error("h.Usecase.getCollectionListing.BaseFilters", err.Error(), err)
	// 	h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
	// 	return
	// }

	// filter := &algolia.AlgoliaFilter{
	// 	Page: int(bf.Page), Limit: int(bf.Limit),
	// }

	// listings, t, tp, err := h.Usecase.AlgoliaSearchProjectListing(filter)
	// if err != nil {
	// 	h.Logger.Error("h.Usecase.getCollectionListing", err.Error(), err)
	// 	h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
	// 	return
	// }

	// result := &entity.Pagination{}
	// result.Result = listings
	// result.Page = int64(filter.Page)
	// result.PageSize = int64(filter.Limit)
	// result.TotalPage = int64(tp)
	// result.Total = int64(t)

	// h.Response.RespondSuccess(w, http.StatusOK, response.Success, h.PaginationResp(result, result.Result), "")
	//

	baseF, err := h.BaseFilters(r)
	if err != nil {
		h.Logger.Error("BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	f := structure.FilterProjects{}
	f.BaseFilters = *baseF
	hidden := false
	f.IsHidden = &hidden
	f.Sort = -1
	f.SortBy = "stats.trending_score"

	uProjects, err := h.Usecase.GetProjects(f)
	if err != nil {
		h.Logger.Error("h.Usecase.GetProjects", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	iProjects := uProjects.Result
	projects := iProjects.([]entity.Projects)
	listings := []*response.ProjectListing{}
	mainW := &sync.WaitGroup{}

	address := []string{}
	for _, project := range projects {
		if project.CreatorAddrr == "" {
			continue
		}
		address = append(address, project.CreatorAddrr)
	}

	users, _ := h.Usecase.Repo.FindUserByAddresses(address)
	usersMap := make(map[string]entity.Users)
	for _, u := range users {
		usersMap[u.WalletAddress] = u
	}

	for _, project := range projects {
		mainW.Add(1)
		go func() {
			defer mainW.Done()
			projectID := project.TokenID
			floorPrice, err := h.Usecase.Repo.RetrieveFloorPriceOfCollection(projectID)
			if err != nil {
				h.Logger.Error(" h.Usecase.Repo.RetrieveFloorPriceOfCollection", err.Error(), err)
				return
			}

			// if floorPrice <= 0 {
			// 	continue
			// }

			currentListing, err := h.Usecase.Repo.ProjectGetCurrentListingNumber(projectID)
			if err != nil {
				h.Logger.Error(" h.Usecase.Repo.ProjectGetCurrentListingNumber", err.Error(), err)
				h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
				return
			}

			volume, err := h.Usecase.Repo.ProjectGetListingVolume(projectID)
			if err != nil {
				h.Logger.Error(" h.Usecase.Repo.ProjectGetListingVolume", err.Error(), err)
				return
			}

			mintVolume, err := h.Usecase.Repo.ProjectGetMintVolume(projectID)
			var result response.ProjectMarketplaceData

			result.FloorPrice = floorPrice
			result.Listed = currentListing
			result.TotalVolume = volume + mintVolume
			result.MintVolume = mintVolume

			data := &response.ProjectListing{
				Project: &response.ProjectInfo{
					Name:            project.Name,
					TokenId:         projectID,
					Thumbnail:       project.Thumbnail,
					ContractAddress: project.ContractAddress,
					CreatorAddress:  project.CreatorAddrr,
				},
				ProjectMarketplaceData: &result,
			}

			if user, ok := usersMap[project.CreatorAddrr]; ok {
				data.Owner = &response.OwnerInfo{
					DisplayName:             user.DisplayName,
					WalletAddress:           user.WalletAddress,
					WalletAddressPayment:    user.WalletAddressPayment,
					WalletAddressBTC:        user.WalletAddressBTC,
					WalletAddressBTCTaproot: user.WalletAddressBTCTaproot,
					Avatar:                  user.Avatar,
				}
			}

			listings = append(listings, data)
		}()

		mainW.Wait()
	}
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, h.PaginationResp(uProjects, listings), "")
}
