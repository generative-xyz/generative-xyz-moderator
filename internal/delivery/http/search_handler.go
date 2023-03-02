package http

import (
	"errors"
	"net/http"

	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/algolia"
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
	if search != "" && len(search) < 3 {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("Term search minimum is 3 characters"))
		return
	}

	bf, err := h.BaseAlgoliaFilters(r)
	if err != nil {
		h.Logger.Error("h.Usecase.getProfileNfts.BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	dataResp := []*response.SearchResponse{}
	filter := &algolia.AlgoliaFilter{
		SearchStr: search, ObjType: objType,
		Page: int(bf.Page), Limit: int(bf.Limit),
	}
	total := 0

	resp, t, err := h.Usecase.AlgoliaSearchProject(filter)
	if err != nil {
		h.Logger.Error("h.Usecase.AlgoliaSearchProject", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	dataResp = append(dataResp, resp...)
	total += t

	resp, t, err = h.Usecase.AlgoliaSearchInscription(filter)
	if err != nil {
		h.Logger.Error("h.Usecase.AlgoliaSearchInscription", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	dataResp = append(dataResp, resp...)
	total += t

	resp, t, err = h.Usecase.AlgoliaSearchArtist(filter)
	if err != nil {
		h.Logger.Error("h.Usecase.AlgoliaSearchArtist", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	dataResp = append(dataResp, resp...)
	total += t

	resp, t, err = h.Usecase.AlgoliaSearchTokenUri(filter)
	if err != nil {
		h.Logger.Error("h.Usecase.AlgoliaSearchTokenUri", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	dataResp = append(dataResp, resp...)
	total += t

	result := &entity.Pagination{}
	result.Result = dataResp
	result.Page = int64(filter.Page)
	result.PageSize = int64(filter.Limit)
	result.Total = int64(total)

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, h.PaginationResp(result, result.Result), "")
}
