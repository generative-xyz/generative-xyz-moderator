package usecase

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	"rederinghub.io/utils/redis"
)

type uploadFileChan struct {
	FileURL *string
	Err error
}

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

func (u Usecase) CreateBTCProject(rootSpan opentracing.Span, req structure.CreateBtcProjectReq) (*entity.Projects, error) {
	span, log := u.StartSpan("CreateBTCProject", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	pe := &entity.Projects{}
	err := copier.Copy(pe, req)
	if err != nil {
		log.Error("copier.Copy", err.Error(), err)
		return nil, err
	}

	pow := math.Pow10(8)
	powBig := big.NewFloat(0).SetFloat64(pow)

	mintPrice := big.NewFloat(0)
	mintPrice.SetString(req.MintPrice)

	minMint, ok := big.NewFloat(0).SetString("0.01")
	if !ok {
		err = errors.New("Cannot convert number")
		log.Error("validate.convert.number", err.Error(), err)
		return  nil, err
	}

	if mintPrice.Cmp(minMint) == -1 {
		err = errors.New("mintPrice must be greater than 0.01")
		log.Error("validate.Min.Fee.Fail", err.Error(), err)
		return  nil, err
	}
	mintPrice.Mul(mintPrice, powBig)

	maxID, err := u.Repo.GetMaxBtcProjectID()
	if err != nil {
		log.Error("u.Repo.GetMaxBtcProjectID", err.Error(), err)
		return nil, err
	}
	maxID = maxID + 1
	pe.TokenIDInt =  maxID
	pe.TokenID =  fmt.Sprintf("%d", maxID)
	pe.ContractAddress = os.Getenv("GENERATIVE_PROJECT")
	pe.MintPrice = mintPrice.String()
	pe.IsHidden = true
	pe.Status = true
	pe.IsSynced = true
	nftTokenURI := make(map[string]interface{})
	nftTokenURI["name"] = pe.Name
	nftTokenURI["description"] = pe.Description
	nftTokenURI["image"] = pe.Thumbnail
	nftTokenURI["animation_url"] = ""
	nftTokenURI["attributes"] = []string{}

	creatorAddrr, err := u.Repo.FindUserByWalletAddress(req.CreatorAddrr)
	if err != nil {
		log.Error("u.Repo.FindUserByWalletAddress", err.Error(), err)
		return nil, err
	}

	if creatorAddrr.WalletAddressBTC == "" {
		creatorAddrr.WalletAddressBTC = req.CreatorAddrrBTC
		updated, err := u.Repo.UpdateUserByID(creatorAddrr.UUID, creatorAddrr)
		if err != nil {
			log.Error("u.Repo.UpdateUserByID", err.Error(), err)
			
		}else{
			log.SetData("updated.creatorAddrr", creatorAddrr)
			log.SetData("updated", updated)
		}
	}

	isPubsub := false
	animationURL := ""
	zipLink := req.ZipLink
	if zipLink != nil {
		if *zipLink != "" {
			pe.IsHidden = true
			isPubsub = true
			pe.IsSynced = false
			pe.Status = false
		}
	}else{
		if req.AnimationURL != nil {
			animationURL = *req.AnimationURL
			nftTokenURI["animation_url"] = animationURL
		}
	}

	bytes, err := json.Marshal(nftTokenURI)
	if err != nil {
		log.Error("json.Marshal.nftTokenURI", err.Error(), err)
		return nil, err
	}
	nftToken := helpers.Base64Encode(bytes)
	now := time.Now().UTC()

	pe.NftTokenUri = fmt.Sprintf("data:application/json;base64,%s",nftToken)
	pe.ProcessingImages = []string{}
	pe.MintedImages = nil
	pe.MintedTime = &now
	pe.CreatorProfile = *creatorAddrr
	pe.CreatorAddrrBTC = req.CreatorAddrrBTC
	pe.LimitSupply = 0
	pe.GenNFTAddr = pe.TokenID
	
	err = u.Repo.CreateProject(pe)
	if err != nil {
		log.Error("u.Repo.CreateProject", err.Error(), err)
		return nil, err
	}

	if isPubsub {
		err = u.PubSub.ProducerWithTrace(span, utils.PUBSUB_PROJECT_UNZIP, redis.PubSubPayload{ Data: structure.ProjectUnzipPayload{ProjectID: pe.TokenID, ZipLink:  *zipLink}})
		if err != nil {
			log.Error("u.Repo.CreateProject", err.Error(), err)
			//return nil, err
		}
	}

	log.SetData("pe", pe)
	log.SetData("pe.isPubsub", isPubsub)
	return pe, nil
}

func (u Usecase) UpdateProject(rootSpan opentracing.Span, req structure.UpdateProjectReq) (*entity.Projects, error) {
	span, log := u.StartSpan("UpdateProject", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	p, err := u.Repo.FindProjectBy(req.ContracAddress, req.TokenID)
	if err != nil {
		log.Error("UpdateProject.FindProjectBy", err.Error(), err)
		return nil, err
	}

	if req.Priority != nil {
		priority  := 0
		p.Priority = &priority
	}
	
	updated, err := u.Repo.UpdateProject(p.UUID, p)
	if err != nil {
		log.Error("UpdateProject.UpdateProject", err.Error(), err)
		return nil, err
	}

	log.SetData("updated", updated)
	return p, nil
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

	log.SetData("projects", projects.Total)
	return projects, nil
}

func (u Usecase) GetRandomProject(rootSpan opentracing.Span) (*entity.Projects, error) {
	span, log := u.StartSpan("GetRandomProject", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	caddr := os.Getenv("RANDOM_PR_CONTRACT")
	pID := os.Getenv("RANDOM_PR_PROJECT")

	if caddr != "" && pID != "" {
		return u.GetProjectDetail(span, structure.GetProjectDetailMessageReq{
	ContractAddress: caddr,
	ProjectID:       pID,
		})
	}

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

	log.SetData("projects", projects.Total)
	return projects, nil
}

func (u Usecase) GetProjectDetail(rootSpan opentracing.Span, req structure.GetProjectDetailMessageReq) (*entity.Projects, error) {
	span, log := u.StartSpan("GetProjectDetail", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	// defer func  ()  {
	// 	//alway update project in a separated process
	// 	go func(rootSpan opentracing.Span) {
	// 		span, log := u.StartSpan("GetProjectDetail.GetProjectFromChain", rootSpan)
	// 		defer u.Tracer.FinishSpan(span, log)

	// 		_, err := u.UpdateProjectFromChain(span, req.ContractAddress, req.ProjectID)
	// 		if err != nil {
	// 	log.Error("u.Repo.FindProjectBy", err.Error(), err)
	// 	return
	// 		}

	// 	}(span)	
	// }()

	log.SetTag("ProjectID", req.ProjectID)
	log.SetTag("ContractAddress", req.ContractAddress)

	c, _ := u.Repo.FindProjectWithoutCache(req.ContractAddress, req.ProjectID)
	if (c == nil) || (c != nil && !c.IsSynced) || c.MintedTime == nil {
		// p, err := u.UpdateProjectFromChain(span, req.ContractAddress, req.ProjectID)
		// if err != nil {
		// 	log.Error("u.Repo.FindProjectBy", err.Error(), err)
		// 	return nil, err
		// }
		// return p, nil
		return nil, errors.New("project is not found")
	}
	mintPriceInt, err := strconv.ParseInt(c.MintPrice, 10, 64)
	if err != nil {
		return nil, err
	}
	ethPrice, err := u.convertBTCToETH(span, fmt.Sprintf("%f", float64(mintPriceInt) / 1e8))
	if err != nil {
		return nil, err
	}
	c.MintPriceEth = ethPrice

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

	log.SetData("projects", projects.Total)
	return projects, nil
}

func (u Usecase) GetUpdatedProjectStats(rootSpan opentracing.Span, req structure.GetProjectReq) (*entity.ProjectStat, []entity.TraitStat,error) {
	span, log := u.StartSpan("SyncProjectStats", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	project, err := u.Repo.FindProjectBy(req.ContractAddr, req.TokenID)
	if err != nil {
		return nil, nil, err
	}

	// do not resync
	if project.Stats.LastTimeSynced != nil && project.Stats.LastTimeSynced.Unix() + int64(u.Config.TimeResyncProjectStat) > time.Now().Unix() {
		return &project.Stats, project.TraitsStat,nil
	}

	allTokenFromDb, err := u.Repo.GetAllTokensByProjectID(project.TokenID)
	if err != nil {
		return nil, nil, err
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
		return nil, nil, err
	}

	allOffers, err = u.Repo.GetAllOfferByCollectionContract(project.GenNFTAddr)
	if err != nil {
		log.Error("u.Repo.GetAllOfferByCollectionContract", err.Error(), err)
		return nil, nil, err
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
		// update listing percent
		if !listing.Closed && (time.Now().Unix() < durationTime || durationTime == 0) {
	listingSet[listing.TokenId] = true
		}

		// update floor price
		if listing.Finished {
	if floorPrice == nil {
		floorPrice = price
	} else {
		if floorPrice.Cmp(price) > 0 {
			floorPrice = price
		}
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

		// update floor price
		if offer.Finished {
	if floorPrice == nil {
		floorPrice = price
	} else {
		if floorPrice.Cmp(price) > 0 {
			floorPrice = price
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

	// update trait stats
	traitToCnt := make(map[string]int32)
	traitValueToCnt := make(map[string]map[string]int32)
	for _, token := range allTokenFromDb {
		for _, attribute := range token.ParsedAttributes {
	traitToCnt[attribute.TraitType] += 1
	if traitValueToCnt[attribute.TraitType] == nil {
		traitValueToCnt[attribute.TraitType] = make(map[string]int32)
	}
	traitValueToCnt[attribute.TraitType][fmt.Sprintf("%v", attribute.Value)] += 1
		}
	}

	traitsStat := make([]entity.TraitStat, 0)
	for k, cnt := range traitToCnt {
		traitValueStat := make([]entity.TraitValueStat, 0)
		for value, cntValue := range traitValueToCnt[k] {
	traitValueStat = append(traitValueStat, entity.TraitValueStat{
		Value: value,
		Rarity: int32(cntValue * 100 / cnt),
	})
		}
		traitsStat = append(traitsStat, entity.TraitStat{
	TraitName: k,
	TraitValuesStat: traitValueStat,
		})
	}

	now := time.Now()

	project, err = u.Repo.FindProjectBy(req.ContractAddr, req.TokenID)
	if err != nil {
		return nil, nil, err
	}

	return &entity.ProjectStat{
		LastTimeSynced: &now,
		UniqueOwnerCount: uint32(len(owners)),
		TotalTradingVolumn: totalTradingVolumn.String(),
		FloorPrice: floorPrice.String(),
		BestMakeOfferPrice: bestMakeOfferPrice.String(),
		ListedPercent: listedPercent,
		MintedCount: project.Stats.MintedCount,
	}, traitsStat, nil
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
	contractDetail, err := u.getNftContractDetailInternal(span, client, addr, *projectID)
	if err != nil {
		log.Error("u.getNftContractDetailInternal", err.Error(), err)
		return nil, err
	}
	//log.SetData("contractDetail", contractDetail)
	u.Cache.SetData(contractDataKey, contractDetail)
	return contractDetail, nil
}

// Get from chain with cache
func (u Usecase) getProjectDetailFromChain(rootSpan opentracing.Span, req structure.GetProjectDetailMessageReq) (*structure.ProjectDetail, error) {
	span, log := u.StartSpan("getProjectDetailFromChain", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	contractDataKey := helpers.ProjectDetailKey(req.ContractAddress, req.ProjectID)
	
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
		contractDetail, err := u.getNftContractDetailInternal(span, client, addr, *projectID)
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
func (u Usecase) getNftContractDetailInternal(rootSpan opentracing.Span, client *ethclient.Client, contractAddr common.Address, projectID big.Int) (*structure.ProjectDetail, error) {
	span, log := u.StartSpan("getNftContractDetailInternal", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	log.SetTag("contractAddress", contractAddr)
	log.SetTag("projectID", projectID)
	
	gProject, err := generative_project_contract.NewGenerativeProjectContract(contractAddr, client)
	if err != nil {
		log.Error("generative_project_contract.NewGenerativeProjectContract", err.Error(), err)
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
		log.Error("statusFromChain.Err", statusFromChain.Err.Error(), statusFromChain.Err)
		return nil, statusFromChain.Err
	}

	if tokenFromChain.Err != nil {
		log.Error("tokenFromChain.Err", tokenFromChain.Err.Error(), tokenFromChain.Err)
		return nil, tokenFromChain.Err
	}

	gNftProject, err := generative_nft_contract.NewGenerativeNftContract(detailFromChain.ProjectDetail.GenNFTAddr, client)
	if err != nil {
		log.Error("generative_nft_contract.NewGenerativeNftContract", err.Error(), err)
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

	log.SetData("resp", resp)
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

	log.SetData("resp", resp)
	return resp, nil
}

func (u Usecase) UnzipProjectFile(rootSpan opentracing.Span, zipPayload *structure.ProjectUnzipPayload) (*entity.Projects, error) {
	span, log := u.StartSpan("getNftContractDetailInternal", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	pe, err := u.Repo.FindProjectByTokenID(zipPayload.ProjectID)
	if err != nil {
		log.Error("http.Get", err.Error(), err)
		return nil, err
	}

	nftTokenURI := make(map[string]interface{})
	nftTokenURI["name"] = pe.Name
	nftTokenURI["description"] = pe.Description
	nftTokenURI["image"] = pe.Thumbnail
	nftTokenURI["animation_url"] = ""
	nftTokenURI["attributes"] = []string{}

	images := []string{}
	zipLink := zipPayload.ZipLink
	resp, err := http.Get(zipLink)
	if err != nil {
		log.Error("http.Get", err.Error(), err)
		return nil, err
	}
	defer resp.Body.Close()

	fileName := fmt.Sprintf("%s.zip", helpers.GenerateSlug(pe.Name))

	out, err := os.Create(fileName)
	if err != nil {
		log.Error("os.Create", err.Error(), err)
		return nil, err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Error("os.Create", err.Error(), err)
		return nil, err
	}

	zf, err := zip.OpenReader(fileName)
	if err != nil {
		log.Error("zip.OpenReader", err.Error(), err)
		return nil, err
	}
	
	uploadChan := make(chan uploadFileChan, len(zf.File))
	processed := 1
	processingFiles := []*zip.File{}
	for _, file := range zf.File {
		if strings.Index(file.Name, "__MACOSX") > -1 {
			continue
		}
		
		if strings.Index(file.Name, ".DS_Store") > -1 {
			continue
		}
		
		processingFiles = append(processingFiles, file)
	}
	log.SetData("processingFiles",processingFiles)
	

	for _, file := range processingFiles {
		log.SetData("fileName", file.Name)
		log.SetData("UncompressedSize64", file.UncompressedSize64)
		maxFileSize := uint64(400000)

		log.SetData("maxFileSize", maxFileSize)
		size := file.UncompressedSize64
		if size > maxFileSize {
			err := errors.New(fmt.Sprintf("%s, size: %d Max file size must be less than 400kb", file.Name, size))
			log.Error("maxFileSize",err.Error(), err)
			return nil, err
		}

		fC, err :=  helpers.ReadFile(file)
		if err != nil {
			log.Error("zip.OpenReader", err.Error(), err)
			return nil, err
		}

		if processed % 5 == 0 {
			time.Sleep(5 * time.Second)
		}
		
		go func(rootSpan opentracing.Span, fc []byte, uploadChan chan uploadFileChan) {
			span, log := u.StartSpan("CreateBTCProject.UploadFile", rootSpan)
			defer u.Tracer.FinishSpan(span, log)
			var err error
			var uploadedUrl *string

			defer func  ()  {
		uploadChan <- uploadFileChan{
			FileURL: uploadedUrl,
			Err: err,
		}
			}()

			base64Data := helpers.Base64Encode(fc)
			uploadFileName := fmt.Sprintf("%s/%s", helpers.GenerateSlug(pe.Name), file.Name)
			uploaded, err := u.GCS.UploadBaseToBucket(base64Data, uploadFileName)
			if err != nil {
		log.Error("u.GCS.UploadBaseToBucket", err.Error(), err)
		return
			}

			cdnURL := fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), uploaded.Name)
			uploadedUrl = &cdnURL

		}(span, fC, uploadChan)

		processed ++
	}

	for i := 0 ;i < len(processingFiles) ; i++ {
		dataFromChan := <- uploadChan
		if dataFromChan.Err != nil {
			log.Error("dataFromChan.Err", dataFromChan.Err.Error(), dataFromChan.Err)
			return nil, dataFromChan.Err
		}

		images = append(images, *dataFromChan.FileURL)
		animationURL := *dataFromChan.FileURL
		nftTokenURI["image"] = animationURL
	}

	pe.Images = images	
	pe.IsFullChain = true
	pe.IsHidden = true
	pe.Status = true
	pe.IsSynced = true

	updated, err := u.Repo.UpdateProject(pe.UUID, pe)
	if err != nil {
		log.Error("u.Repo.UpdateProject",err.Error(), err)
		return nil, err
	}

	log.SetData("updated", updated)
	return pe, nil
}