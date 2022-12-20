package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"

	"rederinghub.io/internal/delivery/http/request"
)

// UserCredits godoc
// @Summary User profile
// @Description User profile
// @Tags Profile
// @Accept  json
// @Produce  json
// @Security Authorization
// @Success 200 {object} response.JsonResponse{data=response.ProfileResponse}
// @Router /profile [GET]
func (h *httpDelivery) profile(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("messages.profile", r)
	defer h.Tracer.FinishSpan(span, log )

	ctx := r.Context()
	iUserID := ctx.Value(utils.SIGNED_USER_ID)
	userID, ok := iUserID.(string)
	if !ok {
		err := errors.New( "Token is incorect")
		log.Error("ctx.Value.Token",  err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	profile, err := h.Usecase.UserProfile(span, userID)
	if err != nil {
		log.Error("h.Usecase.GenerateMessage(", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	log.SetData("profile", profile)
	
	resp := h.profileToResp(profile)
	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

// UserCredits godoc
// @Summary Logout
// @Description Logout
// @Tags Profile
// @Accept  json
// @Produce  json
// @Security Authorization
// @Success 200 {object} response.JsonResponse{data=response.LogoutResponse}
// @Router /profile/logout [POST]
func (h *httpDelivery) logout(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("messages.logout", r)
	defer h.Tracer.FinishSpan(span, log )

	ctx := r.Context()
	iToken := ctx.Value(utils.AUTH_TOKEN)
	token, ok := iToken.(string)
	if !ok {
		err := errors.New( "Token is incorect")
		log.Error("ctx.Value.Token",  err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	isLogedOut, err := h.Usecase.Logout(span, token)
	if err != nil {
		log.Error("h.Usecase.GenerateMessage(", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	msg := "Logout successfully"
	if !isLogedOut {
		msg = "Fail! Cannot logout"
	}

	resp := response.LogoutResponse{
		Message:  msg,
	}

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

// UserCredits godoc
// @Summary Edit User profile
// @Description Edit User profile
// @Tags Profile
// @Accept  json
// @Produce  json
// @Security Authorization
// @Param request body request.UpdateProfileRequest true "Update profile request"
// @Success 200 {object} response.JsonResponse{data=response.ProfileResponse}
// @Router /profile [PUT]
func (h *httpDelivery) updateProfile(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("messages.updateProfile", r)
	defer h.Tracer.FinishSpan(span, log )

	ctx := r.Context()
	iUserID := ctx.Value(utils.SIGNED_USER_ID)
	userID, ok := iUserID.(string)
	if !ok {
		err := errors.New( "Token is incorect")
		log.Error("ctx.Value.Token",  err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	var reqBody request.UpdateProfileRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Error("decoder.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	updateProfile := structure.UpdateProfile{
		Bio: reqBody.Bio,
		DisplayName: reqBody.DisplayName,
	}

	log.SetData("input.Data", updateProfile)
	profile, err := h.Usecase.UpdateUserProfile(span, userID, updateProfile)
	if err != nil {
		log.Error("h.Usecase.GenerateMessage(", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	log.SetData("updated.profile", profile)
	resp := h.profileToResp(profile)
	log.SetData("respond.profile", profile)
	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

func (h *httpDelivery) profileToResp(profile *structure.ProfileResponse) response.ProfileResponse {
	resp := response.ProfileResponse{
		ID: profile.ID,
		WalletAddress: profile.WalletAddress,
		DisplayName: profile.DisplayName,
		Bio: profile.Bio,
	}


	return resp
}