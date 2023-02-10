package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
)

func (h *httpDelivery) btcMarketplaceListing(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("httpDelivery.ethGetReceiveWalletAddress", r)
	defer h.Tracer.FinishSpan(span, log)
	h.Response.SetLog(h.Tracer, span)

	var reqBody request.CreateMarketplaceBTCListing
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Error("httpDelivery.btcMint.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	reqUsecase := structure.MarketplaceBTC_ListingInfo{
		InscriptionID:  reqBody.InscriptionID,
		Name:           reqBody.Name,
		Description:    reqBody.Description,
		SellOrdAddress: reqBody.ReceiveAddress,
		Price:          reqBody.Price,
		ServiceFee:     "5000", //5% 50000/10000
	}

	depositAddress, err := h.Usecase.BTCMarketplaceListingNFT(span, reqUsecase)
	if err != nil {
		log.Error("h.Usecase.BTCMarketplaceListingNFT", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp := response.CreateMarketplaceBTCListing{
		ReceiveAddress: depositAddress,
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
		}
		result = append(result, nftInfo)
	}
	// baseF, err := h.BaseFilters(r)
	// if err != nil {
	// 	log.Error("BaseFilters", err.Error(), err)
	// 	h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
	// 	return
	// }

	// h.Usecase.FindBtcNFTListingByNFTID()
	// f := structure.FilterProposalVote{}
	// f.BaseFilters = *baseF

	// f.ProposalID = &proposalID
	// support := r.URL.Query().Get("support")
	// if support != "" {
	// 	supportInt, err := strconv.Atoi(support)
	// 	if err != nil {
	// 		log.Error("strconv.Atoi", err.Error(), err)
	// 		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
	// 		return
	// 	}
	// 	f.Support = &supportInt
	// }

	// voter := r.URL.Query().Get("voter")
	// if voter != "" {
	// 	f.Voter = &voter
	// }

	// paginationData, err := h.Usecase.GetProposalVotes(span, f)
	// if err != nil {
	// 	log.Error("h.Usecase.GetProposal", err.Error(), err)
	// 	h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
	// 	return
	// }

	// pResp := []response.ProposalVotesResp{}
	// iPro := paginationData.Result
	// pro := iPro.([]entity.ProposalVotes)
	// for _, proItem := range pro {

	// 	tmp := &response.ProposalVotesResp{}
	// 	err := response.CopyEntityToRes(tmp, &proItem)
	// 	if err != nil {
	// 		log.Error("copier.Copy", err.Error(), err)
	// 		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
	// 		return
	// 	}

	// 	pResp = append(pResp, *tmp)
	// }

	// //log.SetData("resp.Proposal", resp)
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
	switch inscriptionID {
	case "c0f8acd8f0d91d490ac9c08977b142aa836207d2ee93d111992866cf47a6d2e6i0":
		nft = &entity.MarketplaceBTCListing{
			InscriptionID: "c0f8acd8f0d91d490ac9c08977b142aa836207d2ee93d111992866cf47a6d2e6i0",
			Name:          "Test1",
			Description:   "test1 blah blah blah",
			Price:         "1234567",
			BaseEntity: entity.BaseEntity{
				UUID: "1",
			},
		}
	case "2696948882cc088f2d1c160981501a48b3744d8d5df0e8d9a71557e716c634dci0":
		nft = &entity.MarketplaceBTCListing{
			InscriptionID: "2696948882cc088f2d1c160981501a48b3744d8d5df0e8d9a71557e716c634dci0",
			Name:          "Test2",
			Description:   "test2 blah blah blah",
			Price:         "1234567", BaseEntity: entity.BaseEntity{
				UUID: "2",
			},
		}
	case "95752b856f94d0c60bee700d6df1b47c949c28f2a06859cf6d5a3466843463b8i0":
		nft = &entity.MarketplaceBTCListing{
			InscriptionID: "95752b856f94d0c60bee700d6df1b47c949c28f2a06859cf6d5a3466843463b8i0",
			Name:          "Test3",
			Description:   "test3 blah blah blah",
			Price:         "1234567", BaseEntity: entity.BaseEntity{
				UUID: "3",
			},
		}
	default:
		nft, err = h.Usecase.Repo.FindBtcNFTListingByNFTID(inscriptionID)
		if err != nil {
			log.Error("h.Usecase.Repo.FindBtcNFTListingByNFTID", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}
	}

	nftInfo := response.MarketplaceNFTDetail{
		InscriptionID: nft.InscriptionID,
		Name:          nft.Name,
		Description:   nft.Description,
		Price:         nft.Price,
		OrderID:       nft.UUID,
	}
	//log.SetData("resp.Proposal", resp)
	h.Response.SetLog(h.Tracer, span)
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, nftInfo, "")
}

func (h *httpDelivery) btcMarketplaceCreateBuyOrder(w http.ResponseWriter, r *http.Request) {
	span, log := h.StartSpan("httpDelivery.btcMarketplaceCreateBuyOrder", r)
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

	reqUsecase := &structure.EthWalletAddressData{}
	err = copier.Copy(reqUsecase, reqBody)
	if err != nil {
		log.Error("copier.Copy", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	ethWallet, err := h.Usecase.CreateETHWalletAddress(span, *reqUsecase)
	if err != nil {
		log.Error("h.Usecase.CreateETHWalletAddress", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	log.SetData("ethWallet", ethWallet)
	resp, err := h.EthWalletAddressToResp(ethWallet)
	if err != nil {
		log.Error(" h.proposalToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}
