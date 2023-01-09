package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
)

// UserCredits godoc
// @Summary Get market place listing
// @Description Get market place listing
// @Tags MarketPlace
// @Accept  json
// @Produce  json
// @Param genNFTAddr path string true "genNFTAddrress"
// @Param tokenID path string true "tokenID"
// @Param closed query bool false "true|false, default all"
// @Param finished query bool false "true|false, default all"
// @Param sort_by query string false "sort by field"
// @Param sort query int false "1: ASC, -1: DESC"
// @Param limit query int false "limit default 10"
// @Param page query int false "page start with 1"
// @Success 200 {object} response.JsonResponse{}
// @Router /marketplace/listing/{genNFTAddr}/token/{tokenID} [GET]
func (h *httpDelivery) getListingViaGenAddressTokenID(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("getListingViaGenAddressTokenID", r)
	defer h.Tracer.FinishSpan(span, log )

	vars := mux.Vars(r)
	genNFTAddr := vars["genNFTAddr"]
	tokenID := vars["tokenID"]
	
	bf, err := h.BaseFilters(r)
	if err != nil {
		log.Error("h.Usecase.getListingViaGenAddressTokenID.BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	
	closed := r.URL.Query().Get("closed")
	finished := r.URL.Query().Get("finished")
	f := structure.FilterMkListing{}
	if closed != "" {
		closedBool, err := strconv.ParseBool(closed)
		if err != nil {
			log.Error("strconv.ParseBool.closed", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
			return
		}
		f.Closed = &closedBool
	}
	
	if finished != "" {
		finishedBool, err := strconv.ParseBool(finished)
		if err != nil {
			log.Error("strconv.ParseBool.finished", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
			return
		}
		f.Finished = &finishedBool
	}
	
	f.CollectionContract = &genNFTAddr
	f.TokenId = &tokenID
	f.BaseFilters = *bf
	
	resp, err := h.getMkListings(span, f)
	if err != nil {
		log.Error("h.Usecase.getMkListings.getTokens", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success , resp, "")
}

// UserCredits godoc
// @Summary Get market place offer
// @Description Get market place offer
// @Tags MarketPlace
// @Accept  json
// @Produce  json
// @Param genNFTAddr path string true "genNFTAddrress"
// @Param tokenID path string true "tokenID"
// @Param closed query bool false "true|false, default all"
// @Param finished query bool false "true|false, default all"
// @Param sort_by query string false "sort by field"
// @Param sort query int false "1: ASC, -1: DESC"
// @Param limit query int false "limit default 10"
// @Param page query int false "page start with 1"
// @Success 200 {object} response.JsonResponse{}
// @Router /marketplace/offers/{genNFTAddr}/token/{tokenID} [GET]
func (h *httpDelivery) getOffersViaGenAddressTokenID(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("getOffersViaGenAddressTokenID", r)
	defer h.Tracer.FinishSpan(span, log )

	vars := mux.Vars(r)
	genNFTAddr := vars["genNFTAddr"]
	tokenID := vars["tokenID"]
	
	bf, err := h.BaseFilters(r)
	if err != nil {
		log.Error("h.Usecase.getListingViaGenAddressTokenID.BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	
	closed := r.URL.Query().Get("closed")
	finished := r.URL.Query().Get("finished")
	f := structure.FilterMkOffers{}
	if closed != "" {
		closedBool, err := strconv.ParseBool(closed)
		if err != nil {
			log.Error("strconv.ParseBool.closed", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
			return
		}
		f.Closed = &closedBool
	}
	
	if finished != "" {
		finishedBool, err := strconv.ParseBool(finished)
		if err != nil {
			log.Error("strconv.ParseBool.finished", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
			return
		}
		f.Finished = &finishedBool
	}
	
	f.CollectionContract = &genNFTAddr
	f.TokenId = &tokenID
	f.BaseFilters = *bf
	
	resp, err := h.getMkOffers(span, f)
	if err != nil {
		log.Error("h.Usecase.getMkListings.getTokens", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success , resp, "")
}

// UserCredits godoc
// @Summary User profile's selling nft
// @Description User profile's selling nft
// @Tags MarketPlace
// @Accept  json
// @Produce  json
// @Param walletAddress path string true "Wallet address"
// @Param sort_by query string false "sort by field"
// @Param sort query int false "1: ASC, -1: DESC"
// @Param limit query int false "limit default 10"
// @Param page query int false "page start with 1"
// @Success 200 {object} response.JsonResponse{data=response.InternalTokenURIResp}
// @Router /marketplace/wallet/{walletAddress}/listing [GET]
func (h *httpDelivery) ListingOfAProfile(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("ListingOfAProfile", r)
	defer h.Tracer.FinishSpan(span, log )

	vars := mux.Vars(r)
	walletAddress := vars["walletAddress"]
	
	bf, err := h.BaseFilters(r)
	if err != nil {
		log.Error("h.Usecase.getListingViaGenAddressTokenID.BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	
	closed := r.URL.Query().Get("closed")
	finished := r.URL.Query().Get("finished")
	f := structure.FilterMkListing{}
	if closed != "" {
		closedBool, err := strconv.ParseBool(closed)
		if err != nil {
			log.Error("strconv.ParseBool.closed", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
			return
		}
		f.Closed = &closedBool
	}
	
	if finished != "" {
		finishedBool, err := strconv.ParseBool(finished)
		if err != nil {
			log.Error("strconv.ParseBool.finished", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
			return
		}
		f.Finished = &finishedBool
	}
	
	f.SellerAddress = &walletAddress
	f.BaseFilters = *bf
	
	resp, err := h.getMkListings(span, f)
	if err != nil {
		log.Error("h.Usecase.getMkListings.getTokens", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success , resp, "")

}

// UserCredits godoc
// @Summary User profile's selling nft
// @Description User profile's selling nft
// @Tags MarketPlace
// @Accept  json
// @Produce  json
// @Param walletAddress path string true "Wallet address"
// @Param sort_by query string false "sort by field"
// @Param sort query int false "1: ASC, -1: DESC"
// @Param limit query int false "limit default 10"
// @Param page query int false "page start with 1"
// @Success 200 {object} response.JsonResponse{data=response.InternalTokenURIResp}
// @Router /marketplace/wallet/{walletAddress}/offer [GET]
func (h *httpDelivery) OfferOfAProfile(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("OfferOfAProfile", r)
	defer h.Tracer.FinishSpan(span, log )

	vars := mux.Vars(r)
	walletAddress := vars["walletAddress"]
	
	bf, err := h.BaseFilters(r)
	if err != nil {
		log.Error("h.Usecase.getListingViaGenAddressTokenID.BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	
	closed := r.URL.Query().Get("closed")
	finished := r.URL.Query().Get("finished")
	f := structure.FilterMkOffers{}
	if closed != "" {
		closedBool, err := strconv.ParseBool(closed)
		if err != nil {
			log.Error("strconv.ParseBool.closed", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
			return
		}
		f.Closed = &closedBool
	}
	
	if finished != "" {
		finishedBool, err := strconv.ParseBool(finished)
		if err != nil {
			log.Error("strconv.ParseBool.finished", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
			return
		}
		f.Finished = &finishedBool
	}
	
	f.BuyerAddress = &walletAddress
	f.BaseFilters = *bf
	
	resp, err := h.getMkOffers(span, f)
	if err != nil {
		log.Error("h.Usecase.getMkListings.getTokens", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success , resp, "")

}

func (h *httpDelivery) getMkListings(rootSpan opentracing.Span, f  structure.FilterMkListing) (*response.PaginationResponse, error) {
	span, log := h.StartSpanFromRoot(rootSpan, "httpDelivery.getTokens")
	defer h.Tracer.FinishSpan(span, log )


	pag, err := h.Usecase.FilterMKListing(span, f)
	if err != nil {
		log.Error("h.Usecase.getProfileNfts.FilterTokens", err.Error(), err)
		return nil, err
	}

	respItems := []response.MarketplaceListing{}
	iMkData := pag.Result
	mkData, ok := (iMkData).([]entity.MarketplaceListings)
	if !ok {
		err := errors.New( "Cannot parse MarketplaceListings")
		log.Error("ctx.Value.Token",  err.Error(), err)
		return nil, err
	}

	for _, mk := range mkData {	
		resp, err := h.mkListingToResp(&mk)
		if err != nil {
			err := errors.New( "Cannot parse MarketplaceListin")
			log.Error("tokenToResp",  err.Error(), err)
			return nil, err
		}
		respItems = append(respItems, *resp)
	}

	resp := h.PaginationResp(pag, respItems)
	return &resp, nil
	
}

func (h *httpDelivery) mkListingToResp(input *entity.MarketplaceListings) (*response.MarketplaceListing, error) {
	resp := &response.MarketplaceListing{}
	err := response.CopyEntityToRes(resp, input)
	if err != nil {
		return nil, err
	}

	pResp, err := h.profileToResp(&input.SellerInfo)
	if err == nil {
		resp.SellerInfo = *pResp
	}

	return resp, nil
}




func (h *httpDelivery) getMkOffers(rootSpan opentracing.Span, f  structure.FilterMkOffers) (*response.PaginationResponse, error) {
	span, log := h.StartSpanFromRoot(rootSpan, "httpDelivery.getTokens")
	defer h.Tracer.FinishSpan(span, log )


	pag, err := h.Usecase.FilterMKOffers(span, f)
	if err != nil {
		log.Error("h.Usecase.getProfileNfts.FilterTokens", err.Error(), err)
		return nil, err
	}

	respItems := []response.MarketplaceOffer{}
	iMkData := pag.Result
	mkData, ok := (iMkData).([]entity.MarketplaceOffers)
	if !ok {
		err := errors.New( "Cannot parse MarketplaceOffers")
		log.Error("ctx.Value.Token",  err.Error(), err)
		return nil, err
	}

	for _, mk := range mkData {	
		resp, err := h.mkOfferToResp(&mk)
		if err != nil {
			err := errors.New( "Cannot parse mkOfferToResp")
			log.Error("tokenToResp",  err.Error(), err)
			return nil, err
		}
		respItems = append(respItems, *resp)
	}

	resp := h.PaginationResp(pag, respItems)
	return &resp, nil
	
}

func (h *httpDelivery) mkOfferToResp(input *entity.MarketplaceOffers) (*response.MarketplaceOffer, error) {
	resp := &response.MarketplaceOffer{}
	err := response.CopyEntityToRes(resp, input)
	if err != nil {
		return nil, err
	}

	pResp, err := h.profileToResp(&input.BuyerInfo)
	if err == nil {
		resp.BuyerInfo = *pResp
	}

	return resp, nil
}

// UserCredits godoc
// @Summary get project's detail
// @Description get project's detail
// @Tags MarketPlace
// @Accept  json
// @Produce  json
// @Param genNFTAddr path string true "Gen NFT Addr"
// @Success 200 {object} response.JsonResponse{}
// @Router /marketplace/stats/{genNFTAddr} [GET]
func (h *httpDelivery) getCollectionStats(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("getCollectionStats", r)
	defer h.Tracer.FinishSpan(span, log)
	vars := mux.Vars(r)
	genNFTAddr := vars["genNFTAddr"]
	span.SetTag("genNFTAddr", genNFTAddr)
	project, err := h.Usecase.GetProjectByGenNFTAddr(span, genNFTAddr)
	if  err != nil {
		log.Error(" h.GetProjectByGenNFTAddr", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}
	resp, err := h.projectToStatResp(project)
	if  err != nil {
		log.Error(" h.projectToStatResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}
	log.SetData("project", project)
	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp , "")
}

func (h *httpDelivery) projectToStatResp(project *entity.Projects) (*response.MarketplaceStatResp, error) {
	return &response.MarketplaceStatResp{
		Stats: response.ProjectStatResp{
			UniqueOwnerCount: project.Stats.UniqueOwnerCount,
			TotalTradingVolumn: project.Stats.TotalTradingVolumn,
			FloorPrice: project.Stats.FloorPrice,
			BestMakeOfferPrice: project.Stats.BestMakeOfferPrice,
			ListedPercent: project.Stats.ListedPercent,
		},
	}, nil
}
