package usecase

import (
	"errors"
	"fmt"
	"strings"

	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/mongo"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/contracts/generative_marketplace_lib"
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
	err := u.Repo.PurchaseTokenByOfferingID(offeringID)
	if err != nil {
		return err
	}

	// TODO: @dac add update collection stats here

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

	// check if listing is created or not
	_, err := u.Repo.FindOfferByOfferingID(offer.OfferingId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// if error is no document -> create
			err := u.Repo.CreateMarketplaceOffer(&offer)
	
			if err != nil {
				log.Error("error when create marketplace offer", "", err)
				return err
			}

			// TODO: @dac add update collection stats here

			return nil
		} else {
			return err
		}
	} else {
		// listing is already created
		log.SetData("offer token offeringId", offer.OfferingId)
		return errors.New("offer is already created")
	}
}

func (u Usecase) AcceptMakeOffer(rootSpan opentracing.Span, event *generative_marketplace_lib.GenerativeMarketplaceLibAcceptMakeOffer) error {
	span, log := u.StartSpan("AcceptMakeOffer", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	offeringID := strings.ToLower(fmt.Sprintf("%x", event.OfferingId))
	log.SetData("accept make offer offeringId", offeringID)
	err := u.Repo.AcceptOfferByOfferingID(offeringID)
	if err != nil {
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
	err := u.Repo.CancelListingByOfferingID(offeringID)
	if err != nil {
		return err
	}

	// TODO: @dac add update collection stats here

	return nil
}

func (u Usecase) CancelOffer(rootSpan opentracing.Span, event *generative_marketplace_lib.GenerativeMarketplaceLibCancelMakeOffer) error {
	span, log := u.StartSpan("CancelMakeOffer", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	offeringID := strings.ToLower(fmt.Sprintf("%x", event.OfferingId))
	log.SetData("cancel make offer offeringId", offeringID)
	err := u.Repo.CancelOfferByOfferingID(offeringID)
	if err != nil {
		return err
	}

	// TODO: @dac add update collection stats here

	return nil
}

