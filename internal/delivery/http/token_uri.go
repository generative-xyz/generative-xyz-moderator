package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/usecase/structure"
)

// UserCredits godoc
// @Summary get token uri data
// @Description get token uri data
// @Tags token_uri
// @Accept  json
// @Produce  json
// @Param contractAddress path string true "contract address"
// @Param tokenID path string true "token ID"
// @Success 200 {object} response.JsonResponse{data=response.TokenURIResp}
// @Router /token/{contractAddress}/{tokenID} [GET]
func (h *httpDelivery) tokenURI(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("tokenURI", r)
	defer h.Tracer.FinishSpan(span, log )

	vars := mux.Vars(r)
	contractAddress := vars["contractAddress"]
	span.SetTag("contractAddress", contractAddress)
	
	tokenID := vars["tokenID"]
	span.SetTag("tokenID", tokenID)

	message, err := h.Usecase.GetToken(span, structure.GetTokenMessageReq{
		ContractAddress: contractAddress,
		TokenID: tokenID,
	})

	if err != nil {
		log.Error("h.Usecase.GetToken", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	resp := response.TokenURIResp{
		Name: message.Name,
		Description: message.Description,
		Image: *message.ParsedImage,
		AnimationURL: message.AnimationURL,
		Attributes: message.ParsedAttributes,
	}

	log.SetData("resp.message", message)
	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondWithoutContainer(w, http.StatusOK, resp)
}

// UserCredits godoc
// @Summary get token's traits
// @Description get token's traits
// @Tags token_uri
// @Accept  json
// @Produce  json
// @Param contractAddress path string true "contract address"
// @Param tokenID path string true "token ID"
// @Success 200 {object} response.JsonResponse{data=response.TokenTraitsResp}
// @Router /trait/{contractAddress}/{tokenID} [GET]
func (h *httpDelivery) tokenTrait(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("tokenTrait", r)
	defer h.Tracer.FinishSpan(span, log )

	vars := mux.Vars(r)
	contractAddress := vars["contractAddress"]
	span.SetTag("contractAddress", contractAddress)
	
	tokenID := vars["tokenID"]
	span.SetTag("tokenID", tokenID)

	message, err := h.Usecase.GetTokenTraits(span, structure.GetTokenMessageReq{
		ContractAddress: contractAddress,
		TokenID: tokenID,
	})

	if err != nil {
		log.Error("h.Usecase.GetToken", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	resp := response.TokenTraitsResp{}
	resp.Attributes  = message.ParsedAttributes
	log.SetData("resp.message", message)
	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondWithoutContainer(w, http.StatusOK, resp)
}

// UserCredits godoc
// @Summary get project's detail
// @Description get project's detail
// @Tags token_uri
// @Accept  json
// @Produce  json
// @Param contractAddress path string true "contract address"
// @Param projectID path string true "token ID"
// @Success 200 {object} response.JsonResponse{data=response.TokenTraitsResp}
// @Router /project/{contractAddress}/{projectID} [GET]
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
