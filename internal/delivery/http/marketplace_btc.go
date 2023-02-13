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
)

func (h *httpDelivery) btcMarketplaceListing(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("httpDelivery.btcMarketplaceListing", r)
	defer h.Tracer.FinishSpan(span, log)
	h.Response.SetLog(h.Tracer, span)

	var reqBody request.CreateMarketplaceBTCListing
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Error("httpDelivery.btcMarketplaceListing.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	inscriptionID := reqBody.InscriptionID

	inscriptionIDs := strings.Split(inscriptionID, "https://ordinals.com/inscription/")

	if len(inscriptionIDs) == 2 {
		inscriptionID = inscriptionIDs[1]
		// err := fmt.Errorf("invalid ordinals link")
		// log.Error("httpDelivery.btcMarketplaceListing.Decode", err.Error(), err)
		// h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		// return
	}
	// if reqBody.Name == "" {
	// 	err := fmt.Errorf("invalid name")
	// 	log.Error("httpDelivery.btcMarketplaceListing.Decode", err.Error(), err)
	// 	h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
	// 	return
	// }

	// check valid inscriptionID:
	suffix := "i0"
	if !strings.HasSuffix(inscriptionID, suffix) {
		err := fmt.Errorf("invalid inscriptionID")
		log.Error("httpDelivery.btcMarketplaceListing.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	txHash := strings.TrimSuffix(inscriptionID, suffix)
	_, err = chainhash.NewHashFromStr(txHash)
	if err != nil {
		err := fmt.Errorf("invalid inscriptionID")
		log.Error("httpDelivery.btcMarketplaceListing.NewHashFromStr", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	// check btc address:
	ok, _ := btc.ValidateAddress("btc", reqBody.ReceiveAddress)
	if !ok {
		err := fmt.Errorf("invalid ReceiveAddress")
		log.Error("httpDelivery.btcMarketplaceListing.ValidateAddress", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	ok, _ = btc.ValidateAddress("btc", reqBody.ReceiveOrdAddress)
	if !ok {
		err := fmt.Errorf("invalid ReceiveOrdAddress")
		log.Error("httpDelivery.btcMarketplaceListing.ValidateAddress", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	priceNumber, err := strconv.ParseInt(reqBody.Price, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid price")
		log.Error("httpDelivery.btcMarketplaceListing.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	// check price:
	if priceNumber < utils.MIN_BTC_TO_LIST_BTC {
		err := fmt.Errorf("Minimum price is %.2f BTC", float64(utils.MIN_BTC_TO_LIST_BTC/1e8))
		log.Error("httpDelivery.btcMarketplaceListing.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase := structure.MarketplaceBTC_ListingInfo{
		InscriptionID:  inscriptionID,
		Name:           reqBody.Name,
		Description:    reqBody.Description,
		SellOrdAddress: reqBody.ReceiveOrdAddress,
		SellerAddress:  reqBody.ReceiveAddress,
		Price:          reqBody.Price,
		ServiceFee:     fmt.Sprintf("%v", utils.BUY_NFT_CHARGE),
	}

	nft, err := h.Usecase.Repo.FindBtcNFTListingUnsoldByNFTID(inscriptionID)
	if err == nil {
		if nft != nil {
			err := fmt.Errorf("this inscription is already listed")
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
	}

	listing, err := h.Usecase.BTCMarketplaceListingNFT(span, reqUsecase)
	if err != nil {
		log.Error("h.Usecase.BTCMarketplaceListingNFT", err.Error(), err)
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
	span, log := h.StartSpan("btcMarketplaceListNFTs", r)
	defer h.Tracer.FinishSpan(span, log)

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 20
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		limit = 20
	}

	buyableOnly := false
	if r.URL.Query().Get("buyable-only") == "true" {
		buyableOnly = true
	}

	result, err := h.Usecase.BTCMarketplaceListNFT(span, buyableOnly, int64(limit), int64(offset))
	if err != nil {
		log.Error("h.Usecase.BTCMarketplaceListNFT", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")
}

func (h *httpDelivery) btcMarketplaceNFTDetail(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("btcMarketplaceNFTDetail", r)
	defer h.Tracer.FinishSpan(span, log)

	vars := mux.Vars(r)
	inscriptionID := vars["ID"]
	span.SetTag("ID", inscriptionID)
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
			log.Error("h.Usecase.Repo.FindBtcNFTListingByNFTID", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
		// priceInt, err := strconv.ParseInt(nft.Price, 10, 64)
		// if err != nil {
		// 	log.Error("h.btcMarketplaceNFTDetail.strconv.ParseInt", err.Error(), err)
		// 	h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		// 	return
		// }
		// lastPrice = priceInt
		isCompleted = true
	}

	if !nft.IsSold {
		buyOrders, err := h.Usecase.Repo.GetBTCListingHaveOngoingOrder(nft.UUID)
		if err != nil {
			log.Error("h.Usecase.Repo.GetBTCListingHaveOngoingOrder", err.Error(), err)
		}
		currentTime := time.Now()
		for _, order := range buyOrders {
			expireTime := order.ExpiredAt
			// not expired yet still waiting for btc
			if currentTime.Before(expireTime) && (order.Status == entity.StatusBuy_Pending || order.Status == entity.StatusBuy_NotEnoughBalance) {
				isBuyable = false
				break
			}
			// could be expired but received btc
			if order.Status != entity.StatusBuy_Pending && order.Status != entity.StatusBuy_NotEnoughBalance {
				isBuyable = false
				break
			}
		}
	}

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
	//log.SetData("resp.Proposal", resp)
	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, nftInfo, "")
}

func (h *httpDelivery) btcMarketplaceListingFee(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("httpDelivery.btcMarketplaceListingFee", r)
	defer h.Tracer.FinishSpan(span, log)
	h.Response.SetLog(h.Tracer, span)

	var reqBody request.ListingFee
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Error("httpDelivery.btcMint.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	inscriptionID := reqBody.InscriptionID

	inscriptionIDs := strings.Split(inscriptionID, "https://ordinals.com/inscription/")

	if len(inscriptionIDs) == 2 {
		inscriptionID = inscriptionIDs[1]
	}

	tokenUri, err := h.Usecase.GetTokenByTokenID(span, inscriptionID, 0)
	if err != nil {
		resp := response.ListingFee{
			ServiceFee: fmt.Sprintf("%v", utils.BUY_NFT_CHARGE),
			RoyaltyFee: fmt.Sprintf("%v", 0),
		}
		h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
		return
	}

	projectDetail, err := h.Usecase.GetProjectDetail(span, structure.GetProjectDetailMessageReq{
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
	span, log := h.StartSpan("httpDelivery.ethGetReceiveWalletAddress", r)
	defer h.Tracer.FinishSpan(span, log)
	h.Response.SetLog(h.Tracer, span)

	var reqBody request.CreateMarketplaceBTCBuyOrder
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Error("httpDelivery.btcMint.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	// check btc address:
	ok, _ := btc.ValidateAddress("btc", reqBody.WalletAddress)
	if !ok {
		err := fmt.Errorf("invalid WalletAddress")
		log.Error("httpDelivery.btcMarketplaceListing.WalletAddress", err.Error(), err)
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
		log.Error("h.Usecase.BTCMarketplaceListingNFT", err.Error(), err)
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

	depositAddress, err := h.Usecase.BTCMarketplaceBuyOrder(span, reqUsecase)
	if err != nil {
		log.Error("h.Usecase.BTCMarketplaceListingNFT", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp := response.CreateMarketplaceBTCBuyOrder{
		ReceiveAddress: depositAddress,
		TimeoutAt:      fmt.Sprintf("%d", time.Now().Add(time.Minute*15).Unix()),
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

func (h *httpDelivery) btcTestListen(w http.ResponseWriter, r *http.Request) {

	span, log := h.StartSpan("BtcChecktListNft", r)
	defer h.Tracer.FinishSpan(span, log)

	result := h.Usecase.JobInscribeSendBTCToOrdWallet(span)

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")

	// err := h.Usecase.BtcCheckSendNFTForBuyOrder(span)

	// fmt.Println("len result", len(result))

	// h.Response.RespondSuccess(w, http.StatusOK, response.Success, err, "")
	// h.Response.RespondSuccess(w, http.StatusOK, response.Success, nil, "")
}

func (h *httpDelivery) btcTestTransfer(w http.ResponseWriter, r *http.Request) {

	// span, log := h.StartSpan("btcTestTransfer", r)
	// defer h.Tracer.FinishSpan(span, log)

	// var reqBody request.SendNFT
	// decoder := json.NewDecoder(r.Body)
	// err := decoder.Decode(&reqBody)
	// if err != nil {
	// 	log.Error("httpDelivery.btcTestTransfer.Decode", err.Error(), err)
	// 	h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
	// 	return
	// }

	// resp, err := h.Usecase.SendTokenMKPTest(span, reqBody.WalletName, reqBody.ReceiveOrdAddress, reqBody.InscriptionID)

	// if err != nil {
	// 	h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
	// 	return
	// }

	// h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, nil, "")
}
