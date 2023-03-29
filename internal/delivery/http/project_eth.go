package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/logger"
)

// UserCredits godoc
// @Summary Create project
// @Description Create projects
// @Tags ETH Project
// @Accept  json
// @Produce  json
// @Security Authorization
// @Param request body request.CreateETHProjectReq true "Create eth-project request"
// @Success 200 {object} response.JsonResponse{}
// @Router /project [POST]
func (h *httpDelivery) createEthProjects(w http.ResponseWriter, r *http.Request) {
	var reqBody request.CreateETHProjectReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		logger.AtLog.Logger.Error("decoder.Decode", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	ctx := r.Context()
	iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)
	walletAddress, ok := iWalletAddress.(string)
	if !ok {
		err := errors.New("Wallet address is incorect")
		logger.AtLog.Logger.Error("ctx.Value.Token",  zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase := &structure.CreateProjectReq{}
	err = copier.Copy(reqUsecase, reqBody)
	if err != nil {
		logger.AtLog.Logger.Error("copier.Copy", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase.CreatorAddrr = walletAddress
	message, err := h.Usecase.CreateProject(*reqUsecase)
	if err != nil {
		logger.AtLog.Logger.Error("h.Usecase.CreateProject", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp, err := h.projectToResp(message)
	if err != nil {
		logger.AtLog.Logger.Error("h.projectToResp", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}