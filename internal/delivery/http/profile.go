package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	copierInternal "rederinghub.io/utils/copier"
	"rederinghub.io/utils/logger"
)

// @Summary User profile
// @Description User profile
// @Tags Profile
// @Accept  json
// @Produce  json
// @Security Authorization
// @Success 200 {object} response.JsonResponse{data=response.ProfileResponse}
// @Router /profile [GET]
func (h *httpDelivery) profile(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, vars map[string]string) (interface{}, error) {
			userID := vars[utils.SIGNED_USER_ID]
			if userID == "" {
				return nil, errors.New("Token is empty")
			}
			profile, err := h.Usecase.UserProfile(userID)
			if err != nil {
				return nil, err
			}
			resp := &response.ProfileResponse{}
			if err := copierInternal.Copy(resp, profile); err != nil {
				return nil, err
			}
			if !profile.ProfileSocial.TwitterVerified {
				daoArtist, exists := h.Usecase.CheckDAOArtistAvailableByUser(ctx, profile.WalletAddress)
				resp.CanCreateProposal = !exists
				if daoArtist != nil {
					proposal := &response.DaoArtist{}
					if err := copierInternal.Copy(proposal, daoArtist); err == nil {
						resp.Proposal = proposal
					}
				}
			}
			return resp, nil
		},
	).ServeHTTP(w, r)
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
		logger.AtLog.Logger.Error("ctx.Value.Token", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	isLogedOut, err := h.Usecase.Logout(token)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.GenerateMessage(", zap.Error(err))
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
		logger.AtLog.Logger.Error("ctx.Value.Token", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	var reqBody request.UpdateProfileRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		logger.AtLog.Logger.Error("decoder.Decode", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	updateProfile := &structure.UpdateProfile{}
	err = copier.Copy(updateProfile, reqBody)
	if err != nil {
		logger.AtLog.Logger.Error("copier.Copy.structure.UpdateProfile", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	logger.AtLog.Logger.Info("input.Data", zap.Any("updateProfile", updateProfile))
	profile, err := h.Usecase.UpdateUserProfile(userID, *updateProfile)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.GenerateMessage", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	logger.AtLog.Logger.Info("updated.profile", zap.Any("profile", profile))
	resp, err := h.profileToResp(profile)
	if err != nil {
		logger.AtLog.Logger.Error("h.profileToResp", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	logger.AtLog.Logger.Info("respond.profile", zap.Any("profile", profile))

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

// @Summary Get current user's projects
// @Description Get current user's projects
// @Tags Profile
// @Accept json
// @Produce json
// @Param limit query int false "limit"
// @Param cursor query string false "The cursor returned in the previous response (used for getting the next page)."
// @Success 200 {array} response.ProjectResp{}
// @Router /profile/projects [GET]
// @Security ApiKeyAuth
func (h *httpDelivery) getUserProjects(w http.ResponseWriter, r *http.Request) {
	var err error
	baseF, err := h.BaseFilters(r)
	if err != nil {
		logger.AtLog.Logger.Error("BaseFilters", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	ctx := r.Context()
	iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)
	walletAddress, ok := iWalletAddress.(string)
	if !ok {
		err := errors.New("Wallet address is incorect")
		logger.AtLog.Logger.Error("ctx.Value.Token", zap.Error(err))
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
		logger.AtLog.Logger.Error("h.Usecase.GetProjects", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	pResp := []response.ProjectResp{}
	iProjects := uProjects.Result
	projects := iProjects.([]entity.Projects)
	for _, project := range projects {

		p, err := h.projectToResp(&project)
		if err != nil {
			logger.AtLog.Logger.Error("copier.Copy", zap.Error(err))
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

// @Summary User profile via wallet address
// @Description User profile via wallet address
// @Tags Profile
// @Accept json
// @Produce json
// @Param walletAddress path string true "Wallet address"
// @Success 200 {object} response.JsonResponse{data=response.ProfileResponse}
// @Router /profile/wallet/{walletAddress} [GET]
func (h *httpDelivery) profileByWallet(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, vars map[string]string) (interface{}, error) {
			walletAddress := vars["walletAddress"]
			profile, err := h.Usecase.GetUserProfileByWalletAddress(walletAddress)
			if err != nil {
				profile, err = h.Usecase.GetUserProfileByBtcAddressTaproot(walletAddress)
				if err != nil {
					logger.AtLog.Logger.Error("GetUserProfileByWalletAddress failed", zap.Error(err))
					profile = &entity.Users{}
				}
			}
			resp := &response.ProfileResponse{}
			if err := copierInternal.Copy(resp, profile); err != nil {
				return nil, err
			}
			if !profile.ProfileSocial.TwitterVerified {
				daoArtist, exists := h.Usecase.CheckDAOArtistAvailableByUser(ctx, profile.WalletAddress)
				resp.CanCreateProposal = !exists && strings.EqualFold(profile.WalletAddress, vars[utils.SIGNED_WALLET_ADDRESS])
				if daoArtist != nil {
					proposal := &response.DaoArtist{}
					if err := copierInternal.Copy(proposal, daoArtist); err == nil {
						resp.Proposal = proposal
					}
				}
			}
			return resp, nil
		},
	).ServeHTTP(w, r)
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
		logger.AtLog.Logger.Error("withdraw.walletAddress", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	var reqBody request.WithDrawItemRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		logger.AtLog.Logger.Error("withdraw.Decode", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	err = reqBody.SelfValidate()
	if err != nil {
		logger.AtLog.Logger.Error("withdraw.SelfValidate", zap.Error(err), zap.Any("reqBody", reqBody))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	wdr := &structure.WithDrawItemRequest{}
	err = copier.Copy(wdr, reqBody)
	if err != nil {
		logger.AtLog.Logger.Error("withdraw.Copy", zap.Error(err), zap.Any("wdr", wdr))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	wd, err := h.Usecase.CreateWithdraw(walletAddress, *wdr)
	if err != nil {
		logger.AtLog.Logger.Error("withdraw.CreateWithdraw", zap.Error(err), zap.Any("wdr", wdr))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, wd, "")
}
