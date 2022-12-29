package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	"rederinghub.io/external/nfts"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
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

	var reqBody request.CreateProjectReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Error("decoder.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase := &structure.CreateProjectReq{}
	err = copier.Copy(reqUsecase, reqBody)
	if err != nil {
		log.Error("copier.Copy", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	message, err := h.Usecase.CreateProject(span, *reqUsecase)
	if err != nil {
		log.Error("h.Usecase.GetToken", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}
	
	resp, err  := h.projectToResp(message)
	if err != nil {
		log.Error("h.projectToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
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

	resp, err := h.projectToResp(project)
	if err != nil {
		log.Error(" h.projectToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}
	
	log.SetData("resp.project", resp)
	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp , "")
}

// UserCredits godoc
// @Summary get projects
// @Description get projects
// @Tags Project
// @Accept  json
// @Produce  json
// @Param contractAddress query string false "Filter project via contract address"
// @Param limit query int false "limit"
// @Param cursor query string false "The cursor returned in the previous response (used for getting the next page)."
// @Success 200 {object} response.JsonResponse{}
// @Router /project [GET]
func (h *httpDelivery) getProjects(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("projects", r)
	defer h.Tracer.FinishSpan(span, log )
	var err error
	vars := mux.Vars(r)
	contractAddress := vars["contractAddress"]
	span.SetTag("contractAddress", contractAddress)

	baseF, err := h.BaseFilters(r)
	if err != nil {
		log.Error("BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	f := structure.FilterProjects{}
	f.BaseFilters = *baseF
	uProjects, err := h.Usecase.GetProjects(span, f)
	if err != nil {
		log.Error("h.Usecase.GetProjects", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	pResp :=  []response.ProjectResp{}
	iProjects := uProjects.Result
	projects := iProjects.([]entity.Projects)
	for _, project := range projects {

		p, err := h.projectToResp(&project)
		if err != nil {
			log.Error("copier.Copy", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
			return
		}

		pResp = append(pResp, *p)
	}

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, h.PaginationResp(uProjects, pResp), "")
}

// UserCredits godoc
// @Summary get minted out projects
// @Description  get minted out projects
// @Tags Project
// @Accept  json
// @Produce  json
// @Param contractAddress query string false "Filter project via contract address"
// @Param limit query int false "limit"
// @Param cursor query string false "The cursor returned in the previous response (used for getting the next page)."
// @Success 200 {object} response.JsonResponse{}
// @Router /project/minted-out [GET]
func (h *httpDelivery) getMintedOutProjects(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("projects", r)
	defer h.Tracer.FinishSpan(span, log )
	var err error
	vars := mux.Vars(r)
	contractAddress := vars["contractAddress"]
	span.SetTag("contractAddress", contractAddress)

	baseF, err := h.BaseFilters(r)
	if err != nil {
		log.Error("BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	f := structure.FilterProjects{}
	f.BaseFilters = *baseF
	uProjects, err := h.Usecase.GetMintedOutProjects(span, f)
	if err != nil {
		log.Error("h.Usecase.GetProjects", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	pResp :=  []response.ProjectResp{}
	iProjects := uProjects.Result
	projects := iProjects.([]entity.Projects)
	for _, project := range projects {

		p, err := h.projectToResp(&project)
		if err != nil {
			log.Error("copier.Copy", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
			return
		}

		pResp = append(pResp, *p)
	}

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, h.PaginationResp(uProjects, pResp), "")
}


// UserCredits godoc
// @Summary get the random projects
// @Description get the random projects
// @Tags Project
// @Accept  json
// @Produce  json
// @Success 200 {object} response.JsonResponse{}
// @Router /project/random [GET]
func (h *httpDelivery) getRandomProject(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("getRandomProject", r)
	defer h.Tracer.FinishSpan(span, log )
	var err error

	project, err := h.Usecase.GetRandomProject(span)
	if err != nil {
		log.Error(" h.GetRandomProject", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}
	
	resp, err := h.projectToResp(project)
	if err != nil {
		log.Error(" h.projectToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

// UserCredits godoc
// @Summary get project's tokens
// @Description get tokens by project address
// @Tags Project
// @Accept  json
// @Produce  json
// @Param limit query int false "limit"
// @Param cursor query string false "The cursor returned in the previous response (used for getting the next page)."
// @Param genNFTAddr path string true "This is provided from Project Detail API"
// @Success 200 {object} response.JsonResponse{}
// @Router /project/{genNFTAddr}/tokens [GET]
func (h *httpDelivery) projectTokens(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("projectTokens", r)
	defer h.Tracer.FinishSpan(span, log )
	var err error
	vars := mux.Vars(r)
	genNFTAddr := vars["genNFTAddr"]
	span.SetTag("genNFTAddr", genNFTAddr)
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
	
	data, err := h.Usecase.GetTokensByContract(span, genNFTAddr, nfts.MoralisFilter{
		Limit: &limitInt,
		Cursor: &cursor,
	})
	if err != nil {
		log.Error("h.Usecase.GetTokensByContract", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	respItems := []response.TokenURIResp{}
	iTokensData := data.Result
	tokensData, ok := (iTokensData).([]entity.TokenUri)
	if !ok {
		err := errors.New( "Cannot parse products")
		log.Error("ctx.Value.Token",  err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	for _, token := range tokensData {	
		resp := response.TokenURIResp{
			Name: token.Name,
			Description: token.Description,
			Image: *token.ParsedImage,
			AnimationURL: token.AnimationURL,
			Attributes: token.ParsedAttributes,
		}
		respItems = append(respItems, resp)
	}

	resp := h.PaginationResp(data, respItems)
	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success , resp, "")
}

func (h *httpDelivery) projectToResp(input *entity.Projects) (*response.ProjectResp, error) {
	resp := &response.ProjectResp{}
	social := make(map[string]string)
	response.CopyEntityToRes(resp, input)
	resp.MintPriceAddr = input.MintTokenAddress
	resp.Limit = input.LimitSupply
	resp.Creator = input.CreatorName
	resp.CreatorAddr = input.CreatorAddrr
	resp.Desc = input.Description
	resp.ItemDesc = input.Description
	resp.License = input.License
	resp.Image = input.Thumbnail
	resp.ScriptType = input.ThirdPartyScripts
	resp.NftTokenURI = input.NftTokenUri
	social["web"] = input.SocialWeb
	social["twitter"] = input.SocialTwitter
	social["discord"] = input.SocialDiscord
	social["medium"] = input.SocialMedium
	social["instagram"] = input.SocialInstagram
	resp.MintingInfo = response.NftMintingDetail{
		Index: input.MintingInfo.Index,
		IndexReserve: input.MintingInfo.IndexReverse,
	}
	resp.Royalty = input.Royalty
	resp.Reservers = input.Reservers
	resp.CompleteTime = input.CompleteTime
	return resp, nil
}