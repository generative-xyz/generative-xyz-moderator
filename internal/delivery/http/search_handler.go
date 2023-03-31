package http

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/algolia"
	"rederinghub.io/utils/logger"
)

// UserCredits godoc
// @Summary Search Token
// @Description Search Token
// @Tags Search
// @Accept  json
// @Produce  json
// @Param page query string false "page"
// @Param limit query int false "limit"
// @Param search query string false "search"
// @Param projectID path string true "projectID request"
// @Param contractAddress path string true "contractAddress request"
// @Success 200 {object} response.JsonResponse{}
// @Router /project/{contractAddress}/{projectID}/token [GET]
func (h *httpDelivery) searchToken(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	vars := mux.Vars(r)
	projectId := vars["projectID"]

	result := &entity.Pagination{}
	dataResp := []*response.SearchResponse{}

	bf, err := h.BaseAlgoliaFilters(r)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.getCollectionListing.BaseFilters", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	filter := &algolia.AlgoliaFilter{
		SearchStr: search, ObjType: "token",
		Page: int(bf.Page), Limit: int(bf.Limit),
		FilterStr: fmt.Sprintf("projectID:%s", projectId),
	}
	t, tp := 0, 0

	var uTokens []entity.TokenUri
	uTokens, t, tp, err = h.Usecase.AlgoliaSearchTokenUri(filter)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.AlgoliaSearchTokenUri", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	for _, token := range uTokens {
		r, err := h.tokenToResp(&token)
		if err != nil {
			logger.AtLog.Logger.Error("copier.Copy", zap.Error(err))
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
					logger.AtLog.Logger.Error("GenBuyETHOrder GetBTCToETHRate", zap.Error(err))
				}
				amountBTCRequired := uint64(listingInfo.Amount) + 1000
				amountBTCRequired += (amountBTCRequired / 10000) * 15 // + 0,15%
				// amountBTCRequired += btc.EstimateTxFee(3, 2, uint(15)) + btc.EstimateTxFee(1, 2, uint(15))

				amountETH, _, _, err := h.Usecase.ConvertBTCToETHWithPriceEthBtc(fmt.Sprintf("%f", float64(amountBTCRequired)/1e8), btcRate, ethRate)
				if err != nil {
					logger.AtLog.Logger.Error("GenBuyETHOrder convertBTCToETH", zap.Error(err))
				}
				r.PriceETH = amountETH
			}
		}
		dataResp = append(dataResp, &response.SearchResponse{ObjectType: "token", TokenUri: r})
	}

	result.Result = dataResp
	result.Page = int64(filter.Page)
	result.PageSize = int64(filter.Limit)
	result.TotalPage = int64(tp)
	result.Total = int64(t)

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, h.PaginationResp(result, result.Result), "")
}

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
	fromNumberStr := r.URL.Query().Get("fromNumber")
	toNumberStr := r.URL.Query().Get("toNumber")
	fromNumber := 0
	toNumber := 0

	if len(fromNumberStr) > 0 {
		fromNumber, _ = strconv.Atoi(fromNumberStr)
	}
	if len(toNumberStr) > 0 {
		toNumber, _ = strconv.Atoi(toNumberStr)
	}

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
		logger.AtLog.Logger.Error("h.Usecase.getCollectionListing.BaseFilters", zap.Error(err))
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
			logger.AtLog.Logger.Error("h.Usecase.AlgoliaSearchProject", zap.Error(err))
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
		for _, p := range uProjects {
			result, err := h.Usecase.GetCollectionMarketplaceStats(p.TokenID)
			r, err := h.projectToResp(&p)
			r.BtcFloorPrice = result.FloorPrice
			if err != nil {
				logger.AtLog.Logger.Error("copier.Copy", zap.Error(err))
				h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
				return
			}

			dataResp = append(dataResp, &response.SearchResponse{ObjectType: "project", Project: r})
		}
	case "inscription":
		if fromNumber > 0 && toNumber > 0 {
			filter.FromNumber = fromNumber
			filter.ToNumber = toNumber
		}
		dataResp, t, tp, err = h.Usecase.AlgoliaSearchInscription(filter)
		if err != nil {
			logger.AtLog.Logger.Error("h.Usecase.AlgoliaSearchInscription", zap.Error(err))
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
	case "artist":
		var users []*response.ArtistResponse
		users, t, tp, err = h.Usecase.AlgoliaSearchArtist(filter)
		if err != nil {
			logger.AtLog.Logger.Error("h.Usecase.AlgoliaSearchArtist", zap.Error(err))
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
			logger.AtLog.Logger.Error("h.Usecase.AlgoliaSearchTokenUri", zap.Error(err))
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
		for _, token := range uTokens {
			r, err := h.tokenToResp(&token)
			if err != nil {
				logger.AtLog.Logger.Error("copier.Copy", zap.Error(err))
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
						logger.AtLog.Logger.Error("GenBuyETHOrder GetBTCToETHRate", zap.Error(err))
					}
					amountBTCRequired := uint64(listingInfo.Amount) + 1000
					amountBTCRequired += (amountBTCRequired / 10000) * 15 // + 0,15%
					// amountBTCRequired += btc.EstimateTxFee(3, 2, uint(15)) + btc.EstimateTxFee(1, 2, uint(15))

					amountETH, _, _, err := h.Usecase.ConvertBTCToETHWithPriceEthBtc(fmt.Sprintf("%f", float64(amountBTCRequired)/1e8), btcRate, ethRate)
					if err != nil {
						logger.AtLog.Logger.Error("GenBuyETHOrder convertBTCToETH", zap.Error(err))
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
