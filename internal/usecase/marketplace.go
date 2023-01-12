package usecase

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/copier"
	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/mongo"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/contracts/generative_marketplace_lib"
	"rederinghub.io/utils/helpers"
)

func (u Usecase) ListToken(rootSpan opentracing.Span, event *generative_marketplace_lib.GenerativeMarketplaceLibListingToken) error {
	span, log := u.StartSpan("ListToken", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	listing := entity.MarketplaceListings{
		OfferingId: strings.ToLower(fmt.Sprintf("%x", event.OfferingId)),
		CollectionContract: strings.ToLower(event.Data.CollectionContract.String()),
		TokenId : event.Data.TokenId.String(),
		Seller: strings.ToLower(event.Data.Seller.String()),
		Erc20Token: strings.ToLower(event.Data.Erc20Token.String()),
		Price: event.Data.Price.String(),
		Closed: event.Data.Closed,
		Finished: false,
		DurationTime: event.Data.DurationTime.String(),
	}

	sendMessage := func(rootSpan opentracing.Span, listing entity.MarketplaceListings) {
		span, log := u.StartSpan("MakeListing.sendMessage", rootSpan)
		defer u.Tracer.FinishSpan(span, log)

		preText := fmt.Sprintf("[ListingID %s] has been created by %s", listing.OfferingId, listing.Seller)
		content := fmt.Sprintf("TokenID: %s, collection: %s", listing.TokenId, listing.CollectionContract)
		title := fmt.Sprintf("User %s create listing with %s",listing.Seller, listing.Price)

		if _, _, err := u.Slack.SendMessageToSlack(preText, title, content); err != nil {
			log.Error("s.Slack.SendMessageToSlack err", err.Error(), err)
		}
	
	}

	log.SetTag("offeringID", listing.OfferingId)
	// check if listing is created or not
	_, err := u.Repo.FindListingByOfferingID(listing.OfferingId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// if error is no document -> create
			err := u.Repo.CreateMarketplaceListing(&listing)
	
			if err != nil {
				log.Error("error when create marketplace listing", "", err)
				return err
			}

			sendMessage(span, listing)

			// TODO: @dac add update collection stats here

			return nil
		} else {
			return err
		}
	} else {
		// listing is already created
		log.SetData("list token offeringId", listing.OfferingId)
		return errors.New("listing is already created")
	}
}

