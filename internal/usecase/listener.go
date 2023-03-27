package usecase

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"rederinghub.io/external/nfts"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/contracts/generative_marketplace_lib"
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/logger"
)

type projectChan struct {
	Data *entity.Projects
	Err  error
}

type projectDetailChan struct {
	Data *structure.ProjectDetail
	Err  error
}

type projectStatChan struct {
	Data      *entity.ProjectStat
	DataTrait []entity.TraitStat
	Err       error
}

func (u Usecase) ResolveMarketplaceListTokenEvent(chainLog types.Log) error {
	marketplaceContract, err := generative_marketplace_lib.NewGenerativeMarketplaceLib(chainLog.Address, u.Blockchain.GetClient())
	if err != nil {
		logger.AtLog.Logger.Error("cannot init marketplace contract", zap.Error(err))
		return err
	}
	event, err := marketplaceContract.ParseListingToken(chainLog)
	blocknumber := chainLog.BlockNumber

	if err != nil {
		logger.AtLog.Logger.Error("cannot parse list token event", zap.Error(err))
		return err
	}

	err = u.ListToken(event, blocknumber)

	if err != nil {
		logger.AtLog.Logger.Error("fail when resolve list token event", zap.Error(err))
	}

	return nil
}

func (u Usecase) ResolveMarketplacePurchaseTokenEvent(chainLog types.Log) error {
	marketplaceContract, err := generative_marketplace_lib.NewGenerativeMarketplaceLib(chainLog.Address, u.Blockchain.GetClient())
	if err != nil {
		logger.AtLog.Logger.Error("cannot init marketplace contract", zap.Error(err))
		return err
	}
	event, err := marketplaceContract.ParsePurchaseToken(chainLog)
	if err != nil {
		logger.AtLog.Logger.Error("cannot parse purchase token event", zap.Error(err))
		return err
	}

	err = u.PurchaseToken(event)

	if err != nil {
		logger.AtLog.Logger.Error("fail when resolve purchase token event", zap.Error(err))
	}

	return nil
}

func (u Usecase) ResolveMarketplaceMakeOffer(chainLog types.Log) error {
	marketplaceContract, err := generative_marketplace_lib.NewGenerativeMarketplaceLib(chainLog.Address, u.Blockchain.GetClient())
	if err != nil {
		logger.AtLog.Logger.Error("cannot init marketplace contract", zap.Error(err))
		return err
	}
	event, err := marketplaceContract.ParseMakeOffer(chainLog)
	blocknumber := chainLog.BlockNumber

	if err != nil {
		logger.AtLog.Logger.Error("cannot parse make offer event", zap.Error(err))
		return err
	}

	err = u.MakeOffer(event, blocknumber)

	if err != nil {
		logger.AtLog.Logger.Error("fail when resolve make offer event", zap.Error(err))
	}

	return nil
}

func (u Usecase) ResolveMarketplaceAcceptOfferEvent(chainLog types.Log) error {
	marketplaceContract, err := generative_marketplace_lib.NewGenerativeMarketplaceLib(chainLog.Address, u.Blockchain.GetClient())
	if err != nil {
		logger.AtLog.Logger.Error("cannot init marketplace contract", zap.Error(err))
		return err
	}
	event, err := marketplaceContract.ParseAcceptMakeOffer(chainLog)
	if err != nil {
		logger.AtLog.Logger.Error("cannot parse accept offer event", zap.Error(err))
		return err
	}

	err = u.AcceptMakeOffer(event)

	if err != nil {
		logger.AtLog.Logger.Error("fail when resolve accept offer event", zap.Error(err))
	}

	return nil
}

func (u Usecase) ResolveMarketplaceCancelListing(chainLog types.Log) error {
	marketplaceContract, err := generative_marketplace_lib.NewGenerativeMarketplaceLib(chainLog.Address, u.Blockchain.GetClient())
	if err != nil {
		logger.AtLog.Logger.Error("cannot init marketplace contract", zap.Error(err))
		return err
	}
	event, err := marketplaceContract.ParseCancelListing(chainLog)
	if err != nil {
		logger.AtLog.Logger.Error("cannot parse cancel listing event", zap.Error(err))
		return err
	}

	err = u.CancelListing(event)

	if err != nil {
		logger.AtLog.Logger.Error("fail when resolve cancel listing event", zap.Error(err))
	}

	return nil
}

