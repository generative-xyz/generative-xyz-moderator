package usecase

import (
	"context"
	"fmt"
	"math/big"
	"rederinghub.io/internal/delivery/http/response"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"
	"rederinghub.io/internal/delivery/http/request"
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

	u.UpdateProjectFromChain(contractAddr, tokenIDStr, txnHash)
}

func (u Usecase) UpdateProjectFromChain(contractAddr string, tokenIDStr string, txnHash string) (*entity.Projects, error) {
	var err error
	project := &entity.Projects{}

	defer func() {
		if err != nil {
			logger.AtLog.Logger.Error("UpdateProjectFromChain", zap.String("contractAddr", contractAddr), zap.String("txnHash", txnHash), zap.Error(err))
		} else {
			logger.AtLog.Logger.Info("UpdateProjectFromChain", zap.String("contractAddr", contractAddr), zap.String("txnHash", txnHash), zap.String("projectID", project.TokenID))
		}
	}()

	project, err = u.Repo.FindProjectByTxHash(txnHash)
	if err != nil {
		return nil, err
	}

	if project.TokenID == "" {
		tokenIDInt, err := strconv.Atoi(tokenIDStr)
		if err != nil {
			return nil, err
		}

		project.TokenID = tokenIDStr
		project.GenNFTAddr = tokenIDStr
		project.TokenIDInt = int64(tokenIDInt)
	}

	projectDetail, err := u.getProjectDetailFromChainWithoutCache(structure.GetProjectDetailMessageReq{
		ContractAddress: contractAddr,
		ProjectID:       tokenIDStr,
	})

	if err != nil {
		return nil, err
	}

	project.IsSynced = true
	project.Name = projectDetail.ProjectDetail.Name
	project.CreatorName = projectDetail.ProjectDetail.Creator
	project.CreatorAddrr = strings.ToLower(projectDetail.ProjectDetail.CreatorAddr.String())
	project.Description = projectDetail.ProjectDetail.Desc
	project.Scripts = projectDetail.ProjectDetail.Scripts
	project.ThirdPartyScripts = projectDetail.ProjectDetail.ScriptType
	project.Styles = projectDetail.ProjectDetail.Styles
	project.GenNFTAddr = strings.ToLower(projectDetail.ProjectDetail.GenNFTAddr.String())
	//project.MintPrice = projectDetail.ProjectDetail.MintPrice.String()
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

	// check is full chain
	tokenUri := response.TokenURIResp{}
	err = helpers.Base64DecodeRaw(project.NftTokenUri, &tokenUri)
	if err != nil {
		return nil, err
	}
	if len(tokenUri.AnimationURL) > 0 {
		maxSize := helpers.CalcOrigBinaryLength(tokenUri.AnimationURL)
		project.MaxFileSize = int64(maxSize)
		project.NetworkFee = big.NewInt(u.networkFeeBySize(int64(maxSize / 4))).String()
		htmlContent, err := helpers.Base64Decode(strings.ReplaceAll(tokenUri.AnimationURL, "data:text/html;base64,", ""))
		if err == nil {
			isFullChain, err := helpers.IsFullChain(string(htmlContent))
			if err == nil {
				project.IsFullChain = isFullChain
				logger.AtLog.Logger.Info("UpdateProjectFromChain", zap.Any("isFullChain", zap.Any("isFullChain)", isFullChain)))
			} else {
				logger.AtLog.Error("UpdateProjectFromChain", zap.Any("isFullChain", err))
			}
		} else {
			logger.AtLog.Error("UpdateProjectFromChain", zap.Any("isFullChain", err))
		}
	}

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

	projectStat, traitStat, err := u.GetUpdatedProjectStats(structure.GetProjectReq{
		ContractAddr: contractAddr,
		TokenID:      tokenIDStr,
	})

	if err != nil {
		return nil, err
	}

	project.Stats = *projectStat
	project.TraitsStat = traitStat

	user, err := u.GetUserProfileByWalletAddress(strings.ToLower(project.CreatorAddrr))
	if err != nil {
		return nil, err
	}

	project.CreatorProfile = *user

	_, err = u.Repo.UpdateProject(project.UUID, project)
	if err != nil {
		return nil, err
	}

	//DAO
	ids, err := u.CreateDAOProject(context.TODO(), &request.CreateDaoProjectRequest{
		ProjectIds: []string{project.ID.Hex()},
		CreatedBy:  project.CreatorAddrr,
	})
	if err != nil {
		logger.AtLog.Logger.Error("CreateDAOProject failed", zap.Error(err))
	} else {
		logger.AtLog.Logger.Info("CreateDAOProject success",
			zap.String("project_id", project.ID.Hex()),
			zap.Strings("ids", ids),
		)
	}

	return project, nil
}

// func (u Usecase) GetProjectsFromChain() error {
// contractAddress := os.Getenv("GENERATIVE_PROJECT")
// 	mProjects, err := u.MoralisNft.GetNftByContract(contractAddress, nfts.MoralisFilter{})
// 	if err != nil {
// 		u.Logger.Error("u.MoralisNft.GetNftByContract", err.Error(), err)
// 		return err
// 	}

// 	u.Logger.Info("contractAddress", contractAddress)

// 	for _, mProject := range mProjects.Result {
// 		_, err := u.UpdateProjectFromChain(contractAddress, mProject.TokenID)
// 		if err != nil {
// 			u.Logger.Error("u.Repo.FindProjectBy", err.Error(), err)
// 			return err
// 		}
// 		//resp = append(resp, *p)
// 		//u.Logger.Info("p", *p)
// 	}

// 	return nil
// }
