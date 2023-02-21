package usecase

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/mongo"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/contracts/generative_marketplace_lib"
	"rederinghub.io/utils/helpers"
)

func (u Usecase) ListToken( event *generative_marketplace_lib.GenerativeMarketplaceLibListingToken, blocknumber uint64) error {

	listing := entity.MarketplaceListings{
		OfferingId: strings.ToLower(fmt.Sprintf("%x", event.OfferingId)),
		CollectionContract: strings.ToLower(event.Data.CollectionContract.String()),
		TokenId : event.Data.TokenId.String(),
		Seller: strings.ToLower(event.Data.Seller.String()),
		Erc20Token: strings.ToLower(event.Data.Erc20Token.String()),
		Price: event.Data.Price.String(),
		Closed: event.Data.Closed,
		BlockNumber: blocknumber,
		Finished: false,
		DurationTime: event.Data.DurationTime.String(),
	}

	sendMessage := func( listing entity.MarketplaceListings) {

		profile, err := u.Repo.FindUserByWalletAddress(listing.Seller)
		if err != nil {
			u.Logger.Error("cancelListing.FindUserByWalletAddress", err.Error(), err)
			return 
		}

		token, err := u.Repo.FindTokenByGenNftAddr(listing.CollectionContract, listing.TokenId)
		if err != nil {
			u.Logger.Error("cancelListing.FindTokenByGenNftAddr", err.Error(), err)
			return 
		}

		preText := fmt.Sprintf("[ListingID %s] has been created by %s", listing.OfferingId, listing.Seller)
		content := fmt.Sprintf("TokenID: %s", helpers.CreateTokenLink(token.ProjectID, token.TokenID, token.Name))
		title := fmt.Sprintf("User %s create listing with %s",helpers.CreateProfileLink(profile.WalletAddress, profile.DisplayName), listing.Price)

		if _, _, err := u.Slack.SendMessageToSlack(preText, title, content); err != nil {
			u.Logger.Error("s.Slack.SendMessageToSlack err", err.Error(), err)
		}
}

	
	// check if listing is created or not
	_, err := u.Repo.FindListingByOfferingID(listing.OfferingId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// if error is no document -> create
			err := u.Repo.CreateMarketplaceListing(&listing)
		if err != nil {
				u.Logger.Error("error when create marketplace listing", "", err)
				return err
			}

			sendMessage(listing)

			// TODO: @dac add update collection stats here

			return nil
		} else {
			return err
		}
	} else {
		// listing is already created
		u.Logger.Info("list token offeringId", listing.OfferingId)
		return errors.New("listing is already created")
	}
}

func (u Usecase) PurchaseToken( event *generative_marketplace_lib.GenerativeMarketplaceLibPurchaseToken) error {

	offeringID := strings.ToLower(fmt.Sprintf("%x", event.OfferingId))
	u.Logger.Info("purchase token offeringId", offeringID)
	
	err := u.Repo.PurchaseTokenByOfferingID(offeringID)
	if err != nil {
		u.Logger.Error("u.PurchaseToken.AcceptOfferByOfferingID", err.Error(), err)
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

err = u.UpdateTokenOnwer("purchased", offeringID,  getToken, event.Buyer)
	if err != nil {
		u.Logger.Error("u.PurchaseToken.UpdateTokenOnwer", err.Error(), err)
		return err
	}

	return nil
}

func (u Usecase) MakeOffer( event *generative_marketplace_lib.GenerativeMarketplaceLibMakeOffer, blocknumber uint64) error {

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
		BlockNumber: blocknumber,
	}

	sendMessage := func( offer entity.MarketplaceOffers) {

		profile, err := u.Repo.FindUserByWalletAddress(offer.Buyer)
		if err != nil {
			u.Logger.Error("cancelListing.FindUserByWalletAddress", err.Error(), err)
			return 
		}

		token, err := u.Repo.FindTokenByGenNftAddr(offer.CollectionContract, offer.TokenId)
		if err != nil {
			u.Logger.Error("cancelListing.FindTokenByGenNftAddr", err.Error(), err)
			return 
		}

		preText := fmt.Sprintf("[OfferID %s] has been created by %s", offer.OfferingId, offer.Buyer)
		content := fmt.Sprintf("TokenID: %s", helpers.CreateTokenLink(token.ProjectID, token.TokenID, token.Name))
		title := fmt.Sprintf("User %s made offer with %s",helpers.CreateProfileLink(profile.WalletAddress, profile.DisplayName), offer.Price)

		if _, _, err := u.Slack.SendMessageToSlack(preText, title, content); err != nil {
			u.Logger.Error("s.Slack.SendMessageToSlack err", err.Error(), err)
		}
}

	
	// check if listing is created or not
	_, err := u.Repo.FindOfferByOfferingID(offer.OfferingId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// if error is no document -> create
			err := u.Repo.CreateMarketplaceOffer(&offer)
		if err != nil {
				u.Logger.Error("makeOffer.Repo.CreateMarketplaceOffer", "", err)
				return err
			}

			sendMessage(offer)
			// TODO: @dac add update collection stats here
			return nil
		} else {
			u.Logger.Error("makeOffer.Repo.FindOfferByOfferingID", "", err)
			return err
		}
	} else {
		err := errors.New("offer is already created")
		// listing is already created
		u.Logger.Error("offer token offeringId", err.Error(), err)
		return err
	}

	return nil
}

