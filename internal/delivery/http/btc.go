package http

import (
	"encoding/json"
	"net/http"

	"github.com/jinzhu/copier"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
)

// UserCredits godoc
// @Summary BTC Generate receive wallet address
// @Description Generate receive wallet address
// @Tags BTC
// @Accept  json
// @Produce  json
// @Param request body request.CreateBtcWalletAddressReq true "Create a btc wallet address request"
// @Success 200 {object} response.JsonResponse{}
// @Router /btc/receive-address [POST]
func (h *httpDelivery) btcGetReceiveWalletAddress(w http.ResponseWriter, r *http.Request) {
	var reqBody request.CreateBtcWalletAddressReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		h.Logger.Error("httpDelivery.btcMint.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase := &structure.BctWalletAddressData{}
	err = copier.Copy(reqUsecase, reqBody)
	if err != nil {
		h.Logger.Error("copier.Copy", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	btcWallet, err := h.Usecase.CreateSegwitBTCWalletAddress(*reqUsecase)
	if err != nil {
		h.Logger.Error("h.Usecase.CreateSegwitBTCWalletAddress", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Logger.Info("btcWallet", btcWallet)
	resp, err := h.BtcWalletAddressToResp(btcWallet)
	if err != nil {
		h.Logger.Error(" h.proposalToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

// UserCredits godoc
// @Summary BTC mint
// @Description mint
// @Tags BTC
// @Accept  json
// @Produce  json
// @Param request body request.CheckBalanceAddressReq true "Check balance of ORD_WALLET_ADDRESS"
// @Success 200 {object} response.JsonResponse{}
// @Router /btc/balance [POST]
func (h *httpDelivery) checkBalance(w http.ResponseWriter, r *http.Request) {
	

	var reqBody request.CheckBalanceAddressReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		h.Logger.Error("httpDelivery.checkBalance.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase := &structure.CheckBalance{}
	err = copier.Copy(reqUsecase, reqBody)
	if err != nil {
		h.Logger.Error("copier.Copy", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	btcWallet, err := h.Usecase.CheckBalanceWalletAddress(*reqUsecase)
	if err != nil {
		h.Logger.Error("h.Usecase.CheckBalanceWalletAddress", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Logger.Info("btcWallet", btcWallet)
	resp, err := h.BtcToResp(btcWallet)
	if err != nil {
		h.Logger.Error(" h.proposalToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

func (h *httpDelivery) BtcWalletAddressToResp(input *entity.BTCWalletAddress) (*response.BctReceiveWalletResp, error) {
	resp := &response.BctReceiveWalletResp{}
	resp.Address = input.OrdAddress
	resp.Price = input.Amount
	return resp, nil
}

func (h *httpDelivery) BtcToResp(input *entity.BTCWalletAddress) (*response.BctWalletResp, error) {
	resp := &response.BctWalletResp{}
	err := copier.Copy(resp, input)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
