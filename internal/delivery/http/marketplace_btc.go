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
	if reqBody.Name == "" {
		err := fmt.Errorf("invalid name")
		log.Error("httpDelivery.btcMarketplaceListing.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	// if strings.Contains(reqBody.ReceiveAddress)
	if _, err := strconv.ParseInt(reqBody.Price, 10, 64); err != nil {
		err := fmt.Errorf("invalid price")
		log.Error("httpDelivery.btcMarketplaceListing.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase := structure.MarketplaceBTC_ListingInfo{
		InscriptionID:  inscriptionID[1],
		Name:           reqBody.Name,
		Description:    reqBody.Description,
		SellOrdAddress: reqBody.ReceiveOrdAddress,
		SellerAddress:  reqBody.ReceiveAddress,
		Price:          reqBody.Price,
		ServiceFee:     "10", //10%
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

	// c0f8acd8f0d91d490ac9c08977b142aa836207d2ee93d111992866cf47a6d2e6i0",
	// 		"price": "1234567",
	// 		"name": "Test1",
	// 		"description": "test1 blah blah blah"
	// 	},
	// 	{
	// 		"inscriptionID": "2696948882cc088f2d1c160981501a48b3744d8d5df0e8d9a71557e716c634dci0",
	// 		"price": "1234567",
	// 		"name": "Test2",
	// 		"description": "test2 blah blah blah"
	// 	},
	// 	{
	// 		"inscriptionID": "95752b856f94d0c60bee700d6df1b47c949c28f2a06859cf6d5a3466843463b8i0",
	var nft *entity.MarketplaceBTCListing
	var err error
	// switch inscriptionID {
	// case "c0f8acd8f0d91d490ac9c08977b142aa836207d2ee93d111992866cf47a6d2e6i0":
	// 	nft = &entity.MarketplaceBTCListing{
	// 		InscriptionID: "c0f8acd8f0d91d490ac9c08977b142aa836207d2ee93d111992866cf47a6d2e6i0",
	// 		Name:          "Test1",
	// 		Description:   "test1 blah blah blah",
	// 		Price:         "1234567",
	// 		BaseEntity: entity.BaseEntity{
	// 			UUID: "1",
	// 		},
	// 	}
	// case "2696948882cc088f2d1c160981501a48b3744d8d5df0e8d9a71557e716c634dci0":
	// 	nft = &entity.MarketplaceBTCListing{
	// 		InscriptionID: "2696948882cc088f2d1c160981501a48b3744d8d5df0e8d9a71557e716c634dci0",
	// 		Name:          "Test2",
	// 		Description:   "test2 blah blah blah",
	// 		Price:         "1234567", BaseEntity: entity.BaseEntity{
	// 			UUID: "2",
	// 		},
	// 	}
	// case "95752b856f94d0c60bee700d6df1b47c949c28f2a06859cf6d5a3466843463b8i0":
	// 	nft = &entity.MarketplaceBTCListing{
	// 		InscriptionID: "95752b856f94d0c60bee700d6df1b47c949c28f2a06859cf6d5a3466843463b8i0",
	// 		Name:          "Test3",
	// 		Description:   "test3 blah blah blah",
	// 		Price:         "1234567", BaseEntity: entity.BaseEntity{
	// 			UUID: "3",
	// 		},
	// 	}
	// default:
	nft, err = h.Usecase.Repo.FindBtcNFTListingByNFTID(inscriptionID)
	if err != nil {
		log.Error("h.Usecase.Repo.FindBtcNFTListingByNFTID", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	// }

	nftInfo := response.MarketplaceNFTDetail{
		InscriptionID: nft.InscriptionID,
		Name:          nft.Name,
		Description:   nft.Description,
		Price:         nft.Price,
		OrderID:       nft.UUID,
		IsConfirmed:   nft.IsConfirm,
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
	_, err = h.Usecase.Repo.FindBtcNFTListingByOrderID(reqBody.OrderID)
	if err != nil {
		log.Error("h.Usecase.BTCMarketplaceListingNFT", err.Error(), err)
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

	// result := h.Usecase.BtcChecktListNft(span)

	// h.Response.RespondSuccess(w, http.StatusOK, response.Success, result, "")

	err := h.Usecase.BtcCheckReceivedBuyingNft(span)

	// fmt.Println("len result", len(result))

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, err, "")
}
