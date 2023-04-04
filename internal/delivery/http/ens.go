package http

import (
	"net/http"

	"rederinghub.io/internal/delivery/http/response"
)

func (h *httpDelivery) getListAuction(w http.ResponseWriter, r *http.Request) {

	list, _ := h.Usecase.APIGetListAuction()
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, list, "")
}
