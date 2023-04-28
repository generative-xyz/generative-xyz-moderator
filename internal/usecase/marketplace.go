package usecase

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-resty/resty/v2"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/algolia"
	"rederinghub.io/utils/contracts/generative_marketplace_lib"
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/logger"
)

func (uc Usecase) SubCollectionItem(bf *structure.BaseFilters, numberFrom, numberTo int) ([]*entity.ItemListing, error) {
	filter := &algolia.AlgoliaFilter{
		ObjType: "inscription", Page: int(bf.Page), Limit: int(bf.Limit),
		FromNumber: numberFrom, ToNumber: numberTo,
	}

	inscriptionIds := []string{}
	ids := []string{}
	inscriptions, err := uc.AlgoliaSearchInscriptionFromTo(filter)
	if err != nil {
		return nil, err
	}
	for _, i := range inscriptions {
		inscriptionIds = append(inscriptionIds, i.InscriptionId)
		ids = append(ids, i.ObjectId)
	}

	wG := &sync.WaitGroup{}
	mapTokenUri := make(map[string]entity.TokenUri)

	wG.Add(1)
	go func() {
		defer wG.Done()
		pe := &entity.FilterTokenUris{Ids: ids}
		pe.Page = 1
		pe.Limit = int64(len(ids))
		tokens, _ := uc.Repo.FilterTokenUri(*pe)

		iTokens := tokens.Result
		rTokens := iTokens.([]entity.TokenUri)

		for _, t := range rTokens {
			mapTokenUri[t.TokenID] = t
		}
	}()

	mapVolume := make(map[string]*entity.ItemListing)
	wG.Add(1)
	go func() {
		defer wG.Done()
		bf.Page = 1
		bf.Limit = int64(len(inscriptionIds))
		data, _ := uc.Repo.ListSubClollectionItem(bf, inscriptionIds)
		for _, d := range data {
			mapVolume[d.InscriptionId] = d
		}
	}()

	mapListing := make(map[string][]*entity.DexBTCListing)
	mapUser := make(map[string]entity.Users)
	wG.Add(1)
	go func() {
		defer wG.Done()
		addresses := []string{}
		listings, _ := uc.Repo.GetDexBTCTrackingByInscriptionIds(inscriptionIds)
		for _, d := range listings {
			if d.SellerAddress != "" {
				addresses = append(addresses, d.SellerAddress)
			}
			if d.Buyer != "" {
				addresses = append(addresses, d.Buyer)
			}

			if _, ok := mapListing[d.InscriptionID]; !ok {
				mapListing[d.InscriptionID] = []*entity.DexBTCListing{}
			}
			mapListing[d.InscriptionID] = append(mapListing[d.InscriptionID], d)
		}

		users, _ := uc.Repo.FindUserByAddresses(addresses)
		for _, u := range users {
			mapUser[u.WalletAddressBTCTaproot] = u
		}
	}()

	wG.Wait()

	result := []*entity.ItemListing{}
	waitG := &sync.WaitGroup{}
	mu := sync.Mutex{}
	client := resty.New()
	for _, i := range inscriptions {
		waitG.Add(1)
		go func(i *response.SearhcInscription) {
			defer waitG.Done()
			r := &entity.ItemListing{}
			if v, ok := mapVolume[i.InscriptionId]; ok {
				r.VolumeOneWeek = v.VolumeOneWeek
				r.VolumeOneDay = v.VolumeOneDay
				r.VolumeOneHour = v.VolumeOneHour
				r.VolumeOneMonth = v.VolumeOneMonth
				r.InscriptionId = v.InscriptionId
				r.Image = v.Image
			} else {
				r.InscriptionId = i.InscriptionId
				r.Image = fmt.Sprintf("https://generativeexplorer.com/preview/%s", r.InscriptionId)
			}

			if t, ok := mapTokenUri[r.InscriptionId]; ok {
				r.OrderInscriptionIndex = float64(t.OrderInscriptionIndex)
			}

			if btcs, ok := mapListing[r.InscriptionId]; ok {
				for _, btc := range btcs {
					if r.SellerAddress == "" {
						r.SellerAddress = btc.SellerAddress
					}

					if r.BuyerAddress == "" {
						r.BuyerAddress = btc.Buyer
					}
				}
			}

			if u, ok := mapUser[r.BuyerAddress]; ok {
				r.BuyerDisplayName = u.DisplayName
			}

			if u, ok := mapUser[r.SellerAddress]; ok {
				r.SellerDisplayName = u.DisplayName
			}

			resp := map[string]interface{}{}
			_, err := client.R().
				SetResult(&resp).
				Get(fmt.Sprintf("%s/inscription/%s", uc.Config.GenerativeExplorerApi, r.InscriptionId))

			if err != nil {
				uc.Logger.Error(err)
			} else {
				if resp["content_type"] != nil {
					r.ContentType = resp["content_type"].(string)
				}

				if resp["number"] != nil {
					r.InscriptionIndex = resp["number"].(float64)
				}

			}
			mu.Lock()
			result = append(result, r)
			mu.Unlock()
		}(i)
	}
	waitG.Wait()

	resultMap := map[string]*entity.ItemListing{}
	for _, r := range result {
		resultMap[r.InscriptionId] = r
	}
	response := []*entity.ItemListing{}
	for _, id := range inscriptionIds {
		response = append(response, resultMap[id])
	}
	return response, nil
}

