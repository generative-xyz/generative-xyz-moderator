package usecase

import (
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
		ServiceFee:     listingInfo.ServiceFee,
		IsConfirm:      false,
		IsSold:         false,
		ExpiredAt:      time.Now().Add(time.Hour * 2),
		Name:           listingInfo.Name,
		Description:    listingInfo.Description,
		InscriptionID:  listingInfo.InscriptionID,
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

func (u Usecase) BTCMarketplaceListNFT(rootSpan opentracing.Span) ([]entity.MarketplaceBTCListing, error) {
	span, log := u.StartSpan("BTCMarketplaceListingNFT", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	result := []entity.MarketplaceBTCListing{}

	// test1 := entity.MarketplaceBTCListing{
	// 	InscriptionID: "c0f8acd8f0d91d490ac9c08977b142aa836207d2ee93d111992866cf47a6d2e6i0",
	// 	Name:          "Test1",
	// 	Description:   "test1 blah blah blah",
	// 	Price:         "1234567",
	// 	BaseEntity: entity.BaseEntity{
	// 		UUID: "1",
	// 	},
	// }

	// test2 := entity.MarketplaceBTCListing{
	// 	InscriptionID: "2696948882cc088f2d1c160981501a48b3744d8d5df0e8d9a71557e716c634dci0",
	// 	Name:          "Test2",
	// 	Description:   "test2 blah blah blah",
	// 	Price:         "1234567",
	// 	BaseEntity: entity.BaseEntity{
	// 		UUID: "2",
	// 	},
	// }

	// test3 := entity.MarketplaceBTCListing{
	// 	InscriptionID: "95752b856f94d0c60bee700d6df1b47c949c28f2a06859cf6d5a3466843463b8i0",
	// 	Name:          "Test3",
	// 	Description:   "test3 blah blah blah",
	// 	Price:         "1234567",
	// 	BaseEntity: entity.BaseEntity{
	// 		UUID: "3",
	// 	},
	// }

	// result = append(result, test1)
	// result = append(result, test2)
	// result = append(result, test3)

	nftList, err := u.Repo.RetrieveBTCNFTListings()
	if err != nil {
		return nil, err
	}

	for _, listing := range nftList {
		// err := u.Repo.CheckBTCListingHaveOngoingOrder(listing.UUID)
		// if err != nil {
		// 	continue
		// }
		result = append(result, listing)
	}
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
