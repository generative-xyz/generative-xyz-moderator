package usecase

import (
	"github.com/opentracing/opentracing-go"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
)

func (u Usecase) BTCMarketplaceListingNFT(rootSpan opentracing.Span, listingInfo structure.MarketplaceBTC_ListingInfo) (string, error) {
	span, log := u.StartSpan("BTCMarketplaceListingNFT", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	listing := entity.MarketplaceBTCListing{
		SellOrdAddress: listingInfo.SellOrdAddress,
		HoldOrdAddress: "",
		Price:          listingInfo.Price,
		ServiceFee:     listingInfo.ServiceFee,
		IsConfirm:      false,
		IsSold:         false,
	}
	holdOrdAddress := ""

	//TODO: gen holdOrdAddress
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
	err := u.Repo.CreateMarketplaceListingBTC(&listing)
	if err != nil {
		log.Error("BTCMarketplaceListingNFT.Repo.CreateMarketplaceListingBTC", "", err)
		return "", err
	}
	return holdOrdAddress, nil
}

func (u Usecase) BTCMarketplaceListNFT(rootSpan opentracing.Span) ([]entity.MarketplaceBTCListing, error) {
	span, log := u.StartSpan("BTCMarketplaceListingNFT", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	result := []entity.MarketplaceBTCListing{}
	return result, nil
}