func (u Usecase) PurchaseToken(rootSpan opentracing.Span, event *generative_marketplace_lib.GenerativeMarketplaceLibPurchaseToken) error {
	span, log := u.StartSpan("PurchaseToken", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	offeringID := strings.ToLower(fmt.Sprintf("%x", event.OfferingId))
	log.SetData("purchase token offeringId", offeringID)
	log.SetTag("offeringID", offeringID)
	err := u.Repo.PurchaseTokenByOfferingID(offeringID)
	if err != nil {
		log.Error("u.PurchaseToken.AcceptOfferByOfferingID", err.Error(), err)
		return err
	}

	getToken := func(offeringID string) (*entity.TokenUri, error) {
		listing, err := u.Repo.FindListingByOfferingID(offeringID)
		if err != nil {
			return nil, err
		}

		token, err := u.Repo.FindTokenByGenNftAddr(listing.CollectionContract, listing.TokenId)
		if err != nil {
			return nil, err
		}

		return token, nil
	}

	
	err = u.UpdateTokenOnwer(span, "purchased", offeringID,  getToken, event.Buyer)
	if err != nil {
		log.Error("u.PurchaseToken.UpdateTokenOnwer", err.Error(), err)
		return err
	}

	return nil
}

func (u Usecase) MakeOffer(rootSpan opentracing.Span, event *generative_marketplace_lib.GenerativeMarketplaceLibMakeOffer) error {
	span, log := u.StartSpan("MakeOffer", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	offer := entity.MarketplaceOffers{
		OfferingId: strings.ToLower(fmt.Sprintf("%x", event.OfferingId)),
		CollectionContract: strings.ToLower(event.Data.CollectionContract.String()),
		TokenId : event.Data.TokenId.String(),
		Buyer: strings.ToLower(event.Data.Buyer.String()),
		Erc20Token: strings.ToLower(event.Data.Erc20Token.String()),
		Price: event.Data.Price.String(),
		Closed: event.Data.Closed,
		Finished: false,
		DurationTime: event.Data.DurationTime.String(),
	}

	sendMessage := func(rootSpan opentracing.Span, offer entity.MarketplaceOffers) {
		span, log := u.StartSpan("MakeOffer.sendMessage", rootSpan)
		defer u.Tracer.FinishSpan(span, log)

		preText := fmt.Sprintf("[OfferID %s] has been created by %s", offer.OfferingId, offer.Buyer)
		content := fmt.Sprintf("TokenID: %s, collection: %s, amount: %s", offer.TokenId, offer.CollectionContract, offer.Price)
		title := fmt.Sprintf("User %s made offer",offer.Buyer)

		if _, _, err := u.Slack.SendMessageToSlack(preText, title, content); err != nil {
			log.Error("s.Slack.SendMessageToSlack err", err.Error(), err)
		}
	
	}

	log.SetTag("offeringID", offer.OfferingId)
	// check if listing is created or not
	_, err := u.Repo.FindOfferByOfferingID(offer.OfferingId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// if error is no document -> create
			err := u.Repo.CreateMarketplaceOffer(&offer)
	
			if err != nil {
				log.Error("makeOffer.Repo.CreateMarketplaceOffer", "", err)
				return err
			}

			sendMessage(span, offer)
			// TODO: @dac add update collection stats here
			return nil
		} else {
			log.Error("makeOffer.Repo.FindOfferByOfferingID", "", err)
			return err
		}
	} else {
		
		err := errors.New("offer is already created")
		// listing is already created
		log.Error("offer token offeringId", err.Error(), err)
		return err
	}

	return nil
}

func (u Usecase) AcceptMakeOffer(rootSpan opentracing.Span, event *generative_marketplace_lib.GenerativeMarketplaceLibAcceptMakeOffer) error {
	span, log := u.StartSpan("AcceptMakeOffer", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	offeringID := strings.ToLower(fmt.Sprintf("%x", event.OfferingId))
	log.SetData("accept make offer offeringId", offeringID)
	log.SetTag("offeringID", offeringID)
	err := u.Repo.AcceptOfferByOfferingID(offeringID)
	if err != nil {
		log.Error("u.AcceptMakeOffer.AcceptOfferByOfferingID", err.Error(), err)
		return err
	}

	getToken := func(offeringID string) (*entity.TokenUri, error) {
		
		listing, err := u.Repo.FindOfferByOfferingID(offeringID)
		if err != nil {
			return nil, err
		}

		token, err := u.Repo.FindTokenByGenNftAddr(listing.CollectionContract, listing.TokenId)
		if err != nil {
			return nil, err
		}

		return token, nil
	}
	
	err = u.UpdateTokenOnwer(span,"accepted", offeringID, getToken,event.Buyer)
	if err != nil {
		log.Error("u.UpdateTokenOnwer.UpdateTokenOnwer", err.Error(), err)
		return err
	}

	// TODO: @dac add update collection stats here
	return nil
}

func (u Usecase) CancelListing(rootSpan opentracing.Span, event *generative_marketplace_lib.GenerativeMarketplaceLibCancelListing) error {
	span, log := u.StartSpan("CancelListing", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	offeringID := strings.ToLower(fmt.Sprintf("%x", event.OfferingId))
	log.SetData("cancel listing offeringId", offeringID)
	log.SetTag("offeringID", offeringID)
	err := u.Repo.CancelListingByOfferingID(offeringID)
	if err != nil {
		return err
	}

	offer, err := u.Repo.FindListingByOfferingID(offeringID)
	if err != nil {
		log.Error("cancelListing.FindListingByOfferingID", err.Error(), err)
		return err
	}

	preText := fmt.Sprintf("[OfferID %s] has been cancelled by %s", offer.OfferingId, offer.Seller)
	content := fmt.Sprintf("TokenID: %s, collection: %s", offer.TokenId, offer.CollectionContract)
	title := fmt.Sprintf("User %s cancelled offer %s",offer.Seller, offeringID)

	if _, _, err := u.Slack.SendMessageToSlack(preText, title, content); err != nil {
		log.Error("s.Slack.SendMessageToSlack err", err.Error(), err)
	}

	// TODO: @dac add update collection stats here
	return nil
}

func (u Usecase) CancelOffer(rootSpan opentracing.Span, event *generative_marketplace_lib.GenerativeMarketplaceLibCancelMakeOffer) error {
	span, log := u.StartSpan("CancelMakeOffer", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	offeringID := strings.ToLower(fmt.Sprintf("%x", event.OfferingId))
	log.SetData("cancel make offer offeringId", offeringID)
	log.SetTag("offeringID", offeringID)
	err := u.Repo.CancelOfferByOfferingID(offeringID)
	if err != nil {
		log.Error("s.Repo.CancelOfferByOfferingID", err.Error(), err)
		return err
	}

	offer, err := u.Repo.FindOfferByOfferingID(offeringID)
	if err != nil {
		log.Error("s.Repo.FindOfferByOfferingID", err.Error(), err)
		return err
	}

	preText := fmt.Sprintf("[OfferID %s] has been cancelled by %s", offer.OfferingId, offer.Buyer)
	content := fmt.Sprintf("TokenID: %s, collection: %s", offer.TokenId, offer.CollectionContract)
	title := fmt.Sprintf("User %s cancelled offer %s",offer.Buyer, offeringID)

	if _, _, err := u.Slack.SendMessageToSlack(preText, title, content); err != nil {
		log.Error("s.Slack.SendMessageToSlack err", err.Error(), err)
	}

	// TODO: @dac add update collection stats here
	return nil
}

func (u Usecase) FilterMKListing(rootSpan opentracing.Span, filter structure.FilterMkListing) (*entity.Pagination, error) {
	span, log := u.StartSpan("FilterListing", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	
	fm := &entity.FilterMarketplaceListings{}
	err := copier.Copy(fm, filter)
	if err != nil {
		log.Error("copier.Copy", err.Error(), err)
		return nil, err
	}
	
	ml, err := u.Repo.FilterMarketplaceListings(*fm)
	if err != nil {
		log.Error("u.Repo.FilterMarketplaceListings", err.Error(), err)
		return nil, err
	}

	log.SetData("filtered", ml)

	listingResp := []entity.MarketplaceListings{}
	iListings := ml.Result
	listing := iListings.([]entity.MarketplaceListings) 
	for _, listingItem := range listing {
		
		tok, err := u.Repo.FindTokenByGenNftAddr(listingItem.CollectionContract, listingItem.TokenId)
		if err != nil  {
			log.Error("u.Repo.FindTokenByGenNftAddr", err.Error(), err)
		}else{
			listingItem.Token = *tok
		}

		p, err  := u.Repo.FindUserByWalletAddress(listingItem.Seller)
		if err ==  nil {
			listingItem.SellerInfo =  *p
		}

		listingResp =append(listingResp, listingItem)
	}
	
	ml.Result = listingResp
	return ml, nil
	
}

func (u Usecase) FilterMKOffers(rootSpan opentracing.Span, filter structure.FilterMkOffers) (*entity.Pagination, error) {
	span, log := u.StartSpan("FilterListing", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	
	fm := &entity.FilterMarketplaceOffers{}
	err := copier.Copy(fm, filter)
	if err != nil {
		log.Error("copier.Copy", err.Error(), err)
		return nil, err
	}
	
	ml, err := u.Repo.FilterMarketplaceOffers(*fm)
	if err != nil {
		log.Error("u.Repo.FilterMarketplaceOffers", err.Error(), err)
		return nil, err
	}

	log.SetData("filtered", ml)

	offersResp := []entity.MarketplaceOffers{}
	iOffers := ml.Result
	offers := iOffers.([]entity.MarketplaceOffers) 
	for _, offer := range offers {
		
		tok, err := u.Repo.FindTokenByGenNftAddr(offer.CollectionContract, offer.TokenId)
		if err != nil  {
			log.Error("u.Repo.FindTokenByGenNftAddr", err.Error(), err)
			//continue
		}else{
			offer.Token = *tok
		}

		p, err  := u.Repo.FindUserByWalletAddress(offer.Buyer)
		if err ==  nil {
			offer.BuyerInfo =  *p
		}
		
		
		offersResp =append(offersResp, offer)
	}

	ml.Result = offersResp
	return ml, nil
}

func (u Usecase) GetListingBySeller(rootSpan opentracing.Span, sellerAddress string) ([]entity.MarketplaceListings,[]string,[]string, error) {
	span, log := u.StartSpan("GetListingBySeller", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	cachedKey, cachedContractIDsKey, cachedTokensIDsKey := helpers.ProfileSelingKey(sellerAddress)
	
	listings := []entity.MarketplaceListings{}
	contractIDS := []string{}
	tokenIDS := []string{}
	var err error

	log.SetTag(utils.WALLET_ADDRESS_TAG, sellerAddress)
	log.SetData("cachedKey", cachedKey)
	log.SetData("cachedContractIDsKey", cachedContractIDsKey)
	log.SetData("cachedTokensIDsKey", cachedTokensIDsKey)

	//always reloa data
	liveReload := func (rootSpan opentracing.Span, sellerAddress string, cachedKey string, cachedContractIDsKey string, cachedTokensIDsKey string) ([]entity.MarketplaceListings, []string, []string, error)  {
		span, log := u.StartSpan("FilterListing.Live.Reload", rootSpan)
		defer u.Tracer.FinishSpan(span, log)

		listings, err = u.Repo.GetListingBySeller(sellerAddress)
		if err != nil {
			log.Error("u.Repo.GetListingBySeller", err.Error(), err)
			return nil, nil, nil, err
		}

		contractIDS := []string{}
		tokenIDS := []string{}
		for key, listing := range listings {
			log.SetData(fmt.Sprintf("listing.%d",key),listing)
			contractIDS = append(contractIDS, listing.CollectionContract)
			tokenIDS = append(tokenIDS, listing.TokenId)
		}
		
		
		u.Cache.SetData(cachedKey, listings)
		u.Cache.SetData(cachedContractIDsKey, contractIDS)
		u.Cache.SetData(cachedTokensIDsKey, tokenIDS)
		return listings, contractIDS, tokenIDS, nil
	}

	go liveReload(span, sellerAddress, cachedKey, cachedContractIDsKey, cachedTokensIDsKey)

	cached, err := u.Cache.GetData(cachedKey)
	if err == nil && cached != nil {
		err = helpers.ParseCache(cached, &listings)
		if err != nil  {
			log.Error("helpers.ParseCache.listings", err.Error(), err)
			return nil, nil, nil, err
		}
		
		cached, err := u.Cache.GetData(cachedContractIDsKey)
		err = helpers.ParseCache(cached, &contractIDS)
		if err != nil  {
			log.Error("helpers.ParseCache.cachedContractIDsKey", err.Error(), err)
			return nil, nil, nil, err
		}
		
		cached, err = u.Cache.GetData(cachedTokensIDsKey)
		err = helpers.ParseCache(cached, &tokenIDS)
		if err != nil  {
			log.Error("helpers.ParseCache.tokenIDS", err.Error(), err)
			return nil, nil, nil, err
		}

	}else{
		listings, contractIDS, tokenIDS, err = liveReload(span, sellerAddress, cachedKey, cachedContractIDsKey, cachedTokensIDsKey)
		if err != nil  {
			log.Error("liveReload", err.Error(), err)
			return nil, nil, nil, err
		}
	}

	log.SetData("listings", listings)
	log.SetData("contractIDS", contractIDS)
	log.SetData("tokenIDS", tokenIDS)
	return listings, contractIDS, tokenIDS, nil
}

func (u Usecase) UpdateTokenOnwer(rootSpan opentracing.Span,event string, offeringID string , fn func(offeringID string) (*entity.TokenUri, error), buyer common.Address) error {
	span, log := u.StartSpan("UpdateTokenOnwer", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	owner := strings.ToLower(buyer.String())
	token, err := fn(offeringID)
	if err != nil {
		log.Error("UpdateTokenOnwer.fn", err.Error(), err)
		return err
	}

	log.SetData("tokenID", token.TokenID)
	log.SetTag("tokenID", token.TokenID)
	log.SetTag("contractAddress", token.GenNFTAddr)

	profile, err := u.Repo.FindUserByWalletAddress(owner)
	if err != nil {
		log.Error("UpdateTokenOnwer.FindUserByWalletAddress", err.Error(), err)
		return err
	}

	oldOwner := token.OwnerAddr
	log.SetData("token.Owner", owner)
	log.SetData(utils.WALLET_ADDRESS_TAG, owner)
	token.Owner = profile
	token.OwnerAddr = owner

	updated, err := u.Repo.UpdateOrInsertTokenUri(token.ContractAddress, token.TokenID, token)
	if err != nil {
		log.Error("UpdateTokenOnwer.UpdateOrInsertTokenUri", err.Error(), err)
		return err
	}

	log.SetData("updated",updated)

	//slack
	preText := fmt.Sprintf("[TokenID %s] has been transfered to %s", token.TokenID, token.OwnerAddr)
	content := fmt.Sprintf("From address: %s - to address: %s", oldOwner, owner)
	title := fmt.Sprintf("OfferingID:  %s is %s", offeringID, event)

	if _, _, err := u.Slack.SendMessageToSlack(preText, title, content); err != nil {
		log.Error("s.Slack.SendMessageToSlack err", err.Error(), err)
	}

	// TODO: @dac add update collection stats here
	

	return nil
}