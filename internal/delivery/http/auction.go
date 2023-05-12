package http

import (
	"errors"
	"net/http"

	"go.uber.org/zap"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/utils"
	"rederinghub.io/utils/logger"
)

func (h *httpDelivery) getListAuction(w http.ResponseWriter, r *http.Request) {

	// list, _ := h.Usecase.APIGetListAuction()

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, nil, "")
}

func (h *httpDelivery) checkDeclared(w http.ResponseWriter, r *http.Request) {

	flag := h.Usecase.APIAuctionCheckDeclared()
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, flag, "")
}

func (h *httpDelivery) listSnapshot(w http.ResponseWriter, r *http.Request) {

	// list := h.Usecase.APIAuctionListSnapshot()
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, nil, "")
}

func (h *httpDelivery) listBidWinnerConfig(w http.ResponseWriter, r *http.Request) {

	// list, _ := h.Usecase.GetAuctionListWinnerAddressFromConfig()
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, nil, "")
}

func (h *httpDelivery) listBidWinnerBigList(w http.ResponseWriter, r *http.Request) {

	// list, _ := h.Usecase.GetAuctionListWinnerAddressFromBidList()
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, nil, "")
}

func (h *httpDelivery) shareNow(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)

	userWalletAddr, ok := iWalletAddress.(string)
	if !ok {
		err := errors.New("Wallet address is incorrect")
		logger.AtLog.Logger.Error("ctx.Value.Token", zap.Error(err))
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Usecase.APIShareNow(userWalletAddr)

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, true, "")
}
