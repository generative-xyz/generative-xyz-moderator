package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
)

func (h *httpDelivery) dexBTCListing(w http.ResponseWriter, r *http.Request) {
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
	err = h.Usecase.DexBTCListing(address, reqBody.RawPSBT, reqBody.InscriptionID)
	if err != nil {
		h.Logger.Error("httpDelivery.dexBTCListing.Usecase.DexBTCListing", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, "ok", "")
}

func (h *httpDelivery) cancelBTCListing(w http.ResponseWriter, r *http.Request) {
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

func (h *httpDelivery) retrieveBTCListingOrderInfo(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("order_id")
	if orderID == "" {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("orderID cannot be empty"))
		return
	}
	orderInfo, err := h.Usecase.Repo.GetDexBTCListingOrderByID(orderID)
	if err != nil {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("get order info failed"))
		return
	}
	result := response.DexBTCListingOrderInfo{
		RawPSBT: orderInfo.RawPSBT,
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")
}

func (h *httpDelivery) historyBTCListing(w http.ResponseWriter, r *http.Request) {
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

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")

}

func (h *httpDelivery) dexBTCListingFee(w http.ResponseWriter, r *http.Request) {

	var reqBody request.ListingFee
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		h.Logger.Error("httpDelivery.dexBTCListingFee.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	inscriptionID := reqBody.InscriptionID

	tokenUri, err := h.Usecase.GetTokenByTokenID(inscriptionID, 0)
	if err != nil {
		resp := response.ListingFee{
			ServiceFee:     fmt.Sprintf("%v", utils.BUY_NFT_CHARGE),
			RoyaltyFee:     fmt.Sprintf("%v", 0),
			ServiceAddress: h.Config.MarketBTCServiceFeeAddress,
		}
		h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
		return
	}

	projectDetail, err := h.Usecase.GetProjectDetail(structure.GetProjectDetailMessageReq{
		ContractAddress: tokenUri.ContractAddress,
		ProjectID:       tokenUri.ProjectID,
	})
	if err != nil {
		resp := response.ListingFee{
			ServiceFee:     fmt.Sprintf("%v", utils.BUY_NFT_CHARGE),
			RoyaltyFee:     fmt.Sprintf("%v", 0),
			ServiceAddress: h.Config.MarketBTCServiceFeeAddress,
			ProjectID:      tokenUri.ProjectID,
		}
		h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
		return
	}

	creator, err := h.Usecase.GetUserProfileByWalletAddress(projectDetail.CreatorAddrr)
	if err != nil {
		resp := response.ListingFee{
			ServiceFee:     fmt.Sprintf("%v", utils.BUY_NFT_CHARGE),
			RoyaltyFee:     fmt.Sprintf("%v", 0),
			ProjectID:      tokenUri.ProjectID,
			ServiceAddress: h.Config.MarketBTCServiceFeeAddress,
		}
		h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
		return
	}
	resp := response.ListingFee{
		ServiceFee:     fmt.Sprintf("%v", utils.BUY_NFT_CHARGE),
		RoyaltyFee:     fmt.Sprintf("%v", float64(projectDetail.Royalty)/10000*100),
		RoyaltyAddress: creator.WalletAddressBTCTaproot,
		ServiceAddress: h.Config.MarketBTCServiceFeeAddress,
		ProjectID:      tokenUri.ProjectID,
	}
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}
