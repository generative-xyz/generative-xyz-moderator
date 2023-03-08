package usecase

import (
	"fmt"
	"math"
	"math/big"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/contracts/generative_dao"
	"rederinghub.io/utils/helpers"
)

func (u *Usecase) PrepareData() error {
	allListings, err := u.Repo.GetAllListings()
	if err != nil {
		return err
	}
	allOffers, err := u.Repo.GetAllOffers()
	if err != nil {
		return err
	}
	allTokens, err := u.Repo.GetAllTokensSeletedFields()
	if err != nil {
		return err
	}
	allProfiles, err := u.Repo.GetAllUserProfiles()
	if err != nil {
		return err
	}
	allProjects, err := u.Repo.GetAllProjectsWithSelectedFields()
	if err != nil {
		return err
	}
	u.gData = gData{
		AllListings: allListings,
		AllOffers:   allOffers,
		AllTokens:   allTokens,
		AllProfile:  allProfiles,
		AllProjects: allProjects,
	}
	return nil
}

func roundTo8DecimalPlaces(x float64) float64 {
	return math.Ceil(x*100000000) / 100000000
}

func (u Usecase) SyncUserStats() error {

	addressToCollectionCreated := make(map[string]int32)
	addressToNftMinted := make(map[string]int32)
	addressToOutputMinted := make(map[string]int32)
	addressToVolumeMinted := make(map[string]float64)

	for _, token := range u.gData.AllTokens {
		u.Logger.Info(fmt.Sprintf("tokenId=%s", token.TokenID), token.TokenID)
		if token.MinterAddress != nil {
			addressToNftMinted[*token.MinterAddress]++
		}

		if token.CreatorAddr != "" {
			addressToOutputMinted[token.CreatorAddr]++
		}
	}

	for _, project := range u.gData.AllProjects {
		addressToCollectionCreated[project.CreatorAddrr]++

		if project.MintingInfo.Index == 0 || project.IsHidden {
			continue
		}

		if _, ok := addressToVolumeMinted[project.CreatorAddrr]; !ok {
			addressToVolumeMinted[project.CreatorAddrr] = 0
		}

		mintPrice, _ := strconv.ParseInt(project.MintPrice, 10, 64)
		amount := mintPrice * project.MintingInfo.Index
		addressToVolumeMinted[project.CreatorAddrr] += roundTo8DecimalPlaces(float64(amount))
	}

	wg := new(sync.WaitGroup)

	updateUserStats := func(wg *sync.WaitGroup, address string, stats entity.UserStats) {
		defer wg.Done()
		//u.Logger.Info(fmt.Sprintf("update user stats address=%s", address), stats)
		u.Repo.UpdateUserStats(address, stats)
	}

	processed := 0
	for _, user := range u.gData.AllProfile {
		if user.WalletAddress == "" {
			continue
		}

		update := false
		collectionCreated := addressToCollectionCreated[user.WalletAddress]
		nftMinted := addressToNftMinted[user.WalletAddress]
		outputMinted := addressToOutputMinted[user.WalletAddress]
		volumeMint := addressToVolumeMinted[user.WalletAddress]

		u.Logger.Info(fmt.Sprintf("address %s collectionCreated %v nftMinted %v", user.WalletAddress, collectionCreated, nftMinted), true)
		if collectionCreated != user.Stats.CollectionCreated {
			user.Stats.CollectionCreated = collectionCreated
			update = true
		}

		if nftMinted != user.Stats.NftMinted {
			user.Stats.NftMinted = nftMinted
			update = true
		}

		if outputMinted != user.Stats.OutputMinted {
			user.Stats.OutputMinted = outputMinted
			update = true
		}

		if volumeMint != user.Stats.VolumeMinted {
			user.Stats.VolumeMinted = volumeMint
			update = true
		}

		if update {
			wg.Add(1)

			go updateUserStats(wg, user.WalletAddress, user.Stats)
			if processed%5 == 0 {
				time.Sleep(5 * time.Second)
			}

			processed++
		}
	}

	wg.Wait()

	return nil
}

