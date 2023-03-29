package http

import (
	"encoding/json"
	"net/http"

	"github.com/jinzhu/copier"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/usecase/structure"
)

// UserCredits godoc
// @Summary Create project
// @Description Create projects
// @Tags ETH Project
// @Accept  json
// @Produce  json
// @Param request body request.CreateETHProjectReq true "Create eth-project request"
// @Success 200 {object} response.JsonResponse{}
// @Router /project/eth [POST]
func (h *httpDelivery) createEthProjects(w http.ResponseWriter, r *http.Request) {
	var reqBody request.CreateETHProjectReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		h.Logger.Error("decoder.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase := &structure.CreateProjectReq{}
	err = copier.Copy(reqUsecase, reqBody)
	if err != nil {
		h.Logger.Error("copier.Copy", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Logger.Info("reqUsecase", reqUsecase)

	message, err := h.Usecase.CreateProject(*reqUsecase)
	if err != nil {
		h.Logger.Error("h.Usecase.CreateProject", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp, err := h.projectToResp(message)
	if err != nil {
		h.Logger.Error("h.projectToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

// UserCredits godoc
// @Summary Upload files for ETH-project
// @Description Upload files for ETH-project
// @Tags ETH Project
// @Accept  json
// @Produce  json
// @Success 200 {object} response.JsonResponse{}
// @Router /project/eth/files [POST]
func (h *httpDelivery) uploadEthProjectFiles(w http.ResponseWriter, r *http.Request) {
	
}