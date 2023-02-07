package usecase

import (
	"fmt"
	"math/big"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/opentracing/opentracing-go"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/contracts/generative_dao"
)

func (u *Usecase) PrepareData(rootSpan opentracing.Span) (error) {
	span, log := u.StartSpan("SyncTokenAndMarketplaceData", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
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
		AllOffers: allOffers,
		AllTokens: allTokens,
		AllProfile: allProfiles,
		AllProjects: allProjects,
	}
	return nil
}

func (u Usecase) SyncUserStats(rootSpan opentracing.Span) error {
	span, log := u.StartSpan("SyncUserStats", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	addressToCollectionCreated := make(map[string]int32)
	addressToNftMinted := make(map[string]int32)
	addressToOutputMinted := make(map[string]int32)

	for _, token := range u.gData.AllTokens {
		log.SetData(fmt.Sprintf("tokenId=%s", token.TokenID), token.TokenID)
		if token.MinterAddress != nil {
			addressToNftMinted[*token.MinterAddress]++
		}
		if token.CreatorAddr != "" {
			addressToOutputMinted[token.CreatorAddr]++
		}
	}

	for _, project := range u.gData.AllProjects {
		addressToCollectionCreated[project.CreatorAddrr]++
	}

	wg := new(sync.WaitGroup)

	updateUserStats := func (wg *sync.WaitGroup, address string, stats entity.UserStats) {
		defer wg.Done()
		//log.SetData(fmt.Sprintf("update user stats address=%s", address), stats)
		u.Repo.UpdateUserStats(address, stats)
	}

	processed := 0
	for _, user := range u.gData.AllProfile {
		update := false
		collectionCreated := addressToCollectionCreated[user.WalletAddress]
		nftMinted := addressToNftMinted[user.WalletAddress]
		outputMinted := addressToOutputMinted[user.WalletAddress]
		log.SetData(fmt.Sprintf("address %s collectionCreated %v nftMinted %v", user.WalletAddress, collectionCreated, nftMinted), true)
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
		if update {
			wg.Add(1)

			go updateUserStats(wg, user.WalletAddress, user.Stats)
			if processed % 5 == 0 {
				time.Sleep(5 * time.Second)
			}

			processed ++
		}
	}

	wg.Wait()

	return nil
}

func (u Usecase) SyncTokenAndMarketplaceData(rootSpan opentracing.Span) error {
	span, log := u.StartSpan("SyncTokenAndMarketplaceData", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	gData := u.gData

	var err error

	errChan := make(chan error, 2)
	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func(wg *sync.WaitGroup, errChan chan error) {
		defer wg.Done()
		err := u.syncMarketplaceDurationAndTokenPrice(span, &gData)
		errChan <- err
	}(wg, errChan)
		
	go func(wg *sync.WaitGroup, errChan chan error) {
		defer wg.Done()
		err := u.syncMarketplaceOfferTokenOwner(span, &gData)
		errChan <- err
	}(wg, errChan)
	
	wg.Wait()
	close(errChan)

	for e := range errChan {
		if e != nil {
			err = e
			log.Error("error when sync data", err.Error(), err)
		}
	}

	return err
}

// synchronize token data
func (u Usecase) syncMarketplaceDurationAndTokenPrice(rootSpan opentracing.Span, gData *gData) error {
	span, log := u.StartSpan("syncMarketplaceDurationAndTokenPrice", rootSpan)
	defer u.Tracer.FinishSpan(span, log )

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
	});

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

func (u Usecase) syncMarketplaceOfferTokenOwner(rootSpan opentracing.Span, gData *gData) error {
	span, log := u.StartSpan("syncMarketplaceOfferTokenOwner", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	allListings := gData.AllListings
	allOffers := gData.AllOffers
	allTokens := gData.AllTokens
	
	tokenIdToToken := make(map[string]entity.TokenUri)
	for _, token := range allTokens {
		tokenIdToToken[token.TokenID] = token
	}

	updateListingOwner := func (wg *sync.WaitGroup, offeringID string, ownerAddress string) {
		defer wg.Done()
		log.SetData(fmt.Sprintf("update listing offeringId=%s to ownerAddress %s", offeringID, ownerAddress), true)
		u.Repo.UpdateListingOwnerAddress(offeringID, ownerAddress)
	}

	updateOfferOwner := func (wg *sync.WaitGroup, offeringID string, ownerAddress string) {
		defer wg.Done()
		log.SetData(fmt.Sprintf("update offer offeringId=%s to ownerAddress %s", offeringID, ownerAddress), true)
		u.Repo.UpdateOfferOwnerAddress(offeringID, ownerAddress)
	}

	wg := new(sync.WaitGroup)
	
	counter := 0;

	for _, listing := range allListings {
		token, ok := tokenIdToToken[listing.TokenId]
		if !ok {
			return fmt.Errorf("cannot find token with token id %s", listing.TokenId)
		}
		if listing.OwnerAddress == nil || *listing.OwnerAddress != token.OwnerAddr {
			counter++
			if counter % 20 == 0 {
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
			if counter % 20 == 0 {
				time.Sleep(time.Second)
			}
			wg.Add(1)
			go updateOfferOwner(wg, offer.OfferingId, token.OwnerAddr)
		}
	}

	wg.Wait()

	return nil
}

func (u Usecase) GetTheCurrentBlockNumber(rootSpan opentracing.Span) error {
	span, log := u.StartSpan("Usecase.GetTheCurrentBlockNumber", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	block, err := u.Blockchain.GetBlockNumber()
	if err != nil {
		log.Error("Usecase.GetTheCurrentBlockNumber.GetBlockNumber",err.Error(), err)
		return err
	}

	log.SetData("block",block)
	return nil
}

func (u Usecase) UpdateProposalState(rootSpan opentracing.Span) error {
	span, log := u.StartSpan("Usecase.UpdateProposalState", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	block, err := u.Blockchain.GetBlockNumber()
	if err != nil {
		log.Error("Usecase.GetTheCurrentBlockNumber.GetBlockNumber",err.Error(), err)
		return err
	}

	proposals, err := u.Repo.AllProposals(entity.FilterProposals{})
	if err != nil {
		log.Error("Usecase.GetTheCurrentBlockNumber.AllProposals",err.Error(), err)
		return err
	}

	addr := common.HexToAddress(os.Getenv("DAO_PROPOSAL_CONTRACT"))
	daoContract, err := generative_dao.NewGenerativeDao(addr, u.Blockchain.GetClient())
	if  err != nil {
		log.Error("cannot init DAO contract", err.Error(), err)
		return err
	}

	processed := 0
	processChain := make(chan bool, len(proposals))
	
	for _, proposal := range proposals {

		go func ( proposal entity.Proposal, processChain chan bool)  {
			span, log := u.StartSpan("Usecase.UpdateProposalState.Routine", rootSpan)
			defer u.Tracer.FinishSpan(span, log)

			defer func(){
				processChain <- true
			}()

			n := new(big.Int)
			n, ok := n.SetString(proposal.ProposalID, 10)
			if ok {
				state, err := daoContract.State(nil, n)
				if err == nil {
					proposal.State = state
				}else{
					log.Error("daoContract.State", err.Error(), err)
				}

				vote, err := daoContract.Proposals(nil, n)
				if err != nil {
					log.Error("daoContract.Proposals.vote.Error", err.Error(), err)
				}else{
					//createdProposal.State = state
					log.SetData("daoContract.Proposals.vote", vote)
				}

				forVote :=  vote.ForVotes.Uint64()
				againstVote := vote.AgainstVotes.Uint64()
				abstainVote := vote.AbstainVotes.Uint64()
				percentFor := float64(0)
				percentAgainst := float64(0)
				percentAbstain := float64(0)

				total := forVote + againstVote + abstainVote
				if total != 0 {
					percentFor =  float64((forVote / total ) * 100)
					percentAgainst =  float64((againstVote / total ) * 100)
					percentAbstain =  float64((abstainVote / total ) * 100)
				}
				
				proposal.Vote = entity.ProposalVote{
					For: forVote,
					Against: againstVote,
					Abstain: abstainVote,
					Total: total,
					PercentFor: percentFor,
					PercentAgainst: percentAgainst,
					PercentAbstain: percentAbstain,
				}
			}
	
			proposal.CurrentBlock = block.Int64()
			updated, err := u.Repo.UpdateProposal(proposal.UUID, &proposal)
			if err != nil {
				log.Error("daoContract.State", err.Error(), err)
			}
			log.SetData("Updated", updated)
			
		}(proposal, processChain)

		if processed % 10 == 0{
			time.Sleep(5 * time.Second)
		}
	}

	for i := 0; i< len(proposals) ; i ++ {
		<- processChain
	}

	return nil
}