func (u Usecase) SyncTokenAndMarketplaceData() error {

	gData := u.gData

	var err error

	errChan := make(chan error, 2)
	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func(wg *sync.WaitGroup, errChan chan error) {
		defer wg.Done()
		err := u.syncMarketplaceDurationAndTokenPrice(&gData)
		errChan <- err
	}(wg, errChan)
	go func(wg *sync.WaitGroup, errChan chan error) {
		defer wg.Done()
		err := u.syncMarketplaceOfferTokenOwner(&gData)
		errChan <- err
	}(wg, errChan)
	wg.Wait()
	close(errChan)

	for e := range errChan {
		if e != nil {
			err = e
			u.Logger.Error(err)
		}
	}

	return err
}

// synchronize token data
func (u Usecase) syncMarketplaceDurationAndTokenPrice(gData *gData) error {
	allListings := u.gData.AllListings
	allOffers := u.gData.AllOffers
	allTokens := u.gData.AllTokens

	// update token price by marketplace data
	activeListings := make([]entity.MarketplaceListings, 0)
	for _, listing := range allListings {
		if !listing.Closed {
			activeListings = append(activeListings, listing)
		}
	}
	activeOffers := make([]entity.MarketplaceOffers, 0)
	for _, offer := range allOffers {
		if !offer.Closed {
			activeOffers = append(activeOffers, offer)
		}
	}

	// sort
	sort.SliceStable(activeListings, func(i, j int) bool {
		if activeListings[i].BlockNumber != activeListings[j].BlockNumber {
			return activeListings[i].CreatedAt.After(*activeListings[j].CreatedAt)
		}
		return activeListings[i].BlockNumber > activeListings[j].BlockNumber
	})

	curTime := time.Now().Unix()
	// update listing/offer that closed
	for id, listing := range activeListings {
		if listing.DurationTime != "0" {
			durationTime, err := strconv.Atoi(listing.DurationTime)
			if err != nil {
				return nil
			}
			// listing is passed deadline
			if int64(durationTime) > curTime {
				u.Repo.CancelListingByOfferingID(listing.OfferingId)
				activeListings[id].Closed = true
			}
		}
	}
	for id, offer := range activeOffers {
		if offer.DurationTime != "0" {
			durationTime, err := strconv.Atoi(offer.DurationTime)
			if err != nil {
				return nil
			}
			// listing is passed deadline
			if int64(durationTime) > curTime {
				u.Repo.CancelOfferByOfferingID(offer.OfferingId)
				activeOffers[id].Closed = true
			}
		}
	}
	// map from token id to price
	fromTokenIdToPrice := make(map[string]int64)
	for _, listing := range activeListings {
		if _, ok := fromTokenIdToPrice[listing.TokenId]; !ok && !listing.Closed {
			price, err := strconv.ParseInt(listing.Price, 10, 64)
			if err != nil {
				return err
			}
			fromTokenIdToPrice[listing.TokenId] = price
		}
	}

	tokenWithPrices := make([]entity.TokenUri, 0)
	for _, token := range allTokens {
		if token.Stats.PriceInt != nil {
			tokenWithPrices = append(tokenWithPrices, token)
		}
	}
	// set of tokens that currently has price
	tokenWithPricesSet := make(map[string]bool)
	for _, token := range tokenWithPrices {
		tokenWithPricesSet[token.TokenID] = true
	}

	for k, v := range fromTokenIdToPrice {
		var token *entity.TokenUri
		for _, _token := range allTokens {
			if _token.TokenID == k {
				token = &_token
				break
			}
		}
		if token == nil {
			return fmt.Errorf("can not find token with tokenID %s", k)
		}
		if token.Stats.PriceInt == nil || *token.Stats.PriceInt != v {
			u.Logger.Info(fmt.Sprintf("setTokenPrice%s", k), v)
			u.Repo.UpdateTokenPriceByTokenId(k, v)
		}
		tokenWithPricesSet[k] = false
	}
	for k, v := range tokenWithPricesSet {
		if !v {
			continue
		}
		u.Logger.Info(fmt.Sprintf("unsetTokenPrice%s", k), true)
		u.Repo.UnsetTokenPriceByTokenId(k)
	}
	return nil
}

