package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
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
		Image:        message.ParsedImage,
		AnimationURL: message.AnimationURL,
		Attributes:   message.ParsedAttributes,
	}

	log.SetData("resp.message", message.TokenID)
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

	message, err := h.Usecase.GetToken(span, structure.GetTokenMessageReq{
		ContractAddress: contractAddress,
		TokenID:         tokenID,
	}, 5)

	if err != nil {
		log.Error("h.Usecase.GetToken", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp := response.TokenTraitsResp{}
	resp.Attributes = message.ParsedAttributes
	log.SetData("resp.message", message)
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
	span, log := h.StartSpan("tokenURIWithResp", r)
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

	token, err := h.Usecase.GetToken(span, structure.GetTokenMessageReq{
		ContractAddress: contractAddress,
		TokenID:         tokenID,
	}, captureTimeoutInt)

	if err != nil {
		log.Error("h.Usecase.GetToken", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	log.SetData("h.Usecase.GetToken", token.TokenID)
	log.SetTag("tokenID", token.TokenID)
	resp, err := h.tokenToResp(token)
	if err != nil {
		err := errors.New("Cannot parse products")
		log.Error("tokenToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	log.SetData("resp.token", token.TokenID)
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

	message, err := h.Usecase.GetToken(span, structure.GetTokenMessageReq{
		ContractAddress: contractAddress,
		TokenID:         tokenID,
	}, 5)

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

// UserCredits godoc
// @Summary get project's tokens
// @Description get tokens by project address
// @Tags Project
// @Accept  json
// @Produce  json
// @Param contract_address query string false "contract_address"
// @Param gen_nft_address query string false "gen_nft_address"
// @Param owner_address query string false "owner_address"
// @Param creator_address query string false "creator_address"
// @Param tokenID query string false "Filter via tokenID"
// @Param attributes query []string false "attributes"
// @Param has_price query bool false "has_price"
// @Param from_price query string false "from_price"
// @Param to_price query string false "to_price"
// @Param sort query string false "newest, minted-newest, token-price-asc, token-price-desc"
// @Param limit query int false "limit"
// @Param cursor query string false "The cursor returned in the previous response (used for getting the next page)."
// @Param genNFTAddr path string true "This is provided from Project Detail API"
// @Success 200 {object} response.JsonResponse{}
// @Router /project/{genNFTAddr}/tokens [GET]
func (h *httpDelivery) TokensOfAProject(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("httpDelivery.TokensOfAProfile", r)
	defer h.Tracer.FinishSpan(span, log)

	vars := mux.Vars(r)
	genNFTAddr := vars["genNFTAddr"]
	log.SetData("genNFTAddr", genNFTAddr)
	log.SetTag(utils.GEN_NFT_ADDRESS_TAG, genNFTAddr)
	f := structure.FilterTokens{}
	err := f.CreateFilter(r)
	if err != nil {
		log.Error("f.CreateFilter", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	f.GenNFTAddr = &genNFTAddr
	bf, err := h.BaseFilters(r)
	if err != nil {
		log.Error("h.Usecase.getProfileNfts.BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	f.BaseFilters = *bf
	resp, err := h.getTokens(span, f)
	if err != nil {
		log.Error("h.Usecase.getProfileNfts.getTokens", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	//h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

// UserCredits godoc
// @Summary User profile's nft
// @Description User profile's nft
// @Tags Profile
// @Accept  json
// @Produce  json
// @Param tokenID path string false "Filter via tokenID"
// @Param walletAddress path string true "Wallet address"
// @Param contract_address query string false "contract_address"
// @Param gen_nft_address query string false "gen_nft_address"
// @Param owner_address query string false "owner_address"
// @Param creator_address query string false "creator_address"
// @Param tokenID query string false "Filter via tokenID"
// @Param sort query string false "newest, minted-newest, priority-asc, priority-desc"
// @Success 200 {object} response.JsonResponse{data=response.InternalTokenURIResp}
// @Router /profile/wallet/{walletAddress}/nfts [GET]
func (h *httpDelivery) TokensOfAProfile(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("httpDelivery.TokensOfAProfile", r)
	defer h.Tracer.FinishSpan(span, log)

	vars := mux.Vars(r)
	walletAddress := vars["walletAddress"]
	log.SetData("walletAddress", walletAddress)
	log.SetTag(utils.WALLET_ADDRESS_TAG, walletAddress)
	f := structure.FilterTokens{}
	f.CreateFilter(r)
	f.OwnerAddr = &walletAddress

	bf, err := h.BaseFilters(r)
	if err != nil {
		log.Error("h.Usecase.getProfileNfts.BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	f.BaseFilters = *bf
	tokenID := r.URL.Query().Get("tokenID")
	if tokenID != "" {
		f.TokenIDs = append(f.TokenIDs, tokenID)
	}

	resp, err := h.getTokens(span, f)
	if err != nil {
		log.Error("h.Usecase.getProfileNfts.getTokens", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	//h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")

}

// UserCredits godoc
// @Summary get  projects by wallet
// @Description get  projects by wallet
// @Tags Profile
// @Accept  json
// @Produce  json
// @Param walletAddress path string false "Filter project via wallet address"
// @Param limit query int false "limit"
// @Param cursor query string false "The cursor returned in the previous response (used for getting the next page)."
// @Success 200 {object} response.JsonResponse{}
// @Router /profile/wallet/{walletAddress}/projects [GET]
func (h *httpDelivery) getProjectsByWallet(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("getProjectsByWallet", r)
	defer h.Tracer.FinishSpan(span, log)

	var err error
	vars := mux.Vars(r)
	walletAddress := vars["walletAddress"]
	span.SetTag("walletAddress", walletAddress)

	baseF, err := h.BaseFilters(r)
	if err != nil {
		log.Error("BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	f := structure.FilterProjects{}
	f.BaseFilters = *baseF
	f.WalletAddress = &walletAddress

	uProjects, err := h.Usecase.GetProjects(span, f)
	if err != nil {
		log.Error("h.Usecase.GetProjects", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	pResp := []response.ProjectResp{}
	iProjects := uProjects.Result
	projects := iProjects.([]entity.Projects)
	for _, project := range projects {

		p, err := h.projectToResp(&project)
		if err != nil {
			log.Error("copier.Copy", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}

		pResp = append(pResp, *p)
	}

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, h.PaginationResp(uProjects, pResp), "")
}

func (h *httpDelivery) getTokens(rootSpan opentracing.Span, f structure.FilterTokens) (*response.PaginationResponse, error) {
	span, log := h.StartSpanFromRoot(rootSpan, "httpDelivery.getTokens")
	defer h.Tracer.FinishSpan(span, log)
	pag, err := h.Usecase.FilterTokens(span, f)
	if err != nil {
		log.Error("h.Usecase.getProfileNfts.FilterTokens", err.Error(), err)
		return nil, err
	}

	respItems := []response.InternalTokenURIResp{}
	tokens := []entity.TokenUri{}
	iTokensData := pag.Result

	bytes, err := json.Marshal(iTokensData)
	if err != nil {
		err := errors.New("Cannot parse respItems")
		log.Error("respItems", err.Error(), err)
		return nil, err
	}

	err = json.Unmarshal(bytes, &tokens)
	if err != nil {
		err := errors.New("Cannot Unmarshal")
		log.Error("Unmarshal", err.Error(), err)
		return nil, err
	}

	// get nft listing from marketplace:
	nftListing, _ := h.Usecase.GetAllListListingWithRule(span)

	for _, token := range tokens {
		resp, err := h.tokenToResp(&token)
		if err != nil {
			err := errors.New("Cannot parse products")
			log.Error("tokenToResp", err.Error(), err)
			return nil, err
		}

		for _, v := range nftListing {
			if resp != nil {
				if strings.EqualFold(v.InscriptionID, resp.TokenID) {
					resp.Project.Buyable = v.Buyable
					resp.Project.PriceBTC = v.Price
					resp.Project.OrderID = v.OrderID
					break
				}
			}
		}

		respItems = append(respItems, *resp)
	}

	resp := h.PaginationResp(pag, respItems)
	return &resp, nil

}

func (h *httpDelivery) tokenToResp(input *entity.TokenUri) (*response.InternalTokenURIResp, error) {
	resp := &response.InternalTokenURIResp{}
	err := response.CopyEntityToResNoID(resp, input)
	if err != nil {
		return nil, err
	}
	resp.Attributes = input.ParsedAttributes
	if input.ParsedImage != nil {
		resp.Image = *input.ParsedImage
	} else {
		resp.Image = input.Thumbnail
	}

	if input.Owner != nil {
		ownerResp, err := h.profileToResp(input.Owner)
		if err == nil {
			resp.Owner = ownerResp
		}
	}

	if input.Creator != nil {
		creatorResp, err := h.profileToResp(input.Creator)
		if err == nil {
			resp.Creator = creatorResp
		}
	}

	if input.Project != nil {
		projectResp := &response.ProjectResp{}
		response.CopyEntityToRes(projectResp, input.Project)
		// projectResp, err := h.projectToResp(input.Project)
		// if err == nil {

		// }

		resp.Project = projectResp
		// resp.Stats.Price = input.Stats.PriceInt
		if input.Stats.PriceInt == nil {
			resp.Stats.Price = nil
		} else {
			x := strconv.Itoa(int(*input.Stats.PriceInt))
			resp.Stats.Price = &x
		}
	}

	resp.InscriptionIndex = input.InscriptionIndex

	//resp.Thumbnail = fmt.Sprintf("%s/%s/%s/%s",os.Getenv("DOMAIN"), "api/thumbnail", input.ContractAddress, input.TokenID)

	return resp, nil
}

// UserCredits godoc
// @Summary get tokens
// @Description get tokens
// @Tags Tokens
// @Accept  json
// @Produce  json
// @Param contract_address query string false "contract_address"
// @Param gen_nft_address query string false "gen_nft_address"
// @Param owner_address query string false "owner_address"
// @Param creator_address query string false "creator_address"
// @Param tokenID query string false "Filter via tokenID"
// @Param sort query string false "newest, minted-newest, priority-asc, priority-desc"
// @Param limit query int false "limit"
// @Param page query string false "The cursor returned in the previous response (used for getting the next page)."
// @Success 200 {object} response.JsonResponse{}
// @Router /tokens [GET]
func (h *httpDelivery) Tokens(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("httpDelivery.Tokens", r)
	defer h.Tracer.FinishSpan(span, log)
	f := structure.FilterTokens{}
	f.CreateFilter(r)

	bf, err := h.BaseFilters(r)
	if err != nil {
		log.Error("h.Tokens.BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	f.BaseFilters = *bf
	resp, err := h.getTokens(span, f)
	if err != nil {
		log.Error("h.Tokens.getTokens", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	//h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

// UserCredits godoc
// @Summary get token uri data
// @Description get token uri data
// @Tags Tokens
// @Accept  json
// @Produce  json
// @Param contractAddress path string true "contract address"
// @Param tokenID path string true "token ID"
// @Param request body request.UpdateTokentReq true "Request body"
// @Success 200 {object} response.JsonResponse{data=response.InternalTokenURIResp}
// @Router /tokens/{contractAddress}/{tokenID} [PUT]
func (h *httpDelivery) updatetokenURIWithResp(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("updatetokenURIWithResp", r)
	defer h.Tracer.FinishSpan(span, log)

	vars := mux.Vars(r)
	contractAddress := vars["contractAddress"]
	span.SetTag("contractAddress", contractAddress)

	tokenID := vars["tokenID"]
	span.SetTag("tokenID", tokenID)

	var reqBody request.UpdateTokentReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Error("decoder.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	token, err := h.Usecase.UpdateToken(span, structure.UpdateTokenReq{
		ContracAddress: contractAddress,
		TokenID:        tokenID,
		Priority:       reqBody.Priority,
	})

	if err != nil {
		log.Error("h.Usecase.GetToken", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	log.SetData("h.Usecase.GetToken", token)
	resp, err := h.tokenToResp(token)
	if err != nil {
		err := errors.New("Cannot parse products")
		log.Error("tokenToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	log.SetData("resp.token", token)
	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}
