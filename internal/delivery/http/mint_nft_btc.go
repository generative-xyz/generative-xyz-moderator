package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
)

// UserCredits godoc
// @Summary BTC/ETH Generate receive wallet address
// @Description Generate receive wallet address
// @Tags BTC/ETH
// @Accept  json
// @Produce  json
// @Param request body request.CreateBtcWalletAddressReq true "Create a btc/eth wallet address request"
// @Success 200 {object} response.JsonResponse{}
// @Router /mint-nft-btc/receive-address [POST]
func (h *httpDelivery) createMintReceiveAddress(w http.ResponseWriter, r *http.Request) {

	// verify user:
	ctx := r.Context()
	iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)

	fmt.Println("iWalletAddress", iWalletAddress)

	userWalletAddr, ok := iWalletAddress.(string)
	if !ok {
		err := errors.New("wallet address is incorect")
		h.Logger.Error("ctx.Value.Token", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	fmt.Println("userWalletAddr", userWalletAddr)

	profile, err := h.Usecase.GetUserProfileByWalletAddress(userWalletAddr)
	if err != nil {
		h.Logger.Error("h.Usecase.GetUserProfileByWalletAddress(", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	var reqBody request.CreateMintReceiveAddressReq
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&reqBody)
	if err != nil {
		h.Logger.Error("httpDelivery.MintNftBtc.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	// if !strings.EqualFold(profile.WalletAddressBTCTaproot, reqBody.WalletAddress) {
	// 	err = errors.New("permission dined")
	// 	h.Logger.Error("h.Usecase.createMintReceiveAddress", err.Error(), err)
	// 	h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
	// 	return
	// }

	reqUsecase := &structure.MintNftBtcData{}
	err = copier.Copy(reqUsecase, reqBody)
	if err != nil {
		h.Logger.Error("copier.Copy", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase.UserID = profile.UUID

	mintNftBtcWallet, err := h.Usecase.CreateMintReceiveAddress(*reqUsecase)
	if err != nil {
		h.Logger.Error("h.Usecase.createMintReceiveAddress", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp := h.MintNftBtcToResp(mintNftBtcWallet)

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

// UserCredits godoc
// @Summary cancel the mint request
// @Description cancel the mint request
// @Tags BTC/ETH
// @Accept  json
// @Produce  json
// @Param request body request.CreateBtcWalletAddressReq true "Create a btc/eth wallet address request"
// @Success 200 {object} response.JsonResponse{}
// @Router /mint-nft-btc/receive-address [DELETE]
func (h *httpDelivery) cancelMintNftBt(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)

	fmt.Println("iWalletAddress", iWalletAddress)

	userWalletAddr, ok := iWalletAddress.(string)
	if !ok {
		err := errors.New("wallet address is incorect")
		h.Logger.Error("ctx.Value.Token", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	fmt.Println("userWalletAddr", userWalletAddr)

	profile, err := h.Usecase.GetUserProfileByWalletAddress(userWalletAddr)
	if err != nil {
		h.Logger.Error("h.Usecase.GetUserProfileByWalletAddress(", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	vars := mux.Vars(r)
	uuid := vars["uuid"]

	err = h.Usecase.CancelMintNftBtc(profile.WalletAddressBTCTaproot, uuid)
	if err != nil {
		h.Logger.Error("h.Usecase.CancelMintNftBt", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, true, "")
}

func (h *httpDelivery) getDetailMintNftBtc(w http.ResponseWriter, r *http.Request) {

	// verify user:
	ctx := r.Context()
	iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)

	fmt.Println("iWalletAddress", iWalletAddress)

	userWalletAddr, ok := iWalletAddress.(string)
	if !ok {
		err := errors.New("wallet address is incorect")
		h.Logger.Error("ctx.Value.Token", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	fmt.Println("userWalletAddr", userWalletAddr)

	profile, err := h.Usecase.GetUserProfileByWalletAddress(userWalletAddr)
	if err != nil {
		h.Logger.Error("h.Usecase.GetUserProfileByWalletAddress(", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	vars := mux.Vars(r)
	uuid := vars["uuid"]

	item, err := h.Usecase.GetDetalMintNftBtc(uuid)
	if err != nil {
		h.Logger.Error("h.Usecase.CancelMintNftBt", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	// if !strings.EqualFold(item.OriginUserAddress, profile.WalletAddressBTCTaproot) {
	// 	err := errors.New("permission dined")
	// 	h.Logger.Error("ctx.Value.Token", err.Error(), err)
	// 	h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
	// 	return
	// }
	if item.UserID != profile.UUID {
		err := errors.New("permission dined")
		h.Logger.Error("ctx.Value.Token", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, item, "")
}

func (h *httpDelivery) MintNftBtcToResp(input *entity.MintNftBtc) *response.MintNftBtcReceiveWalletResp {
	return &response.MintNftBtcReceiveWalletResp{
		Address: input.ReceiveAddress,
		Price:   input.Amount,
		PayType: input.PayType,
	}

}

func (h *httpDelivery) MintNftBtcResp(input *entity.MintNftBtc) (*response.MintNftBtcResp, error) {
	resp := &response.MintNftBtcResp{}
	err := copier.Copy(resp, input)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
