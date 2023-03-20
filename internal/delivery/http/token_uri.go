package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"rederinghub.io/utils/helpers"

	"github.com/gorilla/mux"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/algolia"
)

// UserCredits godoc
// @Summary get token uri data
// @Description get token uri data
// @Tags Token for Opensea
// @Accept  json
// @Produce  json
// @Param contractAddress path string true "contract address"
// @Param tokenID path string true "token ID"
// @Param captureTimeout query integer false "Capture timeout"
// @Success 200 {object} response.JsonResponse{data=response.TokenURIResp}
// @Router /token/{contractAddress}/{tokenID} [GET]
func (h *httpDelivery) tokenURI(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	contractAddress := vars["contractAddress"]

	tokenID := vars["tokenID"]

	captureTimeout := r.URL.Query().Get("captureTimeout")
	h.Logger.Info("captureTimeout", captureTimeout)
	captureTimeoutInt, errT := strconv.Atoi(captureTimeout)
	if errT != nil {
		captureTimeoutInt = 5
	}

	message, err := h.Usecase.GetToken(structure.GetTokenMessageReq{
		ContractAddress: contractAddress,
		TokenID:         tokenID,
	}, captureTimeoutInt)

	if err != nil {
		h.Logger.Error("h.Usecase.GetToken", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp := response.TokenURIResp{
		Name:         message.Name,
		Description:  message.Description,
		Image:        message.ParsedImage,
		AnimationURL: message.AnimationURL,
		Attributes:   message.ParsedAttributes,
	}

	h.Logger.Info("resp.message", message.TokenID)
	h.Response.RespondWithoutContainer(w, http.StatusOK, resp)
}

// UserCredits godoc
// @Summary get token's traits
// @Description get token's traits
// @Tags Token for Opensea
// @Accept  json
// @Produce  json
// @Param contractAddress path string true "contract address"
// @Param tokenID path string true "token ID"
// @Success 200 {object} response.JsonResponse{data=response.TokenTraitsResp}
// @Router /trait/{contractAddress}/{tokenID} [GET]
func (h *httpDelivery) tokenTrait(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	contractAddress := vars["contractAddress"]

	tokenID := vars["tokenID"]

	message, err := h.Usecase.GetToken(structure.GetTokenMessageReq{
		ContractAddress: contractAddress,
		TokenID:         tokenID,
	}, 5)

	if err != nil {
		h.Logger.Error("h.Usecase.GetToken", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp := response.TokenTraitsResp{}
	resp.Attributes = message.ParsedAttributes
	h.Logger.Info("resp.message", message)
	h.Response.RespondWithoutContainer(w, http.StatusOK, resp)
}

// UserCredits godoc
// @Summary get token uri data
// @Description get token uri data
// @Tags Tokens
// @Accept  json
// @Produce  json
// @Param contractAddress path string true "contract address"
// @Param tokenID path string true "token ID"
// @Param captureTimeout query integer false "Capture timeout"
// @Success 200 {object} response.JsonResponse{data=response.InternalTokenURIResp}
// @Router /tokens/{contractAddress}/{tokenID} [GET]
func (h *httpDelivery) tokenURIWithResp(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	contractAddress := vars["contractAddress"]

	tokenID := vars["tokenID"]

	captureTimeout := r.URL.Query().Get("captureTimeout")
	h.Logger.Info("captureTimeout", captureTimeout)
	captureTimeoutInt, errT := strconv.Atoi(captureTimeout)
	if errT != nil {
		captureTimeoutInt = 5
	}

	token, err := h.Usecase.GetToken(structure.GetTokenMessageReq{
		ContractAddress: contractAddress,
		TokenID:         tokenID,
	}, captureTimeoutInt)
	if err != nil {
		h.Logger.Error("h.Usecase.GetToken", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	filter := &algolia.AlgoliaFilter{SearchStr: token.TokenID}
	aresp, _, _, err := h.Usecase.AlgoliaSearchInscription(filter)
	if err != nil {
		h.Logger.Error("h.Usecase.AlgoliaSearchInscription", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Logger.Info("h.Usecase.GetToken", token.TokenID)

	resp, err := h.tokenToResp(token)
	for _, i := range aresp {
		if i.Inscription != nil && i.Inscription.ObjectId == token.TokenID {
			resp.OrdinalsData = &response.OrdinalsData{
				Sat:         i.Inscription.Sat,
				ContentType: i.Inscription.ContentType,
				Timestamp:   i.Inscription.Timestamp,
				Block:       i.Inscription.GenesisHeight,
			}
			break
		}
	}

	if err != nil {
		err := errors.New("Cannot parse products")
		h.Logger.Error("tokenToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	fmt.Println("resp, err====>", resp, err)

	if resp != nil {
		// get nft listing detail to check buyable (contact Phuong):
		// nft, _ := h.Usecase.GetListingDetail(tokenID)
		// if nft != nil {
		// 	resp.Buyable = nft.Buyable
		// 	resp.PriceBTC = nft.Price
		// 	resp.OrderID = nft.OrderID
		// 	resp.IsCompleted = nft.IsCompleted

		// 	resp.ListingDetail = nft

		// }
		listingInfo, err := h.Usecase.Repo.GetDexBTCListingOrderPendingByInscriptionID(resp.TokenID)
		if err != nil {
			h.Logger.Error("tokenURIWithResp.Usecase.Repo.GetDexBTCListingOrderPendingByInscriptionID", resp.TokenID, err.Error(), err)
		} else {
			if listingInfo.CancelTx == "" {
				buyEth, _ := h.Usecase.Repo.GetDexBTCBuyETHBuyingByInscriptionID(listingInfo.InscriptionID)
				resp.Buyable = true
				if buyEth != nil {
					resp.Buyable = false
				}
				resp.PriceBTC = fmt.Sprintf("%v", listingInfo.Amount)
				resp.OrderID = listingInfo.UUID
				resp.SellVerified = listingInfo.Verified
				if resp.SellVerified {
					btcRate, ethRate, err := h.Usecase.GetBTCToETHRate()
					if err != nil {
						h.Logger.Error("GenBuyETHOrder GetBTCToETHRate", err.Error(), err)
					}
					// outLen := 0
					// psbt, err := btc.ParsePSBTFromBase64(listingInfo.RawPSBT)
					// if err != nil {
					// 	h.Logger.Error("GenBuyETHOrder ParsePSBTFromBase64", listingInfo.ID, err)
					// } else {
					// 	outLen = len(psbt.UnsignedTx.TxOut)
					// }
					// amountBTCFee := btc.EstimateTxFee(uint(len(listingInfo.Inputs)+3), uint(outLen+2), uint(15)) + btc.EstimateTxFee(1, 2, uint(15))
					amountBTCRequired := uint64(listingInfo.Amount) + 1000
					amountBTCRequired += (amountBTCRequired / 10000) * 15 // + 0,15%
					// amountBTCRequired += amountBTCFee

					amountETH, _, _, err := h.Usecase.ConvertBTCToETHWithPriceEthBtc(fmt.Sprintf("%f", float64(amountBTCRequired)/1e8), btcRate, ethRate)
					if err != nil {
						h.Logger.Error("GenBuyETHOrder convertBTCToETH", err.Error(), err)
					}
					resp.PriceETH = amountETH
				}
			}
		}
	}
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

// UserCredits godoc
// @Summary get token's traits
// @Description get token's traits
// @Tags Tokens
// @Accept  json
// @Produce  json
// @Param contractAddress path string true "contract address"
// @Param tokenID path string true "token ID"
// @Success 200 {object} response.JsonResponse{data=response.InternalTokenTraitsResp}
// @Router /tokens/traits/{contractAddress}/{tokenID} [GET]
func (h *httpDelivery) tokenTraitWithResp(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	contractAddress := vars["contractAddress"]

	tokenID := vars["tokenID"]

	message, err := h.Usecase.GetToken(structure.GetTokenMessageReq{
		ContractAddress: contractAddress,
		TokenID:         tokenID,
	}, 5)

	if err != nil {
		h.Logger.Error("h.Usecase.GetToken", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	resp := response.InternalTokenTraitsResp{}
	resp.Attributes = message.ParsedAttributes
	h.Logger.Info("resp.message", message)

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

// UserCredits godoc
// @Summary get project's tokens
// @Description get tokens by project address
// @Tags Project
// @Accept  json
// @Produce  json
// @Param contract_address query string false "contract_address"
// @Param gen_nft_address query string false "gen_nft_address"
// @Param owner_address query string false "owner_address"
// @Param creator_address query string false "creator_address"
// @Param tokenID query string false "Filter via tokenID"
// @Param attributes query []string false "attributes"
// @Param has_price query bool false "has_price"
// @Param from_price query string false "from_price"
// @Param to_price query string false "to_price"
// @Param sort query string false "newest, minted-newest, token-price-asc, token-price-desc"
// @Param limit query int false "limit"
// @Param cursor query string false "The cursor returned in the previous response (used for getting the next page)."
// @Param genNFTAddr path string true "This is provided from Project Detail API"
// @Success 200 {object} response.JsonResponse{}
// @Router /project/{genNFTAddr}/tokens [GET]
func (h *httpDelivery) TokensOfAProject(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	genNFTAddr := vars["genNFTAddr"]
	h.Logger.Info("genNFTAddr", genNFTAddr)

	f := structure.FilterTokens{}
	err := f.CreateFilter(r)
	if err != nil {
		h.Logger.Error("f.CreateFilter", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	f.GenNFTAddr = &genNFTAddr
	bf, err := h.BaseFilters(r)
	if err != nil {
		h.Logger.Error("h.Usecase.getProfileNfts.BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	f.BaseFilters = *bf
	resp, err := h.getTokens(f)
	if err != nil {
		h.Logger.Error("h.Usecase.getProfileNfts.getTokens", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	//
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

func (h *httpDelivery) TokensOfAProjectNew(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	genNFTAddr := vars["genNFTAddr"]
	h.Logger.Info("genNFTAddr", genNFTAddr)

	f := structure.FilterTokens{}
	err := f.CreateFilter(r)
	if err != nil {
		h.Logger.Error("f.CreateFilter", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	f.GenNFTAddr = &genNFTAddr
	bf, err := h.BaseFilters(r)
	if err != nil {
		h.Logger.Error("h.Usecase.getProfileNfts.BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	f.BaseFilters = *bf
	resp, err := h.getTokensNew(f)
	if err != nil {
		h.Logger.Error("h.Usecase.getProfileNfts.getTokens", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	//
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

// UserCredits godoc
// @Summary User profile's nft
// @Description User profile's nft
// @Tags Profile
// @Accept  json
// @Produce  json
// @Param tokenID path string false "Filter via tokenID"
// @Param walletAddress path string true "Wallet address"
// @Param contract_address query string false "contract_address"
// @Param gen_nft_address query string false "gen_nft_address"
// @Param owner_address query string false "owner_address"
// @Param creator_address query string false "creator_address"
// @Param tokenID query string false "Filter via tokenID"
// @Param sort query string false "newest, minted-newest, priority-asc, priority-desc"
// @Success 200 {object} response.JsonResponse{data=response.InternalTokenURIResp}
// @Router /profile/wallet/{walletAddress}/nfts [GET]
func (h *httpDelivery) TokensOfAProfile(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	walletAddress := vars["walletAddress"]
	h.Logger.Info("walletAddress", walletAddress)

	f := structure.FilterTokens{}
	f.CreateFilter(r)
	f.OwnerAddr = &walletAddress

	bf, err := h.BaseFilters(r)
	if err != nil {
		h.Logger.Error("h.Usecase.getProfileNfts.BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	f.BaseFilters = *bf
	tokenID := r.URL.Query().Get("tokenID")
	if tokenID != "" {
		f.TokenIDs = append(f.TokenIDs, tokenID)
	}

	resp, err := h.getTokens(f)
	if err != nil {
		h.Logger.Error("h.Usecase.getProfileNfts.getTokens", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	//
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")

}

// UserCredits godoc
// @Summary get  projects by wallet
// @Description get  projects by wallet
// @Tags Profile
// @Accept  json
// @Produce  json
// @Param walletAddress path string false "Filter project via wallet address"
// @Param limit query int false "limit"
// @Param cursor query string false "The cursor returned in the previous response (used for getting the next page)."
// @Success 200 {object} response.JsonResponse{}
// @Router /profile/wallet/{walletAddress}/projects [GET]
func (h *httpDelivery) getProjectsByWallet(w http.ResponseWriter, r *http.Request) {

	var err error
	vars := mux.Vars(r)
	walletAddress := vars["walletAddress"]

	hidden := false
	baseF, err := h.BaseFilters(r)
	if err != nil {
		h.Logger.Error("BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	f := structure.FilterProjects{}
	f.BaseFilters = *baseF
	f.WalletAddress = &walletAddress

	ctx := r.Context()
	iWalletAddress := ctx.Value(utils.SIGNED_WALLET_ADDRESS)
	currentUserWalletAddress, ok := iWalletAddress.(string)
	if !ok {
		f.IsHidden = &hidden
	}

	if ok && currentUserWalletAddress != walletAddress {
		f.IsHidden = &hidden
	}

	uProjects, err := h.Usecase.GetAllProjects(f)
	if err != nil {
		h.Logger.Error("h.Usecase.GetProjects", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	pResp := []response.ProjectResp{}
	iProjects := uProjects.Result
	projects := iProjects.([]entity.Projects)
	for _, project := range projects {

		p, err := h.projectToResp(&project)
		if err != nil {
			h.Logger.Error("copier.Copy", err.Error(), err)
			h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
			return
		}

		pResp = append(pResp, *p)
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, h.PaginationResp(uProjects, pResp), "")
}

// UserCredits godoc
// @Summary get list tokenUris
// @Description get tokenUris
// @Tags TokenUri
// @Accept  json
// @Produce  json
// @Param page query string false "page"
// @Param limit query int false "limit"
// @Param search query string false "search"
// @Param cursor query string false "The cursor returned in the previous response (used for getting the next page)."
// @Success 200 {object} response.JsonResponse{}
// @Router /token-uri [GET]
func (h *httpDelivery) getTokenUris(w http.ResponseWriter, r *http.Request) {
	f := structure.FilterTokens{}
	err := f.CreateFilter(r)
	if err != nil {
		h.Logger.Error("f.CreateFilter", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	if len(*f.Search) < 3 {
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, errors.New("Term search minimum is 3 characters"))
	}

	bf, err := h.BaseFilters(r)
	if err != nil {
		h.Logger.Error("h.Usecase.getProfileNfts.BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	f.BaseFilters = *bf
	resp, err := h.getTokensForSearch(f)
	if err != nil {
		h.Logger.Error("h.Usecase.getProfileNfts.getTokens", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	//
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")

}

func (h *httpDelivery) getTokens(f structure.FilterTokens) (*response.PaginationResponse, error) {
	pag, err := h.Usecase.FilterTokens(f)
	if err != nil {
		h.Logger.Error("h.Usecase.getProfileNfts.FilterTokens", err.Error(), err)
		return nil, err
	}

	respItems := []response.InternalTokenURIResp{}
	tokens := []entity.TokenUri{}
	iTokensData := pag.Result

	bytes, err := json.Marshal(iTokensData)
	if err != nil {
		err := errors.New("Cannot parse respItems")
		h.Logger.Error("respItems", err.Error(), err)
		return nil, err
	}

	err = json.Unmarshal(bytes, &tokens)
	if err != nil {
		err := errors.New("Cannot Unmarshal")
		h.Logger.Error("Unmarshal", err.Error(), err)
		return nil, err
	}

	// get nft listing from marketplace to show button buy or not (ask Phuong if you need):
	// nftListing, _ := h.Usecase.GetAllListListingWithRule()

	// get btc, btc rate:
	// btcPrice, err := helpers.GetExternalPrice("BTC")
	// if err != nil {
	// 	h.Logger.ErrorAny("convertBTCToETH", zap.Error(err))
	// 	return nil, err
	// }

	// h.Logger.Info("btcPrice", btcPrice)
	// ethPrice, err := helpers.GetExternalPrice("ETH")
	// if err != nil {
	// 	h.Logger.ErrorAny("convertBTCToETH", zap.Error(err))
	// 	return nil, err
	// }
	// h.Logger.Info("btcPrice", btcPrice)

	for _, token := range tokens {
		resp, err := h.tokenToResp(&token)
		if err != nil {
			err := errors.New("Cannot parse products")
			h.Logger.Error("tokenToResp", err.Error(), err)
			return nil, err
		}

		listingInfo, err := h.Usecase.Repo.GetDexBTCListingOrderPendingByInscriptionID(resp.TokenID)
		if err != nil {
			h.Logger.Error("getTokens.Usecase.Repo.GetDexBTCListingOrderPendingByInscriptionID", resp.TokenID, err.Error(), err)
		} else {
			if listingInfo.CancelTx == "" {
				resp.Buyable = true
				resp.PriceBTC = fmt.Sprintf("%v", listingInfo.Amount)
				resp.OrderID = listingInfo.UUID
			}
		}
		// for _, v := range nftListing {
		// 	if resp != nil {
		// 		if strings.EqualFold(v.InscriptionID, resp.TokenID) {
		// 			resp.Buyable = v.Buyable
		// 			resp.PriceBTC = v.Price
		// 			resp.OrderID = v.OrderID
		// resp.IsCompleted = v.IsCompleted

		// listPaymentInfo, err := h.Usecase.GetListingPaymentInfoWithEthBtcPrice(v.PayType, v.Price, btcPrice, ethPrice)

		// if err != nil {
		// 	continue
		// }
		// v.PaymentListingInfo = listPaymentInfo

		// resp.ListingDetail = &v

		// 			break
		// 		}
		// 	}
		// }

		respItems = append(respItems, *resp)
	}

	resp := h.PaginationResp(pag, respItems)
	return &resp, nil
}

func (h *httpDelivery) getTokensNew(f structure.FilterTokens) (*response.PaginationResponse, error) {
	pag, err := h.Usecase.FilterTokensNew(f)
	if err != nil {
		h.Logger.Error("h.Usecase.getProfileNfts.FilterTokens", err.Error(), err)
		return nil, err
	}
	newList := []entity.TokenUriListingFilter{}
	btcRate, ethRate, err := h.Usecase.GetBTCToETHRate()
	if err != nil {
		h.Logger.Error("GenBuyETHOrder GetBTCToETHRate", err.Error(), err)
	}
	for _, item := range pag.Result.([]entity.TokenUriListingFilter) {

		amountBTCRequired := uint64(item.Price) + 1000
		amountBTCRequired += (amountBTCRequired / 10000) * 15 // + 0,15%
		// amountBTCRequired += btc.EstimateTxFee(4, 3, uint(15)) + btc.EstimateTxFee(1, 2, uint(15))

		amountETH, _, _, err := h.Usecase.ConvertBTCToETHWithPriceEthBtc(fmt.Sprintf("%f", float64(amountBTCRequired)/1e8), btcRate, ethRate)
		if err != nil {
			h.Logger.Error("GenBuyETHOrder convertBTCToETH", err.Error(), err)
		}

		item.PriceETH = amountETH

		if strings.HasSuffix(item.AnimationURL, ".html") {
			client := http.Client{
				CheckRedirect: func(r *http.Request, via []*http.Request) error {
					r.URL.Opaque = r.URL.Path
					return nil
				},
			}
			r, err := client.Get(item.AnimationURL)
			if err != nil {
				h.Usecase.Logger.LogAny("fail")
			}
			defer r.Body.Close()
			b, err := io.ReadAll(r.Body)
			if err != nil {
				h.Usecase.Logger.LogAny("fail")
			}
			base64 := helpers.Base64Encode(b)
			item.AnimationURL = "data:text/html;base64," + base64
		}
		if err != nil {
			return nil, err
		}
		// resp.Attributes = input.ParsedAttributes
		// if input.ParsedImage != nil {
		// 	resp.Image = *input.ParsedImage
		// } else {
		// 	resp.Image = input.Thumbnail
		// }
		if strings.Index(item.Image, "glb") == -1 {
			item.Image = item.Thumbnail
		}
		newList = append(newList, item)
	}
	pag.Result = newList

	// respItems := []response.InternalTokenURIResp{}
	// tokens := []entity.TokenUri{}
	// iTokensData := pag.Result

	// bytes, err := json.Marshal(iTokensData)
	// if err != nil {
	// 	err := errors.New("Cannot parse respItems")
	// 	h.Logger.Error("respItems", err.Error(), err)
	// 	return nil, err
	// }

	// err = json.Unmarshal(bytes, &tokens)
	// if err != nil {
	// 	err := errors.New("Cannot Unmarshal")
	// 	h.Logger.Error("Unmarshal", err.Error(), err)
	// 	return nil, err
	// }

	// get nft listing from marketplace to show button buy or not (ask Phuong if you need):
	// nftListing, _ := h.Usecase.GetAllListListingWithRule()

	// get btc, btc rate:
	// btcPrice, err := helpers.GetExternalPrice("BTC")
	// if err != nil {
	// 	h.Logger.ErrorAny("convertBTCToETH", zap.Error(err))
	// 	return nil, err
	// }

	// h.Logger.Info("btcPrice", btcPrice)
	// ethPrice, err := helpers.GetExternalPrice("ETH")
	// if err != nil {
	// 	h.Logger.ErrorAny("convertBTCToETH", zap.Error(err))
	// 	return nil, err
	// }
	// h.Logger.Info("btcPrice", btcPrice)

	// for _, token := range tokens {
	// 	resp, err := h.tokenToResp(&token)
	// 	if err != nil {
	// 		err := errors.New("Cannot parse products")
	// 		h.Logger.Error("tokenToResp", err.Error(), err)
	// 		return nil, err
	// 	}

	// listingInfo, err := h.Usecase.Repo.GetDexBTCListingOrderPendingByInscriptionID(resp.TokenID)
	// if err != nil {
	// 	h.Logger.Error("getTokens.Usecase.Repo.GetDexBTCListingOrderPendingByInscriptionID", resp.TokenID, err.Error(), err)
	// } else {
	// 	if listingInfo.CancelTx == "" {
	// 		resp.Buyable = true
	// 		resp.PriceBTC = fmt.Sprintf("%v", listingInfo.Amount)
	// 		resp.OrderID = listingInfo.UUID
	// 	}
	// }
	// for _, v := range nftListing {
	// 	if resp != nil {
	// 		if strings.EqualFold(v.InscriptionID, resp.TokenID) {
	// 			resp.Buyable = v.Buyable
	// 			resp.PriceBTC = v.Price
	// 			resp.OrderID = v.OrderID
	// resp.IsCompleted = v.IsCompleted

	// listPaymentInfo, err := h.Usecase.GetListingPaymentInfoWithEthBtcPrice(v.PayType, v.Price, btcPrice, ethPrice)

	// if err != nil {
	// 	continue
	// }
	// v.PaymentListingInfo = listPaymentInfo

	// resp.ListingDetail = &v

	// 			break
	// 		}
	// 	}
	// }

	// 	respItems = append(respItems, *resp)
	// }

	resp := h.PaginationResp(pag, pag.Result)
	return &resp, nil
}

func (h *httpDelivery) getTokensForSearch(f structure.FilterTokens) (*response.PaginationResponse, error) {
	pag, err := h.Usecase.FilterTokens(f)
	if err != nil {
		h.Logger.Error("h.Usecase.getProfileNfts.FilterTokens", err.Error(), err)
		return nil, err
	}

	respItems := []response.ExternalTokenURIResp{}
	tokens := []entity.TokenUri{}
	iTokensData := pag.Result

	bytes, err := json.Marshal(iTokensData)
	if err != nil {
		err := errors.New("Cannot parse respItems")
		h.Logger.Error("respItems", err.Error(), err)
		return nil, err
	}

	err = json.Unmarshal(bytes, &tokens)
	if err != nil {
		err := errors.New("Cannot Unmarshal")
		h.Logger.Error("Unmarshal", err.Error(), err)
		return nil, err
	}

	for _, token := range tokens {
		resp, err := h.tokenExternalToResp(&token)
		if err != nil {
			err := errors.New("Cannot parse products")
			h.Logger.Error("tokenToResp", err.Error(), err)
			return nil, err
		}
		respItems = append(respItems, *resp)
	}

	resp := h.PaginationResp(pag, respItems)
	return &resp, nil
}

func (h *httpDelivery) tokenExternalToResp(input *entity.TokenUri) (*response.ExternalTokenURIResp, error) {
	resp := &response.ExternalTokenURIResp{}
	err := response.CopyEntityToResNoID(resp, input)
	if err != nil {
		return nil, err
	}

	if input.ParsedImage != nil {
		resp.Image = *input.ParsedImage
	} else {
		resp.Image = input.Thumbnail
	}

	resp.InscriptionIndex = input.InscriptionIndex
	if input.Project != nil {
		resp.ProjectName = input.Project.Name
		resp.ProjectID = input.Project.TokenID
	}

	return resp, nil
}

func (h *httpDelivery) tokenToResp(input *entity.TokenUri) (*response.InternalTokenURIResp, error) {
	resp := &response.InternalTokenURIResp{}
	err := response.CopyEntityToResNoID(resp, input)
	if strings.HasSuffix(resp.AnimationURL, ".html") {
		client := http.Client{
			CheckRedirect: func(r *http.Request, via []*http.Request) error {
				r.URL.Opaque = r.URL.Path
				return nil
			},
		}
		r, err := client.Get(resp.AnimationURL)
		if err != nil {
			h.Usecase.Logger.LogAny("fail")
		}
		defer r.Body.Close()
		b, err := io.ReadAll(r.Body)
		if err != nil {
			h.Usecase.Logger.LogAny("fail")
		}
		base64 := helpers.Base64Encode(b)
		resp.AnimationURL = "data:text/html;base64," + base64
	}
	if err != nil {
		return nil, err
	}
	resp.Attributes = input.ParsedAttributes
	// if input.ParsedImage != nil {
	// 	resp.Image = *input.ParsedImage
	// } else {
	// 	resp.Image = input.Thumbnail
	// }
	if strings.Index(resp.Image, "glb") == -1 {
		resp.Image = input.Thumbnail
	}

	if input.Owner != nil {
		ownerResp, err := h.profileToResp(input.Owner)
		if err == nil {
			resp.Owner = ownerResp
		}
	}

	if input.Creator != nil {
		creatorResp, err := h.profileToResp(input.Creator)
		if err == nil {
			resp.Creator = creatorResp
		}
	}

	if input.Project != nil {
		projectResp := &response.ProjectResp{}
		response.CopyEntityToRes(projectResp, input.Project)
		// projectResp, err := h.projectToResp(input.Project)
		// if err == nil {

		// }

		resp.Project = projectResp
		// resp.Stats.Price = input.Stats.PriceInt
		if input.Stats.PriceInt == nil {
			resp.Stats.Price = nil
		} else {
			x := strconv.Itoa(int(*input.Stats.PriceInt))
			resp.Stats.Price = &x
		}
	}

	resp.InscriptionIndex = input.InscriptionIndex
	resp.OrderInscriptionIndex = input.OrderInscriptionIndex

	//resp.Thumbnail = fmt.Sprintf("%s/%s/%s/%s",os.Getenv("DOMAIN"), "api/thumbnail", input.ContractAddress, input.TokenID)

	return resp, nil
}

// UserCredits godoc
// @Summary get tokens
// @Description get tokens
// @Tags Tokens
// @Accept  json
// @Produce  json
// @Param contract_address query string false "contract_address"
// @Param gen_nft_address query string false "gen_nft_address"
// @Param owner_address query string false "owner_address"
// @Param creator_address query string false "creator_address"
// @Param tokenID query string false "Filter via tokenID"
// @Param sort query string false "newest, minted-newest, priority-asc, priority-desc"
// @Param limit query int false "limit"
// @Param page query string false "The cursor returned in the previous response (used for getting the next page)."
// @Success 200 {object} response.JsonResponse{}
// @Router /tokens [GET]
func (h *httpDelivery) Tokens(w http.ResponseWriter, r *http.Request) {
	f := structure.FilterTokens{}
	f.CreateFilter(r)

	bf, err := h.BaseFilters(r)
	if err != nil {
		h.Logger.Error("h.Tokens.BaseFilters", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}
	f.BaseFilters = *bf
	resp, err := h.getTokens(f)
	if err != nil {
		h.Logger.Error("h.Tokens.getTokens", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	//
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

// UserCredits godoc
// @Summary get token uri data
// @Description get token uri data
// @Tags Tokens
// @Accept  json
// @Produce  json
// @Param contractAddress path string true "contract address"
// @Param tokenID path string true "token ID"
// @Param request body request.UpdateTokentReq true "Request body"
// @Success 200 {object} response.JsonResponse{data=response.InternalTokenURIResp}
// @Router /tokens/{contractAddress}/{tokenID} [PUT]
func (h *httpDelivery) updatetokenURIWithResp(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	contractAddress := vars["contractAddress"]

	tokenID := vars["tokenID"]

	var reqBody request.UpdateTokentReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		h.Logger.Error("decoder.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	token, err := h.Usecase.UpdateToken(structure.UpdateTokenReq{
		ContracAddress: contractAddress,
		TokenID:        tokenID,
		Priority:       reqBody.Priority,
	})

	if err != nil {
		h.Logger.Error("h.Usecase.GetToken", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Logger.Info("h.Usecase.GetToken", token)
	resp, err := h.tokenToResp(token)
	if err != nil {
		err := errors.New("Cannot parse products")
		h.Logger.Error("tokenToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Logger.Info("resp.token", token)

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

// UserCredits godoc
// @Summary Update token's thumbnail
// @Description Update token's thumbnail
// @Tags Tokens
// @Accept  json
// @Produce  json
// @Param tokenID path string true "token ID"
// @Param request body request.UpdateTokenThumbnailReq true "Request body"
// @Success 200 {object} response.JsonResponse{data=response.InternalTokenURIResp}
// @Router /tokens/{tokenID}/thumbnail [POST]
func (h *httpDelivery) updateTokenThumbnail(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	tokenID := vars["tokenID"]

	var reqBody request.UpdateTokenThumbnailReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		h.Logger.Error("decoder.Decode", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	token, err := h.Usecase.UpdateTokenThumbnail(structure.UpdateTokenThumbnailReq{
		TokenID:   tokenID,
		Thumbnail: *reqBody.Thumbnail,
	})

	if err != nil {
		h.Logger.Error("h.Usecase.GetToken", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Logger.Info("h.Usecase.GetToken", token)
	resp, err := h.tokenToResp(token)
	if err != nil {
		err := errors.New("Cannot parse products")
		h.Logger.Error("tokenToResp", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Logger.Info("resp.token", token)

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, resp, "")
}

// UserCredits godoc
// @Summary get volume by wallet
// @Description get volume by wallet
// @Tags Profile
// @Accept  json
// @Produce  json
// @Param walletAddress path string false "Filter project via wallet address"
// @Param payType query string false "payType eth|btc"
// @Param limit query int false "limit"
// @Param cursor query string false "The cursor returned in the previous response (used for getting the next page)."
// @Success 200 {object} response.JsonResponse{}
// @Router /profile/wallet/{walletAddress}/volumn [GET]
func (h *httpDelivery) getVolumnByWallet(w http.ResponseWriter, r *http.Request) {

	var err error
	vars := mux.Vars(r)
	walletAddress := vars["walletAddress"]
	paytype := r.URL.Query().Get("payType")
	uProjects, err := h.Usecase.CreatorVolume(walletAddress, paytype)
	if err != nil {
		h.Logger.Error("h.Usecase.GetProjects", err.Error(), err)
		h.Response.RespondWithError(w, http.StatusBadRequest, response.Error, err)
		return
	}

	h.Response.RespondSuccess(w, http.StatusOK, response.Success, uProjects, "")
}
