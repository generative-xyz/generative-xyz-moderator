package http

import (
	"errors"
	"fmt"
	"net/http"

	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/algolia"
	"rederinghub.io/utils/btc"
)

// UserCredits godoc
// @Summary Search
// @Description Search
// @Tags Search
// @Accept  json
// @Produce  json
// @Param page query string false "page"
// @Param limit query int false "limit"
// @Param search query string false "search"
// @Param type query string false "object type"
// @Success 200 {object} response.JsonResponse{}
// @Router /search [GET]
func (h *httpDelivery) search(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	objType := r.URL.Query().Get("type")
	result := &entity.Pagination{}
	dataResp := []*response.SearchResponse{}
	if search != "" && len(search) < 3 {
		result.Result = dataResp
		h.Response.RespondSuccess(w, http.StatusOK, response.Success, h.PaginationResp(result, result.Result), "")
		return
	}

	if objType == "" {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("Missing object type for search"))
		return
	}

	bf, err := h.BaseAlgoliaFilters(r)
	if err != nil {
		h.Logger.Error("h.Usecase.getCollectionListing.BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	filter := &algolia.AlgoliaFilter{
		SearchStr: search, ObjType: objType,
		Page: int(bf.Page), Limit: int(bf.Limit),
	}
	t, tp := 0, 0

	switch objType {
	case "project":
		var uProjects []entity.Projects
		uProjects, t, tp, err = h.Usecase.AlgoliaSearchProject(filter)
		if err != nil {
			h.Logger.Error("h.Usecase.AlgoliaSearchProject", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
		for _, p := range uProjects {
			result, err := h.Usecase.GetCollectionMarketplaceStats(p.TokenID)
			r, err := h.projectToResp(&p)
			r.BtcFloorPrice = result.FloorPrice
			if err != nil {
				h.Logger.Error("copier.Copy", err.Error(), err)
				h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
				return
			}

			dataResp = append(dataResp, &response.SearchResponse{ObjectType: "project", Project: r})
		}
	case "inscription":
		dataResp, t, tp, err = h.Usecase.AlgoliaSearchInscription(filter)
		if err != nil {
			h.Logger.Error("h.Usecase.AlgoliaSearchInscription", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
	case "artist":
		var users []*response.ArtistResponse
		users, t, tp, err = h.Usecase.AlgoliaSearchArtist(filter)
		if err != nil {
			h.Logger.Error("h.Usecase.AlgoliaSearchArtist", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
		for _, user := range users {
			dataResp = append(dataResp, &response.SearchResponse{ObjectType: "artist", Artist: user})
		}
	case "token":
		var uTokens []entity.TokenUri
		uTokens, t, tp, err = h.Usecase.AlgoliaSearchTokenUri(filter)
		if err != nil {
			h.Logger.Error("h.Usecase.AlgoliaSearchTokenUri", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
		for _, token := range uTokens {
			r, err := h.tokenToResp(&token)
			if err != nil {
				h.Logger.Error("copier.Copy", err.Error(), err)
				h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
				return
			}

			listingInfo, err := h.Usecase.Repo.GetDexBTCListingOrderPendingByInscriptionID(token.TokenID)
			if err == nil && listingInfo.CancelTx == "" {
				r.Buyable = true
				r.PriceBTC = fmt.Sprintf("%v", listingInfo.Amount)
				r.OrderID = listingInfo.UUID
				r.SellVerified = listingInfo.Verified
				if r.SellVerified {
					btcRate, ethRate, err := h.Usecase.GetBTCToETHRate()
					if err != nil {
						h.Logger.Error("GenBuyETHOrder GetBTCToETHRate", err.Error(), err)
					}
					amountBTCRequired := uint64(listingInfo.Amount) + 1000
					amountBTCRequired += btc.EstimateTxFee(3, 2, uint(15)) + btc.EstimateTxFee(1, 2, uint(15))

					amountETH, _, _, err := h.Usecase.ConvertBTCToETHWithPriceEthBtc(fmt.Sprintf("%f", float64(amountBTCRequired)/1e8), btcRate, ethRate)
					if err != nil {
						h.Logger.Error("GenBuyETHOrder convertBTCToETH", err.Error(), err)
					}
					r.PriceETH = amountETH
				}
			}
			dataResp = append(dataResp, &response.SearchResponse{ObjectType: "token", TokenUri: r})
		}
	}

	result.Result = dataResp
	result.Page = int64(filter.Page)
	result.PageSize = int64(filter.Limit)
	result.TotalPage = int64(tp)
	result.Total = int64(t)

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, h.PaginationResp(result, result.Result), "")
}
