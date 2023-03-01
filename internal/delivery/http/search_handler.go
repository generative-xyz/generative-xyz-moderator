package http

import (
	"errors"
	"net/http"

	"rederinghub.io/internal/delivery/http/response"
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
// @Success 200 {object} response.JsonResponse{}
// @Router /search [GET]
func (h *httpDelivery) search(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
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
		SearchStr: search,
		Page:      int(bf.Page), Limit: int(bf.Limit),
	}

	resp, err := h.Usecase.AlgoliaSearchProject(filter)
	if err != nil {
		h.Logger.Error("h.Usecase.AlgoliaSearchProject", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	dataResp = append(dataResp, resp...)

	resp, err = h.Usecase.AlgoliaSearchInscription(filter)
	if err != nil {
		h.Logger.Error("h.Usecase.AlgoliaSearchInscription", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	dataResp = append(dataResp, resp...)

	resp, err = h.Usecase.AlgoliaSearchArtist(filter)
	if err != nil {
		h.Logger.Error("h.Usecase.AlgoliaSearchArtist", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	dataResp = append(dataResp, resp...)

	resp, err = h.Usecase.AlgoliaSearchTokenUri(filter)
	if err != nil {
		h.Logger.Error("h.Usecase.AlgoliaSearchTokenUri", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	dataResp = append(dataResp, resp...)

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, dataResp, "")
}
