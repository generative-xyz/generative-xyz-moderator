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

// UserCredits godoc
// @Summary ETH mint
// @Description mint
// @Tags ETH
// @Accept  json
// @Produce  json
// @Param request body request.CreateMintReq true "Mint request via ORD_WALLET_ADDRESS"
// @Success 200 {object} response.JsonResponse{}
// @Router /eth/mint [POST]
// func (h *httpDelivery) mintETH(w http.ResponseWriter, r *http.Request) {
// 	span, log := h.StartSpan("httpDelivery.mintEth", r)
// 	defer h.Tracer.FinishSpan(span, log)
// 	h.Response.SetLog(h.Tracer, span)

// 	var reqBody request.CreateMintReq
// 	decoder := json.NewDecoder(r.Body)
// 	err := decoder.Decode(&reqBody)
// 	if err != nil {
// 		log.Error("httpDelivery.btcMint.Decode", err.Error(), err)
// 		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
// 		return
// 	}

// 	reqUsecase := &structure.EthMintData{}
// 	err = copier.Copy(reqUsecase, reqBody)
// 	if err != nil {
// 		log.Error("copier.Copy", err.Error(), err)
// 		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
// 		return
// 	}

// 	ethWallet, err := h.Usecase.ETHMint(span, *reqUsecase)
// 	if err != nil {
// 		log.Error("h.Usecase.CreateBTCWalletAddress", err.Error(), err)
// 		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
// 		return
// 	}

// 	log.SetData("ethWallet", ethWallet)
// 	resp, err := h.EthToResp(ethWallet)
// 	if err != nil {
// 		log.Error(" h.proposalToResp", err.Error(), err)
// 		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
// 		return
// 	}
// 	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
// }

func (h *httpDelivery) EthWalletAddressToResp(input *entity.ETHWalletAddress) (*response.EthReceiveWalletResp, error) {
	resp := &response.EthReceiveWalletResp{}
	resp.Address = input.OrdAddress
	resp.Pricce = input.Amount
	return resp, nil
}

func (h *httpDelivery) EthToResp(input *entity.ETHWalletAddress) (*response.EthWalletResp, error) {
	resp := &response.EthWalletResp{}
	err := copier.Copy(resp, input)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
