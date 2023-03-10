package usecase

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
)

const (
	INF_TRENDING_SCORE                  int64 = 9223372036854775807 // max int64 value
	SATOSHI_EACH_BTC                    int64 = 100000000
	TRENDING_SCORE_EACH_BTC_VOLUMN      int64 = 1000000
	TRENDING_SCORE_EACH_OPENING_LISTING int64 = 30
	TRENDING_SCORE_EACH_MINT            int64 = 10
	TRENDING_SCORE_EACH_VIEW            int64 = 1
)

func (u Usecase) JobSyncProjectTrending() error {
	u.Logger.Info("JobSyncProjectTrending.StartJobSyncProjectTrending")
	// All btc activities, which include Mint and Buy activity
	// Mapping from projectID to latest two week's volumn in satoshi
	fromProjectIDToRecentVolumn := map[string]int64{}
	for page := int64(1);; page++ {
		u.Logger.Info("JobSyncProjectTrending.StartGetpagingBtcActivities", zap.Any("page", page))
		resp, err := u.Repo.GetRecentBTCActivity(page, 100)
		if err != nil {
			u.Logger.ErrorAny("JobSyncProjectTrending.ErrorWhenGetPagingActitvities", zap.Error(err))
			break
		}
		uActivities := resp.Result
		btcActivities := uActivities.([]entity.Activity)
		u.Logger.Info("JobSyncProjectTrending.GetpagingActivities", zap.Any("page", page), zap.Any("actitvityCount", len(btcActivities)))
		if len(btcActivities) == 0 {
			break
		}
		for _, btcActivity := range btcActivities {
			var value int64
			if btcActivity.Value > 1000000000 { // this is ETH value
				value = int64(float64(btcActivity.Value) * 0.07 / 1e10) // convert from wei to satoshi
			} else {
				value = btcActivity.Value
			}
			fromProjectIDToRecentVolumn[btcActivity.ProjectID] += value
		}
	}

	fromProjectIDToCountListing := map[string]int64{}
	fromProjectIDToListingVolumn := map[string]int64{}

	for page := int64(1); ; page++ {
		u.Logger.Info("SyncProjectTrending.StartGetpagingListings", zap.Any("page", page))
		listings, err := u.Repo.GetDexBtcsAlongWithProjectInfo(entity.GetDexBtcListingWithProjectInfoReq{
			Page:  page,
			Limit: 100,
		})
		if err != nil {
			u.Logger.ErrorAny("SyncProjectTrending.ErrorWhenGetListings", zap.Any("page", page), zap.Error(err))
			break
		}
		u.Logger.Info("SyncProjectTrending.DoneGetpagingListings", zap.Any("page", page), zap.Any("listing_count", len(listings)))
		if len(listings) == 0 {
			break
		}
		for _, listing := range listings {
			if len(listing.ProjectInfo) == 0 {
				continue
			}
			if listing.Cancelled == true {
				continue
			}
			projectId := listing.ProjectInfo[0].ProjectID
			fromProjectIDToCountListing[projectId]++
			if listing.Matched == true {
				fromProjectIDToListingVolumn[projectId] += int64(listing.Amount)
			}
		}
	}

	var processed int64

	for page := int64(1); ; page++ {
		baseFilter := entity.BaseFilters{
			Limit: 10,
			Page:  page,
		}
		f := entity.FilterProjects{}
		f.BaseFilters = baseFilter
		u.Logger.Info("JobSyncProjectTrending.StartGetpagingProjects", zap.Any("page", page))
		resp, err := u.Repo.GetProjects(f)
		if err != nil {
			u.Logger.ErrorAny("JobSyncProjectTrending.ErrorWhenGetPagingProjects", zap.Any("err", err.Error()))
			break
		}
		uProjects := resp.Result
		projects := uProjects.([]entity.Projects)
		u.Logger.Info("JobSyncProjectTrending.GetpagingProjects", zap.Any("page", page), zap.Any("projectCount", len(projects)))
		if len(projects) == 0 {
			break
		}
		for _, project := range projects {
			processed++
			_countView, err := u.Repo.CountViewActivity(project.TokenID)
			if err != nil {
				return err
			}
			var countView int64 = 0
			if _countView != nil {
				countView = *_countView
			}
			volumnInSatoshi := fromProjectIDToRecentVolumn[project.TokenID]
			volumnInBtc := float64(volumnInSatoshi) / float64(SATOSHI_EACH_BTC)
			volumnInBtc += float64(fromProjectIDToListingVolumn[project.TokenID])
			_countMint, err := u.Repo.CountMintActivity(project.TokenID)
			if err != nil {
				return err
			}
			var numActivity int64 = 0
			if _countMint != nil {
				numActivity = *_countMint
			}

			numListings := fromProjectIDToCountListing[project.TokenID]

			if project.MintingInfo.Index == project.MaxSupply && numListings == 0 {
				numActivity = 0
				volumnInBtc = 0
			}
			trendingScore := countView * TRENDING_SCORE_EACH_VIEW
			trendingScore += int64(volumnInBtc * float64(TRENDING_SCORE_EACH_BTC_VOLUMN))
			trendingScore += numActivity * TRENDING_SCORE_EACH_MINT
			trendingScore += numListings * TRENDING_SCORE_EACH_OPENING_LISTING
			if project.MintingInfo.Index == project.MaxSupply {
				if numListings == 0 {
					trendingScore /= 5
				}
			}
			if project.MintingInfo.Index != project.MaxSupply && project.CreatorAddrr != "0x0000000000000000000000000000000000000000" {
				trendingScore += project.MintingInfo.Index * TRENDING_SCORE_EACH_MINT
			}
			if project.MintPrice == "0" {
				trendingScore /= 5
			}
			isWhitelistedProject := false
			isBoostedProject := false
			// check if this project is whitelisted in top of trending
			for _, str := range u.Config.TrendingConfig.WhitelistedProjectID {
				if project.TokenID == str {
					isWhitelistedProject = true
				}
			}

			if project.Categories != nil {
				for _, str := range project.Categories {
					if str == u.Config.TrendingConfig.BoostedCategoryID {
						isBoostedProject = true
					}
				}
			}

			if isWhitelistedProject {
				trendingScore = INF_TRENDING_SCORE
			} else if isBoostedProject {
				trendingScore *= u.Config.TrendingConfig.BoostedWeight
			}

			u.Repo.UpdateTrendingScoreForProject(project.TokenID, trendingScore)
			u.Logger.Info("SyncProjectTrending.UpdateTrendingScoreForProject", zap.Any("projectID", project.TokenID), zap.Any("trendingScore", trendingScore))
			if numListings != 0 {
				u.Logger.Info("SyncProjectTrending.ProjectHasListing", zap.Any("projectID", project.TokenID), zap.Any("trendingScore", trendingScore), zap.Any("numListings", numListings))
			}
		}
	}

	return nil
}

func (u Usecase) JobDeleteOldActivities() error {
	u.Logger.Info("JobDeleteOldActivities.Start")
	err := u.Repo.JobDeleteOldActivities()
	if err != nil {
		return errors.Wrap(err, "u.Repo.JobDeleteOldActivities")
	}
	return nil
}