func (u Usecase) syncMarketplaceOfferTokenOwner(gData *gData) error {

	allListings := gData.AllListings
	allOffers := gData.AllOffers
	allTokens := gData.AllTokens
	tokenIdToToken := make(map[string]entity.TokenUri)
	for _, token := range allTokens {
		tokenIdToToken[token.TokenID] = token
	}

	updateListingOwner := func(wg *sync.WaitGroup, offeringID string, ownerAddress string) {
		defer wg.Done()
		u.Logger.Info(fmt.Sprintf("update listing offeringId=%s to ownerAddress %s", offeringID, ownerAddress), true)
		u.Repo.UpdateListingOwnerAddress(offeringID, ownerAddress)
	}

	updateOfferOwner := func(wg *sync.WaitGroup, offeringID string, ownerAddress string) {
		defer wg.Done()
		u.Logger.Info(fmt.Sprintf("update offer offeringId=%s to ownerAddress %s", offeringID, ownerAddress), true)
		u.Repo.UpdateOfferOwnerAddress(offeringID, ownerAddress)
	}

	wg := new(sync.WaitGroup)
	counter := 0

	for _, listing := range allListings {
		token, ok := tokenIdToToken[listing.TokenId]
		if !ok {
			return fmt.Errorf("cannot find token with token id %s", listing.TokenId)
		}
		if listing.OwnerAddress == nil || *listing.OwnerAddress != token.OwnerAddr {
			counter++
			if counter%20 == 0 {
				time.Sleep(time.Second)
			}
			wg.Add(1)
			go updateListingOwner(wg, listing.OfferingId, token.OwnerAddr)
		}
	}

	for _, offer := range allOffers {
		token, ok := tokenIdToToken[offer.TokenId]
		if !ok {
			return fmt.Errorf("cannot find token with token id %s", offer.TokenId)
		}
		if offer.OwnerAddress == nil || *offer.OwnerAddress != token.OwnerAddr {
			counter++
			if counter%20 == 0 {
				time.Sleep(time.Second)
			}
			wg.Add(1)
			go updateOfferOwner(wg, offer.OfferingId, token.OwnerAddr)
		}
	}

	wg.Wait()

	return nil
}

func (u Usecase) GetTheCurrentBlockNumber() error {
	block, err := u.Blockchain.GetBlockNumber()
	if err != nil {
		u.Logger.Error("Usecase.GetTheCurrentBlockNumber.GetBlockNumber", err.Error(), err)
		return err
	}

	u.Logger.Info("block", block)
	return nil
}

