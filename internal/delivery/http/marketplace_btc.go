package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
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
	inscriptionID := strings.Split(reqBody.InscriptionID, "https://ordinals.com/inscription/")

	if len(inscriptionID) != 2 {
		err := fmt.Errorf("invalid ordinals link")
		log.Error("httpDelivery.btcMarketplaceListing.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	// if reqBody.Name == "" {
	// 	err := fmt.Errorf("invalid name")
	// 	log.Error("httpDelivery.btcMarketplaceListing.Decode", err.Error(), err)
	// 	h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
	// 	return
	// }

	// if strings.Contains(reqBody.ReceiveAddress)
	if priceInt, err := strconv.ParseInt(reqBody.Price, 10, 64); err != nil {
		err := fmt.Errorf("invalid price")
		log.Error("httpDelivery.btcMarketplaceListing.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	} else {
		if priceInt < utils.MIN_BTC_TO_LIST_BTC {
			err := fmt.Errorf("invalid price")
			log.Error("httpDelivery.btcMarketplaceListing.Decode", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
	}

	// check price:

	reqUsecase := structure.MarketplaceBTC_ListingInfo{
		InscriptionID:  inscriptionID[1],
		Name:           reqBody.Name,
		Description:    reqBody.Description,
		SellOrdAddress: reqBody.ReceiveOrdAddress,
		SellerAddress:  reqBody.ReceiveAddress,
		Price:          reqBody.Price,

		//TODO - Tri comment code, becasue type of utils.BUY_NFT_CHARGE and tils.MIN_BTC_TO_LIST_BTC are not defined
		// ServiceFee:     utils.BUY_NFT_CHARGE,
		// Min:            utils.MIN_BTC_TO_LIST_BTC,
	}

	listing, err := h.Usecase.BTCMarketplaceListingNFT(span, reqUsecase)
	if err != nil {
		log.Error("h.Usecase.BTCMarketplaceListingNFT", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp := response.CreateMarketplaceBTCListing{
		ReceiveAddress: listing.HoldOrdAddress,
		TimeoutAt:      fmt.Sprintf("%d", listing.ExpiredAt.Unix()),
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

func (h *httpDelivery) btcMarketplaceListNFTs(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("btcMarketplaceListNFTs", r)
	defer h.Tracer.FinishSpan(span, log)

	nfts, err := h.Usecase.BTCMarketplaceListNFT(span)
	if err != nil {
		log.Error("h.Usecase.BTCMarketplaceListNFT", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	result := []response.MarketplaceNFTDetail{}
	for _, nft := range nfts {
		nftInfo := response.MarketplaceNFTDetail{
			InscriptionID: nft.InscriptionID,
			Name:          nft.Name,
			Description:   nft.Description,
			Price:         nft.Price,
			OrderID:       nft.UUID,
			IsConfirmed:   nft.IsConfirm,
		}
		result = append(result, nftInfo)
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
	lastPrice := int64(0)
	nft, err = h.Usecase.Repo.FindBtcNFTListingUnsoldByNFTID(inscriptionID)
	if err != nil {
		isBuyable = false
		nft, err = h.Usecase.Repo.FindBtcNFTListingSoldByNFTID(inscriptionID)
		if err != nil {
			log.Error("h.Usecase.Repo.FindBtcNFTListingByNFTID", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
		priceInt, err := strconv.ParseInt(nft.Price, 10, 64)
		if err != nil {
			log.Error("h.btcMarketplaceNFTDetail.strconv.ParseInt", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
		lastPrice = priceInt
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

	nftInfo := response.MarketplaceNFTDetail{
		InscriptionID: nft.InscriptionID,
		Name:          nft.Name,
		Description:   nft.Description,
		Price:         nft.Price,
		OrderID:       nft.UUID,
		IsConfirmed:   nft.IsConfirm,
		Buyable:       isBuyable,
		IsCompleted:   isCompleted,
		LastPrice:     lastPrice,
	}
	//log.SetData("resp.Proposal", resp)
	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, nftInfo, "")
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

	result := h.Usecase.BtcChecktListNft(span)

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")

	err := h.Usecase.BtcCheckSendNFTForBuyOrder(span)

	// fmt.Println("len result", len(result))

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, err, "")
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
