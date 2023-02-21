package usecase

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
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

func (u Usecase) ResolveMarketplaceListTokenEvent( chainLog  types.Log) error {
	marketplaceContract, err := generative_marketplace_lib.NewGenerativeMarketplaceLib(chainLog.Address, u.Blockchain.GetClient())
	if  err != nil {
		u.Logger.Error("cannot init marketplace contract", "", err)
		return err
	}
	event, err := marketplaceContract.ParseListingToken(chainLog)
	blocknumber := chainLog.BlockNumber

	if err != nil {
		u.Logger.Error("cannot parse list token event", "", err)
		return err
	}

	
	u.Logger.Info("resolved-listing-event", strings.ToLower(fmt.Sprintf("%x", event.OfferingId)))

	err = u.ListToken(event, blocknumber)

	if err != nil {
		u.Logger.Error("fail when resolve list token event", "", err)
	}

	return nil
}

func (u Usecase) ResolveMarketplacePurchaseTokenEvent( chainLog types.Log) error {
	marketplaceContract, err := generative_marketplace_lib.NewGenerativeMarketplaceLib(chainLog.Address, u.Blockchain.GetClient())
	if  err != nil {
		u.Logger.Error("cannot init marketplace contract", "", err)
		return err
	}
	event, err := marketplaceContract.ParsePurchaseToken(chainLog)
	if err != nil {
		u.Logger.Error("cannot parse purchase token event", "", err)
		return err
	}

	
	u.Logger.Info("resolved-purchase-event", strings.ToLower(fmt.Sprintf("%x", event.OfferingId)))

	err = u.PurchaseToken(event)

	if err != nil {
		u.Logger.Error("fail when resolve purchase token event", "", err)
	}

	return nil
}

func (u Usecase) ResolveMarketplaceMakeOffer( chainLog types.Log) error {
	marketplaceContract, err := generative_marketplace_lib.NewGenerativeMarketplaceLib(chainLog.Address, u.Blockchain.GetClient())
	if  err != nil {
		u.Logger.Error("cannot init marketplace contract", "", err)
		return err
	}
	event, err := marketplaceContract.ParseMakeOffer(chainLog)
	blocknumber := chainLog.BlockNumber

	if err != nil {
		u.Logger.Error("cannot parse make offer event", "", err)
		return err
	}

	
	u.Logger.Info("resolved-make-offer-event", strings.ToLower(fmt.Sprintf("%x", event.OfferingId)))

	err = u.MakeOffer(event, blocknumber)

	if err != nil {
		u.Logger.Error("fail when resolve make offer event", "", err)
	}

	return nil
}

func (u Usecase) ResolveMarketplaceAcceptOfferEvent( chainLog types.Log) error {
	marketplaceContract, err := generative_marketplace_lib.NewGenerativeMarketplaceLib(chainLog.Address, u.Blockchain.GetClient())
	if  err != nil {
		u.Logger.Error("cannot init marketplace contract", "", err)
		return err
	}
	event, err := marketplaceContract.ParseAcceptMakeOffer(chainLog)
	if err != nil {
		u.Logger.Error("cannot parse accept offer event", "", err)
		return err
	}

	u.Logger.Info("resolved-purchase-event", strings.ToLower(fmt.Sprintf("%x", event.OfferingId)))
	
err = u.AcceptMakeOffer(event)

	if err != nil {
		u.Logger.Error("fail when resolve accept offer event", "", err)
	}

	return nil
}

func (u Usecase) ResolveMarketplaceCancelListing( chainLog types.Log) error {
	marketplaceContract, err := generative_marketplace_lib.NewGenerativeMarketplaceLib(chainLog.Address, u.Blockchain.GetClient())
	if  err != nil {
		u.Logger.Error("cannot init marketplace contract", "", err)
		return err
	}
	event, err := marketplaceContract.ParseCancelListing(chainLog)
	if err != nil {
		u.Logger.Error("cannot parse cancel listing event", "", err)
		return err
	}

	u.Logger.Info("resolved-cancel-listing-event", strings.ToLower(fmt.Sprintf("%x", event.OfferingId)))
	
	err = u.CancelListing(event)

	if err != nil {
		u.Logger.Error("fail when resolve cancel listing event", "", err)
	}

	return nil
}

func (u Usecase) ResolveMarketplaceCancelOffer( chainLog types.Log) error {
	marketplaceContract, err := generative_marketplace_lib.NewGenerativeMarketplaceLib(chainLog.Address, u.Blockchain.GetClient())
	if  err != nil {
		u.Logger.Error("cannot init marketplace contract", "", err)
		return err
	}
	event, err := marketplaceContract.ParseCancelMakeOffer(chainLog)
	if err != nil {
		u.Logger.Error("cannot parse cancel offer event", "", err)
		return err
	}

	u.Logger.Info("resolved-cancel-offer-event", strings.ToLower(fmt.Sprintf("%x", event.OfferingId)))
	
	err = u.CancelOffer(event)

	if err != nil {
		u.Logger.Error("fail when resolve cancel offer event", "", err)
	}

	return nil
}

func (u Usecase) UpdateProjectWithListener(chainLog types.Log) {
	txnHash := chainLog.TxHash.String()
	_ =txnHash
	
	u.Logger.Info("chainLog", chainLog)
	topics := chainLog.Topics

	tokenIDStr :=  helpers.HexaNumberToInteger(topics[3].String())
	tokenID, _ := strconv.Atoi(tokenIDStr)
	tokenIDStr = fmt.Sprintf("%d",tokenID)
	contractAddr := strings.ToLower(chainLog.Address.String()) 

	u.UpdateProjectFromChain(contractAddr, tokenIDStr)
}

