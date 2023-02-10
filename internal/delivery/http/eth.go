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
// @Summary ETH Generate receive wallet address
// @Description Generate receive wallet address
// @Tags ETH
// @Accept  json
// @Produce  json
// @Param request body request.CreateEthWalletAddressReq true "Create a eth wallet address request"
// @Success 200 {object} response.JsonResponse{}
// @Router /eth/receive-address [POST]
func (h *httpDelivery) ethGetReceiveWalletAddress(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("httpDelivery.ethGetReceiveWalletAddress", r)
	defer h.Tracer.FinishSpan(span, log)
	h.Response.SetLog(h.Tracer, span)

	var reqBody request.CreateBtcWalletAddressReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Error("httpDelivery.btcMint.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase := &structure.BctWalletAddressData{}
	err = copier.Copy(reqUsecase, reqBody)
	if err != nil {
		log.Error("copier.Copy", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	btcWallet, err := h.Usecase.CreateBTCWalletAddress(span, *reqUsecase)
	if err != nil {
		log.Error("h.Usecase.CreateBTCWalletAddress", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	log.SetData("btcWallet", btcWallet)
	resp, err := h.BtcWalletAddressToResp(btcWallet)
	if err != nil {
		log.Error(" h.proposalToResp", err.Error(), err)
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
// @Param request body request.CreateMintReq true "Mint request via ORD_WALLET_ADDRESS"
// @Success 200 {object} response.JsonResponse{}
// @Router /btc/mint [POST]
func (h *httpDelivery) mintETH(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("httpDelivery.mint", r)
	defer h.Tracer.FinishSpan(span, log)
	h.Response.SetLog(h.Tracer, span)

	var reqBody request.CreateMintReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Error("httpDelivery.btcMint.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase := &structure.BctMintData{}
	err = copier.Copy(reqUsecase, reqBody)
	if err != nil {
		log.Error("copier.Copy", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	btcWallet, err := h.Usecase.BTCMint(span, *reqUsecase)
	if err != nil {
		log.Error("h.Usecase.CreateBTCWalletAddress", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	log.SetData("btcWallet", btcWallet)
	resp, err := h.BtcToResp(btcWallet)
	if err != nil {
		log.Error(" h.proposalToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

func (h *httpDelivery) EthWalletAddressToResp(input *entity.BTCWalletAddress) (*response.BctReceiveWalletResp, error) {
	resp := &response.BctReceiveWalletResp{}
	resp.Address = input.OrdAddress
	resp.Pricce = input.Amount
	return resp, nil
}

func (h *httpDelivery) EthToResp(input *entity.BTCWalletAddress) (*response.BctWalletResp, error) {
	resp := &response.BctWalletResp{}
	err := copier.Copy(resp, input)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
