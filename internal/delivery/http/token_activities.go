package http

import (
	"net/http"

	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
)

// UserCredits godoc
// @Summary get referrals
// @Description get referrals
// @Tags Token Activities
// @Accept  json
// @Produce  json
// @Param inscription_id query string false "token inscription ID"
// @Param project_id query string false "project"
// @Param types query string false "activity types"
// @Param limit query int false "limit"
// @Param page query int false "page"
// @Success 200 {object} response.JsonResponse{}
// @Router /token-activities [GET]
func (h *httpDelivery) getTokenActivities(w http.ResponseWriter, r *http.Request) {
	f := structure.FilterTokenActivities{}
	f.CreateFilter(r)
	bf, err := h.BaseFilters(r)
	if err != nil {
		h.Logger.Error("h.getTokenActivities.BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	f.BaseFilters = *bf
	uTokenActivities, err := h.Usecase.GetTokenActivities(f)
	if err != nil {
		h.Logger.Error("h.Usecase.GetTokenActivities", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

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

	//
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
	if input.TokenInfo != nil {
		token, err := h.tokenToResp(input.TokenInfo)
		if err != nil {
			return nil, err
		}
		resp.TokenInfo = token
	}
	resp.InscriptionID = input.InscriptionID
	resp.ProjectID = input.ProjectID
	return &resp, nil
}
