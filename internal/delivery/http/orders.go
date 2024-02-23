package http

import (
	"encoding/json"
	"github.com/jinzhu/copier"
	"net/http"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/usecase/structure"
)

// UserCredits godoc
// @Summary BTC Generate receive wallet address
// @Description Generate receive wallet address
// @Tags Orders
// @Accept  json
// @Produce  json
// @Param request body request.CreateBtcWalletAddressReq true "Create a btc wallet address request"
// @Success 200 {object} response.JsonResponse{}
// @Router /orders/receive-address [POST]
func (h *httpDelivery) ordersGetReceiveWalletAddress(w http.ResponseWriter, r *http.Request) {
	var reqBody request.CreateOrderReceiveAddressReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase := &structure.OrderBtcData{}
	err = copier.Copy(reqUsecase, reqBody)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp, err := h.Usecase.CreateOrderReceiveAddress(*reqUsecase)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	//resp := h.MintNftBtcToResp(mintNftBtcWallet)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}
