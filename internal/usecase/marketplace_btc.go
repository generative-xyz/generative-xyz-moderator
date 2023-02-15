package usecase

import (
	"sort"
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

func (u Usecase) BTCMarketplaceListNFT(rootSpan opentracing.Span, buyableOnly bool, limit, offset int64) ([]structure.MarketplaceNFTDetail, error) {
	span, log := u.StartSpan("BTCMarketplaceListingNFT", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	result := []structure.MarketplaceNFTDetail{}
	var nftList []entity.MarketplaceBTCListingFilterPipeline
	var err error

	// if buyableOnly {
	nftList, err = u.Repo.RetrieveBTCNFTListingsUnsold(limit, offset)
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
