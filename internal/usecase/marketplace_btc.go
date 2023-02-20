package usecase

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	"rederinghub.io/external/ord_service"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/btc"
)

func (u Usecase) BTCMarketplaceListingNFT(rootSpan opentracing.Span, listingInfo structure.MarketplaceBTC_ListingInfo) (*entity.MarketplaceBTCListing, error) {
	span, log := u.StartSpan("BTCMarketplaceListingNFT", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	listing := entity.MarketplaceBTCListing{
		SellOrdAddress: listingInfo.SellOrdAddress,
		SellerAddress:  listingInfo.SellerAddress,
		HoldOrdAddress: "",
		Price:          listingInfo.Price,
		//ServiceFee:     listingInfo.ServiceFee, //Tri comment: ServiceFee is not existed
		IsConfirm:     false,
		IsSold:        false,
		ExpiredAt:     time.Now().Add(time.Hour * 1),
		Name:          listingInfo.Name,
		Description:   listingInfo.Description,
		InscriptionID: listingInfo.InscriptionID,
	}
	holdOrdAddress := ""
	resp, err := u.OrdService.Exec(ord_service.ExecRequest{
		Args: []string{
			"--wallet",
			"ord_marketplace_master",
			"wallet",
			"receive",
		},
	})
	if err != nil {
		log.Error("u.OrdService.Exec.create.receive", err.Error(), err)
		return &listing, err
	}
	holdOrdAddress = strings.ReplaceAll(resp.Stdout, "\n", "")
	listing.HoldOrdAddress = holdOrdAddress
	// sendMessage := func(rootSpan opentracing.Span, offer entity.MarketplaceOffers) {
	// 	span, log := u.StartSpan("MakeOffer.sendMessage", rootSpan)
	// 	defer u.Tracer.FinishSpan(span, log)

	// 	profile, err := u.Repo.FindUserByWalletAddress(offer.Buyer)
	// 	if err != nil {
	// 		log.Error("cancelListing.FindUserByWalletAddress", err.Error(), err)
	// 		return
	// 	}

	// 	token, err := u.Repo.FindTokenByGenNftAddr(offer.CollectionContract, offer.TokenId)
	// 	if err != nil {
	// 		log.Error("cancelListing.FindTokenByGenNftAddr", err.Error(), err)
	// 		return
	// 	}

	// 	preText := fmt.Sprintf("[OfferID %s] has been created by %s", offer.OfferingId, offer.Buyer)
	// 	content := fmt.Sprintf("TokenID: %s", helpers.CreateTokenLink(token.ProjectID, token.TokenID, token.Name))
	// 	title := fmt.Sprintf("User %s made offer with %s", helpers.CreateProfileLink(profile.WalletAddress, profile.DisplayName), offer.Price)

	// 	if _, _, err := u.Slack.SendMessageToSlack(preText, title, content); err != nil {
	// 		log.Error("s.Slack.SendMessageToSlack err", err.Error(), err)
	// 	}

	// }

	log.SetTag("inscriptionID", listingInfo.InscriptionID)
	// check if listing is created or not
	err = u.Repo.CreateMarketplaceListingBTC(&listing)
	if err != nil {
		log.Error("BTCMarketplaceListingNFT.Repo.CreateMarketplaceListingBTC", "", err)
		return &listing, err
	}
	return &listing, nil
}

func (u Usecase) BTCMarketplaceListNFT(rootSpan opentracing.Span, filter *entity.FilterString, buyableOnly bool, limit, offset int64) ([]structure.MarketplaceNFTDetail, error) {
	span, log := u.StartSpan("BTCMarketplaceListingNFT", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	result := []structure.MarketplaceNFTDetail{}
	var nftList []entity.MarketplaceBTCListingFilterPipeline
	var err error

	// if buyableOnly {
	nftList, err = u.Repo.RetrieveBTCNFTListingsUnsoldForSearch(filter, limit, offset)
	if err != nil {
		return nil, err
	}
	// } else {
	// 	nftList, err = u.Repo.RetrieveBTCNFTListings(limit, offset)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	for _, listing := range nftList {
		// if listing.IsSold {
		// 	if !buyableOnly {
		// 		nftInfo := structure.MarketplaceNFTDetail{
		// 			InscriptionID: listing.InscriptionID,
		// 			Name:          listing.Name,
		// 			Description:   listing.Description,
		// 			Price:         listing.Price,
		// 			OrderID:       listing.UUID,
		// 			IsConfirmed:   listing.IsConfirm,
		// 			Buyable:       false,
		// 			IsCompleted:   listing.IsSold,
		// 			CreatedAt:     listing.CreatedAt,
		// 		}
		// 		result = append(result, nftInfo)
		// 	}
		// 	continue
		// }
		buyOrders, err := u.Repo.GetBTCListingHaveOngoingOrder(listing.UUID)
		if err != nil {
			if !buyableOnly {
				nftInfo := structure.MarketplaceNFTDetail{
					InscriptionID: listing.InscriptionID,
					Name:          listing.Name,
					Description:   listing.Description,
					Price:         listing.Price,
					OrderID:       listing.UUID,
					IsConfirmed:   listing.IsConfirm,
					Buyable:       false,
					IsCompleted:   listing.IsSold,
					CreatedAt:     listing.CreatedAt,

					Inscription:      listing.Inscription,
					InscriptionName:  listing.InscriptionName,
					InscriptionIndex: listing.InscriptionIndex,
					CollectionID:     listing.CollectionID,
				}
				inscribeInfo, err := u.GetInscribeInfo(span, nftInfo.InscriptionID)
				if err != nil {
					log.Error("h.Usecase.GetInscribeInfo", err.Error(), err)
				}
				if inscribeInfo != nil {
					nftInfo.InscriptionNumber = inscribeInfo.Index
					nftInfo.ContentType = inscribeInfo.ContentType
					nftInfo.ContentLength = inscribeInfo.ContentLength
				}
				result = append(result, nftInfo)
				continue
			}
		}
		currentTime := time.Now()
		isAvailable := true
		for _, order := range buyOrders {
			expireTime := order.ExpiredAt
			// not expired yet still waiting for btc
			if currentTime.Before(expireTime) && (order.Status == entity.StatusBuy_Pending || order.Status == entity.StatusBuy_NotEnoughBalance) {
				isAvailable = false
				break
			}
			// could be expired but received btc
			if order.Status != entity.StatusBuy_Pending && order.Status != entity.StatusBuy_NotEnoughBalance {
				isAvailable = false
				break
			}
		}

		nftInfo := structure.MarketplaceNFTDetail{
			InscriptionID: listing.InscriptionID,
			Name:          listing.Name,
			Description:   listing.Description,
			Price:         listing.Price,
			OrderID:       listing.UUID,
			IsConfirmed:   listing.IsConfirm,
			Buyable:       isAvailable,
			IsCompleted:   listing.IsSold,
			CreatedAt:     listing.CreatedAt,

			Inscription: listing.Inscription,

			InscriptionName:  listing.InscriptionName,
			InscriptionIndex: listing.InscriptionIndex,
			CollectionID:     listing.CollectionID,
		}
		inscribeInfo, err := u.GetInscribeInfo(span, nftInfo.InscriptionID)
		if err != nil {
			log.Error("h.Usecase.GetInscribeInfo", err.Error(), err)
		}
		if inscribeInfo != nil {
			nftInfo.InscriptionNumber = inscribeInfo.Index
			nftInfo.ContentType = inscribeInfo.ContentType
			nftInfo.ContentLength = inscribeInfo.ContentLength
		}
		if buyableOnly && isAvailable {
			result = append(result, nftInfo)
		}
		if !buyableOnly {
			result = append(result, nftInfo)
		}
	}

	if !buyableOnly {
		sort.SliceStable(result, func(i, j int) bool {
			if result[i].Buyable && result[j].Buyable {
				if result[i].CreatedAt.After(result[j].CreatedAt) {
					return true
				}
			}
			return result[i].Buyable
		})
	}

	// result := []response.MarketplaceNFTDetail{}
	// for _, nft := range nfts {

	// 	result = append(result, nftInfo)
	// }
	return result, nil
}

func (u Usecase) BTCMarketplaceBuyOrder(rootSpan opentracing.Span, orderInfo structure.MarketplaceBTC_BuyOrderInfo) (string, error) {
	span, log := u.StartSpan("BTCMarketplaceListingNFT", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	order := entity.MarketplaceBTCBuyOrder{
		InscriptionID: orderInfo.InscriptionID,
		ItemID:        orderInfo.OrderID,
		OrdAddress:    orderInfo.BuyOrdAddress,
		ExpiredAt:     time.Now().Add(time.Minute * 30),
	}

	privKey, _, addressSegwit, err := btc.GenerateAddressSegwit()
	if err != nil {
		log.Error("u.OrdService.Exec.create.receive", err.Error(), err)
		return "", err
	}
	order.SegwitAddress = addressSegwit
	order.SegwitKey = privKey

	// order.HoldOrdAddress = holdOrdAddress
	// sendMessage := func(rootSpan opentracing.Span, offer entity.MarketplaceOffers) {
	// 	span, log := u.StartSpan("MakeOffer.sendMessage", rootSpan)
	// 	defer u.Tracer.FinishSpan(span, log)

	// 	profile, err := u.Repo.FindUserByWalletAddress(offer.Buyer)
	// 	if err != nil {
	// 		log.Error("cancelListing.FindUserByWalletAddress", err.Error(), err)
	// 		return
	// 	}

	// 	token, err := u.Repo.FindTokenByGenNftAddr(offer.CollectionContract, offer.TokenId)
	// 	if err != nil {
	// 		log.Error("cancelListing.FindTokenByGenNftAddr", err.Error(), err)
	// 		return
	// 	}

	// 	preText := fmt.Sprintf("[OfferID %s] has been created by %s", offer.OfferingId, offer.Buyer)
	// 	content := fmt.Sprintf("TokenID: %s", helpers.CreateTokenLink(token.ProjectID, token.TokenID, token.Name))
	// 	title := fmt.Sprintf("User %s made offer with %s", helpers.CreateProfileLink(profile.WalletAddress, profile.DisplayName), offer.Price)

	// 	if _, _, err := u.Slack.SendMessageToSlack(preText, title, content); err != nil {
	// 		log.Error("s.Slack.SendMessageToSlack err", err.Error(), err)
	// 	}

	// }

	log.SetTag("inscriptionID", order.InscriptionID)
	log.SetTag("BuyOrdAddress", order.OrdAddress)
	// check if listing is created or not
	err = u.Repo.CreateMarketplaceBuyOrder(&order)
	if err != nil {
		log.Error("BTCMarketplaceListingNFT.Repo.CreateMarketplaceListingBTC", "", err)
		return "", err
	}
	return addressSegwit, nil
}

// get filter info:
type CollectionFilter struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
	Value string `json:"value"`
	From  int64  `json:"from"`
	To    int64  `json:"to"`
}

func (u Usecase) BTCMarketplaceUpdateNftInfo(rootSpan opentracing.Span) error {
	// create data for listing data:
	// get nft list + update collection:
	listingOrder, _ := u.Repo.RetrieveBTCNFTListingsUnsold(99999, 0)

	fmt.Println("len(listingOrder): ", len(listingOrder))

	for _, v := range listingOrder {
		// get nft collection info:
		nftCollectionInfo, _ := u.Repo.FindTokenByTokenID(v.InscriptionID)
		if nftCollectionInfo != nil {
			fmt.Println("can not get nftCollectionInfo with v.InscriptionID: ", v.InscriptionID)
			_, err := u.Repo.UpdateListingCollectionInfo(v.UUID, nftCollectionInfo)
			if err != nil {
				fmt.Println("can not UpdateListingCollectionInfo err: ", err)
			} else {
				fmt.Println("update done: ", nftCollectionInfo.TokenID)
			}
		} else {
			fmt.Println("can not found id", v.InscriptionID)
		}
	}
	return nil
}

func (u Usecase) BTCMarketplaceFilterInfo(rootSpan opentracing.Span) (interface{}, error) {

	listCollectionFilterMap := map[string]CollectionFilter{}
	listPriceFilterMap := map[string]CollectionFilter{}
	listNftIDFilterMap := map[string]CollectionFilter{}

	listPriceFilterMap["range1"] = CollectionFilter{Name: "0 BTC - 2 BTC", Count: 0, Value: "0_200000000", From: 0, To: 200000000}
	listPriceFilterMap["range2"] = CollectionFilter{Name: "2 BTC - 5 BTC", Count: 0, Value: "200000000_500000000", From: 200000000, To: 500000000}
	listPriceFilterMap["range3"] = CollectionFilter{Name: "5 BTC - 7 BTC", Count: 0, Value: "500000000_700000000", From: 500000000, To: 700000000}
	listPriceFilterMap["range4"] = CollectionFilter{Name: "7 BTC - 9 BTC", Count: 0, Value: "700000000_900000000", From: 700000000, To: 900000000}
	listPriceFilterMap["range5"] = CollectionFilter{Name: "> 9 BTC", Count: 0, Value: "900000000_999999999999", From: 900000000, To: 999999999999}

	listNftIDFilterMap["range1"] = CollectionFilter{Name: "#0 - #1000", Count: 0, Value: "0_1000", From: 0, To: 1000}
	listNftIDFilterMap["range2"] = CollectionFilter{Name: "#1000 - #2000", Count: 0, Value: "1000_2000", From: 1000, To: 2000}
	listNftIDFilterMap["range3"] = CollectionFilter{Name: "#2000 - #3000", Count: 0, Value: "2000_3000", From: 2000, To: 3000}
	listNftIDFilterMap["range4"] = CollectionFilter{Name: "#3000 - #4000", Count: 0, Value: "3000_4000", From: 3000, To: 4000}
	listNftIDFilterMap["range5"] = CollectionFilter{Name: "> #4000", Count: 0, Value: "4000_9999999", From: 4000, To: 9999999}

	listCollectionFilterMapReturn := map[string]interface{}{}

	listingOrder, _ := u.Repo.RetrieveBTCNFTListingsUnsold(99999, 1)
	for _, v := range listingOrder {
		collectionFilter := CollectionFilter{
			Name:  v.CollectionName,
			Count: 1,
			Value: v.CollectionID,
		}

		val, ok := listCollectionFilterMap[v.CollectionID]
		// If the key exists
		if ok {
			collectionFilter.Count = val.Count + 1
			listCollectionFilterMap[v.CollectionID] = collectionFilter
		}
		listCollectionFilterMap[v.CollectionID] = collectionFilter

		var price int64
		if len(v.Price) > 0 {
			price, _ = strconv.ParseInt(v.Price, 10, 64)
		}

		var nftID int64
		if len(v.InscriptionIndex) > 0 {
			nftID, _ = strconv.ParseInt(v.InscriptionIndex, 10, 64)
		}

		// TODO important: if len(listPriceFilterMap) != len(listNftIDFilterMap) please use 2 for loop.
		for key, _ := range listPriceFilterMap {
			filterPrice := listPriceFilterMap[key]

			fmt.Println("filterPrice: ", filterPrice.From, filterPrice.To)

			// for price
			if listPriceFilterMap[key].From <= price && price < listPriceFilterMap[key].To {
				filterPrice.Count += 1
			}
			listPriceFilterMap[key] = filterPrice

			// for nftID:
			filterNftID := listNftIDFilterMap[key]
			// for price
			if listNftIDFilterMap[key].From < nftID && nftID <= listNftIDFilterMap[key].To {
				filterNftID.Count += 1
			}
			listNftIDFilterMap[key] = filterNftID
		}
	}

	listCollectionFilterMapReturn["collection"] = listCollectionFilterMap
	listCollectionFilterMapReturn["price"] = listPriceFilterMap
	listCollectionFilterMapReturn["inscriptionID"] = listNftIDFilterMap

	return listCollectionFilterMapReturn, nil
}
