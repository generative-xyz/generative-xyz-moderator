package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/utils"
)

func (h *httpDelivery) dexBTCOffering(w http.ResponseWriter, r *http.Request) {
	var reqBody request.CreateDexBTCListing
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		h.Logger.Error("httpDelivery.dexBTCListing.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	var ok bool
	ctx := r.Context()
	iUserID := ctx.Value(utils.SIGNED_USER_ID)
	userID, ok := iUserID.(string)
	if !ok {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("address or accessToken cannot be empty"))
		return
	}
	userInfo, err := h.Usecase.UserProfile(userID)
	if err != nil {
		h.Logger.Error("httpDelivery.mintStatus.Usecase.UserProfile", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	address := userInfo.WalletAddressBTCTaproot
	listing, err := h.Usecase.DexBTCListing(address, reqBody.RawPSBT, reqBody.InscriptionID, reqBody.SplitTx)
	if err != nil {
		h.Logger.Error("httpDelivery.dexBTCListing.Usecase.DexBTCListing", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	// Discord Notify NEW LISTING
	if listing != nil {
		go h.Usecase.NotifyNewListing(*listing)
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, "ok", "")
}

func (h *httpDelivery) cancelBTCOffer(w http.ResponseWriter, r *http.Request) {
	var reqBody request.CancelDexBTCListing
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		h.Logger.Error("httpDelivery.dexBTCListing.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	var ok bool
	ctx := r.Context()
	iUserID := ctx.Value(utils.SIGNED_USER_ID)
	userID, ok := iUserID.(string)
	if !ok {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("address or accessToken cannot be empty"))
		return
	}
	userInfo, err := h.Usecase.UserProfile(userID)
	if err != nil {
		h.Logger.Error("httpDelivery.mintStatus.Usecase.UserProfile", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	address := userInfo.WalletAddressBTCTaproot
	//find order by inscription_id and user_address
	err = h.Usecase.CancelDexBTCListing(reqBody.Txhash, address, reqBody.InscriptionID, reqBody.OrderID)
	if err != nil {
		h.Logger.Error("httpDelivery.dexBTCListing.Usecase.DexBTCListing", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, "ok", "")
}

func (h *httpDelivery) historyBTCOffer(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 20
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}
	var ok bool
	ctx := r.Context()
	iUserID := ctx.Value(utils.SIGNED_USER_ID)
	userID, ok := iUserID.(string)
	if !ok {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("address or accessToken cannot be empty"))
		return
	}
	userInfo, err := h.Usecase.UserProfile(userID)
	if err != nil {
		h.Logger.Error("httpDelivery.mintStatus.Usecase.UserProfile", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	address := userInfo.WalletAddressBTCTaproot
	listingList, err := h.Usecase.Repo.GetDexBTCListingOrderUser(address, int64(limit), int64(offset))
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("get order info failed"))
		return
	}
	result := []response.DexBTCHistoryListing{}
	for _, listing := range listingList {
		newHistory := response.DexBTCHistoryListing{
			OrderID:       listing.UUID,
			InscriptionID: listing.InscriptionID,
			Timestamp:     listing.CreatedAt.Unix(),
			Amount:        fmt.Sprintf("%v", listing.Amount),
			Type:          "listing",
		}
		result = append(result, newHistory)
		if listing.CancelTx != "" && !listing.Matched {
			newHistory := response.DexBTCHistoryListing{
				OrderID:       listing.UUID,
				InscriptionID: listing.InscriptionID,
				Timestamp:     listing.CancelAt.Unix(),
				Amount:        fmt.Sprintf("%v", listing.Amount),
				Type:          "cancelling",
				Txhash:        listing.CancelTx,
			}
			if listing.Cancelled {
				newHistory.Type = "cancelled"
			}
			result = append(result, newHistory)
		}
		if listing.Matched {
			if listing.Buyer == listing.SellerAddress {
				newHistory := response.DexBTCHistoryListing{
					OrderID:       listing.UUID,
					InscriptionID: listing.InscriptionID,
					Timestamp:     listing.MatchAt.Unix(),
					Amount:        fmt.Sprintf("%v", listing.Amount),
					Type:          "cancelled",
					Txhash:        listing.MatchedTx,
				}
				result = append(result, newHistory)
			} else {
				newHistory := response.DexBTCHistoryListing{
					OrderID:       listing.UUID,
					InscriptionID: listing.InscriptionID,
					Timestamp:     listing.MatchAt.Unix(),
					Amount:        fmt.Sprintf("%v", listing.Amount),
					Type:          "matched",
					Txhash:        listing.MatchedTx,
				}
				result = append(result, newHistory)
			}
		}
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")

}