func (uc Usecase) ListItemListingOnSale(filter *structure.BaseFilters) ([]*entity.ItemListing, error) {
	data, err := uc.Repo.ListItemListingOnSale(filter)
	if err != nil {
		return nil, err
	}
	client := resty.New()
	waitG := &sync.WaitGroup{}
	for _, d := range data {
		waitG.Add(1)
		go func(d *entity.ItemListing) {
			defer waitG.Done()
			resp := map[string]interface{}{}
			_, err := client.R().
				SetResult(&resp).
				Get(fmt.Sprintf("%s/inscription/%s", uc.Config.GenerativeExplorerApi, d.InscriptionId))

			if err != nil {
				uc.Logger.Error(err)
			} else {
				d.ContentType = resp["content_type"].(string)
				d.InscriptionIndex = resp["number"].(float64)
			}
		}(d)
	}
	waitG.Wait()
	return data, nil
}

func (uc Usecase) ListItemListing(filter *structure.BaseFilters) ([]*entity.ItemListing, error) {
	data, err := uc.Repo.FindListItemListing(filter)
	if err != nil {
		return nil, err
	}
	client := resty.New()
	waitG := &sync.WaitGroup{}
	for _, d := range data {
		waitG.Add(1)
		go func(d *entity.ItemListing) {
			defer waitG.Done()
			resp := map[string]interface{}{}
			_, err := client.R().
				SetResult(&resp).
				Get(fmt.Sprintf("%s/inscription/%s", uc.Config.GenerativeExplorerApi, d.InscriptionId))

			if err != nil {
				uc.Logger.Error(err)
			} else {
				d.ContentType = resp["content_type"].(string)
				d.InscriptionIndex = resp["number"].(float64)
			}
		}(d)
	}
	waitG.Wait()
	return data, nil
}

