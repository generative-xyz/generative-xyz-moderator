package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
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
// @Param sort query string false "newest, minted-newest, token-price-asc, token-price-desc"
// @Param limit query int false "limit default 10"
// @Param page query int false "page start with 1"
// @Success 200 {object} response.JsonResponse{}
// @Router /dao/proposals [GET]
func (h *httpDelivery) proposals(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("httpDelivery.proposals", r)
	defer h.Tracer.FinishSpan(span, log )

	baseF, err := h.BaseFilters(r)
	if err != nil {
		log.Error("BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	f := structure.FilterProposal{}
	f.BaseFilters = *baseF

	uProposals, err := h.Usecase.GetProposals(span, f)
	if err != nil {
		log.Error("h.Usecase.GetProjects", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	pResp :=  []response.ProposalResp{}
	iPro := uProposals.Result
	pro := iPro.([]entity.Proposal)
	for _, proItem := range pro {

		p, err := h.proposalToResp(&proItem)
		if err != nil {
			log.Error("copier.Copy", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
			return
		}

		pResp = append(pResp, *p)
	}


	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success,  h.PaginationResp(uProposals, pResp), "")
}

// UserCredits godoc
// @Summary DAO proposal's detail
// @Description DAO proposal's detail
// @Tags DAO
// @Accept  json
// @Produce  json
// @Param proposalID path string true "proposalID"
// @Success 200 {object} response.JsonResponse{}
// @Router /dao/proposals/{proposalID} [GET]
func (h *httpDelivery) getProposal(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("getProposal", r)
	defer h.Tracer.FinishSpan(span, log )

	vars := mux.Vars(r)
	proposalID := vars["proposalID"]
	span.SetTag("proposalID", proposalID)

	proposal, err := h.Usecase.GetProposal(span, proposalID)
	if err != nil {
		log.Error("h.Usecase.GetProposal", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}

	resp, err := h.proposalToResp(proposal)
	if err != nil {
		log.Error(" h.proposalToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest,response.Error, err)
		return
	}
	
	log.SetData("resp.Proposal", resp)
	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp , "")
}


func (h *httpDelivery) proposalToResp(input *entity.Proposal) (*response.ProposalResp, error) {
	resp := &response.ProposalResp{}
	err := copier.Copy(resp, input)
	if err != nil {
		return nil, err
	}
	
	return resp, nil
}
