package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	"rederinghub.io/external/nfts"
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

	project, err := h.Usecase.GetProjectDetail(span, structure.GetProjectDetailMessageReq{
		ContractAddress: contractAddress,
		ProjectID: projectID,
	})

	if err != nil {
		log.Error("h.Usecase.GetToken", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	resp, err := h.projectToResp(*project)
	if err != nil {
		log.Error(" h.projectToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}
	
	log.SetData("resp.project", project)
	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp , "")
}

// UserCredits godoc
// @Summary get projects
// @Description get projects
// @Tags Project
// @Accept  json
// @Produce  json
// @Param limit query int false "limit"
// @Param cursor query string false "The cursor returned in the previous response (used for getting the next page)."
// @Param contractAddress path string true "contract address"
// @Success 200 {object} response.JsonResponse{}
// @Router /project [GET]
func (h *httpDelivery) getProjects(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("projects", r)
	defer h.Tracer.FinishSpan(span, log )
	var err error
	vars := mux.Vars(r)
	contractAddress := vars["contractAddress"]
	span.SetTag("contractAddress", contractAddress)
	limitInt := 10

	limit := r.URL.Query().Get("limit")
	cursor := r.URL.Query().Get("cursor")
	if limit != "" {
		limitInt, err = strconv.Atoi(limit)
		if err != nil {
			log.Error("strconv.Atoi.limit", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
			return
		}
	}

	_ = limitInt
	_ = cursor

	
	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondWithoutContainer(w, http.StatusOK, nil)
}

// UserCredits godoc
// @Summary get project's tokens
// @Description get tokens by project address
// @Tags Project
// @Accept  json
// @Produce  json
// @Param limit query int false "limit"
// @Param cursor query string false "The cursor returned in the previous response (used for getting the next page)."
// @Param contractAddress path string true "contract address"
// @Success 200 {object} response.JsonResponse{}
// @Router /project/{contractAddress}/tokens [GET]
func (h *httpDelivery) projectTokens(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("projectTokens", r)
	defer h.Tracer.FinishSpan(span, log )
	var err error
	vars := mux.Vars(r)
	contractAddress := vars["contractAddress"]
	span.SetTag("contractAddress", contractAddress)
	limitInt := 10

	limit := r.URL.Query().Get("limit")
	cursor := r.URL.Query().Get("cursor")
	if limit != "" {
		limitInt, err = strconv.Atoi(limit)
		if err != nil {
			log.Error("strconv.Atoi.limit", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
			return
		}
	}
	
	data, err := h.Usecase.GetTokensByContract(span, contractAddress, nfts.MoralisFilter{
		Limit: &limitInt,
		Cursor: &cursor,
	})
	if err != nil {
		log.Error("h.Usecase.GetTokensByContract", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	respItems := []response.ProjectResp{}
	iProjectData := data.Result
	projectsData, ok := (iProjectData).([]structure.ProjectDetail)
	if !ok {
		err := errors.New( "Cannot parse products")
		log.Error("ctx.Value.Token",  err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	for _, project := range projectsData {	
		resp, err := h.projectToResp(project)
		if err != nil {
			log.Error(" h.projectToResp", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
			return
		}
		respItems = append(respItems, *resp)
	}

	resp := h.PaginationResp(data, respItems)
	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success , resp, "")
}

func (h *httpDelivery) projectToResp(input structure.ProjectDetail) (*response.ProjectResp, error) {
	resp := &response.ProjectResp{}
	err := copier.Copy(resp, input.ProjectDetail)
	if err != nil {
		return nil, err
	}
	resp.Status = input.Status
	resp.NftTokenURI = input.NftTokenUri
	return resp, nil
}