func (u Usecase) ListToken(event *generative_marketplace_lib.GenerativeMarketplaceLibListingToken, blocknumber uint64) error {
	//TODO - DEBUG
	u.TokenActivites(blocknumber, event.Data.Price.Int64(), strings.ToLower(event.Data.Erc20Token.String()), event.Data.TokenId.String(), strings.ToLower(event.Data.Seller.String()), "", entity.TokenListing, "Listing")

	listing := entity.MarketplaceListings{
		OfferingId:         strings.ToLower(fmt.Sprintf("%x", event.OfferingId)),
		CollectionContract: strings.ToLower(event.Data.CollectionContract.String()),
		TokenId:            event.Data.TokenId.String(),
		Seller:             strings.ToLower(event.Data.Seller.String()),
		Erc20Token:         strings.ToLower(event.Data.Erc20Token.String()),
		Price:              event.Data.Price.String(),
		Closed:             event.Data.Closed,
		BlockNumber:        blocknumber,
		Finished:           false,
		DurationTime:       event.Data.DurationTime.String(),
	}

	sendMessage := func(listing entity.MarketplaceListings) {

		profile, err := u.Repo.FindUserByWalletAddress(listing.Seller)
		if err != nil {
			logger.AtLog.Logger.Error("ListToken", zap.Error(err))
			return
		}

		token, err := u.Repo.FindTokenByGenNftAddr(listing.CollectionContract, listing.TokenId)
		if err != nil {
			logger.AtLog.Logger.Error("ListToken", zap.Error(err))
			return
		}

		preText := fmt.Sprintf("[ListingID %s] has been created by %s", listing.OfferingId, listing.Seller)
		content := fmt.Sprintf("TokenID: %s", helpers.CreateTokenLink(token.ProjectID, token.TokenID, token.Name))
		title := fmt.Sprintf("User %s create listing with %s", helpers.CreateProfileLink(profile.WalletAddress, profile.DisplayName), listing.Price)

		if _, _, err := u.Slack.SendMessageToSlack(preText, title, content); err != nil {
			logger.AtLog.Logger.Error("s.Slack.SendMessageToSlack err", zap.Error(err))
		}
	}

	// check if listing is created or not
	_, err := u.Repo.FindListingByOfferingID(listing.OfferingId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// if error is no document -> create
			err := u.Repo.CreateMarketplaceListing(&listing)
			if err != nil {
				logger.AtLog.Logger.Error("ListToken", zap.Error(err))
				return err
			}

			sendMessage(listing)
			u.TokenActivites(blocknumber, event.Data.Price.Int64(), strings.ToLower(event.Data.Erc20Token.String()), event.Data.TokenId.String(), strings.ToLower(event.Data.Seller.String()), "", entity.TokenListing, "Listing")

			// TODO: @dac add update collection stats here

			return nil
		} else {
			logger.AtLog.Logger.Error("ListToken", zap.String("listing.OfferingId", listing.OfferingId))
			return err
		}
	} else {
		// listing is already created
		logger.AtLog.Logger.Error("ListToken", zap.String("listing.OfferingId", listing.OfferingId))
		return errors.New("listing is already created")
	}
}

func (u Usecase) PurchaseToken(event *generative_marketplace_lib.GenerativeMarketplaceLibPurchaseToken, blockNumber uint64) error {

	offeringID := strings.ToLower(fmt.Sprintf("%x", event.OfferingId))
	logger.AtLog.Logger.Info("PurchaseToken", zap.String("offeringID", offeringID))

	err := u.Repo.PurchaseTokenByOfferingID(offeringID)
	if err != nil {
		logger.AtLog.Logger.Error("PurchaseToken", zap.String("offeringID", offeringID), zap.Error(err))
		return err
	}
	u.TokenActivites(blockNumber, event.Data.Price.Int64(), strings.ToLower(event.Data.Erc20Token.String()), event.Data.TokenId.String(), strings.ToLower(event.Buyer.String()), "", entity.TokenPurchase, "Purchase token")

	getToken := func(offeringID string) (*entity.TokenUri, error) {
		listing, err := u.Repo.FindListingByOfferingID(offeringID)
		if err != nil {
			logger.AtLog.Logger.Error("PurchaseToken", zap.String("offeringID", offeringID), zap.Error(err))
			return nil, err
		}

		token, err := u.Repo.FindTokenByGenNftAddr(listing.CollectionContract, listing.TokenId)
		if err != nil {
			logger.AtLog.Logger.Error("PurchaseToken", zap.String("offeringID", offeringID), zap.Error(err))
			return nil, err
		}

		return token, nil
	}

	err = u.UpdateTokenOnwer("purchased", offeringID, getToken, event.Buyer)
	if err != nil {
		logger.AtLog.Logger.Error("PurchaseToken", zap.Error(err))
		return err
	}

	u.TokenActivites(blockNumber, event.Data.Price.Int64(), strings.ToLower(event.Data.Erc20Token.String()), event.Data.TokenId.String(), strings.ToLower(event.Data.Seller.String()), strings.ToLower(event.Buyer.String()), entity.TokenTransfer, "Transfer token")
	return nil
}

