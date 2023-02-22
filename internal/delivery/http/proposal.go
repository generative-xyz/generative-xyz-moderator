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

// UserCredits godoc
// @Summary DAO list proposal
// @Description DAO list proposal
// @Tags DAO
// @Accept  json
// @Produce  json
// @Param proposer query string false "filter by proposer"
// @Param proposalID query string false "filter by proposalID"
// @Param state query string false "filter by state"
// @Param sort query string false "newest, minted-newest, token-price-asc, token-price-desc"
// @Param limit query int false "limit default 10"
// @Param page query int false "page start with 1"
// @Success 200 {object} response.JsonResponse{}
// @Router /dao/proposals [GET]
func (h *httpDelivery) proposals(w http.ResponseWriter, r *http.Request) {

	baseF, err := h.BaseFilters(r)
	if err != nil {
		h.Logger.Error("BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	f := structure.FilterProposal{}
	f.BaseFilters = *baseF

	proposer := r.URL.Query().Get("proposer")
	proposalID := r.URL.Query().Get("proposalID")
	state := r.URL.Query().Get("state")

	if proposer != "" {
		f.Proposer = &proposer
	}
if proposalID != "" {
		f.ProposalID = &proposalID
	}
if state != "" {
		stateINT, err := strconv.Atoi(state)
		if err ==   nil {
			f.State = &stateINT
		}
	}

	uProposals, err := h.Usecase.GetProposals(f)
	if err != nil {
		h.Logger.Error("h.Usecase.GetProjects", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	pResp :=  []response.ProposalResp{}
	iPro := uProposals.Result
	pro := iPro.([]entity.Proposal)
	for _, proItem := range pro {

		p, err := h.proposalToResp(&proItem)
		if err != nil {
			h.Logger.Error("copier.Copy", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
			return
		}

		pResp = append(pResp, *p)
	}


	
	h.Response.RespondSuccess(w, http.StatusOK, response.Success,  h.PaginationResp(uProposals, pResp), "")
}

// UserCredits godoc
// @Summary DAO proposal's detail
// @Description DAO proposal's detail
// @Tags DAO
// @Accept  json
// @Produce  json
// @Param proposalID path string true "proposalID: the onchain ID"
// @Success 200 {object} response.JsonResponse{}
// @Router /dao/proposals/{proposalID} [GET]
func (h *httpDelivery) getProposal(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	proposalID := vars["proposalID"]
	

	proposal, err := h.Usecase.GetProposal(proposalID)
	if err != nil {
		h.Logger.Error("h.Usecase.GetProposal", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	resp, err := h.proposalToResp(proposal)
	if err != nil {
		h.Logger.Error(" h.proposalToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}
h.Logger.Info("resp.Proposal", resp)
	
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp , "")
}

// UserCredits godoc
// @Summary DAO proposal's votes
// @Description DAO proposal's detail
// @Tags DAO
// @Accept  json
// @Produce  json
// @Param proposalID path string true "proposalID: the onchain ID"
// @Param voter query string false "filter by voter"
// @Param support query string false "filter by support"
// @Param sort query string false "newest, minted-newest, token-price-asc, token-price-desc"
// @Param limit query int false "limit default 10"
// @Param page query int false "page start with 1"
// @Success 200 {object} response.JsonResponse{}
// @Router /dao/proposals/{proposalID}/votes [GET]
func (h *httpDelivery) getProposalVotes(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	proposalID := vars["proposalID"]
	

	baseF, err := h.BaseFilters(r)
	if err != nil {
		h.Logger.Error("BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	f := structure.FilterProposalVote{}
	f.BaseFilters = *baseF

	f.ProposalID = &proposalID
	support := r.URL.Query().Get("support")
	if support != "" {
		supportInt, err := strconv.Atoi(support)
		if err != nil {
			h.Logger.Error("strconv.Atoi", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
			return
		}
		f.Support = &supportInt
	}

	voter := r.URL.Query().Get("voter")
	if voter != "" {
		f.Voter = &voter
	}


	paginationData, err := h.Usecase.GetProposalVotes(f)
	if err != nil {
		h.Logger.Error("h.Usecase.GetProposal", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	pResp :=  []response.ProposalVotesResp{}
	iPro := paginationData.Result
	pro := iPro.([]entity.ProposalVotes)
	for _, proItem := range pro {
		tmp := &response.ProposalVotesResp{}
		err := response.CopyEntityToRes(tmp, &proItem)
		if err != nil {
			h.Logger.Error("copier.Copy", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
			return
		}

		pResp = append(pResp, *tmp)
	}
//h.Logger.Info("resp.Proposal", resp)
	
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, h.PaginationResp(paginationData, pResp) , "")
}

// UserCredits godoc
// @Summary DAO create a draft proposal
// @Description DAO create a draft proposal
// @Tags DAO
// @Accept  json
// @Produce  json
// @Param request body request.CreateProposalReq true "Create a draft proposal request"
// @Success 200 {object} response.JsonResponse{}
// @Router /dao/proposals [POST]
func (h *httpDelivery) createDraftProposals(w http.ResponseWriter, r *http.Request) {

	var reqBody request.CreateProposalReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		h.Logger.Error("decoder.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase := &structure.CreateProposaltReq{}
	err = copier.Copy(reqUsecase, reqBody)
	if err != nil {
		h.Logger.Error("copier.Copy", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	uProposals, err := h.Usecase.CreateDraftProposal(*reqUsecase)
	if err != nil {
		h.Logger.Error("h.Usecase.GetProjects", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	h.Logger.Info("uProposals", uProposals)
	resp, err := h.proposalDetailToResp(uProposals)
	if err != nil {
		h.Logger.Error(" h.proposalToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}


// UserCredits godoc
// @Summary DAO map off and onchain proposal
// @Description DAO off and onchain proposal
// @Tags DAO
// @Accept  json
// @Produce  json
// @Param ID path string true "ID: the offChain ID"
// @Param proposalID path string true "proposalID: the onchain ID"
// @Success 200 {object} response.JsonResponse{}
// @Router /dao/proposals/{ID}/{proposalID} [PUT]
func (h *httpDelivery) mapOffAndOnChainProposal(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	iD := vars["ID"]
	

	proposalID := vars["proposalID"]
	

	uProposals, err := h.Usecase.MapOffToOnChainProposal(iD, proposalID)
	if err != nil {
		h.Logger.Error("h.Usecase.MapOffToOnChainProposal", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	h.Logger.Info("uProposals", uProposals)
	resp, err := h.proposalDetailToResp(uProposals)
	if err != nil {
		h.Logger.Error(" h.proposalToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

func (h *httpDelivery) proposalToResp(input *entity.Proposal) (*response.ProposalResp, error) {
	resp := &response.ProposalResp{}
	err := response.CopyEntityToRes(resp, input)
	if err != nil {
		return nil, err
	}

	resp.Amount = input.ProposalDetail.Amount
	if resp.Title  == "" && input.ProposalDetail.Title != "" {
		resp.Title = input.ProposalDetail.Title
	}
resp.Description = input.ProposalDetail.Description
	resp.TokenType = input.ProposalDetail.TokenType
	resp.ReceiverAddress = input.ProposalDetail.ReceiverAddress
	return resp, nil
}

func (h *httpDelivery) proposalDetailToResp(input *entity.ProposalDetail) (*response.ProposalResp, error) {
	resp := &response.ProposalResp{}
	err := response.CopyEntityToRes(resp, input)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
