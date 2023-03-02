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
		h.Logger.Error("h.Usecase.getProfileNfts.BaseFilters", err.Error(), err)
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
		dataResp, t, tp, err = h.Usecase.AlgoliaSearchProject(filter)
		if err != nil {
			h.Logger.Error("h.Usecase.AlgoliaSearchProject", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
	case "inscription":
		dataResp, t, tp, err = h.Usecase.AlgoliaSearchInscription(filter)
		if err != nil {
			h.Logger.Error("h.Usecase.AlgoliaSearchInscription", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
	case "artist":
		dataResp, t, tp, err = h.Usecase.AlgoliaSearchArtist(filter)
		if err != nil {
			h.Logger.Error("h.Usecase.AlgoliaSearchArtist", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
	case "token":
		dataResp, t, tp, err = h.Usecase.AlgoliaSearchTokenUri(filter)
		if err != nil {
			h.Logger.Error("h.Usecase.AlgoliaSearchTokenUri", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
	}

	result.Result = dataResp
	result.Page = int64(filter.Page)
	result.PageSize = int64(filter.Limit)
	result.TotalPage = int64(tp)
	result.Total = int64(t)

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, h.PaginationResp(result, result.Result), "")
}
