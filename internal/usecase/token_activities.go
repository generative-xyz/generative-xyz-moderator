package usecase

import (
	"sort"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
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

func (u Usecase) JobCreateTokenActivityFromListings() error {
	for page := int64(1);; page++ {
		u.Logger.Info("StartGetPagingNotCreatedActivitiesListings", zap.Any("page", page))
		uListings, err := u.Repo.GetNotCreatedActivitiesListing(page, 100)
		if err != nil {
			return errors.WithStack(err)
		}
		listings := uListings.Result.([]entity.DexBTCListing)
		if len(listings) == 0 {
			break
		}
		u.Logger.Info("StartGetPagingNotCreatedActivitiesListings", zap.Any("page", page))
		
		for _, listing := range listings {
			token, err := u.Repo.FindTokenByTokenID(listing.InscriptionID)
			if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
				return errors.WithStack(err)
			}

			var projectID string
			if token != nil {
				projectID = token.ProjectID
			}

			if listing.Verified && !listing.CreatedVerifiedActivity {
				activity := entity.TokenActivity{
					Type: entity.TokenListing,
					Title: "Listing",
					UserAAddress: listing.SellerAddress,
					Amount: int64(listing.Amount),
					Time: listing.CreatedAt,			
					InscriptionID: listing.InscriptionID,
					ProjectID: projectID,
				}
				err := u.Repo.InsertTokenActivity(&activity)
				if err != nil {
					u.Logger.ErrorAny("JobCreateTokenActivityFromListings.InsertTokenActivity", zap.Error(err))
				} else {
					u.Logger.Info("JobCreateTokenActivityFromListings.InsertTokenActivity", zap.Any("activity", activity))
					u.Repo.UpdateListingCreatedVerifiedActivity(listing.UUID)
				}
			}
			if listing.Cancelled && !listing.CreatedCancelledActivity {
				activity := entity.TokenActivity{
					Type: entity.TokenCancelListing,
					Title: "Cancel Listing",
					UserAAddress: listing.SellerAddress,
					Amount: int64(listing.Amount),
					Time: listing.CancelAt,				
					InscriptionID: listing.InscriptionID,
					ProjectID: projectID,
				}
				err := u.Repo.InsertTokenActivity(&activity)
				if err != nil {
					u.Logger.ErrorAny("JobCreateTokenActivityFromListings.InsertTokenActivity", zap.Error(err))
				} else {
					u.Logger.Info("JobCreateTokenActivityFromListings.InsertTokenActivity", zap.Any("activity", activity))
					u.Repo.UpdateListingCreatedCancelledActivity(listing.UUID)
				}
			}
			if listing.Matched && !listing.CreatedMatchedActivity {
				activity := entity.TokenActivity{
					Type: entity.TokenMatched,
					Title: "Saled",
					UserAAddress: listing.SellerAddress,
					UserBAddress: listing.Buyer,
					Amount: int64(listing.Amount),
					Time: listing.MatchAt,
					InscriptionID: listing.InscriptionID,
					ProjectID: projectID,
				}
				err := u.Repo.InsertTokenActivity(&activity)
				if err != nil {
					u.Logger.ErrorAny("JobCreateTokenActivityFromListings.InsertTokenActivity", zap.Error(err))
				} else {
					u.Logger.Info("JobCreateTokenActivityFromListings.InsertTokenActivity", zap.Any("activity", activity))
					u.Repo.UpdateListingCreatedMatchedActivity(listing.UUID)
				}
			}
		}
	}
	return nil
}

func (u Usecase) JobCreateTokenMintActivityFromTokenUri() error {
	for page := int64(1);; page++ {
		u.Logger.Info("StartGetPagingNotCreatedActivitiesToken", zap.Any("page", page))
		uTokens, err := u.Repo.GetNotCreatedActivitiesToken(page, 100)
		if err != nil {
			return errors.WithStack(err)
		}
		tokens := uTokens.Result.([]entity.TokenUri)
		if len(tokens) == 0 {
			break
		}
		u.Logger.Info("StartGetPagingNotCreatedActivitiesTokens", zap.Any("page", page))
		
		for _, token := range tokens {
			var minterAddress string
			if token.MinterAddress != nil {
				minterAddress = *token.MinterAddress
			}
			activity := entity.TokenActivity{
				Type: entity.TokenMint,
				Title: "Minted",
				UserAAddress: minterAddress,
				Time: token.CreatedAt,
				InscriptionID: token.TokenID,
				ProjectID: token.ProjectID,
			}
			err := u.Repo.InsertTokenActivity(&activity)
			if err != nil {
				u.Logger.ErrorAny("JobCreateTokenMintActivityFromTokenUri.InsertTokenActivity", zap.Error(err))
			} else {
				u.Logger.Info("JobCreateTokenMintActivityFromTokenUri.InsertTokenActivity", zap.Any("activity", activity))
				u.Repo.UpdateTokenCreatedMintActivity(token.TokenID)
			}
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}