func (u Usecase) UpdateProposalState() error {
	block, err := u.Blockchain.GetBlock()
	if err != nil {
		u.Logger.Error("Usecase.GetTheCurrentBlockNumber.GetBlockNumber", err.Error(), err)
		return err
	}

	proposals, err := u.Repo.AllProposals(entity.FilterProposals{})
	if err != nil {
		u.Logger.Error("Usecase.GetTheCurrentBlockNumber.AllProposals", err.Error(), err)
		return err
	}

	addr := common.HexToAddress(os.Getenv("DAO_PROPOSAL_CONTRACT"))
	daoContract, err := generative_dao.NewGenerativeDao(addr, u.Blockchain.GetClient())
	if err != nil {
		u.Logger.Error(err)
		return err
	}

	processed := 0
	processChain := make(chan bool, len(proposals))
	for _, proposal := range proposals {

		go func(proposal entity.Proposal) {

			defer func() {
				processChain <- true
			}()

			n := new(big.Int)
			n, ok := n.SetString(proposal.ProposalID, 10)
			if ok {
				state, err := daoContract.State(nil, n)
				if err == nil {
					proposal.State = state
				} else {
					u.Logger.Error(err)
				}

				vote, err := daoContract.Proposals(nil, n)
				if err != nil {
					u.Logger.Error(err)
				} else {
					//createdProposal.State = state
					u.Logger.Info("daoContract.Proposals.vote", vote)
				}

				forVote := helpers.ParseBigToFloat(vote.ForVotes)
				againstVote := helpers.ParseBigToFloat(vote.AgainstVotes)
				abstainVote := helpers.ParseBigToFloat(vote.AbstainVotes)
				percentFor := float64(0)
				percentAgainst := float64(0)
				percentAbstain := float64(0)

				total := forVote + againstVote + abstainVote
				if total != 0 {
					percentFor = float64((forVote / total) * 100)
					percentAgainst = float64((againstVote / total) * 100)
					percentAbstain = float64((abstainVote / total) * 100)
				}

				proposal.Vote = entity.ProposalVote{
					For:            vote.ForVotes.String(),
					ForNum:         forVote,
					Against:        vote.AgainstVotes.String(),
					Abstain:        vote.AbstainVotes.String(),
					Total:          fmt.Sprintf("%f", total),
					TotalNum:       total,
					PercentFor:     percentFor,
					PercentAgainst: percentAgainst,
					PercentAbstain: percentAbstain,
				}

				if proposal.ProposalID == "35751750717610809166312996604681477486540366891662940411672289868284123500445" {
					test := helpers.ParseBigToFloat(vote.ForVotes)
					spew.Dump(proposal.Vote)
					_ = test
				}
			}
			proposal.CurrentBlock = block.Number.Int64()
			proposal.CurrentBlockTime = helpers.ParseUintToUnixTime(block.Time)

			stB, err := u.Blockchain.GetBlockByNumber(*big.NewInt(proposal.StartBlock))
			if err == nil {
				proposal.StartBlockTime = helpers.ParseUintToUnixTime(stB.Time())
			}

			eBB, err := u.Blockchain.GetBlockByNumber(*big.NewInt(proposal.EndBlock))
			if err == nil {
				proposal.EndBlockTime = helpers.ParseUintToUnixTime(eBB.Time())
			}

			updated, err := u.Repo.UpdateProposal(proposal.UUID, &proposal)
			if err != nil {
				u.Logger.Error(err)
			}
			u.Logger.Info("Updated", updated)
		}(proposal)

		if processed%10 == 0 {
			time.Sleep(5 * time.Second)
		}
	}

	for i := 0; i < len(proposals); i++ {
		<-processChain
	}

	return nil
}

