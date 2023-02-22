package http

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
)

// UserCredits godoc
// @Summary Create referral
// @Description Create referral
// @Tags Referral
// @Accept  json
// @Produce  json
// @Param referrerID path string true "referrerID"
// @Security Authorization
// @Success 200 {object} response.JsonResponse{data=bool}
// @Router /referrals/{referrerID} [POST]
func (h *httpDelivery) createReferral(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	iUserID := ctx.Value(utils.SIGNED_USER_ID)
	referreeID, ok := iUserID.(string)

	if !ok {
		err := errors.New( "Token is incorect")
		h.Logger.Error("ctx.Value.Token",  err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	vars := mux.Vars(r)
	referrerID := vars["referrerID"]
	

	err := h.Usecase.CreateReferral(referrerID, referreeID)

	if err != nil {
		h.Logger.Error("h.Usecase.CreateReferral", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, true, "")
}

// UserCredits godoc
// @Summary get referrals
// @Description get referrals
// @Tags Referral
// @Accept  json
// @Produce  json
// @Param referrerID query string false "Filter by referrerID"
// @Param referreeID query string false "filter project referreeID"
// @Param limit query int false "limit"
// @Param page query int false "page"
// @Security Authorization
// @Success 200 {object} response.JsonResponse{}
// @Router /referrals [GET]
func (h *httpDelivery) getReferrals(w http.ResponseWriter, r *http.Request) {
	var err error
	baseF, err := h.BaseFilters(r)
	if err != nil {
		h.Logger.Error("BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	referrerID := r.URL.Query().Get("referrerID")
	referreeID := r.URL.Query().Get("referreeID")

	ctx := r.Context()
	iUserID := ctx.Value(utils.SIGNED_USER_ID)
	referrerID, ok := iUserID.(string)

	if !ok {
		err := errors.New( "Token is incorect")
		h.Logger.Error("ctx.Value.Token",  err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	f := structure.FilterReferrals{}
	f.BaseFilters = *baseF
	f.ReferrerID = &referrerID
	f.ReferreeID = &referreeID
	uReferrals, err := h.Usecase.GetReferrals(f)
	if err != nil {
		h.Logger.Error("h.Usecase.GetReferrals", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	pResp := []response.ReferralResp{}
	iReferrals := uReferrals.Result
	referrals := iReferrals.([]entity.Referral)
	for _, referral := range referrals {

		p, err := h.referralToResp(&referral)
		if err != nil {
			h.Logger.Error("copier.Copy", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}

		pResp = append(pResp, *p)
	}

	
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, h.PaginationResp(uReferrals, pResp), "")
}

func (h *httpDelivery) referralToResp(input *entity.Referral) (*response.ReferralResp, error) {
	resp := response.ReferralResp{}
	resp.ReferrerID = input.ReferrerID
	resp.ReferreeID = input.ReferreeID
	referree, err := h.profileToResp(input.Referree)
	if err != nil {
		return nil, err
	}
	referrer, err := h.profileToResp(input.Referrer)
	if err != nil {
		return nil, err
	}
	resp.Referree = *referree
	resp.Referrer = *referrer
	return &resp, nil
}
