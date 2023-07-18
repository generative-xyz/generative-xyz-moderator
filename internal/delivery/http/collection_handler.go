package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"go.uber.org/zap"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/algolia"
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/logger"
)

// UserCredits godoc
// @Summary CollectionListing
// @Description get list CollectionListing
// @Tags CollectionListing
// @Accept  json
// @Produce  json
// @Param page query string false "page"
// @Param limit query int false "limit"
// @Param number_from query int false "number_from"
// @Param number_to query int false "number_to"
// @Success 200 {object} response.JsonResponse{}
// @Router /collections/sub-collection-items [GET]
func (h *httpDelivery) getSubCollectionItemListing(w http.ResponseWriter, r *http.Request) {
	bf, err := h.BaseFilters(r)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.getItemListing.BaseFilters", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	numberFromStr := r.URL.Query().Get("number_from")
	numberToStr := r.URL.Query().Get("number_to")
	numberFrom := 0
	numberTo := 0

	if len(numberFromStr) > 0 {
		numberFrom, _ = strconv.Atoi(numberFromStr)
	}
	if len(numberToStr) > 0 {
		numberTo, _ = strconv.Atoi(numberToStr)
	}

	keyByte, _ := json.Marshal(bf)
	key := fmt.Sprintf("sub_collection_%s_%d_%d", helpers.GenerateMd5String(string(keyByte)), numberFrom, numberTo)
	str, err := h.Cache.GetData(key)
	if err == nil {
		data := []*entity.ItemListing{}
		if err := json.Unmarshal([]byte(*str), &data); err == nil {
			result := &entity.Pagination{}
			result.Result = data
			result.Page = int64(bf.Page)
			result.PageSize = int64(bf.Limit)
			h.Response.RespondSuccess(w, http.StatusOK, response.Success, h.PaginationResp(result, data), "")
			return
		}

	}

	dataResp, err := h.Usecase.SubCollectionItem(bf, numberFrom, numberTo)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	result := &entity.Pagination{}
	result.Result = dataResp
	result.Page = int64(bf.Page)
	result.PageSize = int64(bf.Limit)
	if len(dataResp) > 0 {
		h.Cache.SetDataWithExpireTime(key, dataResp, 15*60)
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, h.PaginationResp(result, dataResp), "")
}

// UserCredits godoc
// @Summary Collection's chart
// @Description get list tokens og a collection and respond data for chart
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
		logger.AtLog.Logger.Error("h.Usecase.getCollectionListing.BaseFilters", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	filter := &algolia.AlgoliaFilter{
		Page: int(bf.Page), Limit: int(bf.Limit),
	}

	dataResp, t, tp, err := h.Usecase.DBProjectProtabAPIFormatData(filter)
	result := &entity.Pagination{}
	result.Result = dataResp
	result.Page = int64(filter.Page)
	result.PageSize = int64(filter.Limit)
	result.TotalPage = int64(tp)
	result.Total = int64(t)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, h.PaginationResp(result, dataResp), "")
}

// UserCredits godoc
// @Summary CollectionListing
// @Description get list CollectionListing
// @Tags CollectionListing
// @Accept  json
// @Produce  json
// @Param page query string false "page"
// @Param limit query int false "limit"
// @Success 200 {object} response.JsonResponse{}
// @Router /collections/items [GET]
func (h *httpDelivery) getItemListing(w http.ResponseWriter, r *http.Request) {
	bf, err := h.BaseFilters(r)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.getItemListing.BaseFilters", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	dataResp, err := h.Usecase.ListItemListing(bf)
	result := &entity.Pagination{}
	result.Result = dataResp
	result.Page = int64(bf.Page)
	result.PageSize = int64(bf.Limit)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, h.PaginationResp(result, dataResp), "")
}

// UserCredits godoc
// @Summary CollectionListing
// @Description get list CollectionListing
// @Tags CollectionListing
// @Accept  json
// @Produce  json
// @Param page query string false "page"
// @Param limit query int false "limit"
// @Success 200 {object} response.JsonResponse{}
// @Router /collections/on-sale-items [GET]
func (h *httpDelivery) getItemListingOnSale(w http.ResponseWriter, r *http.Request) {
	bf, err := h.BaseFilters(r)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.getItemListingOnSale.BaseFilters", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	dataResp, err := h.Usecase.ListItemListingOnSale(bf)
	result := &entity.Pagination{}
	result.Result = dataResp
	result.Page = int64(bf.Page)
	result.PageSize = int64(bf.Limit)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, h.PaginationResp(result, dataResp), "")
}
