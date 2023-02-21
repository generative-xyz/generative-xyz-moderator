package http

import (
	"net/http"
	"strconv"

	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
)

// Artist godoc
// @Summary get list Artist
// @Description get list Artist
// @Tags Artist
// @Accept  json
// @Produce  json
// @Param page query string false "page"
// @Param limit query string false "limit"
// @Success 200 {object} response.JsonResponse{}
// @Router /user/artist [GET]
func (h *httpDelivery) listArtist(w http.ResponseWriter, r *http.Request) {

	// baseF, err := h.BaseFilters(r)
	// if err != nil {
	// 	log.Error("BaseFilters", err.Error(), err)
	// 	h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
	// 	return
	// }
	//
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		h.Logger.Error("parse page param to int", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	
	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		h.Logger.Error("parse limit param to int", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	f := entity.FilteArtist{}
	f.BaseFilters.Limit = int64(limit)
	f.BaseFilters.Page = int64(page)

	result, err := h.Usecase.ListArtist(f)
	if err != nil {
		h.Logger.Error("httpDelivery.listArtist.Usecase.ListArtist", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	pagResp := h.PaginationResp(result, result.Result)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, pagResp, "")
}
