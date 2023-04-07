package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
)

// Discord godoc
// @Summary send new bid notification
// @Description send new bid notification
// @Tags Discord
// @Content-Type: application/json
// @Param request body request.SendNewBidNotifyRequest true "Data for minify"
// @Success 200 {object} response.JsonResponse
// @Security Authorization
// @Router /discord/new-bid [POST]
func (h *httpDelivery) sendDiscordNewBid(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, -1, fmt.Errorf("invalid request body"))
		return
	}

	var req request.SendNewBidNotifyRequest
	if err := json.Unmarshal(body, &req); err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, -1, fmt.Errorf("invalid request body"))
		return
	}

	if req.Quantity == 0 {
		req.Quantity = 1
	}

	err = h.Usecase.NotifyNewBid(req.WalletAddress, req.BidPrice, req.Quantity, req.CollectorRedirectTo)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, -1, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, "", "")
}
