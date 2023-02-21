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
	span, log := h.StartSpan("httpDelivery.listArtist", r)
	defer h.Tracer.FinishSpan(span, log)
	h.Response.SetLog(h.Tracer, span)

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
		log.Error("parse page param to int", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	span.SetTag("page", page)
	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Error("parse limit param to int", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	f := entity.FilteArtist{}
	f.BaseFilters.Limit = int64(limit)
	f.BaseFilters.Page = int64(page)

	result, err := h.Usecase.ListArtist(span, f)
	if err != nil {
		log.Error("httpDelivery.listArtist.Usecase.ListArtist", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	pagResp := h.PaginationResp(result, result.Result)
	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, pagResp, "")
}
