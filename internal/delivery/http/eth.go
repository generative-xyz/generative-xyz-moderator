package http

import (
	"net/http"

	"github.com/jinzhu/copier"

	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
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

	// var reqBody request.CreateEthWalletAddressReq
	// decoder := json.NewDecoder(r.Body)
	// err := decoder.Decode(&reqBody)
	// if err != nil {
	// 	logger.AtLog.Logger.Error("httpDelivery.btcMint.Decode", zap.Error(err))
	// 	h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
	// 	return
	// }

	// reqUsecase := &structure.EthWalletAddressData{}
	// err = copier.Copy(reqUsecase, reqBody)
	// if err != nil {
	// 	logger.AtLog.Logger.Error("copier.Copy", zap.Error(err))
	// 	h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
	// 	return
	// }

	// ethWallet, err := h.Usecase.CreateETHWalletAddress(*reqUsecase)
	// if err != nil {
	// 	logger.AtLog.Logger.Error("h.Usecase.CreateETHWalletAddress", zap.Error(err))
	// 	h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
	// 	return
	// }

	// logger.AtLog.Logger.Info("ethWallet", zap.Any("ethWallet", ethWallet))
	// resp, err := h.EthWalletAddressToResp(ethWallet)
	// if err != nil {
	// 	logger.AtLog.Logger.Error(" h.proposalToResp", zap.Error(err))
	// 	h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
	// 	return
	// }

	// h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, true, "")
}

// UserCredits godoc
// @Summary ETH Generate receive wallet address
// @Description Generate receive wallet address
// @Tags ETH
// @Accept  json
// @Produce  json
// @Security Authorization
// @Param request body request.CreateWhitelistedEthWalletAddressReq true "Create a eth wallet address request"
// @Success 200 {object} response.JsonResponse{}
// @Router /eth/receive-address/whitelist [POST]
func (h *httpDelivery) ethGetReceiveWhitelistedWalletAddress(w http.ResponseWriter, r *http.Request) {

	// ctx := r.Context()
	// iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)
	// userWalletAddr, ok := iWalletAddress.(string)
	// if !ok {
	// 	err := errors.New("Wallet address is incorect")
	// 	logger.AtLog.Logger.Error("ctx.Value.Token", zap.Error(err))
	// 	h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
	// 	return
	// }

	// var reqBody request.CreateWhitelistedEthWalletAddressReq
	// decoder := json.NewDecoder(r.Body)
	// err := decoder.Decode(&reqBody)
	// if err != nil {
	// 	logger.AtLog.Logger.Error("httpDelivery.btcMint.Decode", zap.Error(err))
	// 	h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
	// 	return
	// }

	// reqUsecase := &structure.EthWalletAddressData{}
	// err = copier.Copy(reqUsecase, reqBody)
	// if err != nil {
	// 	logger.AtLog.Logger.Error("copier.Copy", zap.Error(err))
	// 	h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
	// 	return
	// }
	// ethWallet, err := h.Usecase.CreateWhitelistedETHWalletAddress(ctx, userWalletAddr, *reqUsecase)
	// if err != nil {
	// 	logger.AtLog.Logger.Error("h.Usecase.CreateETHWalletAddress", zap.Error(err))
	// 	h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
	// 	return
	// }

	// logger.AtLog.Logger.Info("ethWallet", zap.Any("ethWallet", ethWallet))
	// resp, err := h.EthWalletAddressToResp(ethWallet)
	// if err != nil {
	// 	logger.AtLog.Logger.Error(" h.proposalToResp", zap.Error(err))
	// 	h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
	// 	return
	// }

	// h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, true, "")
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
//
//
//

// 	var reqBody request.CreateMintReq
// 	decoder := json.NewDecoder(r.Body)
// 	err := decoder.Decode(&reqBody)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("httpDelivery.btcMint.Decode", zap.Error(err))
// 		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
// 		return
// 	}

// 	reqUsecase := &structure.EthMintData{}
// 	err = copier.Copy(reqUsecase, reqBody)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("copier.Copy", zap.Error(err))
// 		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
// 		return
// 	}

// 	ethWallet, err := h.Usecase.ETHMint(*reqUsecase)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("h.Usecase.CreateOrdBTCWalletAddress", zap.Error(err))
// 		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
// 		return
// 	}

// 	logger.AtLog.Logger.Info("ethWallet", zap.Any("ethWallet", ethWallet))
// 	resp, err := h.EthToResp(ethWallet)
// 	if err != nil {
// 		logger.AtLog.Logger.Error(" h.proposalToResp", zap.Error(err))
// 		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
// 		return
// 	}
// 	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
// }

func (h *httpDelivery) EthWalletAddressToResp(input *entity.ETHWalletAddress) (*response.EthReceiveWalletResp, error) {
	resp := &response.EthReceiveWalletResp{}
	resp.Address = input.OrdAddress
	resp.Price = input.Amount
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
