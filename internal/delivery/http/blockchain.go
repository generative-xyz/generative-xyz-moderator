package http

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/usecase/structure"
)

// UserCredits godoc
// @Summary get nft transactions
// @Description get nft transactions
// @Tags Blockchain
// @Accept  json
// @Produce  json
// @Param contractAddress path string true "contract address"
// @Param tokenID path string true "token ID"
// @Success 200 {object} response.JsonResponse{}
// @Router /nfts/{contractAddress}/transactions/{tokenID} [GET]
func (h *httpDelivery) getNftTransactions(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("messages.nft", r)
	defer h.Tracer.FinishSpan(span, log )

	vars := mux.Vars(r)
	contractAddress := vars["contractAddress"]
	span.SetTag("contractAddress", contractAddress)
	tokenID := vars["tokenID"]
	span.SetTag("tokenID", tokenID)

	covalentResp, err := h.Usecase.GetNftTransactions(span, structure.GetNftTransactionsReq{
		ContractAddress: contractAddress,
		TokenID: tokenID,
	})

	if err != nil {
		log.Error("getNftTransactions", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp := covalentResp.Data

	log.SetData("resp", resp);
	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")	
}

// UserCredits godoc
// @Summary get token holder
// @Description get token holder
// @Tags Blockchain
// @Accept  json
// @Produce  json
// @Param contractAddress path string true "contract address"
// @Param page query string false "page"
// @Param limit query string false "limit"
// @Success 200 {object} response.JsonResponse{}
// @Router /nfts/{contractAddress}/nft_holders [GET]
func (h *httpDelivery) getTokenHolder(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("messages.token-holder", r)
	defer h.Tracer.FinishSpan(span, log )

	vars := mux.Vars(r)
	contractAddress := vars["contractAddress"]
	span.SetTag("contractAddress", contractAddress)
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		log.Error("parse page param to int", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	span.SetTag("page", page)
	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Error("parse limit param to int", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp, err := h.Usecase.GetTokenHolders(span, structure.GetTokenHolderRequest{
		ContractAddress: contractAddress,
		Page: int32(page),
		Limit: int32(limit),
	})

	if err != nil {
		log.Error("parse limit param to int", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	pagResp := h.PaginationResp(resp, resp.Result)

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, pagResp, "")
}
