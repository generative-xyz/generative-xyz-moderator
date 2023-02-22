package http

import (
	"net/http"
	"strconv"

	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
)

// UserCredits godoc
// @Summary get users
// @Description get users
// @Tags User
// @Accept  json
// @Produce  json
// @Param search query string false "Filter project via contract address"
// @Param limit query int false "limit"
// @Param page query int false "limit"
// @Success 200 {object} response.JsonResponse{}
// @Router /user [GET]
func (h *httpDelivery) getUsers(w http.ResponseWriter, r *http.Request) {
	searchStr := r.URL.Query().Get("search")
	baseF, err := h.BaseFilters(r)
	if err != nil {
		h.Logger.Error("BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	f := structure.FilterUsers{}
	f.BaseFilters = *baseF
	f.Search = &searchStr

	uUsers, err := h.Usecase.ListUsers(f)
	if err != nil {
		h.Logger.Error("h.Usecase.ListUsers", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, h.PaginationResp(uUsers, uUsers), "")
}

// Artist godoc
// @Summary get list Artist
// @Description get list Artist
// @Tags User
// @Accept  json
// @Produce  json
// @Param page query string false "page"
// @Param limit query string false "limit"
// @Success 200 {object} response.JsonResponse{}
// @Router /user/artist [GET]
func (h *httpDelivery) listArtist(w http.ResponseWriter, r *http.Request) {
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
