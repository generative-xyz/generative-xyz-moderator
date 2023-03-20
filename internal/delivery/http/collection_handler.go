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
		h.Logger.Error("h.Usecase.getItemListing.BaseFilters", err.Error(), err)
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
// @Router /collections/not-matched-items [GET]
func (h *httpDelivery) getItemListingNotMatched(w http.ResponseWriter, r *http.Request) {
	bf, err := h.BaseFilters(r)
	if err != nil {
		h.Logger.Error("h.Usecase.getItemListing.BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	dataResp, err := h.Usecase.ListItemListingNotMatched(bf)
	result := &entity.Pagination{}
	result.Result = dataResp
	result.Page = int64(bf.Page)
	result.PageSize = int64(bf.Limit)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, h.PaginationResp(result, dataResp), "")
}
