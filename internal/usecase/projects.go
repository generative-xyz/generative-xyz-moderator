package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jinzhu/copier"
	"github.com/opentracing/opentracing-go"

	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/contracts/generative_nft_contract"
	"rederinghub.io/utils/contracts/generative_project_contract"
	"rederinghub.io/utils/helpers"
)

func (u Usecase) CreateProject(rootSpan opentracing.Span, req structure.CreateProjectReq) (*entity.Projects, error) {
	span, log := u.StartSpan("CreateProject", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	pe := &entity.Projects{}
	err := copier.Copy(pe, req)
	if err != nil {
		log.Error("copier.Copy", err.Error(), err)
		return nil, err
	}

	err = u.Repo.CreateProject(pe)
	if err != nil {
		log.Error("u.Repo.CreateProject", err.Error(), err)
		return nil, err
	}

	log.SetData("pe", pe)
	return pe, nil
}

func (u Usecase) UpdateProject(rootSpan opentracing.Span, req structure.UpdateProjectReq) (*entity.Projects, error) {
	span, log := u.StartSpan("UpdateProject", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	resp := &entity.Projects{}

	return resp, nil
}

func (u Usecase) GetProjectByGenNFTAddr(rootSpan opentracing.Span, genNFTAddr string) (*entity.Projects, error) {
	span, log := u.StartSpan("GetProjects", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	project, err := u.Repo.FindProjectByGenNFTAddr(genNFTAddr)
	return project, err
}

func (u Usecase) GetProjects(rootSpan opentracing.Span, req structure.FilterProjects) (*entity.Pagination, error) {
	span, log := u.StartSpan("GetProjects", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	pe := &entity.FilterProjects{}
	err := copier.Copy(pe, req)
	if err != nil {
		log.Error("copier.Copy", err.Error(), err)
		return nil, err
	}
	
	projects, err := u.Repo.GetProjects(*pe)
	if err != nil {
		log.Error("u.Repo.GetProjects", err.Error(), err)
		return nil, err
	}

	log.SetData("projects", projects)
	return projects, nil
}

func (u Usecase) GetRandomProject(rootSpan opentracing.Span) (*entity.Projects, error) {
	span, log := u.StartSpan("GetRandomProject", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	key := helpers.ProjectRandomKey()

	//always reload data
	go func() {
		p, err := u.Repo.GetAllProjects(entity.FilterProjects{})
		if err != nil {
			return
		}
		u.Cache.SetData(key, p)
	}()

	cached, err := u.Cache.GetData(key)
	if err != nil {
		p, err := u.Repo.GetAllProjects(entity.FilterProjects{})
		if err != nil {
			log.Error("u.Repo.GetProjects", err.Error(), err)
			return nil, err
		}
		u.Cache.SetData(key, p)
	}

	cached, err = u.Cache.GetData(key)
	projects := []entity.Projects{}
	bytes := []byte(*cached)
	err = json.Unmarshal(bytes, &projects)
	if err != nil {
		log.Error("json.Unmarshal", err.Error(), err)
		return nil, err
	}

	if len(projects) == 0 {
		err := errors.New("Project are not found")
		log.Error("Projects.are.not.found", err.Error(), err)
		return nil, err 
	}

	timeNow := time.Now().UTC().Nanosecond()
	rand := int(timeNow) % len(projects)

	//TODO - cache will be applied here

	projectRand := projects[rand]
	return u.GetProjectDetail(span, structure.GetProjectDetailMessageReq{
		ContractAddress: projectRand.ContractAddress,
		ProjectID:       projectRand.TokenID,
	})
}

func (u Usecase) GetMintedOutProjects(rootSpan opentracing.Span, req structure.FilterProjects) (*entity.Pagination, error) {
	span, log := u.StartSpan("GetMintedOutProjects", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	pe := &entity.FilterProjects{}
	err := copier.Copy(pe, req)
	if err != nil {
		log.Error("copier.Copy", err.Error(), err)
		return nil, err
	}

	pe.WalletAddress = req.WalletAddress
	projects, err := u.Repo.GetMintedOutProjects(*pe)
	if err != nil {
		log.Error("u.Repo.GetMintedOutProjects", err.Error(), err)
		return nil, err
	}

	log.SetData("projects", projects)
	return projects, nil
}

func (u Usecase) GetProjectDetail(rootSpan opentracing.Span, req structure.GetProjectDetailMessageReq) (*entity.Projects, error) {
	span, log := u.StartSpan("GetProjectDetail", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	defer func  ()  {
		//alway update project in a separated process
		go func(rootSpan opentracing.Span) {
			span, log := u.StartSpan("GetProjectDetail.GetProjectFromChain", rootSpan)
			defer u.Tracer.FinishSpan(span, log)

			_, err := u.UpdateProjectFromChain(span, req.ContractAddress, req.ProjectID)
			if err != nil {
				log.Error("u.Repo.FindProjectBy", err.Error(), err)
				return
			}

		}(span)	
	}()

	log.SetTag("ProjectID", req.ProjectID)
	log.SetTag("ContractAddress", req.ContractAddress)

	c, _ := u.Repo.FindProjectBy(req.ContractAddress, req.ProjectID)
	if (c == nil) || (c != nil && !c.IsSynced) || c.MintedTime == nil {
		p, err := u.UpdateProjectFromChain(span, req.ContractAddress, req.ProjectID)
		if err != nil {
			log.Error("u.Repo.FindProjectBy", err.Error(), err)
			return nil, err
		}
		return p, nil
	}
	return c, nil
}

func (u Usecase) GetRecentWorksProjects(rootSpan opentracing.Span, req structure.FilterProjects) (*entity.Pagination, error) {
	span, log := u.StartSpan("GetRecentWorksProjects", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	pe := &entity.FilterProjects{}
	err := copier.Copy(pe, req)
	if err != nil {
		log.Error("copier.Copy", err.Error(), err)
		return nil, err
	}

	pe.WalletAddress = req.WalletAddress
	projects, err := u.Repo.GetRecentWorksProjects(*pe)
	if err != nil {
		log.Error("u.Repo.GetRecentWorksProjects", err.Error(), err)
		return nil, err
	}

	log.SetData("projects", projects)
	return projects, nil
}

func (u Usecase) GetUpdatedProjectStats(rootSpan opentracing.Span, req structure.GetProjectReq) (*entity.ProjectStat, error) {
	span, log := u.StartSpan("SyncProjectStats", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	project, err := u.Repo.FindProjectBy(req.ContractAddr, req.TokenID)
	if err != nil {
		return nil, err
	}

	// do not resync
	if project.Stats.LastTimeSynced != nil && project.Stats.LastTimeSynced.Unix() + int64(u.Config.TimeResyncProjectStat) > time.Now().Unix() {
		return &project.Stats, nil
	}

	allTokenFromDb, err := u.Repo.GetAllTokensByProjectID(project.TokenID)
	if err != nil {
		return nil, err
	}
	owners := make(map[string]bool)
	for _, token := range allTokenFromDb {
		owners[token.OwnerAddr] = true
	}

	var allListings []entity.MarketplaceListings
	var allOffers []entity.MarketplaceOffers

	allListings, err = u.Repo.GetAllListingByCollectionContract(project.GenNFTAddr)
	if err != nil {
		log.Error("u.Repo.GetAllListingByCollectionContract", err.Error(), err)
		return nil, err
	}

	allOffers, err = u.Repo.GetAllOfferByCollectionContract(project.GenNFTAddr)
	if err != nil {
		log.Error("u.Repo.GetAllOfferByCollectionContract", err.Error(), err)
		return nil, err
	}

	var totalTradingVolumn *big.Int
	var floorPrice *big.Int
	var bestMakeOfferPrice *big.Int
	var listedPercent int32
	listingSet := make(map[string]bool)

	for _, listing := range allListings {
		if listing.Erc20Token != utils.EVM_NULL_ADDRESS {
			continue
		}
		price := new(big.Int)
		price, ok := price.SetString(listing.Price, 10)
		if !ok {
			err := errors.New("fail to convert price to big int")
			log.Error("fail to convert price to big int", err.Error(), err)
			continue
		}
		durationTime, err := strconv.ParseInt(listing.DurationTime, 10, 64)
		if err != nil {
			log.Error("fail to parse duration time", err.Error(), err)
			continue
		}

		// update total volumn trading
		if listing.Finished {
			if totalTradingVolumn == nil {
				totalTradingVolumn = new(big.Int)
			}
			totalTradingVolumn.Add(totalTradingVolumn, price)
		}
		// update floor price
		if !listing.Closed && (time.Now().Unix() < durationTime || durationTime == 0) {
			if floorPrice == nil {
				floorPrice = price
			} else {
				if floorPrice.Cmp(price) > 0 {
					floorPrice = price
				}
				listingSet[listing.CollectionContract] = true
			}
		}
	}

	for _, offer := range allOffers {
		price := new(big.Int)
		price, ok := price.SetString(offer.Price, 10)
		if !ok {
			err := errors.New("fail to convert price to big int")
			log.Error("fail to convert price to big int", err.Error(), err)
			continue
		}
		durationTime, err := strconv.ParseInt(offer.DurationTime, 10, 64)
		if err != nil {
			log.Error("fail to parse duration time", err.Error(), err)
			continue
		}

		// update total volumn trading
		if offer.Finished {
			if totalTradingVolumn == nil {
				totalTradingVolumn = new(big.Int)
			}
			totalTradingVolumn.Add(totalTradingVolumn, price)
		}

		// update floor price
		if !offer.Closed && (time.Now().Unix() < durationTime || durationTime == 0) {
			if bestMakeOfferPrice == nil {
				bestMakeOfferPrice = price
			} else {
				if bestMakeOfferPrice.Cmp(price) < 0 {
					bestMakeOfferPrice = price
				}
			}
		}

	}

	if len(allTokenFromDb) > 0 {
		listedPercent = int32(len(listingSet) * 100 / len(allTokenFromDb))
	} else {
		listedPercent = 0
	}
	

	if totalTradingVolumn == nil {
		totalTradingVolumn = new(big.Int)
	}
	if floorPrice == nil {
		floorPrice = new(big.Int)
	}
	if bestMakeOfferPrice == nil {
		bestMakeOfferPrice = new(big.Int)
	}

	now := time.Now()
	return &entity.ProjectStat{
		LastTimeSynced: &now,
		UniqueOwnerCount: uint32(len(owners)),
		TotalTradingVolumn: totalTradingVolumn.String(),
		FloorPrice: floorPrice.String(),
		BestMakeOfferPrice: bestMakeOfferPrice.String(),
		ListedPercent: listedPercent,
	}, nil
}


func (u Usecase) getProjectDetailFromChainWithoutCache(rootSpan opentracing.Span, req structure.GetProjectDetailMessageReq) (*structure.ProjectDetail, error) {
	span, log := u.StartSpan("getProjectDetailFromChainWithoutCache", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	contractDataKey := fmt.Sprintf("detail.%s.%s", req.ContractAddress, req.ProjectID)

	log.SetData("req", req)
	
	addr := common.HexToAddress(req.ContractAddress)
	// call to contract to get emotion
	client, err := helpers.EthDialer()
	if err != nil {
		log.Error("ethclient.Dial", err.Error(), err)
		return nil, err
	}

	projectID := new(big.Int)
	projectID, ok := projectID.SetString(req.ProjectID, 10)
	if !ok {
		return nil, errors.New("cannot convert tokenID")
	}
	contractDetail, err := u.getNftContractDetailInternal(client, addr, *projectID)
	if err != nil {
		log.Error("u.getNftContractDetailInternal", err.Error(), err)
		return nil, err
	}
	log.SetData("contractDetail", contractDetail)
	u.Cache.SetData(contractDataKey, contractDetail)
	return contractDetail, nil
}

// Get from chain with cache
func (u Usecase) getProjectDetailFromChain(rootSpan opentracing.Span, req structure.GetProjectDetailMessageReq) (*structure.ProjectDetail, error) {
	span, log := u.StartSpan("getProjectDetailFromChain", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	contractDataKey := fmt.Sprintf("detail.%s.%s", req.ContractAddress, req.ProjectID)

	//u.Cache.Delete(contractDataKey)
	data, err := u.Cache.GetData(contractDataKey)
	if err != nil {
		log.SetData("req", req)
		
		addr := common.HexToAddress(req.ContractAddress)
		// call to contract to get emotion
		client, err := helpers.EthDialer()
		if err != nil {
			log.Error("ethclient.Dial", err.Error(), err)
			return nil, err
		}

		projectID := new(big.Int)
		projectID, ok := projectID.SetString(req.ProjectID, 10)
		if !ok {
			return nil, errors.New("cannot convert tokenID")
		}
		contractDetail, err := u.getNftContractDetailInternal(client, addr, *projectID)
		if err != nil {
			log.Error("u.getNftContractDetail", err.Error(), err)
			return nil, err
		}
		log.SetData("contractDetail", contractDetail)
		u.Cache.SetData(contractDataKey, contractDetail)
		return contractDetail, nil
	}

	
	contractDetail := &structure.ProjectDetail{}
	err = helpers.ParseCache(data, contractDetail)
	if err != nil {
		log.Error("helpers.ParseCache", err.Error(), err)
		return nil, err
	}

	return contractDetail, nil
}

// Internal get project detail
func (u Usecase) getNftContractDetailInternal(client *ethclient.Client, contractAddr common.Address, projectID big.Int) (*structure.ProjectDetail, error) {
	gProject, err := generative_project_contract.NewGenerativeProjectContract(contractAddr, client)
	if err != nil {
		return nil, err
	}

	pDchan := make(chan structure.ProjectDetailChan, 1)
	pStatuschan := make(chan structure.ProjectStatusChan, 1)
	pTokenURIchan := make(chan structure.ProjectNftTokenUriChan, 1)

	go func(pDchan chan structure.ProjectDetailChan, projectID *big.Int) {
		proDetail := &generative_project_contract.NFTProjectProject{}
		var err error

		defer func() {
			pDchan <- structure.ProjectDetailChan{
				ProjectDetail: proDetail,
				Err:           err,
			}
		}()

		proDetailReps, err := gProject.ProjectDetails(nil, projectID)
		if err != nil {
			return
		}

		proDetail = &proDetailReps

	}(pDchan, &projectID)

	go func(pDchan chan structure.ProjectStatusChan, projectID *big.Int) {
		var status *bool
		var err error

		defer func() {
			pDchan <- structure.ProjectStatusChan{
				Status: status,
				Err:    err,
			}
		}()

		pStatus, err := gProject.ProjectStatus(nil, projectID)
		if err != nil {
			return
		}

		status = &pStatus

	}(pStatuschan, &projectID)

	go func(pDchan chan structure.ProjectNftTokenUriChan, projectID *big.Int) {
		var tokenURI *string
		var err error

		defer func() {
			pDchan <- structure.ProjectNftTokenUriChan{
				TokenURI: tokenURI,
				Err:      err,
			}
		}()

		pTokenUri, err := gProject.TokenURI(nil, projectID)
		if err != nil {
			return
		}

		tokenURI = &pTokenUri

	}(pTokenURIchan, &projectID)

	detailFromChain := <-pDchan
	statusFromChain := <-pStatuschan
	tokenFromChain := <-pTokenURIchan

	if detailFromChain.Err != nil {
		return nil, detailFromChain.Err
	}

	if statusFromChain.Err != nil {
		return nil, statusFromChain.Err
	}

	if tokenFromChain.Err != nil {
		return nil, tokenFromChain.Err
	}

	gNftProject, err := generative_nft_contract.NewGenerativeNftContract(detailFromChain.ProjectDetail.GenNFTAddr, client)
	if err != nil {
		return nil, err
	}

	//nft project detail chain
	nftProjectDchan := make(chan structure.NftProjectDetailChan, 1)
	go func(nftProjectDchan chan structure.NftProjectDetailChan, gNftProject *generative_nft_contract.GenerativeNftContract) {
		data := &structure.NftProjectDetail{}
		var err error

		defer func() {
			nftProjectDchan <- structure.NftProjectDetailChan{
				Data: data,
				Err:  err,
			}
		}()

		respData, err := gNftProject.Project(nil)
		err = copier.Copy(data, respData)

	}(nftProjectDchan, gNftProject)

	nftRoyaltychan := make(chan structure.RoyaltyChan, 1)
	go func(nftRoyaltychan chan structure.RoyaltyChan, gNftProject *generative_nft_contract.GenerativeNftContract) {
		var data *big.Int
		var err error

		defer func() {
			nftRoyaltychan <- structure.RoyaltyChan{
				Data: data,
				Err:  err,
			}
		}()

		data, err = gNftProject.Royalty(nil)

	}(nftRoyaltychan, gNftProject)

	dataFromNftPChan := <-nftProjectDchan
	dataFromRoyaltyPChan := <-nftRoyaltychan

	resp := &structure.ProjectDetail{
		ProjectDetail: detailFromChain.ProjectDetail,
		Status:        *statusFromChain.Status,
		NftTokenUri:   *tokenFromChain.TokenURI,
	}

	if dataFromNftPChan.Err == nil && dataFromNftPChan.Data != nil {
		resp.NftProjectDetail = *dataFromNftPChan.Data
	} else {
		resp.NftProjectDetail = structure.NftProjectDetail{}
	}

	if dataFromRoyaltyPChan.Err == nil && dataFromRoyaltyPChan.Data != nil {
		resp.Royalty = structure.ProjectRoyalty{
			Data: *dataFromRoyaltyPChan.Data,
		}
	}

	return resp, nil
}