func (u Usecase) UpdateProjectFromChain( contractAddr string, tokenIDStr string) (*entity.Projects, error) {

	pChan := make(chan projectChan, 1)
	pDChan := make(chan projectDetailChan, 1)
	pSChan := make(chan projectStatChan, 1)
	u.Logger.Info("contractAddr", contractAddr)
	u.Logger.Info("tokenIDStr", tokenIDStr)
	
	
	tokenIDInt, err := strconv.Atoi(tokenIDStr)
	if err != nil {
		u.Logger.Error("UpdateProjectFromChain.Atoi.tokenIDStr", err.Error(), err)
		return nil, err
	}

	go func( pChan chan projectChan, contractAddr string, tokenIDStr string) {

		
		

		project := &entity.Projects{}
		var err error

		defer func  ()  {
			pChan <- projectChan {
				Data: project,
				Err:  err,
			}
		}()

		project, err = u.Repo.FindProjectWithoutCache(contractAddr, tokenIDStr)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				project = &entity.Projects{}
				project.ContractAddress = contractAddr 
				project.TokenID = tokenIDStr 

				err = u.Repo.CreateProject(project)
				if err != nil {
					u.Logger.Error("UpdateProjectFromChain.CreateProject", err.Error(), err)
					return
				}

		}else{
			u.Logger.Error("UpdateProjectFromChain.FindProjectBy", err.Error(), err)
			return
		}
	}

	}(pChan, contractAddr, tokenIDStr)

	go func( pDChan chan projectDetailChan, contractAddr string, tokenIDStr string) {
		projectDetail := &structure.ProjectDetail{}
		var err error

		defer func  ()  {
			pDChan <- projectDetailChan {
				Data: projectDetail,
				Err:  err,
			}
		}()

		projectDetail, err = u.getProjectDetailFromChainWithoutCache(structure.GetProjectDetailMessageReq{
			ContractAddress:  contractAddr,
			ProjectID:  tokenIDStr,
		})
		if err != nil {
			u.Logger.Error("UpdateProjectFromChain.getProjectDetailFromChainWithoutCache.GetProjectDetail", err.Error(), err)
			return
		}
	

	}(pDChan, contractAddr, tokenIDStr)

	go func( pDChan chan projectStatChan, contractAddr string, tokenIDStr string) {
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

		projectStat, traitStat, err = u.GetUpdatedProjectStats(structure.GetProjectReq{
			ContractAddr: contractAddr,
			TokenID: tokenIDStr,
		})
		if err != nil {
			u.Logger.Error("UpdateProjectFromChain.GetUpdatedProjectStats.error", err.Error(), err)
			return
		}
	}(pSChan, contractAddr, tokenIDStr)

	projectFChan := <- pChan
	projectDetailFChan := <- pDChan
	projectStatFChan := <- pSChan

	err = projectFChan.Err 
	if err != nil {
		u.Logger.Error("projectFChan.Err ", err.Error(), err)
		return nil, err
	}

	project := projectFChan.Data
	u.Logger.Info("project", project)
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

		user, err = u.GetUserProfileByWalletAddress(strings.ToLower(address))
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
		u.Logger.Error("projectDetailFChan.Err ", err.Error(), err)
		//return nil, err
	}else{

		projectDetail := projectDetailFChan.Data
		//u.Logger.Info("projectDetail", projectDetail)
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
			priority := 0
			project.Priority =  &priority
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
	
			mintedTime, err = u.GetNftMintedTime(structure.GetNftMintedTimeReq{
				ContractAddress: project.ContractAddress,
				TokenID: project.TokenID,
			})
		}(mintedTimeChan)
		mintedTimeFChan := <-mintedTimeChan
		if mintedTimeFChan.Err != nil {
			u.Logger.Error("mintedTimeFChan.Err ", mintedTimeFChan.Err.Error(), mintedTimeFChan.Err)
		} else {
			project.BlockNumberMinted = mintedTimeFChan.NftMintedTime.BlockNumberMinted
			project.MintedTime = mintedTimeFChan.NftMintedTime.MintedTime
		}
	}

project.TokenIDInt = int64(tokenIDInt)

	if usrFromChan.Err != nil {
		u.Logger.Error("usrFromChan.Err", usrFromChan.Err.Error(), usrFromChan.Err)
	}else{
		project.CreatorProfile = *usrFromChan.Data
	}

	if projectStatFChan.Err != nil {
		u.Logger.Error("projectStatFChan.Err", projectStatFChan.Err.Error(), projectStatFChan.Err)
	} else {
		project.Stats = *projectStatFChan.Data
		project.TraitsStat = projectStatFChan.DataTrait
	}

	u.Logger.Info("project",project)
	updated, err := u.Repo.UpdateProject(project.UUID, project)
	if err != nil {
		u.Logger.Error(" u.UpdateProject", err.Error(), err)
		return nil, err
	}
	u.Logger.Info("projectUUID", project.UUID)
	u.Logger.Info("updated",updated)
	return  project, nil
}

func (u Usecase) GetProjectsFromChain() error {
contractAddress := os.Getenv("GENERATIVE_PROJECT")
	mProjects, err := u.MoralisNft.GetNftByContract(contractAddress, nfts.MoralisFilter{})
	if err != nil {
		u.Logger.Error("u.MoralisNft.GetNftByContract", err.Error(), err)
		return err
	}

	u.Logger.Info("contractAddress", contractAddress)
	
	for _, mProject := range mProjects.Result {
		_, err := u.UpdateProjectFromChain(contractAddress, mProject.TokenID)
		if err != nil {
			u.Logger.Error("u.Repo.FindProjectBy", err.Error(), err)
			return err
		}
		//resp = append(resp, *p)
		//u.Logger.Info("p", *p)
	}

	return nil
}