func (u Usecase) MakeOffer(event *generative_marketplace_lib.GenerativeMarketplaceLibMakeOffer, blocknumber uint64) error {

	offer := entity.MarketplaceOffers{
		OfferingId:         strings.ToLower(fmt.Sprintf("%x", event.OfferingId)),
		CollectionContract: strings.ToLower(event.Data.CollectionContract.String()),
		TokenId:            event.Data.TokenId.String(),
		Buyer:              strings.ToLower(event.Data.Buyer.String()),
		Erc20Token:         strings.ToLower(event.Data.Erc20Token.String()),
		Price:              event.Data.Price.String(),
		Closed:             event.Data.Closed,
		Finished:           false,
		DurationTime:       event.Data.DurationTime.String(),
		BlockNumber:        blocknumber,
	}

	sendMessage := func(offer entity.MarketplaceOffers) {

		profile, err := u.Repo.FindUserByWalletAddress(offer.Buyer)
		if err != nil {
			logger.AtLog.Logger.Error("cancelListing.FindUserByWalletAddress", zap.Error(err))
			return
		}

		token, err := u.Repo.FindTokenByGenNftAddr(offer.CollectionContract, offer.TokenId)
		if err != nil {
			logger.AtLog.Logger.Error("cancelListing.FindTokenByGenNftAddr", zap.Error(err))
			return
		}

		preText := fmt.Sprintf("[OfferID %s] has been created by %s", offer.OfferingId, offer.Buyer)
		content := fmt.Sprintf("TokenID: %s", helpers.CreateTokenLink(token.ProjectID, token.TokenID, token.Name))
		title := fmt.Sprintf("User %s made offer with %s", helpers.CreateProfileLink(profile.WalletAddress, profile.DisplayName), offer.Price)

		if _, _, err := u.Slack.SendMessageToSlack(preText, title, content); err != nil {
			logger.AtLog.Logger.Error("s.Slack.SendMessageToSlack err", zap.Error(err))
		}
	}

	// check if listing is created or not
	_, err := u.Repo.FindOfferByOfferingID(offer.OfferingId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// if error is no document -> create
			err := u.Repo.CreateMarketplaceOffer(&offer)
			if err != nil {
				logger.AtLog.Logger.Error("makeOffer.Repo.CreateMarketplaceOffer", zap.Error(err))
				return err
			}

			sendMessage(offer)

			u.TokenActivites(blocknumber, event.Data.Price.Int64(), strings.ToLower(event.Data.Erc20Token.String()), event.Data.TokenId.String(), strings.ToLower(event.Data.Buyer.String()), "", entity.TokenMakeOffer, "Make offer")
			// TODO: @dac add update collection stats here
			return nil
		} else {
			logger.AtLog.Logger.Error("makeOffer.Repo.FindOfferByOfferingID", zap.Error(err))
			return err
		}
	} else {
		err := errors.New("offer is already created")
		// listing is already created
		logger.AtLog.Logger.Error("offer token offeringId", zap.Error(err))
		return err
	}

	return nil
}

