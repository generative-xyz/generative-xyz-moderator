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
// @Param request body request.CreateBtcWalletAddressReqV2 true "Create a btc wallet address request"
// @Success 200 {object} response.JsonResponse{}
// @Router /btc-v2/receive-address [POST]
func (h *httpDelivery) btcGetReceiveWalletAddressV2(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("httpDelivery.btcGetReceiveWalletAddress", r)
	defer h.Tracer.FinishSpan(span, log)
	h.Response.SetLog(h.Tracer, span)

	var reqBody request.CreateBtcWalletAddressReqV2
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Error("httpDelivery.btcMint.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase := &structure.BctWalletAddressDataV2{}
	err = copier.Copy(reqUsecase, reqBody)
	if err != nil {
		log.Error("copier.Copy", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	btcWallet, err := h.Usecase.CreateBTCWalletAddressV2(span, *reqUsecase)
	if err != nil {
		log.Error("h.Usecase.CreateOrdBTCWalletAddress", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	log.SetData("btcWallet", btcWallet)
	resp, err := h.BtcWalletAddressToRespV2(btcWallet)
	if err != nil {
		log.Error(" h.proposalToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

func (h *httpDelivery) BtcWalletAddressToRespV2(input *entity.BTCWalletAddressV2) (*response.BtcWalletRespV2, error) {
	resp := &response.BtcWalletRespV2{}
	resp.UserAddress = input.UserAddress
	resp.Amount = input.Amount
	resp.MintFee = input.MintFee
	resp.SentTokenFee = input.SentTokenFee
	resp.OrdAddress = input.OrdAddress
	resp.FileURI = input.FileURI
	resp.IsConfirm = input.IsConfirm
	resp.InscriptionID = input.InscriptionID
	resp.Balance = input.Balance
	return resp, nil
}
