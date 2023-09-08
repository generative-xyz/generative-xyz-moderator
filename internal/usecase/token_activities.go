package usecase

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/logger"
)

func (u Usecase) GetTokenActivities(filter structure.FilterTokenActivities) (*entity.Pagination, error) {
	pe := &entity.FilterTokenActivities{}
	err := copier.Copy(pe, filter)
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		return nil, errors.WithStack(err)
	}
	activitiesResp, err := u.Repo.GetTokenActivities(*pe)
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		return nil, errors.WithStack(err)
	}
	activities := activitiesResp.Result.([]entity.TokenActivity)

	userAddresses := []string{}
	tokenIds := []string{}

	for _, activity := range activities {
		if activity.UserAAddress != "" {
			userAddresses = append(userAddresses, activity.UserAAddress)
		}
		if activity.UserBAddress != "" {
			userAddresses = append(userAddresses, activity.UserBAddress)
		}
		tokenIds = append(tokenIds, activity.InscriptionID)
	}

	userMap, err := u.GetUsersMap(userAddresses)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	tokenMap, err := u.GetTokensMap(tokenIds)

	for id := range activities {
		if activities[id].UserAAddress != "" {
			activities[id].UserA = userMap[activities[id].UserAAddress]
		}
		if activities[id].UserBAddress != "" {
			activities[id].UserB = userMap[activities[id].UserBAddress]
		}

		tokenInfo := tokenMap[activities[id].InscriptionID]
		tokenInfo.PriceBrc20 = entity.PriceBRC20Obj{
			Value:   fmt.Sprintf("%d", activities[id].Amount),
			Address: activities[id].Erc20Address,
			//OfferingID: activities[id],
		}

		activities[id].TokenInfo = &tokenInfo
	}

	activitiesResp.Result = activities

	logger.AtLog.Logger.Info("activities", zap.Any("activitiesResp.Total", activitiesResp.Total))
	return activitiesResp, nil
}

func (u Usecase) JobCreateTokenActivityFromListings() error {
	for page := int64(1); ; page++ {
		logger.AtLog.Logger.Info("StartGetPagingNotCreatedActivitiesListings", zap.Any("page", zap.Any("page)", page)))
		uListings, err := u.Repo.GetNotCreatedActivitiesListing(page, 100)
		if err != nil {
			return errors.WithStack(err)
		}
		listings := uListings.Result.([]entity.DexBTCListing)
		if len(listings) == 0 {
			break
		}
		logger.AtLog.Logger.Info("StartGetPagingNotCreatedActivitiesListings", zap.Any("page", zap.Any("page)", page)))

		for _, listing := range listings {
			token, err := u.Repo.FindTokenByTokenID(listing.InscriptionID)
			if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
				return errors.WithStack(err)
			}

			var projectID string
			if token != nil {
				projectID = token.ProjectID
			}

			if listing.Verified && !listing.CreatedVerifiedActivity && !listing.FromOtherMkp {
				activity := entity.TokenActivity{
					Type:          entity.TokenListing,
					Title:         "Listing",
					UserAAddress:  listing.SellerAddress,
					Amount:        int64(listing.Amount),
					Time:          listing.CreatedAt,
					InscriptionID: listing.InscriptionID,
					ProjectID:     projectID,
				}
				err := u.Repo.InsertTokenActivity(&activity)
				if err != nil {
					logger.AtLog.Logger.Error("JobCreateTokenActivityFromListings.InsertTokenActivity", zap.Error(err))
				} else {
					logger.AtLog.Logger.Info("JobCreateTokenActivityFromListings.InsertTokenActivity", zap.Any("activity", zap.Any("activity)", activity)))
					u.Repo.UpdateListingCreatedVerifiedActivity(listing.UUID)
				}
			}
			if listing.Cancelled && !listing.CreatedCancelledActivity {
				activity := entity.TokenActivity{
					Type:          entity.TokenCancelListing,
					Title:         "Cancel Listing",
					UserAAddress:  listing.SellerAddress,
					Amount:        int64(listing.Amount),
					Time:          listing.CancelAt,
					InscriptionID: listing.InscriptionID,
					ProjectID:     projectID,
				}
				err := u.Repo.InsertTokenActivity(&activity)
				if err != nil {
					logger.AtLog.Logger.Error("JobCreateTokenActivityFromListings.InsertTokenActivity", zap.Error(err))
				} else {
					logger.AtLog.Logger.Info("JobCreateTokenActivityFromListings.InsertTokenActivity", zap.Any("activity", zap.Any("activity)", activity)))
					u.Repo.UpdateListingCreatedCancelledActivity(listing.UUID)
				}
			}
			if listing.Matched && !listing.CreatedMatchedActivity {
				activity := entity.TokenActivity{
					Type:          entity.TokenMatched,
					Title:         "Saled",
					UserAAddress:  listing.SellerAddress,
					UserBAddress:  listing.Buyer,
					Amount:        int64(listing.Amount),
					Time:          listing.MatchAt,
					InscriptionID: listing.InscriptionID,
					ProjectID:     projectID,
				}
				err := u.Repo.InsertTokenActivity(&activity)
				if err != nil {
					logger.AtLog.Logger.Error("JobCreateTokenActivityFromListings.InsertTokenActivity", zap.Error(err))
				} else {
					logger.AtLog.Logger.Info("JobCreateTokenActivityFromListings.InsertTokenActivity", zap.Any("activity", zap.Any("activity)", activity)))
					u.Repo.UpdateListingCreatedMatchedActivity(listing.UUID)
				}
			}
		}
	}
	return nil
}

