package http

import (
	"net/http"

	"rederinghub.io/internal/delivery/http/response"
)

func (h *httpDelivery) getListAuction(w http.ResponseWriter, r *http.Request) {

	list, _ := h.Usecase.APIGetListAuction()
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, list, "")
}

func (h *httpDelivery) checkDeclared(w http.ResponseWriter, r *http.Request) {

	flag := h.Usecase.APIAuctionCheckDeclared()
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, flag, "")
}

func (h *httpDelivery) listSnapshot(w http.ResponseWriter, r *http.Request) {

	list := h.Usecase.APIAuctionListSnapshot()
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, list, "")
}

func (h *httpDelivery) listBidWinnerConfig(w http.ResponseWriter, r *http.Request) {

	list, _ := h.Usecase.GetAuctionListWinnerAddressFromConfig()
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, list, "")
}

func (h *httpDelivery) listBidWinnerBigList(w http.ResponseWriter, r *http.Request) {

	list, _ := h.Usecase.GetAuctionListWinnerAddressFromBidList()
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, list, "")
}