func (u Usecase) SyncLeaderboard() error {

	allUsers := u.gData.AllProfile
	addressToProfile := make(map[string]entity.Users)
	for _, user := range allUsers {
		addressToProfile[user.WalletAddress] = user
	}

	allTokenHolders, err := u.Repo.GetAllTokenHolders()
	if err != nil {
		return err
	}

	addressToOldRank := make(map[string]*int32)
	addressToOldBalance := make(map[string]*string)
	for _, tokenHolder := range allTokenHolders {
		addressToOldRank[tokenHolder.Address] = tokenHolder.OldRank
		addressToOldBalance[tokenHolder.Address] = tokenHolder.OldBalance
	}

	allNewTokenHolders, err := u.GetAllTokenHolder()
	if err != nil {
		return err
	}

	sort.SliceStable(allNewTokenHolders, func(i, j int) bool {
		lhs := new(big.Int)
		lhs, ok := lhs.SetString(allNewTokenHolders[i].Balance, 10)
		if !ok {
			lhs = big.NewInt(0)
		}
		rhs := new(big.Int)
		rhs, ok = rhs.SetString(allNewTokenHolders[j].Balance, 10)
		if !ok {
			rhs = big.NewInt(0)
		}
		return lhs.Cmp(rhs) > 0
	})

	tokenHolders := make([]entity.TokenHolder, 0, len(allNewTokenHolders))

	addressToProjects := map[string][]entity.Projects{}
	for _, project := range u.gData.AllProjects {
		_, ok := addressToProjects[project.CreatorAddrr]
		if !ok {
			addressToProjects[project.CreatorAddrr] = []entity.Projects{}
		}
		addressToProjects[project.CreatorAddrr] = append(addressToProjects[project.CreatorAddrr], project)
	}

	// map from user's address to set of owner of user's token
	addressToOwners := map[string]map[string]bool{}
	for _, token := range u.gData.AllTokens {
		if token.OwnerAddr == "" {
			continue
		}
		_, ok := addressToOwners[token.CreatorAddr]
		if !ok {
			addressToOwners[token.CreatorAddr] = map[string]bool{}
		}
		addressToOwners[token.CreatorAddr][token.OwnerAddr] = true
	}

	// add rank
	for l, r := 0, 1; l < len(allNewTokenHolders); l = r {
		for r < len(allNewTokenHolders) && allNewTokenHolders[r].Balance == allNewTokenHolders[l].Balance {
			r++
		}
		for i := l; i < r; i++ {
			_tokenHolder := &allNewTokenHolders[i]
			tokenHolder := entity.TokenHolder{
				ContractDecimals:     _tokenHolder.ContractDecimals,
				ContractName:         _tokenHolder.ContractName,
				ContractTickerSymbol: _tokenHolder.ContractTickerSymbol,
				ContractAddress:      _tokenHolder.ContractAddress,
				SupportsErc:          _tokenHolder.SupportsErc,
				LogoURL:              _tokenHolder.LogoURL,
				Address:              _tokenHolder.Address,
				Balance:              _tokenHolder.Balance,
				TotalSupply:          _tokenHolder.TotalSupply,
				BlockHeight:          _tokenHolder.BlockHeight,
				CurrentRank:          int32(l + 1),
				OldRank:              addressToOldRank[_tokenHolder.Address],
				OldBalance:           addressToOldBalance[_tokenHolder.Address],
			}
			profile, ok := addressToProfile[_tokenHolder.Address]
			pProfile := &profile
			if !ok {
				pProfile = nil
			}
			tokenHolder.Profile = pProfile
			tokenHolder.Projects = addressToProjects[tokenHolder.Address]
			var ownerCount int32
			_, ok = addressToOwners[tokenHolder.Address]
			if ok {
				ownerCount = int32(len(addressToOwners[tokenHolder.Address]))
			}
			tokenHolder.OwnerCount = ownerCount
			tokenHolders = append(tokenHolders, tokenHolder)
		}
	}

	err = u.Repo.DeleteAllTokenHolders()

	if err != nil {
		return err
	}

	err = u.Repo.CreateTokenHolders(tokenHolders)

	return err
}

func (u Usecase) SnapShotOldRankAndOldBalance() error {
	_, err := u.Repo.SnapShotOldRankAndOldBalance()
	return err
}

// Currently, this function only syncs projects' nft minted data
// TODO: move all other stats to be synced in this function
func (u Usecase) SyncProjectsStats() error {
	allProjects := u.gData.AllProjects
	allTokens := u.gData.AllTokens
	projectIdToMintedCount := map[string]int32{}
	for _, token := range allTokens {
		projectIdToMintedCount[token.ProjectID]++
	}
	var processed int32
	for _, project := range allProjects {
		mintedCount := projectIdToMintedCount[project.TokenID]
		if mintedCount != project.Stats.MintedCount {
			processed++
			_, err := u.Repo.UpdateProjectMintedCount(project.UUID, mintedCount)
			if err != nil {
				return err
			}
			if processed%10 == 0 {
				time.Sleep(1 * time.Second)
			}
		}
	}

	return nil
}

