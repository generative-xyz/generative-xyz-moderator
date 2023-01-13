package usecase

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/mongo"
	"rederinghub.io/external/nfts"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/contracts/generative_marketplace_lib"
	"rederinghub.io/utils/helpers"
)

type projectChan struct {
	Data *entity.Projects
	Err error
}

type projectDetailChan struct {
	Data *structure.ProjectDetail
	Err error
}

type projectStatChan struct {
	Data *entity.ProjectStat
	DataTrait []entity.TraitStat
	Err error
}

func (u Usecase) ResolveMarketplaceListTokenEvent(rootSpan opentracing.Span, chainLog  types.Log) error {
	span, log := u.StartSpan("ResolveMarketplaceListTokenEvent", rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	marketplaceContract, err := generative_marketplace_lib.NewGenerativeMarketplaceLib(chainLog.Address, u.Blockchain.GetClient())
	if  err != nil {
		log.Error("cannot init marketplace contract", "", err)
		return err
	}
	event, err := marketplaceContract.ParseListingToken(chainLog)

	if err != nil {
		log.Error("cannot parse list token event", "", err)
		return err
	}

	log.SetTag("OfferingID", strings.ToLower(fmt.Sprintf("%x", event.OfferingId)))
	log.SetData("resolved-listing-event", strings.ToLower(fmt.Sprintf("%x", event.OfferingId)))

	err = u.ListToken(span, event)

	if err != nil {
		log.Error("fail when resolve list token event", "", err)
	}

	return nil
}

func (u Usecase) ResolveMarketplacePurchaseTokenEvent(rootSpan opentracing.Span, chainLog types.Log) error {
	span, log := u.StartSpan("ResolveMarketplacePurchaseTokenEvent", rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	marketplaceContract, err := generative_marketplace_lib.NewGenerativeMarketplaceLib(chainLog.Address, u.Blockchain.GetClient())
	if  err != nil {
		log.Error("cannot init marketplace contract", "", err)
		return err
	}
	event, err := marketplaceContract.ParsePurchaseToken(chainLog)
	if err != nil {
		log.Error("cannot parse purchase token event", "", err)
		return err
	}

	log.SetTag("OfferingID", strings.ToLower(fmt.Sprintf("%x", event.OfferingId)))
	log.SetData("resolved-purchase-event", strings.ToLower(fmt.Sprintf("%x", event.OfferingId)))

	err = u.PurchaseToken(span, event)

	if err != nil {
		log.Error("fail when resolve purchase token event", "", err)
	}

	return nil
}

func (u Usecase) ResolveMarketplaceMakeOffer(rootSpan opentracing.Span, chainLog types.Log) error {
	span, log := u.StartSpan("ResolveMarketplaceMakeOffer", rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	marketplaceContract, err := generative_marketplace_lib.NewGenerativeMarketplaceLib(chainLog.Address, u.Blockchain.GetClient())
	if  err != nil {
		log.Error("cannot init marketplace contract", "", err)
		return err
	}
	event, err := marketplaceContract.ParseMakeOffer(chainLog)
	if err != nil {
		log.Error("cannot parse make offer event", "", err)
		return err
	}

	log.SetTag("OfferingID", strings.ToLower(fmt.Sprintf("%x", event.OfferingId)))
	log.SetData("resolved-make-offer-event", strings.ToLower(fmt.Sprintf("%x", event.OfferingId)))

	err = u.MakeOffer(span, event)

	if err != nil {
		log.Error("fail when resolve make offer event", "", err)
	}

	return nil
}

func (u Usecase) ResolveMarketplaceAcceptOfferEvent(rootSpan opentracing.Span, chainLog types.Log) error {
	span, log := u.StartSpan("ResolveMarketplaceAcceptOfferEvent", rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	marketplaceContract, err := generative_marketplace_lib.NewGenerativeMarketplaceLib(chainLog.Address, u.Blockchain.GetClient())
	if  err != nil {
		log.Error("cannot init marketplace contract", "", err)
		return err
	}
	event, err := marketplaceContract.ParseAcceptMakeOffer(chainLog)
	if err != nil {
		log.Error("cannot parse accept offer event", "", err)
		return err
	}

	log.SetData("resolved-purchase-event", strings.ToLower(fmt.Sprintf("%x", event.OfferingId)))
	log.SetTag("OfferingID", strings.ToLower(fmt.Sprintf("%x", event.OfferingId)))
	
	err = u.AcceptMakeOffer(span, event)

	if err != nil {
		log.Error("fail when resolve accept offer event", "", err)
	}

	return nil
}

func (u Usecase) ResolveMarketplaceCancelListing(rootSpan opentracing.Span, chainLog types.Log) error {
	span, log := u.StartSpan("ResolveMarketplaceCancelListing", rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	marketplaceContract, err := generative_marketplace_lib.NewGenerativeMarketplaceLib(chainLog.Address, u.Blockchain.GetClient())
	if  err != nil {
		log.Error("cannot init marketplace contract", "", err)
		return err
	}
	event, err := marketplaceContract.ParseCancelListing(chainLog)
	if err != nil {
		log.Error("cannot parse cancel listing event", "", err)
		return err
	}

	log.SetData("resolved-cancel-listing-event", strings.ToLower(fmt.Sprintf("%x", event.OfferingId)))
	log.SetTag("OfferingID", strings.ToLower(fmt.Sprintf("%x", event.OfferingId)))
	err = u.CancelListing(span, event)

	if err != nil {
		log.Error("fail when resolve cancel listing event", "", err)
	}

	return nil
}

func (u Usecase) ResolveMarketplaceCancelOffer(rootSpan opentracing.Span, chainLog types.Log) error {
	span, log := u.StartSpan("ResolveMarketplaceCancelOffer", rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	marketplaceContract, err := generative_marketplace_lib.NewGenerativeMarketplaceLib(chainLog.Address, u.Blockchain.GetClient())
	if  err != nil {
		log.Error("cannot init marketplace contract", "", err)
		return err
	}
	event, err := marketplaceContract.ParseCancelMakeOffer(chainLog)
	if err != nil {
		log.Error("cannot parse cancel offer event", "", err)
		return err
	}

	log.SetData("resolved-cancel-offer-event", strings.ToLower(fmt.Sprintf("%x", event.OfferingId)))
	log.SetTag("OfferingID", strings.ToLower(fmt.Sprintf("%x", event.OfferingId)))
	err = u.CancelOffer(span, event)

	if err != nil {
		log.Error("fail when resolve cancel offer event", "", err)
	}

	return nil
}

func (u Usecase) UpdateProjectWithListener(chainLog types.Log) {
	txnHash := chainLog.TxHash.String()
	span, log := u.StartSpanWithoutRoot("UpdateProjectWithListener.GetProjectDetail")
	defer u.Tracer.FinishSpan(span, log )
	log.SetTag("transaction_hash", txnHash)
	log.SetData("chainLog", chainLog)
	topics := chainLog.Topics

	tokenIDStr :=  helpers.HexaNumberToInteger(topics[3].String())
	tokenID, _ := strconv.Atoi(tokenIDStr)
	tokenIDStr = fmt.Sprintf("%d",tokenID)
	contractAddr := strings.ToLower(chainLog.Address.String()) 

	u.UpdateProjectFromChain(span, contractAddr, tokenIDStr)
}

func (u Usecase) UpdateProjectFromChain(rootSpan opentracing.Span, contractAddr string, tokenIDStr string) (*entity.Projects, error) {
	span, log := u.StartSpan("UpdateProjectWithListener.GetProjectDetail", rootSpan)
	defer u.Tracer.FinishSpan(span, log )

	pChan := make(chan projectChan, 1)
	pDChan := make(chan projectDetailChan, 1)
	pSChan := make(chan projectStatChan, 1)
	
	log.SetData("contractAddr", contractAddr)
	log.SetData("tokenIDStr", tokenIDStr)
	log.SetTag("contractAddr", contractAddr)
	log.SetTag("tokenIDStr", tokenIDStr)
	tokenIDInt, err := strconv.Atoi(tokenIDStr)
	if err != nil {
		log.Error("strconv.Atoi.tokenIDStr", err.Error(), err)
		return nil, err
	}

	go func(rootSpan opentracing.Span, pChan chan projectChan, contractAddr string, tokenIDStr string) {
		span, log := u.StartSpan("GetProjectDetail.Project", rootSpan)
		defer u.Tracer.FinishSpan(span, log )

		log.SetTag("contractAddr",contractAddr)
		log.SetTag("tokenIDStr",tokenIDStr)

		project := &entity.Projects{}
		var err error

		defer func  ()  {
			pChan <- projectChan {
				Data: project,
				Err:  err,
			}
		}()

		project, err = u.Repo.FindProjectBy(contractAddr, tokenIDStr)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				project = &entity.Projects{}
				project.ContractAddress = contractAddr 
				project.TokenID = tokenIDStr 

				err = u.Repo.CreateProject(project)
				if err != nil {
					log.Error("u.Repo.CreateProject", err.Error(), err)
					return
				}

		}else{
			log.Error("u.Repo.FindProjectBy", err.Error(), err)
			return
		}
	}

	}(span, pChan, contractAddr, tokenIDStr)

	go func(rootSpan opentracing.Span, pDChan chan projectDetailChan, contractAddr string, tokenIDStr string) {
		span, log := u.StartSpan("GetProjectDetail.Project", rootSpan)
		defer u.Tracer.FinishSpan(span, log )
		
		projectDetail := &structure.ProjectDetail{}
		var err error

		defer func  ()  {
			pDChan <- projectDetailChan {
				Data: projectDetail,
				Err:  err,
			}
		}()

		projectDetail, err = u.getProjectDetailFromChainWithoutCache(span, structure.GetProjectDetailMessageReq{
			ContractAddress:  contractAddr,
			ProjectID:  tokenIDStr,
		})
		if err != nil {
			log.Error(" u.GetProjectDetail", err.Error(), err)
			return
		}
	

	}(span, pDChan, contractAddr, tokenIDStr)

	go func(rootSpan opentracing.Span, pDChan chan projectStatChan, contractAddr string, tokenIDStr string) {
		span, log := u.StartSpan("GetProjectDetail.ProjectStat", rootSpan)
		defer u.Tracer.FinishSpan(span, log )
		projectStat := &entity.ProjectStat{}
		traitStat := make([]entity.TraitStat, 0)
		var err error

		defer func  ()  {
			pDChan <- projectStatChan {
				Data: projectStat,
				DataTrait: traitStat,
				Err:  err,
			}
		}()

		projectStat, traitStat, err = u.GetUpdatedProjectStats(span, structure.GetProjectReq{
			ContractAddr: contractAddr,
			TokenID: tokenIDStr,
		})
		if err != nil {
			log.Error(" u.SyncProjectStats", err.Error(), err)
			return
		}
	}(span, pSChan, contractAddr, tokenIDStr)

	projectFChan := <- pChan
	projectDetailFChan := <- pDChan
	projectStatFChan := <- pSChan

	err = projectFChan.Err 
	if err != nil {
		log.Error("projectFChan.Err ", err.Error(), err)
		return nil, err
	}

	project := projectFChan.Data
	log.SetData("project", project)
	
	//get creator profile
	getProfile := func(profileChan chan structure.ProfileChan, address string) {
		var user *entity.Users
		var err error

		defer func() {
			profileChan <- structure.ProfileChan{
				Data: user,
				Err: err,
			}
		}()

		user, err = u.GetUserProfileByWalletAddress(span, strings.ToLower(address))
		if err != nil {
			return
		}
	}

	profileChan := make(chan structure.ProfileChan, 1)
	go getProfile(profileChan, project.CreatorAddrr)

	usrFromChan := <- profileChan

	project.MintingInfo = entity.ProjectMintingInfo{
		Index: 0,
		IndexReverse: 0,
	}

	err = projectDetailFChan.Err
	if err != nil {
		log.Error("projectDetailFChan.Err ", err.Error(), err)
		//return nil, err
	}else{

		
		projectDetail := projectDetailFChan.Data
		//log.SetData("projectDetail", projectDetail)
		project.IsSynced = true
		project.Name = projectDetail.ProjectDetail.Name
		project.CreatorName = projectDetail.ProjectDetail.Creator
		project.CreatorAddrr = strings.ToLower(projectDetail.ProjectDetail.CreatorAddr.String())
		project.Description = projectDetail.ProjectDetail.Desc
		project.Scripts= projectDetail.ProjectDetail.Scripts
		project.ThirdPartyScripts= projectDetail.ProjectDetail.ScriptType
		project.Styles= projectDetail.ProjectDetail.Styles
		project.GenNFTAddr= strings.ToLower( projectDetail.ProjectDetail.GenNFTAddr.String())
		project.MintPrice = projectDetail.ProjectDetail.MintPrice.String()
		project.MaxSupply = projectDetail.ProjectDetail.MaxSupply.Int64()
		project.LimitSupply = projectDetail.ProjectDetail.Limit.Int64()
		project.MintTokenAddress = strings.ToLower(string(projectDetail.ProjectDetail.MintPriceAddr.String()))
		project.License = projectDetail.ProjectDetail.License
		project.Status = projectDetail.Status
		project.SocialWeb = projectDetail.ProjectDetail.Social.Web
		project.SocialTwitter = projectDetail.ProjectDetail.Social.Twitter
		project.SocialDiscord = projectDetail.ProjectDetail.Social.Discord
		project.SocialMedium = projectDetail.ProjectDetail.Social.Medium
		project.SocialInstagram = projectDetail.ProjectDetail.Social.Instagram
		project.Thumbnail = projectDetail.ProjectDetail.Image
		project.NftTokenUri = projectDetail.NftTokenUri
		project.Royalty = int(projectDetail.Royalty.Data.Int64())
		project.CompleteTime = projectDetail.ProjectDetail.CompleteTime.Int64()
		for _, reserve := range projectDetail.ProjectDetail.Reserves {
			project.Reservers = append(project.Reservers, strings.ToLower(reserve.String()) )
		}

		if projectDetail.NftProjectDetail.Index != nil && projectDetail.NftProjectDetail.IndexReserve != nil {
			project.MintingInfo = entity.ProjectMintingInfo{
				Index: projectDetail.NftProjectDetail.Index.Int64(),
				IndexReverse: projectDetail.NftProjectDetail.IndexReserve.Int64(),
			}
		}

		if project.Priority ==  nil {
			*project.Priority = 0
		}
	}

	// get minted time 
	if project.BlockNumberMinted == nil || project.MintedTime == nil {
		mintedTimeChan := make (chan structure.NftMintedTimeChan, 1)
		go func(mintedTimeChan chan structure.NftMintedTimeChan) {
			var mintedTime *structure.NftMintedTime
			var err error
			defer func() {
				mintedTimeChan <- structure.NftMintedTimeChan{
					NftMintedTime: mintedTime,
					Err: err,
				}
			}()
			span, _ := u.StartSpanWithoutRoot("getNftContractDetail.GetNftMintedTime")
			mintedTime, err = u.GetNftMintedTime(span, structure.GetNftMintedTimeReq{
				ContractAddress: project.ContractAddress,
				TokenID: project.TokenID,
			})
		}(mintedTimeChan)
		mintedTimeFChan := <-mintedTimeChan
		if mintedTimeFChan.Err != nil {
			log.Error("mintedTimeFChan.Err ", mintedTimeFChan.Err.Error(), mintedTimeFChan.Err)
		} else {
			project.BlockNumberMinted = mintedTimeFChan.NftMintedTime.BlockNumberMinted
			project.MintedTime = mintedTimeFChan.NftMintedTime.MintedTime
		}
	}

	
	project.TokenIDInt = int64(tokenIDInt)

	if usrFromChan.Err != nil {
		log.Error("usrFromChan.Err", usrFromChan.Err.Error(), usrFromChan.Err)
	}else{
		project.CreatorProfile = *usrFromChan.Data
	}

	if projectStatFChan.Err != nil {
		log.Error("projectStatFChan.Err", projectStatFChan.Err.Error(), projectStatFChan.Err)
	} else {
		project.Stats = *projectStatFChan.Data
		project.TraitsStat = projectStatFChan.DataTrait
	}

	log.SetData("project",project)
	updated, err := u.Repo.UpdateProject(project.UUID, project)
	if err != nil {
		log.Error(" u.UpdateProject", err.Error(), err)
		return nil, err
	}
	log.SetData("projectUUID", project.UUID)
	log.SetData("updated",updated)
	return  project, nil
}

func (u Usecase) GetProjectsFromChain(rootSpan opentracing.Span) error {
	span, log := u.StartSpan("Usecase.GetProjectsFromChain", rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	
	contractAddress := os.Getenv("GENERATIVE_PROJECT")
	mProjects, err := u.MoralisNft.GetNftByContract(contractAddress, nfts.MoralisFilter{})
	if err != nil {
		log.Error("u.MoralisNft.GetNftByContract", err.Error(), err)
		return err
	}

	log.SetData("contractAddress", contractAddress)
	log.SetTag("contractAddress", contractAddress)
	for _, mProject := range mProjects.Result {
		_, err := u.UpdateProjectFromChain(span, contractAddress, mProject.TokenID)
		if err != nil {
			log.Error("u.Repo.FindProjectBy", err.Error(), err)
			return err
		}
		//resp = append(resp, *p)
		//log.SetData("p", *p)
	}

	return nil
}