func (u Usecase) AcceptMakeOffer( event *generative_marketplace_lib.GenerativeMarketplaceLibAcceptMakeOffer) error {

	offeringID := strings.ToLower(fmt.Sprintf("%x", event.OfferingId))
	u.Logger.Info("accept make offer offeringId", offeringID)
	
	err := u.Repo.AcceptOfferByOfferingID(offeringID)
	if err != nil {
		u.Logger.Error("u.AcceptMakeOffer.AcceptOfferByOfferingID", err.Error(), err)
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
	err = u.UpdateTokenOnwer("accepted", offeringID, getToken,event.Buyer)
	if err != nil {
		u.Logger.Error("u.UpdateTokenOnwer.UpdateTokenOnwer", err.Error(), err)
		return err
	}

	// TODO: @dac add update collection stats here
	return nil
}

func (u Usecase) CancelListing( event *generative_marketplace_lib.GenerativeMarketplaceLibCancelListing) error {

	offeringID := strings.ToLower(fmt.Sprintf("%x", event.OfferingId))
	u.Logger.Info("cancel listing offeringId", offeringID)
	
	err := u.Repo.CancelListingByOfferingID(offeringID)
	if err != nil {
		return err
	}

	done := make(chan bool)
	go func ( done chan bool)  {
		defer func ()  {
			done <- true
		}()

		listing, err := u.Repo.FindListingByOfferingID(offeringID)
		if err != nil {
			u.Logger.Error("cancelListing.FindListingByOfferingID", err.Error(), err)
			return 
		}

		profile, err := u.Repo.FindUserByWalletAddress(listing.Seller)
		if err != nil {
			u.Logger.Error("cancelListing.FindUserByWalletAddress", err.Error(), err)
			return 
		}

		token, err := u.Repo.FindTokenByGenNftAddr(listing.CollectionContract, listing.TokenId)
		if err != nil {
			u.Logger.Error("cancelListing.FindTokenByGenNftAddr", err.Error(), err)
			return 
		}

		preText := fmt.Sprintf("[Listing %s] has been cancelled by %s", listing.OfferingId, listing.Seller)
		content := fmt.Sprintf("TokenID: %s", helpers.CreateTokenLink(token.ProjectID, token.TokenID, token.Name))
		title := fmt.Sprintf("User %s cancelled offer %s", helpers.CreateProfileLink(profile.WalletAddress, profile.DisplayName) , offeringID)

		if _, _, err := u.Slack.SendMessageToSlack(preText, title, content); err != nil {
			u.Logger.Error("s.Slack.SendMessageToSlack err", err.Error(), err)
		}}(done)
	<- done

	// TODO: @dac add update collection stats here
	return nil
}

func (u Usecase) CancelOffer( event *generative_marketplace_lib.GenerativeMarketplaceLibCancelMakeOffer) error {

	offeringID := strings.ToLower(fmt.Sprintf("%x", event.OfferingId))
	u.Logger.Info("cancel make offer offeringId", offeringID)
	
	err := u.Repo.CancelOfferByOfferingID(offeringID)
	if err != nil {
		u.Logger.Error("s.Repo.CancelOfferByOfferingID", err.Error(), err)
		return err
	}

	done := make(chan bool)
	go func ( done chan bool)  {
		defer func ()  {
			done <- true
		}()

		offer, err := u.Repo.FindOfferByOfferingID(offeringID)
		if err != nil {
			u.Logger.Error("s.Repo.FindOfferByOfferingID", err.Error(), err)
			return 
		}

		profile, err := u.Repo.FindUserByWalletAddress(offer.Buyer)
		if err != nil {
			u.Logger.Error("cancelListing.FindUserByWalletAddress", err.Error(), err)
			return 
		}

		token, err := u.Repo.FindTokenByGenNftAddr(offer.CollectionContract, offer.TokenId)
		if err != nil {
			u.Logger.Error("cancelListing.FindTokenByGenNftAddr", err.Error(), err)
			return 
		}

		preText := fmt.Sprintf("[Listing %s] has been cancelled by %s", offer.OfferingId, offer.Buyer)
		content := fmt.Sprintf("TokenID: %s",  helpers.CreateTokenLink(token.ProjectID, token.TokenID, token.Name))
		title := fmt.Sprintf("User %s cancelled offer %s", helpers.CreateProfileLink(profile.WalletAddress, profile.DisplayName) , offeringID)

		if _, _, err := u.Slack.SendMessageToSlack(preText, title, content); err != nil {
			u.Logger.Error("s.Slack.SendMessageToSlack err", err.Error(), err)
		}	

	}(done)
	<- done

	

	// TODO: @dac add update collection stats here
	return nil
}

func (u Usecase) FilterMKListing( filter structure.FilterMkListing) (*entity.Pagination, error) {

fm := &entity.FilterMarketplaceListings{}
	err := copier.Copy(fm, filter)
	if err != nil {
		u.Logger.Error("copier.Copy", err.Error(), err)
		return nil, err
	}
ml, err := u.Repo.FilterMarketplaceListings(*fm)
	if err != nil {
		u.Logger.Error("u.Repo.FilterMarketplaceListings", err.Error(), err)
		return nil, err
	}

	u.Logger.Info("filtered", ml)

	listingResp := []entity.MarketplaceListings{}
	iListings := ml.Result
	listing := iListings.([]entity.MarketplaceListings) 
	for _, listingItem := range listing {
		tok, err := u.Repo.FindTokenByGenNftAddr(listingItem.CollectionContract, listingItem.TokenId)
		if err != nil  {
			u.Logger.Error("u.Repo.FindTokenByGenNftAddr", err.Error(), err)
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

func (u Usecase) FilterMKOffers( filter structure.FilterMkOffers) (*entity.Pagination, error) {

fm := &entity.FilterMarketplaceOffers{}
	err := copier.Copy(fm, filter)
	if err != nil {
		u.Logger.Error("copier.Copy", err.Error(), err)
		return nil, err
	}
ml, err := u.Repo.FilterMarketplaceOffers(*fm)
	if err != nil {
		u.Logger.Error("u.Repo.FilterMarketplaceOffers", err.Error(), err)
		return nil, err
	}

	u.Logger.Info("filtered", ml)

	offersResp := []entity.MarketplaceOffers{}
	iOffers := ml.Result
	offers := iOffers.([]entity.MarketplaceOffers) 
	for _, offer := range offers {
		tok, err := u.Repo.FindTokenByGenNftAddr(offer.CollectionContract, offer.TokenId)
		if err != nil  {
			u.Logger.Error("u.Repo.FindTokenByGenNftAddr", err.Error(), err)
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

func (u Usecase) GetListingBySeller( sellerAddress string) ([]entity.MarketplaceListings,[]string,[]string, error) {

	cachedKey, cachedContractIDsKey, cachedTokensIDsKey := helpers.ProfileSelingKey(sellerAddress)
listings := []entity.MarketplaceListings{}
	contractIDS := []string{}
	tokenIDS := []string{}
	var err error

	
	u.Logger.Info("cachedKey", cachedKey)
	u.Logger.Info("cachedContractIDsKey", cachedContractIDsKey)
	u.Logger.Info("cachedTokensIDsKey", cachedTokensIDsKey)

	//always reloa data
	liveReload := func ( sellerAddress string, cachedKey string, cachedContractIDsKey string, cachedTokensIDsKey string) ([]entity.MarketplaceListings, []string, []string, error)  {

		listings, err = u.Repo.GetListingBySeller(sellerAddress)
		if err != nil {
			u.Logger.Error("u.Repo.GetListingBySeller", err.Error(), err)
			return nil, nil, nil, err
		}

		contractIDS := []string{}
		tokenIDS := []string{}
		for key, listing := range listings {
			u.Logger.Info(fmt.Sprintf("listing.%d",key),listing)
			contractIDS = append(contractIDS, listing.CollectionContract)
			tokenIDS = append(tokenIDS, listing.TokenId)
		}
		u.Cache.SetData(cachedKey, listings)
		u.Cache.SetData(cachedContractIDsKey, contractIDS)
		u.Cache.SetData(cachedTokensIDsKey, tokenIDS)
		return listings, contractIDS, tokenIDS, nil
	}

	go liveReload(sellerAddress, cachedKey, cachedContractIDsKey, cachedTokensIDsKey)

	cached, err := u.Cache.GetData(cachedKey)
	if err == nil && cached != nil {
		err = helpers.ParseCache(cached, &listings)
		if err != nil  {
			u.Logger.Error("helpers.ParseCache.listings", err.Error(), err)
			return nil, nil, nil, err
		}
		cached, err := u.Cache.GetData(cachedContractIDsKey)
		err = helpers.ParseCache(cached, &contractIDS)
		if err != nil  {
			u.Logger.Error("helpers.ParseCache.cachedContractIDsKey", err.Error(), err)
			return nil, nil, nil, err
		}
		cached, err = u.Cache.GetData(cachedTokensIDsKey)
		err = helpers.ParseCache(cached, &tokenIDS)
		if err != nil  {
			u.Logger.Error("helpers.ParseCache.tokenIDS", err.Error(), err)
			return nil, nil, nil, err
		}

	}else{
		listings, contractIDS, tokenIDS, err = liveReload(sellerAddress, cachedKey, cachedContractIDsKey, cachedTokensIDsKey)
		if err != nil  {
			u.Logger.Error("liveReload", err.Error(), err)
			return nil, nil, nil, err
		}
	}

	u.Logger.Info("listings", listings)
	u.Logger.Info("contractIDS", contractIDS)
	u.Logger.Info("tokenIDS", tokenIDS)
	return listings, contractIDS, tokenIDS, nil
}

func (u Usecase) UpdateTokenOnwer(event string, offeringID string , fn func(offeringID string) (*entity.TokenUri, error), buyer common.Address) error {


	owner := strings.ToLower(buyer.String())
	token, err := fn(offeringID)
	if err != nil {
		u.Logger.Error("UpdateTokenOnwer.fn", err.Error(), err)
		return err
	}

	u.Logger.Info("tokenID", token.TokenID)
	
	

	profile, err := u.Repo.FindUserByWalletAddress(owner)
	if err != nil {
		// if can not find user profile in db, set owner to nil
		if errors.Is(err, mongo.ErrNoDocuments) {
			profile = nil
		} else {
			u.Logger.Error("UpdateTokenOnwer.FindUserByWalletAddress", err.Error(), err)
			return err
		}
	}

	u.Logger.Info("token.Owner", owner)
	
	token.Owner = profile
	token.OwnerAddr = owner

	updated, err := u.Repo.UpdateOrInsertTokenUri(token.ContractAddress, token.TokenID, token)
	if err != nil {
		u.Logger.Error("UpdateTokenOnwer.UpdateOrInsertTokenUri", err.Error(), err)
		return err
	}

	u.Logger.Info("updated",updated)

	//slack
	preText := fmt.Sprintf("[TokenID %s] has been transfered to %s", token.TokenID, token.OwnerAddr)
	content := fmt.Sprintf("To user: %s. Token: %s", helpers.CreateProfileLink(owner,  profile.DisplayName),  helpers.CreateTokenLink( token.ProjectID, token.TokenID,  token.Name))
	title := fmt.Sprintf("OfferingID:  %s is %s", offeringID, event)

	if _, _, err := u.Slack.SendMessageToSlack(preText, title, content); err != nil {
		u.Logger.Error("s.Slack.SendMessageToSlack err", err.Error(), err)
	}

	// TODO: @dac add update collection stats here
	

	return nil
}
