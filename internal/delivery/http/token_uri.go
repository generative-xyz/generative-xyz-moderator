package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/usecase/structure"
)

// UserCredits godoc
// @Summary get token uri data
// @Description get token uri data
// @Tags Token for Opensea
// @Accept  json
// @Produce  json
// @Param contractAddress path string true "contract address"
// @Param tokenID path string true "token ID"
// @Param captureTimeout query integer false "Capture timeout"
// @Success 200 {object} response.JsonResponse{data=response.TokenURIResp}
// @Router /token/{contractAddress}/{tokenID} [GET]
func (h *httpDelivery) tokenURI(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("tokenURI", r)
	defer h.Tracer.FinishSpan(span, log)

	vars := mux.Vars(r)
	contractAddress := vars["contractAddress"]
	span.SetTag("contractAddress", contractAddress)

	tokenID := vars["tokenID"]
	span.SetTag("tokenID", tokenID)

	captureTimeout := r.URL.Query().Get("captureTimeout")
	log.SetData("captureTimeout", captureTimeout)
	captureTimeoutInt, errT := strconv.Atoi(captureTimeout)
	if errT != nil {
		captureTimeoutInt = 5
	}

	message, err := h.Usecase.GetToken(span, structure.GetTokenMessageReq{
		ContractAddress: contractAddress,
		TokenID:         tokenID,
	}, captureTimeoutInt)

	if err != nil {
		log.Error("h.Usecase.GetToken", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp := response.TokenURIResp{
		Name:         message.Name,
		Description:  message.Description,
		Image:        *message.ParsedImage,
		AnimationURL: message.AnimationURL,
		Attributes:   message.ParsedAttributes,
	}

	log.SetData("resp.message", message)
	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondWithoutContainer(w, http.StatusOK, resp)
}

// UserCredits godoc
// @Summary get token's traits
// @Description get token's traits
// @Tags Token for Opensea
// @Accept  json
// @Produce  json
// @Param contractAddress path string true "contract address"
// @Param tokenID path string true "token ID"
// @Success 200 {object} response.JsonResponse{data=response.TokenTraitsResp}
// @Router /trait/{contractAddress}/{tokenID} [GET]
func (h *httpDelivery) tokenTrait(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("tokenTrait", r)
	defer h.Tracer.FinishSpan(span, log)

	vars := mux.Vars(r)
	contractAddress := vars["contractAddress"]
	span.SetTag("contractAddress", contractAddress)

	tokenID := vars["tokenID"]
	span.SetTag("tokenID", tokenID)

	message, err := h.Usecase.GetTokenTraits(span, structure.GetTokenMessageReq{
		ContractAddress: contractAddress,
		TokenID:         tokenID,
	})

	if err != nil {
		log.Error("h.Usecase.GetToken", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp := response.TokenTraitsResp{}
	resp.Attributes = message.ParsedAttributes
	log.SetData("resp.message", message)
	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondWithoutContainer(w, http.StatusOK, resp)
}

// UserCredits godoc
// @Summary get token uri data
// @Description get token uri data
// @Tags Tokens
// @Accept  json
// @Produce  json
// @Param contractAddress path string true "contract address"
// @Param tokenID path string true "token ID"
// @Param captureTimeout query integer false "Capture timeout"
// @Success 200 {object} response.JsonResponse{data=response.InternalTokenURIResp}
// @Router /tokens/{contractAddress}/{tokenID} [GET]
func (h *httpDelivery) tokenURIWithResp(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("tokenURI", r)
	defer h.Tracer.FinishSpan(span, log)

	vars := mux.Vars(r)
	contractAddress := vars["contractAddress"]
	span.SetTag("contractAddress", contractAddress)

	tokenID := vars["tokenID"]
	span.SetTag("tokenID", tokenID)

	captureTimeout := r.URL.Query().Get("captureTimeout")
	log.SetData("captureTimeout", captureTimeout)
	captureTimeoutInt, errT := strconv.Atoi(captureTimeout)
	if errT != nil {
		captureTimeoutInt = 5
	}

	message, err := h.Usecase.GetToken(span, structure.GetTokenMessageReq{
		ContractAddress: contractAddress,
		TokenID:         tokenID,
	}, captureTimeoutInt)

	if err != nil {
		log.Error("h.Usecase.GetToken", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	// find project by projectID and contract address
	project, err := h.Usecase.GetProjectDetail(span, structure.GetProjectDetailMessageReq{ContractAddress: contractAddress, ProjectID: message.ProjectID})
	if err != nil {
		log.Error("h.Usecase.GetToken", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	// find owner address on moralis
	nft, err := h.Usecase.MoralisNft.GetNftByContractAndTokenID(project.GenNFTAddr, tokenID)
	if err != nil {
		log.Error("h.Usecase.GetToken", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	profile, _ := h.Usecase.GetUserProfileByWalletAddress(span, strings.ToLower(nft.Owner))
	var profileResp *response.ProfileResponse
	if profile != nil {
		profileResp, _ = h.profileToResp(profile)
	} else {
		profileResp = nil
	}

	projectResp, err := h.projectToResp(project)
	if err != nil {
		log.Error("h.Usecase.GetToken", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp := response.InternalTokenURIResp{
		Name:         message.Name,
		Description:  message.Description,
		Image:        *message.ParsedImage,
		AnimationURL: message.AnimationURL,
		Attributes:   message.ParsedAttributes,
		OwnerAddr:    nft.Owner,
		Owner:        profileResp,
		MintedTime:   *message.MintedTime,
		Project:      projectResp,
	}

	log.SetData("resp.message", message)
	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

// UserCredits godoc
// @Summary get token's traits
// @Description get token's traits
// @Tags Tokens
// @Accept  json
// @Produce  json
// @Param contractAddress path string true "contract address"
// @Param tokenID path string true "token ID"
// @Success 200 {object} response.JsonResponse{data=response.InternalTokenTraitsResp}
// @Router /tokens/traits/{contractAddress}/{tokenID} [GET]
func (h *httpDelivery) tokenTraitWithResp(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("tokenTrait", r)
	defer h.Tracer.FinishSpan(span, log)

	vars := mux.Vars(r)
	contractAddress := vars["contractAddress"]
	span.SetTag("contractAddress", contractAddress)

	tokenID := vars["tokenID"]
	span.SetTag("tokenID", tokenID)

	message, err := h.Usecase.GetTokenTraits(span, structure.GetTokenMessageReq{
		ContractAddress: contractAddress,
		TokenID:         tokenID,
	})

	if err != nil {
		log.Error("h.Usecase.GetToken", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp := response.InternalTokenTraitsResp{}
	resp.Attributes = message.ParsedAttributes
	log.SetData("resp.message", message)
	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}
