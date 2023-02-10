package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
)

func (h *httpDelivery) btcMarketplaceListing(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("httpDelivery.ethGetReceiveWalletAddress", r)
	defer h.Tracer.FinishSpan(span, log)
	h.Response.SetLog(h.Tracer, span)

	var reqBody request.CreateEthWalletAddressReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Error("httpDelivery.btcMint.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase := &structure.EthWalletAddressData{}
	err = copier.Copy(reqUsecase, reqBody)
	if err != nil {
		log.Error("copier.Copy", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	ethWallet, err := h.Usecase.CreateETHWalletAddress(span, *reqUsecase)
	if err != nil {
		log.Error("h.Usecase.CreateETHWalletAddress", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	log.SetData("ethWallet", ethWallet)
	resp, err := h.EthWalletAddressToResp(ethWallet)
	if err != nil {
		log.Error(" h.proposalToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

func (h *httpDelivery) btcMarketplaceListNFTs(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("btcMarketplaceListNFTs", r)
	defer h.Tracer.FinishSpan(span, log)

	// baseF, err := h.BaseFilters(r)
	// if err != nil {
	// 	log.Error("BaseFilters", err.Error(), err)
	// 	h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
	// 	return
	// }

	// h.Usecase.FindBtcNFTListingByNFTID()
	// f := structure.FilterProposalVote{}
	// f.BaseFilters = *baseF

	// f.ProposalID = &proposalID
	// support := r.URL.Query().Get("support")
	// if support != "" {
	// 	supportInt, err := strconv.Atoi(support)
	// 	if err != nil {
	// 		log.Error("strconv.Atoi", err.Error(), err)
	// 		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
	// 		return
	// 	}
	// 	f.Support = &supportInt
	// }

	// voter := r.URL.Query().Get("voter")
	// if voter != "" {
	// 	f.Voter = &voter
	// }

	// paginationData, err := h.Usecase.GetProposalVotes(span, f)
	// if err != nil {
	// 	log.Error("h.Usecase.GetProposal", err.Error(), err)
	// 	h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
	// 	return
	// }

	// pResp := []response.ProposalVotesResp{}
	// iPro := paginationData.Result
	// pro := iPro.([]entity.ProposalVotes)
	// for _, proItem := range pro {

	// 	tmp := &response.ProposalVotesResp{}
	// 	err := response.CopyEntityToRes(tmp, &proItem)
	// 	if err != nil {
	// 		log.Error("copier.Copy", err.Error(), err)
	// 		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
	// 		return
	// 	}

	// 	pResp = append(pResp, *tmp)
	// }

	// //log.SetData("resp.Proposal", resp)
	// h.Response.SetLog(h.Tracer, span)
	// h.Response.RespondSuccess(w, http.StatusOK, response.Success, h.PaginationResp(paginationData, pResp), "")
}

func (h *httpDelivery) btcMarketplaceNFTDetail(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("getProposalVotes", r)
	defer h.Tracer.FinishSpan(span, log)

	vars := mux.Vars(r)
	proposalID := vars["proposalID"]
	span.SetTag("proposalID", proposalID)

	baseF, err := h.BaseFilters(r)
	if err != nil {
		log.Error("BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	f := structure.FilterProposalVote{}
	f.BaseFilters = *baseF

	f.ProposalID = &proposalID
	support := r.URL.Query().Get("support")
	if support != "" {
		supportInt, err := strconv.Atoi(support)
		if err != nil {
			log.Error("strconv.Atoi", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
		f.Support = &supportInt
	}

	voter := r.URL.Query().Get("voter")
	if voter != "" {
		f.Voter = &voter
	}

	paginationData, err := h.Usecase.GetProposalVotes(span, f)
	if err != nil {
		log.Error("h.Usecase.GetProposal", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	pResp := []response.ProposalVotesResp{}
	iPro := paginationData.Result
	pro := iPro.([]entity.ProposalVotes)
	for _, proItem := range pro {

		tmp := &response.ProposalVotesResp{}
		err := response.CopyEntityToRes(tmp, &proItem)
		if err != nil {
			log.Error("copier.Copy", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}

		pResp = append(pResp, *tmp)
	}

	//log.SetData("resp.Proposal", resp)
	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, h.PaginationResp(paginationData, pResp), "")
}

func (h *httpDelivery) btcMarketplaceCreateBuyOrder(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("getProposalVotes", r)
	defer h.Tracer.FinishSpan(span, log)

	vars := mux.Vars(r)
	proposalID := vars["proposalID"]
	span.SetTag("proposalID", proposalID)

	baseF, err := h.BaseFilters(r)
	if err != nil {
		log.Error("BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	f := structure.FilterProposalVote{}
	f.BaseFilters = *baseF

	f.ProposalID = &proposalID
	support := r.URL.Query().Get("support")
	if support != "" {
		supportInt, err := strconv.Atoi(support)
		if err != nil {
			log.Error("strconv.Atoi", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
		f.Support = &supportInt
	}

	voter := r.URL.Query().Get("voter")
	if voter != "" {
		f.Voter = &voter
	}

	paginationData, err := h.Usecase.GetProposalVotes(span, f)
	if err != nil {
		log.Error("h.Usecase.GetProposal", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	pResp := []response.ProposalVotesResp{}
	iPro := paginationData.Result
	pro := iPro.([]entity.ProposalVotes)
	for _, proItem := range pro {

		tmp := &response.ProposalVotesResp{}
		err := response.CopyEntityToRes(tmp, &proItem)
		if err != nil {
			log.Error("copier.Copy", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}

		pResp = append(pResp, *tmp)
	}

	//log.SetData("resp.Proposal", resp)
	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, h.PaginationResp(paginationData, pResp), "")
}