func (u Usecase) JobCreateTokenMintActivityFromTokenUri() error {
	for page := int64(1); ; page++ {
		logger.AtLog.Logger.Info("StartGetPagingNotCreatedActivitiesToken", zap.Any("page", zap.Any("page)", page)))
		uTokens, err := u.Repo.GetNotCreatedActivitiesToken(page, 100)
		if err != nil {
			return errors.WithStack(err)
		}
		tokens := uTokens.Result.([]entity.TokenUri)
		if len(tokens) == 0 {
			break
		}
		logger.AtLog.Logger.Info("StartGetPagingNotCreatedActivitiesTokens", zap.Any("page", zap.Any("page)", page)))

		for _, token := range tokens {
			var minterAddress string
			var amount int64
			mintNftBtc, err := u.Repo.FindMintNftBtcByInscriptionID(token.TokenID)
			if err != nil {
				if !errors.Is(err, mongo.ErrNoDocuments) {
					return errors.WithStack(err)
				}
			}
			if mintNftBtc != nil {
				minterAddress = mintNftBtc.OriginUserAddress
				amount, err = strconv.ParseInt(mintNftBtc.EstFeeInfo["btc"].MintPrice, 10, 64)
				if err != nil {
					logger.AtLog.Logger.Error("JobCreateTokenMintActivityFromTokenUri.FailedParseMintPrice")
				}
			}

			activity := entity.TokenActivity{
				Type:          entity.TokenMint,
				Title:         "Minted",
				UserAAddress:  minterAddress,
				Time:          token.CreatedAt,
				InscriptionID: token.TokenID,
				ProjectID:     token.ProjectID,
				Amount:        amount,
			}
			err = u.Repo.InsertTokenActivity(&activity)
			if err != nil {
				logger.AtLog.Logger.Error("JobCreateTokenMintActivityFromTokenUri.InsertTokenActivity", zap.Error(err))
			} else {
				logger.AtLog.Logger.Info("JobCreateTokenMintActivityFromTokenUri.InsertTokenActivity", zap.Any("activity", zap.Any("activity)", activity)))
				u.Repo.UpdateTokenCreatedMintActivity(token.TokenID)
			}
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}
