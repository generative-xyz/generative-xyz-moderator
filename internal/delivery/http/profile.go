package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
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

	ctx := r.Context()
	iUserID := ctx.Value(utils.SIGNED_USER_ID)
	userID, ok := iUserID.(string)
	if !ok {
		err := errors.New("Token is incorect")
		h.Logger.Error("ctx.Value.Token", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	profile, err := h.Usecase.UserProfile(userID)
	if err != nil {
		h.Logger.Error("h.Usecase.GenerateMessage(", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Logger.Info("profile", profile)

	resp, err := h.profileToResp(profile)
	if err != nil {
		h.Logger.Error("h.profileToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

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

	ctx := r.Context()
	iToken := ctx.Value(utils.AUTH_TOKEN)
	token, ok := iToken.(string)
	if !ok {
		err := errors.New("Token is incorect")
		h.Logger.Error("ctx.Value.Token", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	isLogedOut, err := h.Usecase.Logout(token)
	if err != nil {
		h.Logger.Error("h.Usecase.GenerateMessage(", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	msg := "Logout successfully"
	if !isLogedOut {
		msg = "Fail! Cannot logout"
	}

	resp := response.LogoutResponse{
		Message: msg,
	}

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

	ctx := r.Context()
	iUserID := ctx.Value(utils.SIGNED_USER_ID)
	userID, ok := iUserID.(string)
	if !ok {
		err := errors.New("Token is incorect")
		h.Logger.Error("ctx.Value.Token", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	var reqBody request.UpdateProfileRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		h.Logger.Error("decoder.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	updateProfile := &structure.UpdateProfile{}
	err = copier.Copy(updateProfile, reqBody)
	if err != nil {
		h.Logger.Error("copier.Copy.structure.UpdateProfile", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Logger.Info("input.Data", updateProfile)
	profile, err := h.Usecase.UpdateUserProfile(userID, *updateProfile)
	if err != nil {
		h.Logger.Error("h.Usecase.GenerateMessage", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Logger.Info("updated.profile", profile)
	resp, err := h.profileToResp(profile)
	if err != nil {
		h.Logger.Error("h.profileToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Logger.Info("respond.profile", profile)

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

// UserCredits godoc
// @Summary get current user's projects
// @Description get current user's projects
// @Tags Profile
// @Accept  json
// @Produce  json
// @Security Authorization
// @Param contractAddress query string false "Filter project via contract address"
// @Param limit query int false "limit"
// @Param cursor query string false "The cursor returned in the previous response (used for getting the next page)."
// @Success 200 {object} response.JsonResponse{}
// @Router /profile/projects [GET]
func (h *httpDelivery) getUserProjects(w http.ResponseWriter, r *http.Request) {
	var err error
	baseF, err := h.BaseFilters(r)
	if err != nil {
		h.Logger.Error("BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	ctx := r.Context()
	iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)
	walletAddress, ok := iWalletAddress.(string)
	if !ok {
		err := errors.New("Wallet address is incorect")
		h.Logger.Error("ctx.Value.Token", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	f := structure.FilterProjects{}
	f.BaseFilters = *baseF
	f.WalletAddress = &walletAddress
	f.SortBy = "created_at"
	f.Sort = 1
	uProjects, err := h.Usecase.GetProjects(f)
	if err != nil {
		h.Logger.Error("h.Usecase.GetProjects", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	pResp := []response.ProjectResp{}
	iProjects := uProjects.Result
	projects := iProjects.([]entity.Projects)
	for _, project := range projects {

		p, err := h.projectToResp(&project)
		if err != nil {
			h.Logger.Error("copier.Copy", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}

		pResp = append(pResp, *p)
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, h.PaginationResp(uProjects, pResp), "")
}

func (h *httpDelivery) profileToResp(profile *entity.Users) (*response.ProfileResponse, error) {
	resp := &response.ProfileResponse{}

	err := response.CopyEntityToRes(resp, profile)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// UserCredits godoc
// @Summary User profile via wallet address
// @Description User profile via wallet address
// @Tags Profile
// @Accept  json
// @Produce  json
// @Param walletAddress path string true "Wallet address"
// @Success 200 {object} response.JsonResponse{data=response.ProfileResponse}
// @Router /profile/wallet/{walletAddress} [GET]
func (h *httpDelivery) profileByWallet(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	walletAddress := vars["walletAddress"]
	profile, err := h.Usecase.GetUserProfileByWalletAddress(walletAddress)
	if err != nil {
		profile, err = h.Usecase.GetUserProfileByBtcAddressTaproot(walletAddress)
		if err != nil {
			h.Logger.Error("h.Usecase.GetUserProfileByWalletAddress(", err.Error(), err)
			profile = &entity.Users{}
		}
	}

	h.Logger.Info("profile", profile)
	resp, err := h.profileToResp(profile)
	if err != nil {
		h.Logger.Error("h.profileToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

// UserCredits godoc
// @Summary User profile via wallet address
// @Description User profile via wallet address
// @Tags Profile
// @Accept  json
// @Produce  json
// @Security Authorization
// @Param WithDrawItemRequest body request.WithDrawItemRequest true "Withdraw request"
// @Success 200 {object} response.JsonResponse{data=response.ProfileResponse}
// @Router /profile/withdraw [POST]
func (h *httpDelivery) withdraw(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)
	walletAddress, ok := iWalletAddress.(string)
	if !ok {
		err := errors.New("WalletAddress is incorect")
		h.Logger.ErrorAny("withdraw.walletAddress", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	var reqBody request.WithDrawItemRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		h.Logger.ErrorAny("withdraw.Decode", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	err = reqBody.SelfValidate()
	if err != nil {
		h.Logger.ErrorAny("withdraw.SelfValidate", zap.Error(err), zap.Any("reqBody", reqBody))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	wdr := &structure.WithDrawItemRequest{}
	err = copier.Copy(wdr, reqBody)
	if err != nil {
		h.Logger.ErrorAny("withdraw.Copy", zap.Error(err), zap.Any("wdr", wdr))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	// TODO: 0x2525
	wd, err := h.Usecase.CreateWithdraw(walletAddress, *wdr)
	if err != nil {
		h.Logger.ErrorAny("withdraw.CreateWithdraw", zap.Error(err), zap.Any("wdr", wdr))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, wd, "")
}