func (u Usecase) SyncTokenInscribeIndex() error {
	notSyncedTokens, err := u.Repo.GetAllNotSyncInscriptionIndexToken()
	if err != nil {
		return err
	}
	processed := 0
	for _, token := range notSyncedTokens {
		processed++
		inscribeInfo, err := u.GetInscribeInfo(token.TokenID)
		if err != nil {
			u.Logger.Error("FailedToGetInscribeInfo", zap.String("tokenID", token.TokenID), zap.Error(err))
			continue
		}
		u.Repo.UpdateTokenInscriptionIndex(token.TokenID, inscribeInfo.Index)

		if token.OwnerAddr != inscribeInfo.Address {
			u.Repo.UpdateTokenOwnerAddr(inscribeInfo.ID, inscribeInfo.Address)
		}
		// try to find user with address as btc address
		profile, err := u.Repo.FindUserByBtcAddress(inscribeInfo.Address)
		if err == nil && profile != nil {
			u.Repo.UpdateTokenOwner(inscribeInfo.ID, profile)
		}
		if processed%10 == 0 {
			time.Sleep(time.Second)
		}
	}
	return nil
}

const (
	INF_TRENDING_SCORE             int64 = 9223372036854775807 // max int64 value
	SATOSHI_EACH_BTC               int64 = 100000000
	TRENDING_SCORE_EACH_BTC_VOLUMN int64 = 1000000
	TRENDING_SCORE_EACH_OPENING_LISTING int64 = 5000
	TRENDING_SCORE_EACH_MINT       int64 = 10
	TRENDING_SCORE_EACH_VIEW       int64 = 1
)

func (u Usecase) SyncProjectTrending() error {
	u.Logger.Info("SyncProjectTrending.StartSyncProjectTrending")
	// All btc activities, which include Mint and Buy activity
	btcActivites, err := u.Repo.GetRecentBTCActivity()
	if err != nil {
		u.Logger.ErrorAny("SyncProjectTrending.ErrorWhenGetBtcActivities", zap.Any("err", err.Error()))
		return err
	}

	u.Logger.Info("SyncProjectTrending.DoneGetBtcActivities", zap.Any("act_len", len(btcActivites)))

	// Mapping from projectID to latest 24h's volumn in satoshi
	fromProjectIDToRecentVolumn := map[string]int64{}
	for _, btcActivity := range btcActivites {
		var value int64
		if btcActivity.Value > 1000000000 { // this is ETH value
			value = int64(float64(btcActivity.Value) * 0.07 / 1e10) // convert from wei to satoshi
		} else {
			value = btcActivity.Value
		}
		fromProjectIDToRecentVolumn[btcActivity.ProjectID] += value
	}

	fromProjectIDToCountListing := map[string]int64{}

	for page := int64(1);; page++ {
		u.Logger.Info("SyncProjectTrending.StartGetpagingListings", zap.Any("page", page))
		listings, err := u.Repo.GetDexBtcsAlongWithProjectInfo(entity.GetDexBtcListingWithProjectInfoReq{
			Page: page,
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
		u.Logger.Info("SyncProjectTrending.StartGetpagingProjects", zap.Any("page", page))
		resp, err := u.Repo.GetProjects(f)
		if err != nil {
			u.Logger.ErrorAny("SyncProjectTrending.ErrorWhenGetPagingProjects", zap.Any("err", err.Error()))
			break
		}
		uProjects := resp.Result
		projects := uProjects.([]entity.Projects)
		u.Logger.Info("SyncProjectTrending.GetpagingProjects", zap.Any("page", page), zap.Any("projectCount", len(projects)))
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
			numActivity := int64(len(btcActivites))

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
				} else {
					trendingScore *= 2
				}
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

func (u Usecase) DeleteOldActivities() error {
	u.Logger.Info("DeleteOldActivities.Start")
	err := u.Repo.DeleteOldActivities()
	if err != nil {
		return errors.Wrap(err, "u.Repo.DeleteOldActivities")
	}
	return nil
}
