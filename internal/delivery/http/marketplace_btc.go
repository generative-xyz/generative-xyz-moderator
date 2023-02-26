package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/gorilla/mux"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/btc"
	"rederinghub.io/utils/eth"
)

func (h *httpDelivery) btcMarketplaceListing(w http.ResponseWriter, r *http.Request) {

	var reqBody request.CreateMarketplaceBTCListing
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		h.Logger.Error("httpDelivery.btcMarketplaceListing.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	inscriptionID := reqBody.InscriptionID

	inscriptionIDs := strings.Split(inscriptionID, "https://ordinals.com/inscription/")

	if len(inscriptionIDs) == 2 {
		inscriptionID = inscriptionIDs[1]
	}

	// TODO: check exists:

	// check valid inscriptionID:
	suffix := "i0"
	if !strings.HasSuffix(inscriptionID, suffix) {
		err := fmt.Errorf("invalid inscriptionID")
		h.Logger.Error("httpDelivery.btcMarketplaceListing.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	txHash := strings.TrimSuffix(inscriptionID, suffix)
	_, err = chainhash.NewHashFromStr(txHash)
	if err != nil {
		err := fmt.Errorf("invalid inscriptionID")
		h.Logger.Error("httpDelivery.btcMarketplaceListing.NewHashFromStr", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	// check btc address:
	ok, _ := btc.ValidateAddress("btc", reqBody.OrdWalletAddress)
	if !ok {
		err := fmt.Errorf("invalid ordWalletAddress")
		h.Logger.Error("httpDelivery.btcMarketplaceListing.ValidateAddress", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	// check paytype:
	btcPaymentAddress, okBtc := reqBody.PayType["btc"]
	ethPaymentAddress, okEth := reqBody.PayType["eth"]
	if !okBtc && !okEth {
		err := fmt.Errorf("payment type is requied")
		h.Logger.Error("httpDelivery.btcMarketplaceListing.Validate Payment type", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	if okBtc {
		ok, _ = btc.ValidateAddress("btc", btcPaymentAddress)
		if !ok {
			err := fmt.Errorf("invalid btcPaymentAddress")
			h.Logger.Error("httpDelivery.btcMarketplaceListing.ValidateAddress.btcPaymentAddress", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
	}

	if okEth {
		ok = eth.ValidateAddress(ethPaymentAddress)
		if !ok {
			err := fmt.Errorf("invalid ethPaymentAddress")
			h.Logger.Error("httpDelivery.btcMarketplaceListing.ValidateAddress.ethPaymentAddress", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
	}

	priceNumber, err := strconv.ParseInt(reqBody.Price, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid price")
		h.Logger.Error("httpDelivery.btcMarketplaceListing.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	// check price:
	if priceNumber < utils.MIN_BTC_TO_LIST_BTC {
		err := fmt.Errorf("Minimum price is %.2f BTC", float64(utils.MIN_BTC_TO_LIST_BTC/1e8))
		h.Logger.Error("httpDelivery.btcMarketplaceListing.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase := structure.MarketplaceBTC_ListingInfo{
		InscriptionID: inscriptionID,
		Name:          reqBody.Name,
		Description:   reqBody.Description,

		SellOrdAddress: reqBody.OrdWalletAddress,

		Price:      reqBody.Price,
		ServiceFee: fmt.Sprintf("%v", utils.BUY_NFT_CHARGE),

		PayType: reqBody.PayType,
	}

	nft, err := h.Usecase.Repo.FindBtcNFTListingUnsoldByNFTID(inscriptionID)
	if err == nil {
		if nft != nil {
			err := fmt.Errorf("this inscription is already listed")
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
	}

	listing, err := h.Usecase.BTCMarketplaceListingNFT(reqUsecase)
	if err != nil {
		h.Logger.Error("h.Usecase.BTCMarketplaceListingNFT", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp := response.CreateMarketplaceBTCListing{
		ReceiveAddress: listing.HoldOrdAddress,
		TimeoutAt:      fmt.Sprintf("%d", time.Now().Add(time.Minute*90).Unix()),
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

func (h *httpDelivery) btcMarketplaceListNFTs(w http.ResponseWriter, r *http.Request) {

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 20
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}

	buyableOnly := false
	if r.URL.Query().Get("buyable-only") == "true" {
		buyableOnly = true
	}

	keyword := r.URL.Query().Get("keyword")
	listCollectionIDs := r.URL.Query().Get("listCollectionIDs") // collection id
	listPrices := r.URL.Query().Get("listPrices")               // price
	listIDs := r.URL.Query().Get("listIDs")                     // nft id

	filterObject := &entity.FilterString{
		Keyword:           keyword,
		ListCollectionIDs: listCollectionIDs,
		ListPrices:        listPrices,
		ListIDs:           listIDs,
	}

	result, err := h.Usecase.BTCMarketplaceListNFT(filterObject, buyableOnly, int64(limit), int64(offset))
	if err != nil {
		h.Logger.Error("h.Usecase.BTCMarketplaceListNFT", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")
}

func (h *httpDelivery) btcMarketplaceNFTDetail(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	inscriptionID := vars["ID"]

	var nft *entity.MarketplaceBTCListing
	var err error
	isBuyable := true
	isCompleted := false
	// lastPrice := int64(0)

	nft, err = h.Usecase.Repo.FindBtcNFTListingUnsoldByNFTID(inscriptionID)
	if err != nil {
		isBuyable = false
		nft, err = h.Usecase.Repo.FindBtcNFTListingLastSoldByNFTID(inscriptionID)
		if err != nil {
			h.Logger.Error("h.Usecase.Repo.FindBtcNFTListingByNFTID", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}

		isCompleted = true
	}

	// if !nft.IsSold {
	// 	buyOrders, err := h.Usecase.Repo.GetBTCListingHaveOngoingOrder(nft.UUID)
	// 	if err != nil {
	// 		h.Logger.Error("h.Usecase.Repo.GetBTCListingHaveOngoingOrder", err.Error(), err)
	// 	}
	// 	currentTime := time.Now()
	// 	for _, order := range buyOrders {
	// 		expireTime := order.ExpiredAt
	// 		// not expired yet still waiting for btc
	// 		if currentTime.Before(expireTime) && (order.Status == entity.StatusBuy_Pending || order.Status == entity.StatusBuy_NotEnoughBalance) {
	// 			isBuyable = false
	// 			break
	// 		}
	// 		// could be expired but received btc
	// 		if order.Status != entity.StatusBuy_Pending && order.Status != entity.StatusBuy_NotEnoughBalance {
	// 			isBuyable = false
	// 			break
	// 		}
	// 	}
	// }

	nftInfo := structure.MarketplaceNFTDetail{
		InscriptionID: nft.InscriptionID,
		Name:          nft.Name,
		Description:   nft.Description,
		Price:         nft.Price,
		OrderID:       nft.UUID,
		IsConfirmed:   nft.IsConfirm,
		Buyable:       isBuyable,
		IsCompleted:   isCompleted,
		// LastPrice:     lastPrice,
	}
	inscribeInfo, err := h.Usecase.GetInscribeInfo(nftInfo.InscriptionID)
	if err != nil {
		h.Logger.Error("h.Usecase.GetInscribeInfo", err.Error(), err)
	}
	if inscribeInfo != nil {
		nftInfo.InscriptionNumber = inscribeInfo.Index
		nftInfo.ContentType = inscribeInfo.ContentType
		nftInfo.ContentLength = inscribeInfo.ContentLength
	}
	//h.Logger.Info("resp.Proposal", resp)

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, nftInfo, "")
}

func (h *httpDelivery) btcMarketplaceListingFee(w http.ResponseWriter, r *http.Request) {

	var reqBody request.ListingFee
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		h.Logger.Error("httpDelivery.btcMint.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	inscriptionID := reqBody.InscriptionID

	inscriptionIDs := strings.Split(inscriptionID, "https://ordinals.com/inscription/")

	if len(inscriptionIDs) == 2 {
		inscriptionID = inscriptionIDs[1]
	}

	tokenUri, err := h.Usecase.GetTokenByTokenID(inscriptionID, 0)
	if err != nil {
		resp := response.ListingFee{
			ServiceFee: fmt.Sprintf("%v", utils.BUY_NFT_CHARGE),
			RoyaltyFee: fmt.Sprintf("%v", 0),
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
			ServiceFee: fmt.Sprintf("%v", utils.BUY_NFT_CHARGE),
			RoyaltyFee: fmt.Sprintf("%v", 0),
		}
		h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
		return
	}

	resp := response.ListingFee{
		ServiceFee: fmt.Sprintf("%v", utils.BUY_NFT_CHARGE),
		RoyaltyFee: fmt.Sprintf("%v", float64(projectDetail.Royalty)/10000*100),
	}
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

func (h *httpDelivery) btcMarketplaceCreateBuyOrder(w http.ResponseWriter, r *http.Request) {

	var reqBody request.CreateMarketplaceBTCBuyOrder
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		h.Logger.Error("httpDelivery.btcMint.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	// check btc address:
	ok, _ := btc.ValidateAddress("btc", reqBody.WalletAddress)
	if !ok {
		err := fmt.Errorf("invalid WalletAddress")
		h.Logger.Error("httpDelivery.btcMarketplaceListing.WalletAddress", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase := structure.MarketplaceBTC_BuyOrderInfo{
		InscriptionID: reqBody.InscriptionID,
		OrderID:       reqBody.OrderID,
		BuyOrdAddress: reqBody.WalletAddress,
	}
	//TODO: lam uncomment
	listing, err := h.Usecase.Repo.FindBtcNFTListingByOrderID(reqBody.OrderID)
	if err != nil {
		h.Logger.Error("h.Usecase.BTCMarketplaceListingNFT", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("Inscription not available to buy"))
		return
	}
	if !listing.IsConfirm {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("Inscription not available to buy"))
		return
	}

	if listing.IsSold {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("Inscription not available to buy"))
		return
	}

	depositAddress, err := h.Usecase.BTCMarketplaceBuyOrder(reqUsecase)
	if err != nil {
		h.Logger.Error("h.Usecase.BTCMarketplaceListingNFT", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp := response.CreateMarketplaceBTCBuyOrder{
		ReceiveAddress: depositAddress.ReceiveAddress,
		TimeoutAt:      fmt.Sprintf("%d", time.Now().Add(time.Minute*30).Unix()),
		Price:          depositAddress.Price,
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

func (h *httpDelivery) btcTestListen(w http.ResponseWriter, r *http.Request) {

	result := h.Usecase.JobMint_CheckTxMasterAndRefund()

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")

	// err := h.Usecase.BtcCheckSendNFTForBuyOrder()

	// fmt.Println("len result", len(result))

	// h.Response.RespondSuccess(w, http.StatusOK, response.Success, err, "")
	// h.Response.RespondSuccess(w, http.StatusOK, response.Success, nil, "")
}

func (h *httpDelivery) btcTestTransfer(w http.ResponseWriter, r *http.Request) {

	//
	//

	// var reqBody request.SendNFT
	// decoder := json.NewDecoder(r.Body)
	// err := decoder.Decode(&reqBody)
	// if err != nil {
	// 	h.Logger.Error("httpDelivery.btcTestTransfer.Decode", err.Error(), err)
	// 	h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
	// 	return
	// }

	// resp, err := h.Usecase.SendTokenMKPTest(reqBody.WalletName, reqBody.ReceiveOrdAddress, reqBody.InscriptionID)

	// if err != nil {
	// 	h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
	// 	return
	// }

	// h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, nil, "")
}

func (h *httpDelivery) btcMarketplaceFilterInfo(w http.ResponseWriter, r *http.Request) {

	result, err := h.Usecase.BTCMarketplaceFilterInfo()

	if err != nil {
		h.Logger.Error("h.Usecase.BTCMarketplaceListNFT", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")
}

func (h *httpDelivery) btcMarketplaceRunFilterInfo(w http.ResponseWriter, r *http.Request) {

	err := h.Usecase.BTCMarketplaceUpdateNftInfo()

	if err != nil {
		h.Logger.Error("h.Usecase.BTCMarketplaceListNFT", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, true, "")
}

func (h *httpDelivery) btcMarketplaceCollectionStats(w http.ResponseWriter, r *http.Request) {
	collectionID := r.URL.Query().Get("collection_id")
	_ = collectionID // not use for now

	projectID := r.URL.Query().Get("project_id")

	result, err := h.Usecase.GetCollectionMarketplaceStats(projectID)
	if err != nil {
		h.Logger.Error("h.Usecase.BTCMarketplaceListNFT", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")
}
