package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/usecase/structure"
)

// UserCredits godoc
// @Summary Create project
// @Description Create projects
// @Tags Project
// @Accept  json
// @Produce  json
// @Param request body request.CreateProjectReq true "Create profile request"
// @Success 200 {object} response.JsonResponse{}
// @Router /project [POST]
func (h *httpDelivery) createProjects(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("messages.projects", r)
	defer h.Tracer.FinishSpan(span, log )

	message, err := h.Usecase.CreateProject(span, structure.CreateProjectReq{
	
	})

	if err != nil {
		log.Error("h.Usecase.GetToken", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	resp := &response.ProjectResp{}
	err = copier.Copy(resp, message)
	if err != nil {
		log.Error("copier.Copy", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}
	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, nil, "")
}

// UserCredits godoc
// @Summary get project's detail
// @Description get project's detail
// @Tags Project
// @Accept  json
// @Produce  json
// @Param contractAddress path string true "contract address"
// @Param projectID path string true "token ID"
// @Success 200 {object} response.JsonResponse{}
// @Router /project/{contractAddress}/tokens/{projectID} [GET]
func (h *httpDelivery) projectDetail(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("tokenTrait", r)
	defer h.Tracer.FinishSpan(span, log )

	vars := mux.Vars(r)
	contractAddress := vars["contractAddress"]
	span.SetTag("contractAddress", contractAddress)
	
	projectID := vars["projectID"]
	span.SetTag("projectID", projectID)

	message, err := h.Usecase.GetProjectDetail(span, structure.GetProjectDetailMessageReq{
		ContractAddress: contractAddress,
		ProjectID: projectID,
	})

	if err != nil {
		log.Error("h.Usecase.GetToken", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	resp := &response.ProjectResp{}
	err = copier.Copy(resp, message.ProjectDetail)
	if err != nil {
		log.Error("copier.Copy", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}
	resp.Status = message.Status
	resp.NftTokenURI = message.NftTokenUri
	
	log.SetData("resp.message", message)
	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondWithoutContainer(w, http.StatusOK, resp)
}

// UserCredits godoc
// @Summary get project's tokens
// @Description get tokens by project address
// @Tags Project
// @Accept  json
// @Produce  json
// @Param contractAddress path string true "contract address"
// @Success 200 {object} response.JsonResponse{}
// @Router /project/{contractAddress}/tokens [GET]
func (h *httpDelivery) projectTokens(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("projectTokens", r)
	defer h.Tracer.FinishSpan(span, log )

	vars := mux.Vars(r)
	contractAddress := vars["contractAddress"]
	span.SetTag("contractAddress", contractAddress)
	
	

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondWithoutContainer(w, http.StatusOK, nil)
}
