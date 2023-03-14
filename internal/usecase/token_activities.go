package usecase

import (
	"sort"

	"github.com/pkg/errors"
	"rederinghub.io/internal/entity"
)

func (u Usecase) GetTokenActivities(page int64, limit int64, inscriptionID string) (*entity.Pagination, error) {
	token, err := u.Repo.FindTokenByTokenID(inscriptionID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var minterAddress string
	if token.MinterAddress != nil {
		minterAddress = *token.MinterAddress
	}
	allDexBtcListings, err := u.Repo.GetAllDexBTCListingByInscriptionID(inscriptionID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	
	activities := []entity.TokenActivity{
		{
			Type: entity.TokenMint,
			Title: "Minted",
			UserAAddress: "",
			UserBAddress: minterAddress,
			Time: token.CreatedAt,
			Amount: 0,
		},
	}

	for _, listing := range allDexBtcListings {
		if listing.Verified {
			activities = append(activities, entity.TokenActivity{
				Type: entity.TokenListing,
				Title: "Listing",
				UserAAddress: listing.SellerAddress,
				Amount: int64(listing.Amount),
				Time: listing.CreatedAt,			
			})
		}
		if listing.Cancelled {
			activities = append(activities, entity.TokenActivity{
				Type: entity.TokenCancelListing,
				Title: "Cancel Listing",
				UserAAddress: listing.SellerAddress,
				Amount: int64(listing.Amount),
				Time: listing.CancelAt,				
			})
		}
		if listing.Matched {
			activities = append(activities, entity.TokenActivity{
				Type: entity.TokenMatched,
				Title: "Saled",
				UserAAddress: listing.SellerAddress,
				UserBAddress: listing.Buyer,
				Amount: int64(listing.Amount),
				Time: listing.MatchAt,	
			})
		}
	}

	sort.SliceStable(activities, func (i, j int) bool {
		if activities[j].Time == nil {
			return false
		}
		if activities[i].Time == nil {
			return true
		}
		return activities[i].Time.Before(*activities[j].Time)
	})

	// paginating on slice
	skip := (page - 1) * limit
	if skip > int64(len(activities)) {
		skip = int64(len(activities))
	}
	end := skip + limit
	if end > int64(len(activities)) {
		end = int64(len(activities))
	}
	
	pagedActivities := activities[skip:end]
	addresses := []string{}
	for _, activity := range pagedActivities {
		if activity.UserAAddress != "" {
			addresses = append(addresses, activity.UserAAddress)
		}
		if activity.UserBAddress != "" {
			addresses = append(addresses, activity.UserBAddress)
		}
	}

	userMap, err := u.GetUsersMap(addresses)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, activity := range pagedActivities {
		if activity.UserAAddress != "" {
			activity.UserA = userMap[activity.UserAAddress]
		}
		if activity.UserBAddress != "" {
			activity.UserB = userMap[activity.UserBAddress]
		}
	}
	
	return &entity.Pagination{
		Result: pagedActivities,
		Page: page,
		PageSize: limit, 
		Total: int64(len(activities)),
		TotalPage: int64(float64(len(activities)) / float64(limit)),
	}, nil
}