func (u Usecase) ResolveMarketplaceCancelOffer(chainLog types.Log) error {
	marketplaceContract, err := generative_marketplace_lib.NewGenerativeMarketplaceLib(chainLog.Address, u.Blockchain.GetClient())
	if err != nil {
		logger.AtLog.Logger.Error("cannot init marketplace contract", zap.Error(err))
		return err
	}
	event, err := marketplaceContract.ParseCancelMakeOffer(chainLog)
	if err != nil {
		logger.AtLog.Logger.Error("cannot parse cancel offer event", zap.Error(err))
		return err
	}

	err = u.CancelOffer(event)

	if err != nil {
		logger.AtLog.Logger.Error("fail when resolve cancel offer event", zap.Error(err))
	}

	return nil
}

func (u Usecase) UpdateProjectWithListener(chainLog types.Log) {
	txnHash := chainLog.TxHash.String()
	_ = txnHash

	logger.AtLog.Logger.Info("chainLog", zap.Any("chainLog", chainLog))
	topics := chainLog.Topics

	tokenIDStr := helpers.HexaNumberToInteger(topics[3].String())
	tokenID, _ := strconv.Atoi(tokenIDStr)
	tokenIDStr = fmt.Sprintf("%d", tokenID)
	contractAddr := strings.ToLower(chainLog.Address.String())

	u.UpdateProjectFromChain(contractAddr, tokenIDStr)
}

