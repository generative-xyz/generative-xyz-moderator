package usecase

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/bson"
)

// synchronize token data
func (u Usecase) SyncTokenAndMarketplaceData(rootSpan opentracing.Span) error {
	span, log := u.StartSpan("SyncTokenAndMarketplaceData", rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	// update token price by marketplace data
	listings, err := u.Repo.GetAllActiveListings() 
	if err != nil {
		return err
	}
	offers, err := u.Repo.GetAllActiveOffers() 
	if err != nil {
		return err
	}
	// sort
	sort.SliceStable(listings, func(i, j int) bool {
		if listings[i].BlockNumber != listings[j].BlockNumber {
			return listings[i].CreatedAt.After(*listings[j].CreatedAt)
		}
		return listings[i].BlockNumber > listings[j].BlockNumber
	});

	curTime := time.Now().Unix()
	// update listing/offer that closed
	for id, listing := range listings {
		if listing.DurationTime != "0" {
			durationTime, err := strconv.Atoi(listing.DurationTime)
			if err != nil {
				return nil
			}
			// listing is passed deadline
			if int64(durationTime) > curTime {
				u.Repo.CancelListingByOfferingID(listing.OfferingId)
				listings[id].Closed = true
			}
		}
	}
	for id, offer := range offers {
		if offer.DurationTime != "0" {
			durationTime, err := strconv.Atoi(offer.DurationTime)
			if err != nil {
				return nil
			}
			// listing is passed deadline
			if int64(durationTime) > curTime {
				u.Repo.CancelOfferByOfferingID(offer.OfferingId)
				offers[id].Closed = true
			}
		}
	}
	
	// map from token id to price
	fromTokenIdToPrice := make(map[string]string)
	for _, listing := range listings {
		if _, ok := fromTokenIdToPrice[listing.TokenId]; !ok && !listing.Closed {
			fromTokenIdToPrice[listing.TokenId] = listing.Price
		}
	}

	tokenWithPrices, err := u.Repo.FindTokenUrisWithoutCache(bson.M{"stats.price": bson.M{"$exists": true}})
	// set of tokens that currently has price
	tokenWithPricesSet := make(map[string]bool)
	for _, token := range tokenWithPrices {
		tokenWithPricesSet[token.TokenID] = true
	}

	if err != nil {
		return nil
	}
	for k, v := range fromTokenIdToPrice {
		token, err := u.Repo.FindTokenUriWithoutCache(bson.D{{Key: "token_id", Value: k}})
		if err != nil {
			return nil
		}
		if token.Stats.Price != v {
			log.SetData(fmt.Sprintf("setTokenPrice%s", k), v)
			u.Repo.UpdateTokenPriceByTokenId(k, v)
		}
		tokenWithPricesSet[k] = false
	}
	for k, v := range tokenWithPricesSet {
		if !v {
			continue
		}
		log.SetData(fmt.Sprintf("unsetTokenPrice%s", k), true)
		u.Repo.UnsetTokenPriceByTokenId(k)
	}
	return nil
}
