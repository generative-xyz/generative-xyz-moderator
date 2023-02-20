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
// @Summary BTC/ETH Generate receive wallet address
// @Description Generate receive wallet address
// @Tags BTC/ETH
// @Accept  json
// @Produce  json
// @Param request body request.CreateBtcWalletAddressReq true "Create a btc/eth wallet address request"
// @Success 200 {object} response.JsonResponse{}
// @Router /mint-nft-btc/receive-address [POST]
func (h *httpDelivery) createMintReceiveAddress(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("httpDelivery.createMintReceiveAddress", r)
	defer h.Tracer.FinishSpan(span, log)
	h.Response.SetLog(h.Tracer, span)

	var reqBody request.CreateMintReceiveAddressReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Error("httpDelivery.MintNftBtc.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase := &structure.MintNftBtcData{}
	err = copier.Copy(reqUsecase, reqBody)
	if err != nil {
		log.Error("copier.Copy", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	mintNftBtcWallet, err := h.Usecase.CreateMintReceiveAddress(span, *reqUsecase)
	if err != nil {
		log.Error("h.Usecase.createMintReceiveAddress", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	log.SetData("mintNftBtcWallet", mintNftBtcWallet)
	resp := h.MintNftBtcToResp(mintNftBtcWallet)

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
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