func (u Usecase) UpdateProjectFromChain(contractAddr string, tokenIDStr string) (*entity.Projects, error) {

	pChan := make(chan projectChan, 1)
	pDChan := make(chan projectDetailChan, 1)
	pSChan := make(chan projectStatChan, 1)
	logger.AtLog.Logger.Info("contractAddr", zap.Any("contractAddr", contractAddr))
	logger.AtLog.Logger.Info("tokenIDStr", zap.Any("tokenIDStr", tokenIDStr))

	tokenIDInt, err := strconv.Atoi(tokenIDStr)
	if err != nil {
		logger.AtLog.Logger.Error("UpdateProjectFromChain.Atoi.tokenIDStr", zap.Error(err))
		return nil, err
	}

	go func(pChan chan projectChan, contractAddr string, tokenIDStr string) {

		project := &entity.Projects{}
		var err error

		defer func() {
			pChan <- projectChan{
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
					logger.AtLog.Logger.Error("UpdateProjectFromChain.CreateProject", zap.Error(err))
					return
				}

			} else {
				logger.AtLog.Logger.Error("UpdateProjectFromChain.FindProjectBy", zap.Error(err))
				return
			}
		}

	}(pChan, contractAddr, tokenIDStr)

	go func(pDChan chan projectDetailChan, contractAddr string, tokenIDStr string) {
		projectDetail := &structure.ProjectDetail{}
		var err error

		defer func() {
			pDChan <- projectDetailChan{
				Data: projectDetail,
				Err:  err,
			}
		}()

		projectDetail, err = u.getProjectDetailFromChainWithoutCache(structure.GetProjectDetailMessageReq{
			ContractAddress: contractAddr,
			ProjectID:       tokenIDStr,
		})
		if err != nil {
			logger.AtLog.Logger.Error("UpdateProjectFromChain.getProjectDetailFromChainWithoutCache.GetProjectDetail", zap.Error(err))
			return
		}

	}(pDChan, contractAddr, tokenIDStr)

	go func(pDChan chan projectStatChan, contractAddr string, tokenIDStr string) {
		projectStat := &entity.ProjectStat{}
		traitStat := make([]entity.TraitStat, 0)
		var err error

		defer func() {
			pDChan <- projectStatChan{
				Data:      projectStat,
				DataTrait: traitStat,
				Err:       err,
			}
		}()

		projectStat, traitStat, err = u.GetUpdatedProjectStats(structure.GetProjectReq{
			ContractAddr: contractAddr,
			TokenID:      tokenIDStr,
		})
		if err != nil {
			logger.AtLog.Logger.Error("UpdateProjectFromChain.GetUpdatedProjectStats.error", zap.Error(err))
			return
		}
	}(pSChan, contractAddr, tokenIDStr)

	projectFChan := <-pChan
	projectDetailFChan := <-pDChan
	projectStatFChan := <-pSChan

	err = projectFChan.Err
	if err != nil {
		logger.AtLog.Logger.Error("projectFChan.Err ", zap.Error(err))
		return nil, err
	}

	project := projectFChan.Data
	logger.AtLog.Logger.Info("project", zap.Any("project", project))
	//get creator profile
	getProfile := func(profileChan chan structure.ProfileChan, address string) {
		var user *entity.Users
		var err error

		defer func() {
			profileChan <- structure.ProfileChan{
				Data: user,
				Err:  err,
			}
		}()

		user, err = u.GetUserProfileByWalletAddress(strings.ToLower(address))
		if err != nil {
			return
		}
	}

	profileChan := make(chan structure.ProfileChan, 1)
	go getProfile(profileChan, project.CreatorAddrr)

	usrFromChan := <-profileChan

	project.MintingInfo = entity.ProjectMintingInfo{
		Index:        0,
		IndexReverse: 0,
	}

	err = projectDetailFChan.Err
	if err != nil {
		logger.AtLog.Logger.Error("projectDetailFChan.Err ", zap.Error(err))
		//return nil, err
	} else {

		projectDetail := projectDetailFChan.Data
		//logger.AtLog.Logger.Info("projectDetail", zap.Any("projectDetail", projectDetail))
		project.IsSynced = true
		project.Name = projectDetail.ProjectDetail.Name
		project.CreatorName = projectDetail.ProjectDetail.Creator
		project.CreatorAddrr = strings.ToLower(projectDetail.ProjectDetail.CreatorAddr.String())
		project.Description = projectDetail.ProjectDetail.Desc
		project.Scripts = projectDetail.ProjectDetail.Scripts
		project.ThirdPartyScripts = projectDetail.ProjectDetail.ScriptType
		project.Styles = projectDetail.ProjectDetail.Styles
		project.GenNFTAddr = strings.ToLower(projectDetail.ProjectDetail.GenNFTAddr.String())
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
			project.Reservers = append(project.Reservers, strings.ToLower(reserve.String()))
		}

		if projectDetail.NftProjectDetail.Index != nil && projectDetail.NftProjectDetail.IndexReserve != nil {
			project.MintingInfo = entity.ProjectMintingInfo{
				Index:        projectDetail.NftProjectDetail.Index.Int64(),
				IndexReverse: projectDetail.NftProjectDetail.IndexReserve.Int64(),
			}
		}

		if project.Priority == nil {
			priority := 0
			project.Priority = &priority
		}
	}

	// get minted time
	if project.BlockNumberMinted == nil || project.MintedTime == nil {
		mintedTimeChan := make(chan structure.NftMintedTimeChan, 1)
		go func(mintedTimeChan chan structure.NftMintedTimeChan) {
			var mintedTime *structure.NftMintedTime
			var err error
			defer func() {
				mintedTimeChan <- structure.NftMintedTimeChan{
					NftMintedTime: mintedTime,
					Err:           err,
				}
			}()

			mintedTime, err = u.GetNftMintedTime(structure.GetNftMintedTimeReq{
				ContractAddress: project.ContractAddress,
				TokenID:         project.TokenID,
			})
		}(mintedTimeChan)
		mintedTimeFChan := <-mintedTimeChan
		if mintedTimeFChan.Err != nil {
			logger.AtLog.Logger.Error("mintedTimeFChan.Err ", zap.Error(mintedTimeFChan.Err))
		} else {
			project.BlockNumberMinted = mintedTimeFChan.NftMintedTime.BlockNumberMinted
			project.MintedTime = mintedTimeFChan.NftMintedTime.MintedTime
		}
	}

	project.TokenIDInt = int64(tokenIDInt)

	if usrFromChan.Err != nil {
		logger.AtLog.Logger.Error("usrFromChan.Err", zap.Error(usrFromChan.Err))
	} else {
		project.CreatorProfile = *usrFromChan.Data
	}

	if projectStatFChan.Err != nil {
		logger.AtLog.Logger.Error("projectStatFChan.Err", zap.Error(projectStatFChan.Err))
	} else {
		project.Stats = *projectStatFChan.Data
		project.TraitsStat = projectStatFChan.DataTrait
	}

	logger.AtLog.Logger.Info("project", zap.Any("project",project))
	_, err = u.Repo.UpdateProject(project.UUID, project)
	if err != nil {
		logger.AtLog.Logger.Error(" u.UpdateProject", zap.Error(err))
		return nil, err
	}
	logger.AtLog.Logger.Info("projectUUID", zap.Any("project.UUID", project.UUID))
	return project, nil
}

func (u Usecase) GetProjectsFromChain() error {
	contractAddress := os.Getenv("GENERATIVE_PROJECT")
	mProjects, err := u.MoralisNft.GetNftByContract(contractAddress, nfts.MoralisFilter{})
	if err != nil {
		logger.AtLog.Logger.Error("u.MoralisNft.GetNftByContract", zap.Error(err))
		return err
	}

	logger.AtLog.Logger.Info("contractAddress", zap.Any("contractAddress", contractAddress))

	for _, mProject := range mProjects.Result {
		_, err := u.UpdateProjectFromChain(contractAddress, mProject.TokenID)
		if err != nil {
			logger.AtLog.Logger.Error("u.Repo.FindProjectBy", zap.Error(err))
			return err
		}
		//resp = append(resp, *p)
		//logger.AtLog.Logger.Info("p", zap.Any("*p", *p))
	}

	return nil
}
