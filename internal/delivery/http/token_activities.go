package http

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
)

// UserCredits godoc
// @Summary get referrals
// @Description get referrals
// @Tags Tokens
// @Accept  json
// @Produce  json
// @Param inscriptionID path string true "token inscription ID"
// @Param limit query int false "limit"
// @Param page query int false "page"
// @Success 200 {object} response.JsonResponse{}
// @Router /tokens/activities/{inscriptionID} [GET]
func (h *httpDelivery) getTokenActivities(w http.ResponseWriter, r *http.Request) {
	limitInt := 10
	pageInt := 1
	var err error

	limit := r.URL.Query().Get("limit")
	if limit != "" {
		limitInt, err = strconv.Atoi(limit)
		if err != nil {
			h.Logger.ErrorAny("getTokenActivities.FailedParseLimit",  zap.Error(err))
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
	}

	page := r.URL.Query().Get("page")
	if page != "" {
		pageInt, err = strconv.Atoi(page)
		if err != nil {
			h.Logger.ErrorAny("getTokenActivities.FailedParsePage",  zap.Error(err))
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
	}

	vars := mux.Vars(r)
	inscriptionID := vars["inscriptionID"]

	uTokenActivities, err := h.Usecase.GetTokenActivities(int64(pageInt), int64(limitInt), inscriptionID)

	pResp := []response.TokenActivityResp{}
	iTokenActivities := uTokenActivities.Result
	tokenActivities := iTokenActivities.([]entity.TokenActivity)
	for _, tokenActivity := range tokenActivities {

		p, err := h.tokenActivityToResp(&tokenActivity)
		if err != nil {
			h.Logger.Error("copier.Copy", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}

		pResp = append(pResp, *p)
	}

	
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, h.PaginationResp(uTokenActivities, pResp), "")
}

func (h *httpDelivery) tokenActivityToResp(input *entity.TokenActivity) (*response.TokenActivityResp, error) {
	resp := response.TokenActivityResp{}
	resp.Type = int(input.Type)
	resp.Title = input.Title
	resp.UserAAddress = input.UserAAddress
	resp.UserBAddress = input.UserBAddress
	resp.Amount = input.Amount
	resp.Time = input.Time
	if input.UserA != nil {
		userA, err := h.profileToResp(input.UserA)
		if err != nil {
			return nil, err
		}
		resp.UserA = userA
	}
	if input.UserB != nil {
		userB, err := h.profileToResp(input.UserB)
		if err != nil {
			return nil, err
		}
		resp.UserB = userB
	}
	return &resp, nil
}