func (u Usecase) AcceptMakeOffer(event *generative_marketplace_lib.GenerativeMarketplaceLibAcceptMakeOffer, blockNumber uint64) error {

	offeringID := strings.ToLower(fmt.Sprintf("%x", event.OfferingId))
	logger.AtLog.Logger.Info("accept make offer offeringId", zap.Any("offeringID", offeringID))

	err := u.Repo.AcceptOfferByOfferingID(offeringID)
	if err != nil {
		logger.AtLog.Logger.Error("u.AcceptMakeOffer.AcceptOfferByOfferingID", zap.Error(err))
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
	err = u.UpdateTokenOnwer("accepted", offeringID, getToken, event.Buyer)
	if err != nil {
		logger.AtLog.Logger.Error("u.UpdateTokenOnwer.UpdateTokenOnwer", zap.Error(err))
		return err
	}

	// TODO: @dac add update collection stats here
	u.TokenActivites(blockNumber, event.Data.Price.Int64(), strings.ToLower(event.Data.Erc20Token.String()), event.Data.TokenId.String(), "", strings.ToLower(event.Buyer.String()), entity.TokenAcceptOffer, "Accept offer")

	return nil
}

func (u Usecase) CancelListing(event *generative_marketplace_lib.GenerativeMarketplaceLibCancelListing) error {

	offeringID := strings.ToLower(fmt.Sprintf("%x", event.OfferingId))
	logger.AtLog.Logger.Info("CancelListing", zap.String("offeringID", offeringID))

	err := u.Repo.CancelListingByOfferingID(offeringID)
	if err != nil {
		logger.AtLog.Logger.Error("CancelListing", zap.String("offeringID", offeringID), zap.Error(err))
		return err
	}

	done := make(chan bool)
	go func(done chan bool) {
		defer func() {
			done <- true
		}()

		listing, err := u.Repo.FindListingByOfferingID(offeringID)
		if err != nil {
			logger.AtLog.Logger.Error("CancelListing", zap.String("offeringID", offeringID), zap.Error(err))
			return
		}

		profile, err := u.Repo.FindUserByWalletAddress(listing.Seller)
		if err != nil {
			logger.AtLog.Logger.Error("CancelListing", zap.String("offeringID", offeringID), zap.Error(err))
			return
		}

		token, err := u.Repo.FindTokenByGenNftAddr(listing.CollectionContract, listing.TokenId)
		if err != nil {
			logger.AtLog.Logger.Error("CancelListing", zap.String("offeringID", offeringID), zap.Error(err))
			return
		}

		preText := fmt.Sprintf("[Listing %s] has been cancelled by %s", listing.OfferingId, listing.Seller)
		content := fmt.Sprintf("TokenID: %s", helpers.CreateTokenLink(token.ProjectID, token.TokenID, token.Name))
		title := fmt.Sprintf("User %s cancelled offer %s", helpers.CreateProfileLink(profile.WalletAddress, profile.DisplayName), offeringID)

		if _, _, err := u.Slack.SendMessageToSlack(preText, title, content); err != nil {
			logger.AtLog.Logger.Error("CancelListing", zap.String("offeringID", offeringID), zap.Error(err))
		}
	}(done)
	<-done

	// TODO: @dac add update collection stats here
	return nil
}

func (u Usecase) CancelOffer(event *generative_marketplace_lib.GenerativeMarketplaceLibCancelMakeOffer, blockNumber uint64) error {

	offeringID := strings.ToLower(fmt.Sprintf("%x", event.OfferingId))
	logger.AtLog.Logger.Info("cancel make offer offeringId", zap.Any("offeringID", offeringID))

	err := u.Repo.CancelOfferByOfferingID(offeringID)
	if err != nil {
		logger.AtLog.Logger.Error("s.Repo.CancelOfferByOfferingID", zap.Error(err))
		return err
	}

	done := make(chan bool)
	go func(done chan bool) {
		defer func() {
			done <- true
		}()

		offer, err := u.Repo.FindOfferByOfferingID(offeringID)
		if err != nil {
			logger.AtLog.Logger.Error("s.Repo.FindOfferByOfferingID", zap.Error(err))
			return
		}

		profile, err := u.Repo.FindUserByWalletAddress(offer.Buyer)
		if err != nil {
			logger.AtLog.Logger.Error("cancelListing.FindUserByWalletAddress", zap.Error(err))
			return
		}

		token, err := u.Repo.FindTokenByGenNftAddr(offer.CollectionContract, offer.TokenId)
		if err != nil {
			logger.AtLog.Logger.Error("cancelListing.FindTokenByGenNftAddr", zap.Error(err))
			return
		}

		preText := fmt.Sprintf("[Listing %s] has been cancelled by %s", offer.OfferingId, offer.Buyer)
		content := fmt.Sprintf("TokenID: %s", helpers.CreateTokenLink(token.ProjectID, token.TokenID, token.Name))
		title := fmt.Sprintf("User %s cancelled offer %s", helpers.CreateProfileLink(profile.WalletAddress, profile.DisplayName), offeringID)

		if _, _, err := u.Slack.SendMessageToSlack(preText, title, content); err != nil {
			logger.AtLog.Logger.Error("s.Slack.SendMessageToSlack err", zap.Error(err))
		}

		u.TokenActivites(blockNumber, event.Data.Price.Int64(), strings.ToLower(event.Data.Erc20Token.String()), event.Data.TokenId.String(), strings.ToLower(event.Data.Buyer.String()), "", entity.TokenCancelOffer, "Cancel offer")

	}(done)
	<-done

	// TODO: @dac add update collection stats here
	return nil
}

func (u Usecase) FilterMKListing(filter structure.FilterMkListing) (*entity.Pagination, error) {

	fm := &entity.FilterMarketplaceListings{}
	err := copier.Copy(fm, filter)
	if err != nil {
		logger.AtLog.Logger.Error("copier.Copy", zap.Error(err))
		return nil, err
	}
	ml, err := u.Repo.FilterMarketplaceListings(*fm)
	if err != nil {
		logger.AtLog.Logger.Error("u.Repo.FilterMarketplaceListings", zap.Error(err))
		return nil, err
	}

	logger.AtLog.Logger.Info("filtered", zap.Any("ml", ml))

	listingResp := []entity.MarketplaceListings{}
	iListings := ml.Result
	listing := iListings.([]entity.MarketplaceListings)
	for _, listingItem := range listing {
		tok, err := u.Repo.FindTokenByGenNftAddr(listingItem.CollectionContract, listingItem.TokenId)
		if err != nil {
			logger.AtLog.Logger.Error("u.Repo.FindTokenByGenNftAddr", zap.Error(err))
		} else {
			listingItem.Token = *tok
		}

		p, err := u.Repo.FindUserByWalletAddress(listingItem.Seller)
		if err == nil {
			listingItem.SellerInfo = *p
		}

		listingResp = append(listingResp, listingItem)
	}
	ml.Result = listingResp
	return ml, nil
}

func (u Usecase) FilterMKOffers(filter structure.FilterMkOffers) (*entity.Pagination, error) {

	fm := &entity.FilterMarketplaceOffers{}
	err := copier.Copy(fm, filter)
	if err != nil {
		logger.AtLog.Logger.Error("copier.Copy", zap.Error(err))
		return nil, err
	}
	ml, err := u.Repo.FilterMarketplaceOffers(*fm)
	if err != nil {
		logger.AtLog.Logger.Error("u.Repo.FilterMarketplaceOffers", zap.Error(err))
		return nil, err
	}

	logger.AtLog.Logger.Info("filtered", zap.Any("ml", ml))

	offersResp := []entity.MarketplaceOffers{}
	iOffers := ml.Result
	offers := iOffers.([]entity.MarketplaceOffers)
	for _, offer := range offers {
		tok, err := u.Repo.FindTokenByGenNftAddr(offer.CollectionContract, offer.TokenId)
		if err != nil {
			logger.AtLog.Logger.Error("u.Repo.FindTokenByGenNftAddr", zap.Error(err))
			//continue
		} else {
			offer.Token = *tok
		}

		p, err := u.Repo.FindUserByWalletAddress(offer.Buyer)
		if err == nil {
			offer.BuyerInfo = *p
		}
		offersResp = append(offersResp, offer)
	}

	ml.Result = offersResp
	return ml, nil
}

func (u Usecase) GetListingBySeller(sellerAddress string) ([]entity.MarketplaceListings, []string, []string, error) {

	cachedKey, cachedContractIDsKey, cachedTokensIDsKey := helpers.ProfileSelingKey(sellerAddress)
	listings := []entity.MarketplaceListings{}
	contractIDS := []string{}
	tokenIDS := []string{}
	var err error

	logger.AtLog.Logger.Info("cachedKey", zap.Any("cachedKey", cachedKey))
	logger.AtLog.Logger.Info("cachedContractIDsKey", zap.Any("cachedContractIDsKey", cachedContractIDsKey))
	logger.AtLog.Logger.Info("cachedTokensIDsKey", zap.Any("cachedTokensIDsKey", cachedTokensIDsKey))

	//always reloa data
	liveReload := func(sellerAddress string, cachedKey string, cachedContractIDsKey string, cachedTokensIDsKey string) ([]entity.MarketplaceListings, []string, []string, error) {

		listings, err = u.Repo.GetListingBySeller(sellerAddress)
		if err != nil {
			logger.AtLog.Logger.Error("u.Repo.GetListingBySeller", zap.Error(err))
			return nil, nil, nil, err
		}

		contractIDS := []string{}
		tokenIDS := []string{}
		for key, listing := range listings {
			logger.AtLog.Logger.Info(fmt.Sprintf("listing.%d", key), zap.Any("listing", listing))
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
		if err != nil {
			logger.AtLog.Logger.Error("helpers.ParseCache.listings", zap.Error(err))
			return nil, nil, nil, err
		}
		cached, err := u.Cache.GetData(cachedContractIDsKey)
		err = helpers.ParseCache(cached, &contractIDS)
		if err != nil {
			logger.AtLog.Logger.Error("helpers.ParseCache.cachedContractIDsKey", zap.Error(err))
			return nil, nil, nil, err
		}
		cached, err = u.Cache.GetData(cachedTokensIDsKey)
		err = helpers.ParseCache(cached, &tokenIDS)
		if err != nil {
			logger.AtLog.Logger.Error("helpers.ParseCache.tokenIDS", zap.Error(err))
			return nil, nil, nil, err
		}

	} else {
		listings, contractIDS, tokenIDS, err = liveReload(sellerAddress, cachedKey, cachedContractIDsKey, cachedTokensIDsKey)
		if err != nil {
			logger.AtLog.Logger.Error("liveReload", zap.Error(err))
			return nil, nil, nil, err
		}
	}

	logger.AtLog.Logger.Info("listings", zap.Any("listings", listings))
	logger.AtLog.Logger.Info("contractIDS", zap.Any("contractIDS", contractIDS))
	logger.AtLog.Logger.Info("tokenIDS", zap.Any("tokenIDS", tokenIDS))
	return listings, contractIDS, tokenIDS, nil
}

func (u Usecase) UpdateTokenOnwer(event string, offeringID string, fn func(offeringID string) (*entity.TokenUri, error), buyer common.Address) error {

	owner := strings.ToLower(buyer.String())
	token, err := fn(offeringID)
	if err != nil {
		logger.AtLog.Logger.Error("UpdateTokenOnwer.fn", zap.Error(err))
		return err
	}

	logger.AtLog.Logger.Info("tokenID", zap.Any("token.TokenID", token.TokenID))

	profile, err := u.Repo.FindUserByWalletAddress(owner)
	if err != nil {
		// if can not find user profile in db, set owner to nil
		if errors.Is(err, mongo.ErrNoDocuments) {
			profile = nil
		} else {
			logger.AtLog.Logger.Error("UpdateTokenOnwer.FindUserByWalletAddress", zap.Error(err))
			return err
		}
	}

	logger.AtLog.Logger.Info("token.Owner", zap.Any("owner", owner))

	token.Owner = profile
	token.OwnerAddr = owner

	updated, err := u.Repo.UpdateOrInsertTokenUri(token.ContractAddress, token.TokenID, token)
	if err != nil {
		logger.AtLog.Logger.Error("UpdateTokenOnwer.UpdateOrInsertTokenUri", zap.Error(err))
		return err
	}

	logger.AtLog.Logger.Info("updated", zap.Any("updated", updated))

	//slack
	preText := fmt.Sprintf("[TokenID %s] has been transfered to %s", token.TokenID, token.OwnerAddr)
	content := fmt.Sprintf("To user: %s. Token: %s", helpers.CreateProfileLink(owner, profile.DisplayName), helpers.CreateTokenLink(token.ProjectID, token.TokenID, token.Name))
	title := fmt.Sprintf("OfferingID:  %s is %s", offeringID, event)

	if _, _, err := u.Slack.SendMessageToSlack(preText, title, content); err != nil {
		logger.AtLog.Logger.Error("s.Slack.SendMessageToSlack err", zap.Error(err))
	}

	// TODO: @dac add update collection stats here
	return nil
}

func (u Usecase) TokenActivites(blocknumber uint64, amount int64, erc20Address string, tokenID string, fromWallet string, toWallet string, action entity.TokenActivityType, title string) {
	bn := big.NewInt(int64(blocknumber))
	blockInfo, err := u.TcClientPublicNode.BlockByNumber(context.Background(), bn)

	blockTime := blockInfo.Header().Time
	tm := time.Unix(int64(blockTime), 0).UTC()

	tok, err := u.Repo.FindTokenByTokenID(tokenID)
	if err == nil {
		//token activities here
		err = u.Repo.InsertTokenActivity(&entity.TokenActivity{
			Type:          action,
			Title:         title,
			UserAAddress:  fromWallet,
			UserBAddress:  toWallet,
			Amount:        amount,
			Erc20Address:  erc20Address,
			InscriptionID: tokenID,
			ProjectID:     tok.ProjectID,
			Time:          &tm,
			BlockNumber:   blocknumber,
		})
	} else {
		logger.AtLog.Logger.Error("TokenActivites", zap.String("FindTokenByTokenID", tokenID), zap.Error(err))
	}
}
